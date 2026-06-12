package services

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransporterPayrollService struct {
	notificationService *UINotificationService
}

func NewTransporterPayrollService() *TransporterPayrollService {
	return &TransporterPayrollService{
		notificationService: NewUINotificationService(),
	}
}

func (s *TransporterPayrollService) Create(req dtos.CreateTransporterPayrollRequest, userID uint64) (*models.TransporterPayroll, error) {
	var pdr models.MemberPayDateRange
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", req.PayDateRangeID).First(&pdr).Error; err != nil {
		return nil, fmt.Errorf("invalid pay date range: %w", err)
	}

	var existing models.TransporterPayroll
	err := db.DB.Where("pay_date_range_id = ? AND deleted_at IS NULL", req.PayDateRangeID).First(&existing).Error
	if err == nil && existing.Status != "draft" {
		return nil, fmt.Errorf("payroll for this period is already %s and cannot be regenerated", existing.Status)
	}

	dateOpened := utils.ParseDate(req.DateOpened)
	var payroll models.TransporterPayroll

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if existing.ID != 0 {
			if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.TransporterPayslip{}).Error; err != nil {
				return err
			}
			if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.TransporterPayrollBenefit{}).Error; err != nil {
				return err
			}
			// Assuming TransporterPayrollDeduction exists and needs to be cleared
			// if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.TransporterPayrollDeduction{}).Error; err != nil {
			// 	return err
			// }
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		}

		payroll = models.TransporterPayroll{
			BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			PayDateRangeID: &req.PayDateRangeID,
			DateOpened:     &dateOpened,
			Description:    req.Description,
			FiscalPeriod:   req.FiscalPeriod,
			Status:         "processing",
		}
		return tx.Create(&payroll).Error
	})

	if err != nil {
		return nil, err
	}

	go s.generatePayrollInBackground(payroll.ID, pdr, req, userID)

	return &payroll, nil
}

