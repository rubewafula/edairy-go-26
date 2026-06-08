package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberPayrollService struct {
	notificationService *UINotificationService
}

func NewMemberPayrollService() *MemberPayrollService {
	return &MemberPayrollService{
		notificationService: NewUINotificationService(),
	}
}

func (s *MemberPayrollService) Create(req dtos.CreateMemberPayrollRequest, userID uint64) (*models.MemberPayroll, error) {
	// 1. Verify Date Range
	var pdr models.MemberPayDateRange
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", req.PayDateRangeID).First(&pdr).Error; err != nil {
		return nil, fmt.Errorf("invalid pay date range: %w", err)
	}

	// Check if payroll already exists for this period
	var existing models.MemberPayroll
	err := db.DB.Where("pay_date_range_id = ? AND deleted_at IS NULL", req.PayDateRangeID).First(&existing).Error
	allowedStatuses := map[string]bool{
		"draft":      true,
		"incomplete": true,
		"DRAFT":      true,
		"INCOMPLTE":  true,
	}

	if err == nil && !allowedStatuses[existing.Status] {
		return nil, fmt.Errorf("payroll for this period is already %s and cannot be regenerated", existing.Status)
	}

	dateOpened := utils.ParseDate(req.DateOpened)
	var payroll models.MemberPayroll

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Clean up existing draft if it exists
		if existing.ID != 0 {
			if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.MemberPayslip{}).Error; err != nil {
				return err
			}
			if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.MemberPayrollDeduction{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		}

		payroll = models.MemberPayroll{
			BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			DateOpened:     &dateOpened,
			Description:    req.Description,
			Status:         "processing",
			PayDateRangeID: &req.PayDateRangeID,
			FiscalPeriod:   req.FiscalPeriod,
		}
		return tx.Create(&payroll).Error
	})

	if err != nil {
		return nil, err
	}

	// Defer heavy processing to background
	go s.generatePayrollInBackground(payroll.ID, pdr, req, userID)

	return &payroll, nil
}

