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

type MemberPayrollService struct{}

func NewMemberPayrollService() *MemberPayrollService {
	return &MemberPayrollService{}
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
	if err == nil && existing.Status != "draft" {
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
			PhysicalPeriod: req.PhysicalPeriod,
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
	db.DB.Table("milk_journal_entries").
		Select("member_id, SUM(quantity) as kilos").
		Joins("JOIN milk_journals ON milk_journal_entries.milk_journal_id = milk_journals.id").
		Where("milk_journals.confirmed = ? AND milk_journals.journal_date BETWEEN ? AND ?", true, pdr.StartDate, pdr.EndDate).
		Group("member_id").Scan(&milkCollections)

	if len(milkCollections) == 0 {
		db.DB.Model(&models.MemberPayroll{}).Where("id = ?", payrollID).Update("status", "draft")
		return
	}

	// 2. Fetch Aggregated Rejects
	var milkRejects []milkSum
	db.DB.Table("milk_rejects").
		Select("member_id, SUM(quantity) as kilos").
		Where("transaction_date BETWEEN ? AND ?", pdr.StartDate, pdr.EndDate).
		Group("member_id").Scan(&milkRejects)

	rejectMap := make(map[uint64]float64)
	for _, r := range milkRejects {
		rejectMap[r.MemberID] = r.Kilos
	}

	// 3. Pre-fetch Rate Resolution Maps
	var specialRates []models.MilkSpecialRate
	db.DB.Where("pay_date_range_id = ? AND deleted_at IS NULL", req.PayDateRangeID).Find(&specialRates)
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
	db.DB.Where("deleted_at IS NULL").Find(&defaultRates)
	defaultRouteRateMap := make(map[uint64]float64)
	var globalDefault float64
	for _, r := range defaultRates {
		if r.RouteID != nil {
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
	db.DB.Select("id, route_id").Where("id IN ?", memberIDs).Find(&members)
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
							db.DB.Create(&models.MemberPayrollGenerationError{
								BaseModel: models.BaseModel{CreatedBy: userID},
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
							mGross, mNet, mDeductions = 0, 0, 0
							rejectKilos := rejectMap[collection.MemberID]
							netKilos := collection.Kilos - rejectKilos
							if netKilos < 0 {
								netKilos = 0
							}

							routeID := memberRouteMap[collection.MemberID]
							var rate float64
							if r, ok := memberRateMap[collection.MemberID]; ok {
								rate = r
							} else if r, ok := routePeriodRateMap[routeID]; ok {
								rate = r
							} else if r, ok := defaultRouteRateMap[routeID]; ok {
								rate = r
							} else {
								rate = globalDefault
							}

							mGross = netKilos * rate

							// Fetch deductions for this member with a row-level lock
							var memberDeductions []models.RecurrentDeduction
							if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
								Where("customer_id = ? AND settled = 0 AND customer_type = 'member'", collection.MemberID).
								Order("created_at ASC").Find(&memberDeductions).Error; err != nil {
								return err
							}

							var payslipDeductions float64
							for _, rd := range memberDeductions {
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
									MemberID:        int64(collection.MemberID),
									PayrollID:       payrollID,
									DeductionTypeID: rd.DeductionTypeID,
									Amount:          deductAmount,
									TransactionDate: dateOpened,
									Reference:       rd.Reference,
								}
								if err := tx.Create(&mpd).Error; err != nil {
									return err
								}

								payslipDeductions += deductAmount
							}

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
								NetPay:          mNet,
								PhysicalPeriod:  req.PhysicalPeriod,
								PayDateRangeID:  &req.PayDateRangeID,
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
						log.Printf("[MemberPayrollService] Failed to generate payroll for member %d: %v", collection.MemberID, err)
						db.DB.Create(&models.MemberPayrollGenerationError{
							BaseModel: models.BaseModel{CreatedBy: userID},
							MemberID:  collection.MemberID,
							PayrollID: payrollID,
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
	db.DB.Where("transaction_type = ?", "MILK_PAYMENT").First(&grossRule)
	db.DB.Where("transaction_type = ?", "MEMBER_REPAYMENT_DEDUCTION").First(&loanRule)
	db.DB.Where("transaction_type = ?", "SHARES_CONTRIBUTION").First(&shareRule)

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
		log.Printf("[MemberPayrollService] Failed to fetch payslips for payroll %d: %v", payrollID, err)
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
							db.DB.Create(&models.MemberPayrollGenerationError{
								BaseModel: models.BaseModel{CreatedBy: userID},
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
								TransactionName: fmt.Sprintf("Payslip Approval - %s (Member %d)", payroll.PhysicalPeriod, ps.MemberID),
								TransactionType: "PAYROLL_PAYSLIP",
								TransactionDate: now,
								Description:     fmt.Sprintf("Payslip for Member %d in Payroll %s", ps.MemberID, payroll.PhysicalPeriod),
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

								// DR Member Payables / CR Loan/Share Receivable
								glDeduct := []models.GeneralLedgerEntry{
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       grossRule.CreditAccountID,
										Debit:           d.Amount,
										TransactionDate: now,
										Description:     fmt.Sprintf("Deduction DR: Member %d - %s", ps.MemberID, d.Reference),
									},
									{
										BaseModel:       models.BaseModel{CreatedBy: userID},
										TransactionID:   transaction.ID,
										AccountID:       creditAcc,
										Credit:          d.Amount,
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
							PayrollID: payrollID,
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