func (s *TransporterPayrollService) generatePayrollInBackground(payrollID uint64, pdr models.MemberPayDateRange, req dtos.CreateTransporterPayrollRequest, userID uint64) {
	type transporterMilkSummary struct {
		TransporterID uint64
		TotalQuantity float64
		TotalRejects  float64
	}

	var transporterSummaries []transporterMilkSummary
	if err := db.DB.Table("milk_transporter_cost").
		Select("transporter_id, SUM(quantity) as total_quantity, SUM(rejects) as total_rejects").
		Where("pay_date_range_id = ? ", pdr.ID).
		Group("transporter_id").Scan(&transporterSummaries).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Generation Failed", "Failed to fetch aggregated milk summaries", err, "incomplete")
		return
	}

	if len(transporterSummaries) == 0 {
		db.DB.Model(&models.TransporterPayroll{}).Where("id = ?", payrollID).Update("status", "draft")
		return
	}

	transporterIDs := make([]uint64, len(transporterSummaries))
	for i, ts := range transporterSummaries {
		transporterIDs[i] = ts.TransporterID
	}

	// Pre-fetch Transporter details (especially RouteID for rate calculation)
	var transporters []models.TransporterRouteAssignment
	if err := db.DB.Select("id, route_id").Where("transporter_id IN ?", transporterIDs).Find(&transporters).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Generation Failed", "Failed to fetch transporter assignments", err, "incomplete")
		return
	}

	transporterRouteMap := make(map[uint64]uint64)
	for _, t := range transporters {
		transporterRouteMap[t.ID] = t.RouteID
	}

	// Pre-fetch Recurrent Deductions for transporters
	var recurrentDeductions []models.RecurrentDeduction
	if err := db.DB.Where("customer_id IN ? AND settled = 0 AND customer_type = 'transporter'", transporterIDs).
		Order("created_at ASC").Find(&recurrentDeductions).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Generation Failed", "Failed to fetch recurrent deductions", err, "incomplete")
		return
	}
	deductionMap := make(map[uint64][]models.RecurrentDeduction)
	for _, rd := range recurrentDeductions {
		deductionMap[rd.CustomerID] = append(deductionMap[rd.CustomerID], rd)
	}

	// Pre-fetch Transporter Benefits
	var transporterBenefits []models.TransporterBenefit
	if err := db.DB.Where("status = 'active' AND (route_id IS NULL OR route_id IN ?)", utils.GetUniqueRouteIDs(transporterRouteMap)).Find(&transporterBenefits).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Generation Failed", "Failed to fetch transporter benefits", err, "incomplete")
		return
	}
	benefitMap := make(map[uint64][]models.TransporterBenefit) // Map routeID to benefits, or 0 for global
	for _, b := range transporterBenefits {
		if b.RouteID != nil {
			benefitMap[*b.RouteID] = append(benefitMap[*b.RouteID], b)
		} else {
			benefitMap[0] = append(benefitMap[0], b) // Global benefits
		}
	}

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
	jobs := make(chan transporterMilkSummary, len(transporterSummaries))
	var mu sync.Mutex
	var totalGross, totalNet, totalDeductions, totalBenefits, totalKilos, totalRejects float64
	var failedCount int64

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for summary := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[TransporterPayrollService] Worker panicked during payroll generation for transporter %d: %v", summary.TransporterID, r)
							db.DB.Create(&models.TransporterPayrollGenerationError{
								BaseModel:     models.BaseModel{CreatedBy: userID},
								TransporterID: summary.TransporterID,
								PayrollID:     payrollID,
								Error:         fmt.Sprintf("Panic during generation: %v", r),
							})
							db.DB.Model(&models.TransporterPayslip{}).Where("transporter_id = ?", summary.TransporterID).Update("status", "incomplete")
							mu.Lock()
							failedCount++
							mu.Unlock()
						}
					}()
					var tGross, tNet, tDeductions, tBenefits float64
					var err error

					for attempt := 1; attempt <= 3; attempt++ {
						err = db.DB.Transaction(func(tx *gorm.DB) error {
							tGross, tNet, tDeductions, tBenefits = 0, 0, 0, 0

							// 1. Determine Transport Rate
							var rate float64
							routeID := transporterRouteMap[summary.TransporterID]

							// Try to find specific rate for (route_id, member_id, transporter_id) - assuming member_id is relevant here
							var transportRate models.TransportRate
							if err := tx.Where("route_id = ? AND member_id = ? AND transporter_id = ?", routeID, 0, summary.TransporterID).First(&transportRate).Error; err == nil {
								rate = transportRate.TransportRate
							} else if err := tx.Where("route_id = ? AND transporter_id = ?", routeID, summary.TransporterID).First(&transportRate).Error; err == nil {
								rate = transportRate.TransportRate
							} else if err := tx.Where("route_id = ? AND member_id IS NULL AND transporter_id IS NULL", routeID).First(&transportRate).Error; err == nil {
								rate = transportRate.TransportRate
							} else if err := tx.Where("route_id IS NULL AND member_id IS NULL AND transporter_id IS NULL").First(&transportRate).Error; err == nil { // Global default
								rate = transportRate.TransportRate
							} else {
								return fmt.Errorf("no transport rate found for transporter %d on route %d", summary.TransporterID, routeID)
							}

							tGross = summary.TotalQuantity * rate

							// 2. Process Benefits
							applicableBenefits := []models.TransporterBenefit{}
							if routeBenefits, ok := benefitMap[routeID]; ok {
								applicableBenefits = append(applicableBenefits, routeBenefits...)
							}
							if globalBenefits, ok := benefitMap[0]; ok {
								applicableBenefits = append(applicableBenefits, globalBenefits...)
							}

							for _, benefit := range applicableBenefits {
								if summary.TotalQuantity >= benefit.MinQuantity {
									benefitAmount := benefit.Rate
									tBenefits += benefitAmount
									// Create TransporterPayrollBenefit record
									tpb := models.TransporterPayrollBenefit{
										TransporterID:        summary.TransporterID,
										TransporterBenefitID: benefit.ID,
										Amount:               benefitAmount,
										PayrollID:            payrollID,
										BaseModel:            models.BaseModel{CreatedBy: userID},
									}
									if err := tx.Create(&tpb).Error; err != nil {
										return err
									}
								}
							}
							tGross += tBenefits

							// 3. Process Deductions
							if tDeds, ok := deductionMap[summary.TransporterID]; ok {
								for _, rd := range tDeds {
									remaining := rd.TotalAmount - rd.PaidAmount
									deductAmount := rd.RecurrentAmount
									if deductAmount > remaining {
										deductAmount = remaining
									}
									if tDeductions+deductAmount > tGross {
										deductAmount = tGross - tDeductions
									}
									if deductAmount <= 0 {
										continue
									}

									// Assuming a TransporterPayrollDeduction model exists
									// mpd := models.TransporterPayrollDeduction{
									// 	TransporterID:   summary.TransporterID,
									// 	PayrollID:       payrollID,
									// 	DeductionTypeID: rd.DeductionTypeID,
									// 	Amount:          deductAmount,
									// 	TransactionDate: dateOpened,
									// 	Reference:       rd.Reference,
									// 	BaseModel:       models.BaseModel{CreatedBy: userID},
									// }
									// if err := tx.Create(&mpd).Error; err != nil {
									// 	return err
									// }
									tDeductions += deductAmount
								}
							}

							tNet = tGross - tDeductions

							payslip := models.TransporterPayslip{
								TransporterID:   summary.TransporterID,
								PayrollID:       payrollID,
								PayDateRangeID:  &pdr.ID,
								FiscalPeriod:    req.FiscalPeriod,
								TotalKilos:      summary.TotalQuantity,
								TotalRejects:    summary.TotalRejects,
								GrossPay:        tGross,
								TotalBenefits:   tBenefits,
								TotalDeductions: tDeductions,
								NetPay:          tNet,
								Status:          "draft",
								BaseModel:       models.BaseModel{CreatedBy: userID},
							}
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
						log.Printf("[TransporterPayrollService] Failed to generate payroll for transporter %d: %v", summary.TransporterID, err)
						db.DB.Create(&models.TransporterPayrollGenerationError{
							BaseModel:     models.BaseModel{CreatedBy: userID},
							TransporterID: summary.TransporterID,
							PayrollID:     payrollID,
							Error:         err.Error(),
						})
					} else {
						mu.Lock()
						totalGross += tGross
						totalNet += tNet
						totalDeductions += tDeductions
						totalBenefits += tBenefits
						totalKilos += summary.TotalQuantity
						totalRejects += summary.TotalRejects
						mu.Unlock()
					}
				}()
			}
		}()
	}

	for _, ts := range transporterSummaries {
		jobs <- ts
	}
	close(jobs)
	wg.Wait()

	status := "draft"
	if failedCount > 0 {
		status = "incomplete"
	}

	db.DB.Model(&models.TransporterPayroll{}).Where("id = ?", payrollID).Updates(map[string]interface{}{
		"total_kilos":      totalKilos,
		"total_rejects":    totalRejects,
		"total_deductions": totalDeductions,
		"total_benefits":   totalBenefits,
		"gross_pay":        totalGross,
		"net_pay":          totalNet,
		"status":           status,
		"updated_by":       userID,
		"updated_at":       time.Now(),
	})

	// Send UI notification
	message := fmt.Sprintf("Transporter payroll generation completed. Success: %d, Failed: %d out of %d records.", int64(len(transporterSummaries))-failedCount, failedCount, len(transporterSummaries))
	notificationType := "SUCCESS"
	errorLink := ""

	if failedCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/transporter-payrolls/generation-errors/%d", payrollID)
		log.Printf("[TransporterPayrollService.generatePayrollInBackground] Transporter payroll generation completed with %d failures for payroll %d.", failedCount, payrollID)
	} else if len(transporterSummaries) == 0 {
		message = "Transporter payroll generation completed. No records were processed."
		notificationType = "INFO"
		log.Printf("[TransporterPayrollService.generatePayrollInBackground] Transporter payroll generation completed. No records were processed for payroll %d.", payrollID)
	} else {
		log.Printf("[TransporterPayrollService.generatePayrollInBackground] Transporter payroll generation completed successfully for payroll %d.", payrollID)
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Transporter Payroll Generation Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("TRANSPORTER_PAYROLL"),
	})
}