func (s *MemberPayrollService) generatePayrollInBackground(payrollID uint64, pdr models.MemberPayDateRange, req dtos.CreateMemberPayrollRequest, userID uint64) {
	type milkSum struct {
		MemberID uint64
		Kilos    float64
	}

	// 1. Fetch Aggregated Milk Collections
	var milkCollections []milkSum
	if err := db.DB.Table("milk_journal_entries").
		Select("member_id, SUM(quantity) as kilos").
		Joins("JOIN milk_journals ON milk_journal_entries.milk_journal_id = milk_journals.id").
		Where("milk_journals.confirmed = ? AND milk_journals.journal_date BETWEEN ? AND ?", true, pdr.StartDate, pdr.EndDate).
		Group("member_id").Scan(&milkCollections).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Member Payroll Generation Failed", "Failed to fetch aggregated milk collections", err, "incomplete")
		return
	}

	if len(milkCollections) == 0 {
		db.DB.Model(&models.MemberPayroll{}).Where("id = ?", payrollID).Update("status", "draft")
		return
	}

	// 2. Fetch Aggregated Rejects
	var milkRejects []milkSum
	if err := db.DB.Table("milk_rejects").
		Select("member_id, SUM(quantity) as kilos").
		Where("transaction_date BETWEEN ? AND ?", pdr.StartDate, pdr.EndDate).
		Group("member_id").Scan(&milkRejects).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Member Payroll Generation Failed", "Failed to fetch aggregated milk rejects", err, "incomplete")
		return
	}

	rejectMap := make(map[uint64]float64)
	for _, r := range milkRejects {
		rejectMap[r.MemberID] = r.Kilos
	}

	// 3. Pre-fetch Rate Resolution Maps
	var specialRates []models.MilkSpecialRate
	if err := db.DB.Where("pay_date_range_id = ? AND deleted_at IS NULL", req.PayDateRangeID).Find(&specialRates).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Member Payroll Generation Failed", "Failed to fetch milk special rates", err, "incomplete")
		return
	}
	memberRateMap := make(map[uint64]float64)
	routePeriodRateMap := make(map[uint64]float64)
	for _, r := range specialRates {
		if r.MemberID != nil {
			memberRateMap[*r.MemberID] = r.Rate
		} else if r.RouteID != nil {
			routePeriodRateMap[*r.RouteID] = r.Rate
		}
	}

	var defaultRates []models.DefaultMilkRate
	if err := db.DB.Where("deleted_at IS NULL").Find(&defaultRates).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Member Payroll Generation Failed", "Failed to fetch default milk rates", err, "incomplete")
		return
	}
	defaultRouteRateMap := make(map[uint64]float64)
	defaultMemberRateMap := make(map[uint64]float64)
	var globalDefault float64
	for _, r := range defaultRates {
		if r.MemberID != nil {
			defaultMemberRateMap[*r.MemberID] = r.Rate
		} else if r.RouteID != nil {
			defaultRouteRateMap[*r.RouteID] = r.Rate
		} else {
			globalDefault = r.Rate
		}
	}

	// 4. Pre-fetch Member metadata and Deductions in bulk
	memberIDs := make([]uint64, len(milkCollections))
	for i, c := range milkCollections {
		memberIDs[i] = c.MemberID
	}

	var members []models.Member
	if err := db.DB.Select("id, route_id").Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Member Payroll Generation Failed", "Failed to fetch member metadata", err, "incomplete")
		return
	}
	memberRouteMap := make(map[uint64]uint64)
	for _, m := range members {
		memberRouteMap[m.ID] = m.RouteID
	}

	// 5. Worker Pool Setup
	dateOpened := utils.ParseDate(req.DateOpened)
	numWorkers := runtime.NumCPU() * 2
	if sqlDB, err := db.DB.DB(); err == nil {
		stats := sqlDB.Stats()
		if stats.MaxOpenConnections > 0 && numWorkers > stats.MaxOpenConnections/2 {
			numWorkers = stats.MaxOpenConnections / 2
		}
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	var wg sync.WaitGroup
	jobs := make(chan milkSum, len(milkCollections))
	var mu sync.Mutex
	var totalGross, totalNet, totalDeductions float64
	var failedCount int64

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for collection := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[MemberPayrollService] Worker panicked during payroll generation for member %d: %v", collection.MemberID, r)
							db.DB.Create(&models.MemberPayrollGenerationError{ // Corrected PayrollID type
								BaseModel: models.BaseModel{CreatedBy: userID}, // Corrected MemberID type
								MemberID:  collection.MemberID,
								PayrollID: payrollID,
								Error:     fmt.Sprintf("Panic during generation: %v", r),
							})
							db.DB.Model(&models.MemberPayslip{}).Where("member_id = ?", collection.MemberID).Update("status", "incomplete")
							mu.Lock()
							failedCount++
							mu.Unlock()
						}
					}()
					var mGross, mNet, mDeductions float64
					var err error

					// Retry logic: try up to 3 times per member
					for attempt := 1; attempt <= 3; attempt++ {
						err = db.DB.Transaction(func(tx *gorm.DB) error {
							// Idempotency: Delete existing payslip and its deductions for this member and payroll
							// to ensure a clean re-creation on re-run.
							if err := tx.Where("member_id = ? AND payroll_id = ?", collection.MemberID, payrollID).Delete(&models.MemberPayslip{}).Error; err != nil {
								return fmt.Errorf("failed to delete existing payslip for member %d, payroll %d: %w", collection.MemberID, payrollID, err)
							}
							// Also delete associated deductions. This ensures that if a previous attempt created deductions
							// but failed before the payslip was finalized, or if we are updating an existing payslip,
							// old deductions are cleared and re-created based on current state.
							// Note: The previous code already had this deduction deletion, but it's good to keep it
							// explicitly here for idempotency of the entire payslip-deduction unit.
							if err := tx.Where("member_id = ? AND payroll_id = ?", collection.MemberID, payrollID).Delete(&models.MemberPayrollDeduction{}).Error; err != nil {
								return fmt.Errorf("failed to delete existing deductions: %w", err)
							}

							// Fetch current recurrent deductions for this member with a row-level lock
							var memberRecurrentDeductions []models.RecurrentDeduction
							if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
								Where("customer_id = ? AND settled = 0 AND customer_type = 'member'", collection.MemberID).
								Order("created_at ASC").Find(&memberRecurrentDeductions).Error; err != nil {
								return fmt.Errorf("failed to fetch recurrent deductions with lock: %w", err)
							}
							mGross, mNet, mDeductions = 0, 0, 0
							rejectKilos := rejectMap[collection.MemberID]
							netKilos := collection.Kilos - rejectKilos
							if netKilos < 0 {
								netKilos = 0
							}
							rateUsed := "SPECIAL"
							routeID := memberRouteMap[collection.MemberID]
							var rate float64
							if r, ok := memberRateMap[collection.MemberID]; ok {
								rate = r
							} else if r, ok := defaultMemberRateMap[collection.MemberID]; ok {
								rate = r
								rateUsed = "DEFAULT"
							} else if r, ok := routePeriodRateMap[routeID]; ok {
								rate = r
							} else if r, ok := defaultRouteRateMap[routeID]; ok {
								rateUsed = "DEFAULT"
								rate = r
							} else {
								rate = globalDefault
								rateUsed = "DEFAULT"
							}

							mGross = netKilos * rate

							var payslipDeductions float64                  // Initialize payslipDeductions here
							for _, rd := range memberRecurrentDeductions { // Corrected from memberDeductions
								remaining := rd.TotalAmount - rd.PaidAmount
								deductAmount := rd.RecurrentAmount
								if deductAmount > remaining {
									deductAmount = remaining
								}
								if payslipDeductions+deductAmount > mGross {
									deductAmount = mGross - payslipDeductions
								}
								if deductAmount <= 0 {
									continue
								}

								mpd := models.MemberPayrollDeduction{
									MemberID:        collection.MemberID,
									PayrollID:       payrollID,
									DeductionTypeID: rd.DeductionTypeID,
									Amount:          fmt.Sprintf("%.2f", deductAmount),
									TransactionDate: &dateOpened,
									DateCaptured:    utils.NowPtr(),
									Reference:       rd.Reference,
								}
								if err := tx.Create(&mpd).Error; err != nil { //
									return fmt.Errorf("failed to create member payroll deduction: %w", err)
								}

								payslipDeductions += deductAmount
							}
							// 4. Create Payslip (after deleting any existing one)
							mDeductions = payslipDeductions
							mNet = mGross - mDeductions
							payslip := models.MemberPayslip{
								MemberID:        collection.MemberID,
								PayrollID:       payrollID,
								DateOpened:      &dateOpened,
								Status:          "draft",
								GrossKilos:      collection.Kilos,
								RejectKilos:     rejectKilos,
								NetKilos:        netKilos,
								GrossPay:        mGross,
								TotalDeductions: mDeductions,
								RateUsed:        rateUsed,
								NetPay:          mNet,
								FiscalPeriod:    req.FiscalPeriod,
								PayDateRangeID:  &req.PayDateRangeID,
							}

							// Always create a new payslip, as any existing one for this member/payroll was deleted.
							return tx.Create(&payslip).Error
						})

						if err == nil {
							break
						}
						time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
					}

					if err != nil {
						mu.Lock()
						failedCount++
						mu.Unlock()
						log.Printf("[MemberPayrollService] Failed to generate payroll for member %d: %v", collection.MemberID, err)
						db.DB.Create(&models.MemberPayrollGenerationError{
							BaseModel: models.BaseModel{CreatedBy: userID},
							MemberID:  collection.MemberID,
							PayrollID: payrollID, // Corrected PayrollID type
							Error:     err.Error(),
						})
					} else {
						mu.Lock()
						totalGross += mGross
						totalNet += mNet
						totalDeductions += mDeductions
						mu.Unlock()
					}
				}()
			}
		}()
	}

	for _, c := range milkCollections {
		jobs <- c
	}
	close(jobs)
	wg.Wait()

	// 6. Finalize Payroll Header
	status := "draft"
	if failedCount > 0 {
		status = "incomplete"
	}

	db.DB.Model(&models.MemberPayroll{}).Where("id = ?", payrollID).Updates(map[string]interface{}{
		"gross_pay":        totalGross,
		"net_pay":          totalNet,
		"total_deductions": totalDeductions,
		"status":           status,
		"updated_by":       userID,
		"updated_at":       utils.Now(),
	})

	// Send UI notification
	message := fmt.Sprintf("Member payroll generation completed. Success: %d, Failed: %d out of %d records.", len(milkCollections)-int(failedCount), failedCount, len(milkCollections))
	notificationType := "SUCCESS"
	errorLink := ""

	if failedCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/member-payrolls/generation-errors/%d", payrollID)
		log.Printf("[MemberPayrollService.generatePayrollInBackground] Member payroll generation completed with %d failures for payroll %d.", failedCount, payrollID)
	} else if len(milkCollections) == 0 {
		message = "Member payroll generation completed. No milk collections were processed."
		notificationType = "INFO"
		log.Printf("[MemberPayrollService.generatePayrollInBackground] Member payroll generation completed. No milk collections were processed for payroll %d.", payrollID)
	} else {
		log.Printf("[MemberPayrollService.generatePayrollInBackground] Member payroll generation completed successfully for payroll %d.", payrollID)
	}

	_, err := s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Member Payroll Generation Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("MEMBER_PAYROLL"),
	})
	if err != nil {
		log.Printf("[MemberPayrollService.generatePayrollInBackground] Failed to create UI notification: %v", err)
	}
}