func (s *TransporterPayrollService) List(page, limit int) ([]dtos.TransporterPayrollResponse, int64, error) {
	var results []dtos.TransporterPayrollResponse
	var total int64
	db.DB.Model(&models.TransporterPayroll{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT
			tp.*,
			tpdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM transporter_payrolls tp
		LEFT JOIN member_pay_date_ranges tpdr ON tp.pay_date_range_id = tpdr.id
		LEFT JOIN users u_post ON tp.posted_by = u_post.id
		LEFT JOIN users u_conf ON tp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON tp.approved_by = u_appr.id
		WHERE tp.deleted_at IS NULL
		ORDER BY tp.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *TransporterPayrollService) Get(id string) (*dtos.TransporterPayrollResponse, error) {
	var result dtos.TransporterPayrollResponse
	query := `
		SELECT
			tp.*,
			tpdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM transporter_payrolls tp
		LEFT JOIN member_pay_date_ranges tpdr ON tp.pay_date_range_id = tpdr.id
		LEFT JOIN users u_post ON tp.posted_by = u_post.id
		LEFT JOIN users u_conf ON tp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON tp.approved_by = u_appr.id
		WHERE tp.id = ? AND tp.deleted_at IS NULL
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

func (s *TransporterPayrollService) Confirm(payrollID uint64, userID uint64) (*models.TransporterPayroll, error) {
	var payroll models.TransporterPayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		if payroll.Status != "draft" {
			return fmt.Errorf("payroll must be in 'draft' status to be confirmed, current status: %s", payroll.Status)
		}

		now := time.Now()
		if err := tx.Model(&payroll).Updates(map[string]interface{}{
			"status":       "confirmed",
			"confirmed_at": &now,
			"confirmed_by": &userID,
			"updated_by":   userID,
			"updated_at":   now,
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.TransporterPayslip{}).Where("payroll_id = ?", payrollID).Update("status", "confirmed").Error; err != nil {
			return err
		}
		return nil
	})
	return &payroll, err
}

func (s *TransporterPayrollService) Approve(payrollID uint64, userID uint64) (*models.TransporterPayroll, error) {
	var payroll models.TransporterPayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		if payroll.Status != "confirmed" {
			return fmt.Errorf("payroll must be in 'confirmed' status to be approved, current status: %s", payroll.Status)
		}

		var errCount int64
		tx.Model(&models.TransporterPayrollGenerationError{}).Where("payroll_id = ?", payrollID).Count(&errCount)
		if errCount > 0 {
			return fmt.Errorf("cannot approve payroll with %d generation errors", errCount)
		}

		var incompleteCount int64
		tx.Model(&models.TransporterPayslip{}).Where("payroll_id = ? AND status = 'incomplete'", payrollID).Count(&incompleteCount)
		if incompleteCount > 0 {
			return fmt.Errorf("cannot approve payroll with %d incomplete payslips", incompleteCount)
		}

		return tx.Model(&payroll).Update("status", "approving").Error
	})

	if err != nil {
		return nil, err
	}

	go s.approvePayrollInBackground(payrollID, userID)

	return &payroll, nil
}

func (s *TransporterPayrollService) approvePayrollInBackground(payrollID uint64, userID uint64) {
	var payroll models.TransporterPayroll
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
		log.Printf("[TransporterPayrollService] Background approval failed to find payroll %d: %v", payrollID, err)
		return
	}

	now := time.Now()

	var grossRule, loanRule, shareRule models.TransactionPostingRule
	db.DB.Where("transaction_type = ?", "TRANSPORTER_PAYMENT_GROSS").First(&grossRule)
	db.DB.Where("transaction_type = ?", "TRANSPORTER_REPAYMENT_DEDUCTION").First(&loanRule)
	db.DB.Where("transaction_type = ?", "TRANSPORTER_BENEFIT_EXPENSE").First(&shareRule) // Using shareRule for benefits for now, needs proper rule

	if grossRule.ID == 0 {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Approval Failed", "Missing posting rule for TRANSPORTER_PAYMENT_GROSS. Please check accounting configurations.", nil, "confirmed")
		return
	}
	if loanRule.ID == 0 {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Approval Failed", "Missing posting rule for TRANSPORTER_REPAYMENT_DEDUCTION. Please check accounting configurations.", nil, "confirmed")
		return
	}
	if shareRule.ID == 0 {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Approval Failed", "Missing posting rule for TRANSPORTER_BENEFIT_EXPENSE. Please check accounting configurations.", nil, "confirmed")
		return
	}

	var payslips []models.TransporterPayslip
	if err := db.DB.Where("payroll_id = ? AND status != 'approved'", payrollID).Find(&payslips).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Transporter Payroll Approval Failed", "Failed to fetch payslips from database", err, "confirmed")
		return
	}

	numWorkers := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	jobs := make(chan models.TransporterPayslip, len(payslips))
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
							log.Printf("[TransporterPayrollService] Worker panicked during payroll approval for payslip %d: %v", ps.ID, r)
							db.DB.Create(&models.TransporterPayrollGenerationError{
								BaseModel:     models.BaseModel{CreatedBy: userID},
								TransporterID: ps.TransporterID,
								PayrollID:     payrollID,
								Error:         fmt.Sprintf("Panic during approval: %v", r),
							})
							db.DB.Model(&models.TransporterPayslip{}).Where("id = ?", ps.ID).Update("status", "incomplete")
							mu.Lock()
							failedCount++
							mu.Unlock()
						}
					}()
					var err error
					for attempt := 1; attempt <= 3; attempt++ {
						err = db.DB.Transaction(func(tx *gorm.DB) error {
							var currentPS models.TransporterPayslip
							if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&currentPS, ps.ID).Error; err != nil {
								return err
							}
							if currentPS.Status == "approved" {
								return nil
							}

							transaction := models.Transaction{
								BaseModel:       models.BaseModel{CreatedBy: userID},
								Reference:       fmt.Sprintf("TPSL-%d", ps.ID),
								TransactionName: fmt.Sprintf("Transporter Payslip Approval - %s (Transporter %d)", payroll.FiscalPeriod, ps.TransporterID),
								TransactionType: "TRANSPORTER_PAYSLIP",
								TransactionDate: now,
								Description:     fmt.Sprintf("Payslip for Transporter %d in Payroll %s", ps.TransporterID, payroll.FiscalPeriod),
								Status:          "approved",
							}
							if err := tx.Create(&transaction).Error; err != nil {
								return err
							}

							// GL: DR Transporter Expense / CR Transporter Payables (Gross Pay)
							glGross := []models.GeneralLedgerEntry{
								{
									BaseModel:       models.BaseModel{CreatedBy: userID},
									TransactionID:   transaction.ID,
									AccountID:       grossRule.DebitAccountID,
									Debit:           ps.GrossPay,
									TransactionDate: now,
									Description:     fmt.Sprintf("Gross Pay: Transporter %d", ps.TransporterID),
								},
								{
									BaseModel:       models.BaseModel{CreatedBy: userID},
									TransactionID:   transaction.ID,
									AccountID:       grossRule.CreditAccountID,
									Credit:          ps.GrossPay,
									TransactionDate: now,
									Description:     fmt.Sprintf("Gross Pay Payable: Transporter %d", ps.TransporterID),
								},
							}
							if err := tx.Create(&glGross).Error; err != nil {
								return err
							}

							// Process Benefits (if any)
							if ps.TotalBenefits > 0 {
								glBenefit := []models.GeneralLedgerEntry{
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       shareRule.DebitAccountID, // Using shareRule for benefits
										Debit:           ps.TotalBenefits,
										TransactionDate: now,
										Description:     fmt.Sprintf("Benefit Expense: Transporter %d", ps.TransporterID),
									},
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       grossRule.CreditAccountID, // CR Transporter Payables
										Credit:          ps.TotalBenefits,
										TransactionDate: now,
										Description:     fmt.Sprintf("Benefit Payable: Transporter %d", ps.TransporterID),
									},
								}
								if err := tx.Create(&glBenefit).Error; err != nil {
									return err
								}
							}

							// Process Individual Deductions
							var deductions []models.TransporterPayrollDeduction // Assuming same model for now
							if err := tx.Where("payroll_id = ? AND transporter_id = ?", payrollID, ps.TransporterID).Find(&deductions).Error; err != nil {
								return err
							}

							for _, d := range deductions {
								var dType models.DeductionType
								if err := tx.First(&dType, d.DeductionTypeID).Error; err != nil {
									return err
								}

								// Map deduction type to rules (Loan/Other)
								rule := loanRule
								// Add more specific rules if needed for other deduction types
								creditAcc := rule.CreditAccountID
								if creditAcc == 0 {
									return fmt.Errorf("no accounting mapping for deduction type %s", dType.Code)
								}

								// GL: DR Transporter Payables / CR Loan/Other Receivable
								glDeduct := []models.GeneralLedgerEntry{
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       grossRule.CreditAccountID, // DR the payable account
										Debit:           d.Amount,
										TransactionDate: now,
										Description:     fmt.Sprintf("Deduction DR: Transporter %d - %s", ps.TransporterID, d.Reference),
									},
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       creditAcc,
										Credit:          d.Amount,
										TransactionDate: now,
										Description:     fmt.Sprintf("Deduction CR: Transporter %d - %s", ps.TransporterID, d.Reference),
									},
								}
								if err := tx.Create(&glDeduct).Error; err != nil {
									return err
								}

								// Update Recurrent Deduction balances
								var rd models.RecurrentDeduction
								if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
									Where("customer_id = ? AND reference = ? AND customer_type = 'transporter'", d.TransporterID, d.Reference).
									First(&rd).Error; err == nil {

									newPaid := rd.PaidAmount + d.Amount
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
						log.Printf("[TransporterPayrollService] Approval failed for payslip %d: %v", ps.ID, err)
						db.DB.Create(&models.TransporterPayrollGenerationError{
							BaseModel:     models.BaseModel{CreatedBy: userID},
							TransporterID: ps.TransporterID,
							PayrollID:     payrollID,
							Error:         fmt.Sprintf("Approval Error: %s", err.Error()),
						})
						db.DB.Model(&models.TransporterPayslip{}).Where("id = ?", ps.ID).Update("status", "incomplete")
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

	finalStatus := "approved"
	notificationTitle := "Transporter Payroll Approval Status"
	notificationMsg := fmt.Sprintf("Transporter payroll for period %s has been approved successfully.", payroll.FiscalPeriod)
	notificationType := "SUCCESS"
	errorLink := ""

	if failedCount > 0 {
		finalStatus = "incomplete"
		notificationTitle = "Transporter Payroll Approval Incomplete"
		notificationMsg = fmt.Sprintf("Transporter payroll for %s completed with %d errors. Please check the generation errors tab.", payroll.FiscalPeriod, failedCount)
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/transporter-payrolls/approval-errors/%d", payrollID)
	}

	db.DB.Model(&payroll).Updates(map[string]interface{}{
		"status":      finalStatus,
		"approved_at": &now,
		"approved_by": &userID,
		"updated_by":  userID,
		"posted_at":   &now,
		"posted_by":   &userID,
		"updated_at":  time.Now(),
	})

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            notificationTitle,
		Message:          notificationMsg,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("TRANSPORTER_PAYROLL"),
	})
}

// handleProcessingError logs a database error, updates the payroll status, and sends a UI notification.
func (s *TransporterPayrollService) handleProcessingError(payrollID uint64, userID uint64, title, message string, err error, status string) {
	msg := message
	if err != nil {
		msg = fmt.Sprintf("%s: %v", message, err)
	}
	log.Printf("[%s] %s for payroll %d", title, msg, payrollID)
	db.DB.Model(&models.TransporterPayroll{}).Where("id = ?", payrollID).Update("status", status)
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            title,
		Message:          msg,
		NotificationType: "ERROR",
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("TRANSPORTER_PAYROLL"),
	})
}