func (s *MemberPayrollService) getRawPayroll(id string) (*models.MemberPayroll, error) {
	var payroll models.MemberPayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (s *MemberPayrollService) List(page, limit int) ([]dtos.MemberPayrollResponse, int64, error) {
	var results []dtos.MemberPayrollResponse
	var total int64
	db.DB.Model(&models.MemberPayroll{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			mp.*, 
			pdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM member_payrolls mp
		LEFT JOIN customer_pay_date_ranges pdr ON mp.pay_date_range_id = pdr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.deleted_at IS NULL
		ORDER BY mp.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *MemberPayrollService) Get(id string) (*dtos.MemberPayrollResponse, error) {
	var result dtos.MemberPayrollResponse
	query := `
		SELECT 
			mp.*, 
			pdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM member_payrolls mp
		LEFT JOIN customer_pay_date_ranges pdr ON mp.pay_date_range_id = pdr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.id = ? AND mp.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

// GetPayslips retrieves a list of member payslips with pagination and filtering.
func (s *MemberPayrollService) GetPayslips(payrollID string, memberID string, routeID string, page, limit int) ([]dtos.MemberPayslipResponse, int64, error) {
	var results []dtos.MemberPayslipResponse
	var total int64

	queryBuilder := db.DB.Model(&models.MemberPayslip{})

	if payrollID != "" {
		queryBuilder = queryBuilder.Where("payroll_id = ?", payrollID)
	}
	if memberID != "" {
		queryBuilder = queryBuilder.Where("member_id = ?", memberID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT
			mp.id, mp.member_id,
			mr.member_no,
			CONCAT(mr.first_name, ' ', mr.last_name) as member_name,
			mp.payroll_id, mp.pay_date_range_id, mp.date_opened, mp.description, mp.status, mp.rate_used,
			mp.posted_at, u_post.name as posted_by_name,
			mp.confirmed_at, u_conf.name as confirmed_by_name,
			mp.approved_at, u_appr.name as approved_by_name,
			mp.gross_kilos, mp.reject_kilos, mp.net_kilos, mp.gross_pay, mp.total_deductions, mp.net_pay, mp.fiscal_period,
			mp.created_at, mp.updated_at
		FROM member_payslips mp
		LEFT JOIN member_registrations mr ON mp.member_id = mr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.deleted_at IS NULL
	`
	var args []interface{}
	if payrollID != "" {
		query += " AND mp.payroll_id = ?"
		args = append(args, payrollID)
	}
	if memberID != "" {
		query += " AND mp.member_id = ?"
		args = append(args, memberID)
	}
	if routeID != "" {
		query += " AND mr.route_id = ?"
		args = append(args, routeID)
	}
	query += " ORDER BY mp.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := db.DB.Raw(query, args...).Scan(&results).Error
	return results, total, err
}

// GetPayslip retrieves a single member payslip by ID.
func (s *MemberPayrollService) GetPayslip(id string) (*dtos.MemberPayslipResponse, error) {
	var result dtos.MemberPayslipResponse
	query := `
		SELECT
			mp.id, mp.member_id, mr.member_no, CONCAT(mr.first_name, ' ', mr.last_name) as member_name,
			mp.payroll_id, mp.pay_date_range_id, mp.date_opened, mp.description, mp.status, mp.rate_used,
			mp.posted_at, u_post.name as posted_by_name, mp.confirmed_at, u_conf.name as confirmed_by_name,
			mp.approved_at, u_appr.name as approved_by_name, mp.gross_kilos, mp.reject_kilos, mp.net_kilos,
			mp.gross_pay, mp.total_deductions, mp.net_pay, mp.fiscal_period, mp.created_at, mp.updated_at
		FROM member_payslips mp
		LEFT JOIN member_registrations mr ON mp.member_id = mr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.id = ? AND mp.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

// ExportPayslips initiates a background process to export member payslips.
func (s *MemberPayrollService) ExportPayslips(userID uint64, payrollID, memberID, routeID, reportType string) error {
	go s.processPayslipExportInBackground(userID, payrollID, memberID, routeID, reportType)
	return nil
}

func (s *MemberPayrollService) processPayslipExportInBackground(userID uint64, payrollID, memberID, routeID, reportType string) {
	var results []dtos.MemberPayslipResponse

	baseQuery := `
		SELECT
			mp.id, mp.member_id,
			mr.member_no,
			CONCAT(mr.first_name, ' ', mr.last_name) as member_name,
			mp.payroll_id, mp.pay_date_range_id, mp.date_opened, mp.description, mp.status, mp.rate_used,
			mp.posted_at, u_post.name as posted_by_name,
			mp.confirmed_at, u_conf.name as confirmed_by_name,
			mp.approved_at, u_appr.name as approved_by_name,
			mp.gross_kilos, mp.reject_kilos, mp.net_kilos, mp.gross_pay, mp.total_deductions, mp.net_pay, mp.fiscal_period,
			mp.created_at, mp.updated_at
		FROM member_payslips mp
		LEFT JOIN member_registrations mr ON mp.member_id = mr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.deleted_at IS NULL
	`

	var args []interface{}
	if payrollID != "" {
		baseQuery += " AND mp.payroll_id = ?"
		args = append(args, payrollID)
	}
	if memberID != "" {
		baseQuery += " AND mp.member_id = ?"
		args = append(args, memberID)
	}
	if routeID != "" {
		baseQuery += " AND mr.route_id = ?"
		args = append(args, routeID)
	}

	err := db.DB.Raw(baseQuery+" ORDER BY mp.id DESC", args...).Scan(&results).Error
	if err != nil {
		log.Printf("[MemberPayrollService.processPayslipExportInBackground] Error: %v", err)
		return
	}

	var fileData []byte
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generatePayslipPDF(results)
	} else {
		fileData, err = s.generatePayslipCSV(results)
	}

	if err != nil {
		log.Printf("[MemberPayrollService.processPayslipExportInBackground] Generation error: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("payslips_export_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[MemberPayrollService.processPayslipExportInBackground] File error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Payslip Export Complete",
		Message:          fmt.Sprintf("Your payslip data export (%s) is ready for download.", strings.ToUpper(ext)),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/member-payslips/export/download/%s", filename),
	})
}

// ExportPayslipStatements initiates a background process to export detailed payslip statements.
func (s *MemberPayrollService) ExportPayslipStatements(userID uint64, payslipID, memberID, reportType string) error {
	go s.processPayslipStatementsInBackground(userID, payslipID, memberID, reportType)
	return nil
}

func (s *MemberPayrollService) processPayslipStatementsInBackground(userID uint64, payslipID, memberID, reportType string) {
	var payslips []dtos.MemberPayslipResponse

	// 1. Fetch Payslips using filters
	baseQuery := `
		SELECT
			mp.id, mp.member_id, mr.member_no, CONCAT(mr.first_name, ' ', mr.last_name) as member_name,
			mp.payroll_id, mp.pay_date_range_id, mp.date_opened, mp.description, mp.status, mp.rate_used,
			mp.gross_kilos, mp.reject_kilos, mp.net_kilos, mp.gross_pay, mp.total_deductions, mp.net_pay, mp.fiscal_period,
			mp.created_at, mp.updated_at
		FROM member_payslips mp
		LEFT JOIN member_registrations mr ON mp.member_id = mr.id
		WHERE mp.deleted_at IS NULL AND mp.id = ? AND mp.member_id = ?
	`
	if err := db.DB.Raw(baseQuery, payslipID, memberID).Scan(&payslips).Error; err != nil {
		log.Printf("[MemberPayrollService.processPayslipStatementsInBackground] Query error: %v", err)
		return
	}

	if len(payslips) == 0 {
		return
	}

	// 2. Fetch Deductions for this specific payslip
	var deductions []dtos.MemberPayrollDeductionResponse
	p := payslips[0]
	db.DB.Table("member_payroll_deductions mpd").
		Select("mpd.*, dt.description as deduction_type_name").
		Joins("LEFT JOIN deduction_types dt ON mpd.deduction_type_id = dt.id").
		Where("mpd.payroll_id = ? AND mpd.member_id = ? AND mpd.deleted_at IS NULL", p.PayrollID, p.MemberID).
		Scan(&deductions)

	// Map deductions to payslip IDs
	deductMap := make(map[uint64][]dtos.MemberPayrollDeductionResponse)
	for _, d := range deductions {
		deductMap[d.PayrollID] = append(deductMap[d.PayrollID], d)
	}

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateDetailedPayslipPDF(payslips, deductMap)
	} else {
		fileData, err = s.generateDetailedPayslipCSV(payslips, deductMap)
	}

	if err != nil {
		log.Printf("[MemberPayrollService] Generation error: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("payslip_statements_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[MemberPayrollService] File write error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Payslip Statements Ready",
		Message:          fmt.Sprintf("Your detailed payslip statements (%s) are ready for download.", strings.ToUpper(ext)),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/member-payslips/export/download/%s", filename),
	})
}

func (s *MemberPayrollService) generateDetailedPayslipCSV(payslips []dtos.MemberPayslipResponse, deductMap map[uint64][]dtos.MemberPayrollDeductionResponse) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"Member No", "Name", "Period", "Category", "Description", "Amount"})

	for _, p := range payslips {
		// Earnings row
		writer.Write([]string{p.MemberNo, p.MemberName, p.FiscalPeriod, "EARNING", "Gross Milk Pay", fmt.Sprintf("%.2f", p.GrossPay)})

		// Deduction rows
		for _, d := range deductMap[p.PayrollID] {
			writer.Write([]string{p.MemberNo, p.MemberName, p.FiscalPeriod, "DEDUCTION", d.DeductionTypeName, d.Amount})
		}

		// Summary row
		writer.Write([]string{p.MemberNo, p.MemberName, p.FiscalPeriod, "SUMMARY", "NET PAY", fmt.Sprintf("%.2f", p.NetPay)})
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *MemberPayrollService) generateDetailedPayslipPDF(payslips []dtos.MemberPayslipResponse, deductMap map[uint64][]dtos.MemberPayrollDeductionResponse) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
	}
	db.DB.Table("organization_details").First(&org)

	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, p := range payslips {
		pdf.AddPage()
		// Header
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 10, "MEMBER PAYSLIP STATEMENT", "B", 1, "C", false, 0, "")
		pdf.Ln(5)

		// Member Details
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 7, "Member No:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 7, p.MemberNo, "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 7, "Member Name:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 7, p.MemberName, "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 7, "Fiscal Period:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 7, p.FiscalPeriod, "", 1, "L", false, 0, "")
		pdf.Ln(5)

		// Breakdown Table
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(140, 8, "Description", "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 8, "Amount", "1", 1, "C", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		// Earnings
		itemDesc := fmt.Sprintf("Gross Milk Pay (%.2f Kgs)", p.NetKilos)
		pdf.CellFormat(140, 8, itemDesc, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 8, fmt.Sprintf("%.2f", p.GrossPay), "1", 1, "R", false, 0, "")

		// Deductions
		for _, d := range deductMap[p.PayrollID] {
			pdf.CellFormat(140, 8, d.DeductionTypeName, "1", 0, "L", false, 0, "")
			pdf.CellFormat(50, 8, "-"+d.Amount, "1", 1, "R", false, 0, "")
		}

		// Footer
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(140, 10, "NET PAYABLE", "1", 0, "R", false, 0, "")
		pdf.CellFormat(50, 10, fmt.Sprintf("%.2f", p.NetPay), "1", 1, "R", false, 0, "")
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func (s *MemberPayrollService) generatePayslipCSV(results []dtos.MemberPayslipResponse) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Member No", "Name", "Period", "Status", "G-Kilos", "N-Kilos", "Gross Pay", "Deductions", "Net Pay"}
	writer.Write(headers)

	var totalGK, totalNK, totalGross, totalDeduct, totalNet float64 // Roll-up totals

	for _, ps := range results {
		totalGK += ps.GrossKilos
		totalNK += ps.NetKilos
		totalGross += ps.GrossPay
		totalDeduct += ps.TotalDeductions
		totalNet += ps.NetPay

		writer.Write([]string{
			ps.MemberNo,
			ps.MemberName,
			ps.FiscalPeriod,
			ps.Status,
			fmt.Sprintf("%.2f", ps.GrossKilos),
			fmt.Sprintf("%.2f", ps.NetKilos),
			fmt.Sprintf("%.2f", ps.GrossPay),
			fmt.Sprintf("%.2f", ps.TotalDeductions),
			fmt.Sprintf("%.2f", ps.NetPay),
		})
	}

	// Add totals row
	writer.Write([]string{
		"TOTALS", "", "", "",
		fmt.Sprintf("%.2f", totalGK),
		fmt.Sprintf("%.2f", totalNK),
		fmt.Sprintf("%.2f", totalGross),
		fmt.Sprintf("%.2f", totalDeduct),
		fmt.Sprintf("%.2f", totalNet),
	})

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *MemberPayrollService) generatePayslipPDF(results []dtos.MemberPayslipResponse) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
		Phone          string `gorm:"column:phone"`
		Email          string `gorm:"column:email"`
	}
	db.DB.Table("organization_details").First(&org)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Phone: %s | Email: %s", org.Phone, org.Email), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "MEMBER PAYSLIPS REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"M-No", "Name", "Period", "Status", "G-Kilos", "N-Kilos", "Gross", "Deductions", "Net"}
	widths := []float64{25, 65, 25, 25, 25, 25, 30, 30, 30}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	var totalGK, totalNK, totalGross, totalDeduct, totalNet float64 // Roll-up totals

	pdf.SetFont("Arial", "", 8)
	for _, ps := range results {
		totalGK += ps.GrossKilos
		totalNK += ps.NetKilos
		totalGross += ps.GrossPay
		totalDeduct += ps.TotalDeductions
		totalNet += ps.NetPay

		pdf.CellFormat(widths[0], 8, ps.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, ps.MemberName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, ps.FiscalPeriod, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, ps.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 8, fmt.Sprintf("%.2f", ps.GrossKilos), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[5], 8, fmt.Sprintf("%.2f", ps.NetKilos), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[6], 8, fmt.Sprintf("%.2f", ps.GrossPay), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[7], 8, fmt.Sprintf("%.2f", ps.TotalDeductions), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[8], 8, fmt.Sprintf("%.2f", ps.NetPay), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	// Add totals row
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(widths[0], 8, "TOTALS", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[1], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[2], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[3], 8, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(widths[4], 8, fmt.Sprintf("%.2f", totalGK), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[5], 8, fmt.Sprintf("%.2f", totalNK), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[6], 8, fmt.Sprintf("%.2f", totalGross), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[7], 8, fmt.Sprintf("%.2f", totalDeduct), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[8], 8, fmt.Sprintf("%.2f", totalNet), "1", 0, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func (s *MemberPayrollService) Approve(payrollID uint64, userID uint64) (*models.MemberPayroll, error) {
	var payroll models.MemberPayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validate state and Lock payroll record
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		// 2. Ensure payroll is in 'confirmed' status
		if payroll.Status != "confirmed" {
			return fmt.Errorf("payroll must be in 'confirmed' status to be approved, current status: %s", payroll.Status)
		}

		// 3. Validate all payslips are complete and no generation errors exist
		var errCount int64
		tx.Model(&models.MemberPayrollGenerationError{}).Where("payroll_id = ?", payrollID).Count(&errCount)
		if errCount > 0 {
			return fmt.Errorf("cannot approve payroll with %d generation errors", errCount)
		}

		var incompleteCount int64
		tx.Model(&models.MemberPayslip{}).Where("payroll_id = ? AND status = 'incomplete'", payrollID).Count(&incompleteCount)
		if incompleteCount > 0 {
			return fmt.Errorf("cannot approve payroll with %d incomplete payslips", incompleteCount)
		}

		// Update status to prevent re-triggering
		return tx.Model(&payroll).Update("status", "approving").Error
	})

	if err != nil {
		return nil, err
	}

	go s.approvePayrollInBackground(payrollID, userID)

	return &payroll, nil
}

func (s *MemberPayrollService) approvePayrollInBackground(payrollID uint64, userID uint64) {
	var payroll models.MemberPayroll
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
		log.Printf("[MemberPayrollService] Background approval failed to find payroll %d: %v", payrollID, err)
		return
	}

	now := utils.Now()

	// 2. Fetch Account Posting Rules
	var grossRule, loanRule, shareRule models.TransactionPostingRule
	if err := db.DB.Where("transaction_type = ?", "MILK_PAYMENT").First(&grossRule).Error; err != nil && err != gorm.ErrRecordNotFound {
		s.handleProcessingError(payrollID, userID, "Member Payroll Approval Failed", "Database error fetching MILK_PAYMENT rule", err, "confirmed")
		return
	}
	if err := db.DB.Where("transaction_type = ?", "MEMBER_REPAYMENT_DEDUCTION").First(&loanRule).Error; err != nil && err != gorm.ErrRecordNotFound {
		s.handleProcessingError(payrollID, userID, "Member Payroll Approval Failed", "Database error fetching MEMBER_REPAYMENT_DEDUCTION rule", err, "confirmed")
		return
	}
	if err := db.DB.Where("transaction_type = ?", "SHARES_CONTRIBUTION").First(&shareRule).Error; err != nil && err != gorm.ErrRecordNotFound {
		s.handleProcessingError(payrollID, userID, "Member Payroll Approval Failed", "Database error fetching SHARES_CONTRIBUTION rule", err, "confirmed")
		return
	}

	if grossRule.ID == 0 {
		log.Printf("[MemberPayrollService] Missing posting rule for MILK_PAYROLL_GROSS, aborting payroll %d", payrollID)
		db.DB.Model(&payroll).Update("status", "confirmed")
		return
	}

	if loanRule.ID == 0 {
		log.Printf("[MemberPayrollService] Missing posting rule for MEMBER_REPAYMENT_DEDUCTION, aborting payroll %d", payrollID)
		db.DB.Model(&payroll).Update("status", "confirmed")
		return
	}

	if shareRule.ID == 0 {
		log.Printf("[MemberPayrollService] Missing posting rule for SHARES_CONTRIBUTION, aborting payroll %d", payrollID)
		db.DB.Model(&payroll).Update("status", "confirmed")
		return
	}
	// 3. Load Payslips
	var payslips []models.MemberPayslip
	if err := db.DB.Where("payroll_id = ? AND status != 'approved'", payrollID).Find(&payslips).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Member Payroll Approval Failed", "Failed to fetch payslips from database", err, "confirmed")
		return
	}

	// 4. Worker Pool Setup
	numWorkers := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	jobs := make(chan models.MemberPayslip, len(payslips))
	var mu sync.Mutex
	var failedCount int64

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ps := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[MemberPayrollService] Worker panicked during payroll approval for payslip %d: %v", ps.ID, r)
							db.DB.Create(&models.MemberPayrollGenerationError{ // Corrected PayrollID type
								BaseModel: models.BaseModel{CreatedBy: userID}, // Corrected MemberID type
								MemberID:  ps.MemberID,
								PayrollID: payrollID,
								Error:     fmt.Sprintf("Panic during approval: %v", r),
							})
							// Mark the affected payslip as incomplete
							db.DB.Model(&models.MemberPayslip{}).Where("id = ?", ps.ID).Update("status", "incomplete")
							mu.Lock()
							failedCount++
							mu.Unlock()
						}
					}()
					var err error
					// Retry logic: try up to 3 times per payslip
					for attempt := 1; attempt <= 3; attempt++ {
						err = db.DB.Transaction(func(tx *gorm.DB) error {
							// Lock and check current payslip status
							var currentPS models.MemberPayslip
							if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&currentPS, ps.ID).Error; err != nil {
								return err
							}
							if currentPS.Status == "approved" {
								return nil // Already processed
							}

							// 1. Create Independent Transaction Header for THIS payslip
							transaction := models.Transaction{
								BaseModel:       models.BaseModel{CreatedBy: userID},
								Reference:       fmt.Sprintf("PSL-%d", ps.ID),
								TransactionName: fmt.Sprintf("Payslip Approval - %s (Member %d)", payroll.FiscalPeriod, ps.MemberID),
								TransactionType: "PAYROLL_PAYSLIP",
								TransactionDate: now,
								Description:     fmt.Sprintf("Payslip for Member %d in Payroll %s", ps.MemberID, payroll.FiscalPeriod),
								Status:          "approved",
							}
							if err := tx.Create(&transaction).Error; err != nil {
								return err
							}

							// DR Milk Purchase Expense / CR Member Payables
							glGross := []models.GeneralLedgerEntry{
								{
									BaseModel:       models.BaseModel{CreatedBy: userID},
									TransactionID:   transaction.ID,
									AccountID:       grossRule.DebitAccountID,
									Debit:           ps.GrossPay,
									TransactionDate: now,
									Description:     fmt.Sprintf("Gross Pay: Member %d", ps.MemberID),
								},
								{
									BaseModel:       models.BaseModel{CreatedBy: userID},
									TransactionID:   transaction.ID,
									AccountID:       grossRule.CreditAccountID,
									Credit:          ps.GrossPay,
									TransactionDate: now,
									Description:     fmt.Sprintf("Gross Pay Payable: Member %d", ps.MemberID),
								},
							}
							if err := tx.Create(&glGross).Error; err != nil {
								return err
							}

							// Process Individual Deductions
							var deductions []models.MemberPayrollDeduction
							if err := tx.Where("payroll_id = ? AND member_id = ?", payrollID, ps.MemberID).Find(&deductions).Error; err != nil {
								return err
							}

							for _, d := range deductions {
								var dType models.DeductionType
								if err := tx.First(&dType, d.DeductionTypeID).Error; err != nil {
									return err
								}

								rule := loanRule
								if dType.Code == "SHARE" {
									rule = shareRule
								}
								creditAcc := rule.CreditAccountID
								if creditAcc == 0 {
									return fmt.Errorf("no accounting mapping for deduction type %s", dType.Code)
								}

								// Parse deduction amount from string to float64 for accounting entries
								parsedDeductionAmount, err := strconv.ParseFloat(d.Amount, 64)
								if err != nil {
									parsedDeductionAmount = 0
								}

								// DR Member Payables / CR Loan/Share Receivable
								glDeduct := []models.GeneralLedgerEntry{
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       grossRule.CreditAccountID,
										Debit:           parsedDeductionAmount,
										TransactionDate: now,
										Description:     fmt.Sprintf("Deduction DR: Member %d - %s", ps.MemberID, d.Reference),
									},
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       creditAcc,
										Credit:          parsedDeductionAmount,
										TransactionDate: now,
										Description:     fmt.Sprintf("Deduction CR: Member %d - %s", ps.MemberID, d.Reference),
									},
								}
								if err := tx.Create(&glDeduct).Error; err != nil {
									return err
								}

								// Update Recurrent Deduction balances
								var rd models.RecurrentDeduction
								if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
									Where("customer_id = ? AND reference = ? AND customer_type = 'member'", d.MemberID, d.Reference).
									First(&rd).Error; err == nil {

									newPaid := rd.PaidAmount + parsedDeductionAmount
									settled := 0
									if newPaid >= rd.TotalAmount {
										settled = 1
									}
									if err := tx.Model(&rd).Updates(map[string]interface{}{
										"paid_amount": newPaid,
										"settled":     settled,
										"updated_by":  userID,
									}).Error; err != nil {
										return err
									}
								}
							}

							// 5. Finalize Payslip status
							return tx.Model(&currentPS).Updates(map[string]interface{}{
								"status":      "approved",
								"approved_at": &now,
								"approved_by": &userID,
								"posted_at":   &now,
								"posted_by":   &userID,
								"updated_by":  userID,
							}).Error
						})

						if err == nil {
							break
						}
						time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
					}

					if err != nil {
						mu.Lock()
						failedCount++
						mu.Unlock()
						log.Printf("[MemberPayrollService] Approval failed for payslip %d: %v", ps.ID, err)
						db.DB.Create(&models.MemberPayrollGenerationError{
							BaseModel: models.BaseModel{CreatedBy: userID},
							MemberID:  ps.MemberID,
							PayrollID: payrollID, // Corrected PayrollID type
							Error:     fmt.Sprintf("Approval Error: %s", err.Error()),
						})
						db.DB.Model(&models.MemberPayslip{}).Where("id = ?", ps.ID).Update("status", "incomplete")
					}
				}()
			}
		}()
	}

	for _, ps := range payslips {
		jobs <- ps
	}
	close(jobs)
	wg.Wait()

	// 5. Finalize Payroll Header
	finalStatus := "approved"
	if failedCount > 0 {
		finalStatus = "incomplete"
	}

	db.DB.Model(&payroll).Updates(map[string]interface{}{
		"status":      finalStatus,
		"approved_at": &now,
		"approved_by": &userID,
		"updated_by":  userID,
		"posted_at":   &now,
		"posted_by":   &userID,
		"updated_at":  utils.Now(),
	})

	// Send UI notification
	message := fmt.Sprintf("Member payroll approval completed. Success: %d, Failed: %d out of %d payslips.", len(payslips)-int(failedCount), failedCount, len(payslips))
	notificationType := "SUCCESS"
	errorLink := ""

	if failedCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/member-payrolls/approval-errors/%d", payrollID)
		log.Printf("[MemberPayrollService.approvePayrollInBackground] Member payroll approval completed with %d failures for payroll %d.", failedCount, payrollID)
	} else if len(payslips) == 0 {
		message = "Member payroll approval completed. No payslips were processed."
		notificationType = "INFO"
		log.Printf("[MemberPayrollService.approvePayrollInBackground] Member payroll approval completed. No payslips were processed for payroll %d.", payrollID)
	} else {
		log.Printf("[MemberPayrollService.approvePayrollInBackground] Member payroll approval completed successfully for payroll %d.", payrollID)
	}

	_, err := s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Member Payroll Approval Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("MEMBER_PAYROLL"),
	})
	if err != nil {
		log.Printf("[MemberPayrollService.approvePayrollInBackground] Failed to create UI notification: %v", err)
	}
}

func (s *MemberPayrollService) Confirm(payrollID uint64, userID uint64) (*models.MemberPayroll, error) {

	var payroll models.MemberPayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		if payroll.Status != "draft" {
			return fmt.Errorf("payroll must be in 'draft' status to be confirmed, current status: %s", payroll.Status)
		}

		now := utils.Now()
		// Update the main payroll record
		if err := tx.Model(&payroll).Updates(map[string]interface{}{
			"status":       "confirmed",
			"confirmed_at": &now,
			"confirmed_by": &userID,
			"updated_by":   userID,
			"updated_at":   now,
		}).Error; err != nil {
			return err
		}

		// Update all associated payslips to confirmed status
		if err := tx.Model(&models.MemberPayslip{}).Where("payroll_id = ?", payrollID).Update("status", "confirmed").Error; err != nil {
			return err
		}
		return nil
	})
	return &payroll, err
}

// GetGenerationErrors retrieves the list of errors encountered during a specific payroll generation.
func (s *MemberPayrollService) GetGenerationErrors(payrollID uint64) ([]models.MemberPayrollGenerationError, error) {
	var errors []models.MemberPayrollGenerationError
	err := db.DB.Where("payroll_id = ?", payrollID).Order("id DESC").Find(&errors).Error
	return errors, err
}

// handleProcessingError logs a database error, updates the payroll status, and sends a UI notification.
func (s *MemberPayrollService) handleProcessingError(payrollID uint64, userID uint64, title, message string, err error, status string) {
	log.Printf("[%s] %s for payroll %d: %v", title, message, payrollID, err)
	db.DB.Model(&models.MemberPayroll{}).Where("id = ?", payrollID).Update("status", status)
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            title,
		Message:          fmt.Sprintf("%s: %v", message, err),
		NotificationType: "ERROR",
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("MEMBER_PAYROLL"),
	})
}

// GetApprovalErrors retrieves the list of errors encountered during a specific payroll approval.
func (s *MemberPayrollService) GetApprovalErrors(payrollID uint64) ([]models.MemberPayrollGenerationError, error) {
	var errors []models.MemberPayrollGenerationError
	err := db.DB.Where("payroll_id = ?", payrollID).Order("id DESC").Find(&errors).Error
	return errors, err
}
