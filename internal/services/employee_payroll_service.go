package services

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeePayrollService struct{}

func NewEmployeePayrollService() *EmployeePayrollService {
	return &EmployeePayrollService{}
}

func (s *EmployeePayrollService) CreatePayroll(req dtos.CreateEmployeePayrollRequest, userID uint64) error {
	log.Printf("[EmployeePayrollService.CreatePayroll] User %d initiated payroll generation for %s %s", userID, req.PayrollMonth, req.PayrollYear)

	// IDEMPOTENCY CHECK: Ensure we don't overwrite confirmed/approved payrolls
	var existing models.EmployeePayroll
	err := db.DB.Where("payroll_month = ? AND payroll_year = ? AND deleted_at IS NULL", req.PayrollMonth, req.PayrollYear).First(&existing).Error
	if err == nil && existing.Status != "draft" && existing.Status != "incomplete" {
		return fmt.Errorf("payroll for %s %s is already %s and cannot be regenerated", req.PayrollMonth, req.PayrollYear, existing.Status)
	}

	// IDEMPOTENCY PREPARATION: Clear existing draft data for this month
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var payrollIDs []uint64
		tx.Model(&models.EmployeePayroll{}).Where("payroll_month = ? AND payroll_year = ? AND status IN ?", req.PayrollMonth, req.PayrollYear, []string{"draft", "incomplete", "processing"}).Pluck("id", &payrollIDs)

		if len(payrollIDs) > 0 {
			tx.Where("payroll_id IN ?", payrollIDs).Delete(&models.EmployeePayslip{})
			tx.Where("payroll_id IN ?", payrollIDs).Delete(&models.EmployeePayrollDeduction{})
			tx.Where("payroll_id IN ?", payrollIDs).Delete(&models.EmployeePayrollBenefit{})
			tx.Where("payroll_id IN ?", payrollIDs).Delete(&models.EmployeePayrollRelief{})
			tx.Where("id IN ?", payrollIDs).Delete(&models.EmployeePayroll{})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to clear existing payroll data: %w", err)
	}

	go s.generatePayrollInBackground(req, userID)
	return nil
}

func (s *EmployeePayrollService) generatePayrollInBackground(req dtos.CreateEmployeePayrollRequest, userID uint64) {
	log.Printf("[EmployeePayrollService.generatePayrollInBackground] Starting payroll generation for %s %s", req.PayrollMonth, req.PayrollYear)

	// Compute Monthly date boundaries for benefit/deduction pre-loading
	monthTime, _ := time.Parse("January", req.PayrollMonth)
	yearInt, _ := strconv.Atoi(req.PayrollYear)
	startDate := time.Date(yearInt, monthTime.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	// 1. Fetch all active employees
	var employees []models.Employee
	if err := db.DB.Where("status = 'ACTIVE' AND deleted_at IS NULL").Find(&employees).Error; err != nil {
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Error fetching employees: %v", err)
		return
	}

	// 2. Preload necessary data for all employees
	employeeIDs := make([]uint64, len(employees))
	for i, emp := range employees {
		employeeIDs[i] = emp.ID
	}

	var salaries []models.EmployeeSalary
	db.DB.Where("employee_id IN ? AND status = 'ACTIVE' AND deleted_at IS NULL", employeeIDs).Find(&salaries)
	salaryMap := make(map[uint64]float64)
	for _, sal := range salaries {
		salaryMap[sal.EmployeeID] = sal.BasicSalary
	}

	var benefits []models.EmployeeBenefit
	db.DB.Where("employee_id IN ? AND status = 'ACTIVE' AND deleted_at IS NULL AND (start_date <= ? OR start_date IS NULL) AND (end_date >= ? OR end_date IS NULL)",
		employeeIDs, endDate, startDate).Find(&benefits)
	benefitMap := make(map[uint64][]models.EmployeeBenefit)
	for _, b := range benefits {
		benefitMap[b.EmployeeID] = append(benefitMap[b.EmployeeID], b)
	}

	var deductions []models.EmployeeDeduction
	db.DB.Where("employee_id IN ? AND status = 1 AND deleted_at IS NULL AND (start_date <= ? OR start_date IS NULL) AND (end_date >= ? OR end_date IS NULL)",
		employeeIDs, endDate, startDate).Find(&deductions)
	deductionMap := make(map[uint64][]models.EmployeeDeduction)
	for _, d := range deductions {
		deductionMap[d.EmployeeID] = append(deductionMap[d.EmployeeID], d)
	}

	var statutoryConfigs []models.StatutoryDeductionConfiguration
	db.DB.Where("deleted_at IS NULL").Find(&statutoryConfigs)
	statutoryConfigMap := make(map[uint64]models.StatutoryDeductionConfiguration) // DeductionTypeID -> Config
	for _, sc := range statutoryConfigs {
		statutoryConfigMap[sc.DeductionID] = sc
	}

	// 3. Create Payroll Header
	payrollHeader := models.EmployeePayroll{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		PayrollMonth: req.PayrollMonth,
		PayrollYear:  req.PayrollYear,
		DateOpened:   time.Now(),
		Status:       "processing",
	}
	if err := db.DB.Create(&payrollHeader).Error; err != nil {
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Error creating payroll header: %v", err)
		return
	}

	// 4. Concurrently process payroll for each employee
	var wg sync.WaitGroup
	errorChan := make(chan error, len(employees))
	failedEmployeeCount := 0

	numWorkers := runtime.NumCPU() * 2
	if numWorkers < 1 {
		numWorkers = 1
	}
	jobs := make(chan models.Employee, len(employees))

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for emp := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[EmployeePayrollService] Worker panicked during payroll generation for employee %d: %v", emp.ID, r)
							errorMessage := fmt.Sprintf("[%s %s] Panic during generation: %v", req.PayrollMonth, req.PayrollYear, r)
							db.DB.Create(&models.EmployeePayrollGenerationError{
								BaseModel:  models.BaseModel{CreatedBy: userID},
								EmployeeID: emp.ID,
								Error:      errorMessage,
							})
							errorChan <- fmt.Errorf("panic during payroll generation for employee %d: %v", emp.ID, r)
						}
					}()

					err := s.processSingleEmployeePayroll(
						emp,
						payrollHeader.ID,
						req.PayrollMonth,
						req.PayrollYear,
						userID,
						salaryMap[emp.ID],
						benefitMap[emp.ID],
						deductionMap[emp.ID],
						statutoryConfigMap,
					)
					if err != nil {
						log.Printf("[EmployeePayrollService.generatePayrollInBackground] Error processing payroll for employee %d: %v", emp.ID, err)
						errorMessage := fmt.Sprintf("[%s %s] Generation Error: %v", req.PayrollMonth, req.PayrollYear, err)
						db.DB.Create(&models.EmployeePayrollGenerationError{
							BaseModel:  models.BaseModel{CreatedBy: userID},
							EmployeeID: emp.ID,
							Error:      errorMessage,
						})
						errorChan <- err
					}
				}()
			}
		}()
	}

	for _, emp := range employees {
		jobs <- emp
	}
	close(jobs)
	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			failedEmployeeCount++
		}
	}

	// 5. Finalize Payroll Header and Pay Date Range status
	finalPayrollStatus := "draft"
	if failedEmployeeCount > 0 {
		finalPayrollStatus = "incomplete"
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Completed payroll generation for %s %s with %d failures.", req.PayrollMonth, req.PayrollYear, failedEmployeeCount)
	} else {
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Completed payroll generation for %s %s successfully.", req.PayrollMonth, req.PayrollYear)
	}

	db.DB.Model(&payrollHeader).Where("id = ?", payrollHeader.ID).Updates(map[string]interface{}{
		"status":     finalPayrollStatus,
		"updated_by": userID,
	})
}

// processSingleEmployeePayroll handles the payroll generation for a single employee.
func (s *EmployeePayrollService) processSingleEmployeePayroll(
	employee models.Employee,
	payrollID uint64,
	payMonth, payYear string,
	userID uint64,
	basicSalary float64,
	employeeBenefits []models.EmployeeBenefit,
	employeeDeductions []models.EmployeeDeduction,
	statutoryConfigMap map[uint64]models.StatutoryDeductionConfiguration,
) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var grossPay, totalBenefits, totalDeductions, totalTax, totalRelief, netPay float64

		// 1. Basic Salary
		grossPay = basicSalary

		// 2. Process Benefits
		for _, benefit := range employeeBenefits {
			totalBenefits += benefit.Amount
			// Create EmployeePayrollBenefit record
			epb := models.EmployeePayrollBenefit{
				BaseModel:  models.BaseModel{CreatedBy: userID},
				EmployeeID: employee.ID,
				BenefitID:  benefit.ID,
				Amount:     benefit.Amount,
				PayrollID:  payrollID,
				// Month and Year can be derived from payDateRangeID if needed
			}
			if err := tx.Create(&epb).Error; err != nil {
				return err
			}
		}
		grossPay += totalBenefits

		// 3. Process Deductions
		for _, deduction := range employeeDeductions {
			deductionAmount := deduction.Amount

			// Check if it's a statutory deduction and apply rules
			if config, ok := statutoryConfigMap[deduction.DeductionTypeID]; ok {
				// Example: NSSF/NHIF/PAYE calculation logic
				// This is a simplified example; real statutory calculations are complex.
				if config.FixedAmount > 0 {
					deductionAmount = config.FixedAmount
				} else if config.EmployeeDeductionRate > 0 {
					deductionAmount = grossPay * (config.EmployeeDeductionRate / 100)
				}
				// Apply min/max amounts, bands etc.
				if deductionAmount < config.MinAmount {
					deductionAmount = config.MinAmount
				}
				if config.MaxAmount > 0 && deductionAmount > config.MaxAmount {
					deductionAmount = config.MaxAmount
				}
				// For PAYE, you'd need tax bands and reliefs, which is beyond this scope but would fit here.
				// For simplicity, let's assume statutory deductions are directly applied or calculated based on basic rules.
			}

			totalDeductions += deductionAmount

			// Create EmployeePayrollDeduction record
			epd := models.EmployeePayrollDeduction{
				BaseModel:   models.BaseModel{CreatedBy: userID},
				EmployeeID:  employee.ID,
				DeductionID: deduction.ID,
				Amount:      deductionAmount,
				PayrollID:   payrollID,
				// Month and Year can be derived from payDateRangeID
			}
			if err := tx.Create(&epd).Error; err != nil {
				return err
			}
		}

		// 4. Process Reliefs (if any)
		// This would involve fetching EmployeeReliefs and applying them to reduce tax.
		// For now, let's assume totalRelief is 0 or pre-calculated.

		// 5. Calculate Net Pay
		netPay = grossPay - totalDeductions - totalTax + totalRelief // totalTax and totalRelief would be calculated from grossPay and statutory rules

		// 6. Create EmployeePayslip
		payslip := models.EmployeePayslip{
			BaseModel:       models.BaseModel{CreatedBy: userID},
			EmployeeID:      employee.ID,
			PayrollID:       payrollID,
			PayrollMonth:    payMonth,
			PayrollYear:     payYear,
			BasicSalary:     basicSalary,
			GrossPay:        grossPay,
			TotalBenefits:   totalBenefits,
			TotalDeductions: totalDeductions,
			TotalTax:        totalTax,
			TotalRelief:     totalRelief,
			NetPay:          netPay,
			Status:          "draft",
		}
		if err := tx.Create(&payslip).Error; err != nil {
			return err
		}

		// Update payslip_id for deductions and benefits
		if err := tx.Model(&models.EmployeePayrollDeduction{}).Where("payroll_id = ? AND employee_id = ?", payrollID, employee.ID).Update("payslip_id", payslip.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.EmployeePayrollBenefit{}).Where("payroll_id = ? AND employee_id = ?", payrollID, employee.ID).Update("payslip_id", payslip.ID).Error; err != nil {
			return err
		}
		// If reliefs are created as separate records, update their payslip_id too.

		log.Printf("[EmployeePayrollService.processSingleEmployeePayroll] Successfully processed payroll for employee %d for %s %s", employee.ID, payMonth, payYear)
		return nil
	})
}

func (s *EmployeePayrollService) GetEmployeePayrolls(page, limit int) ([]dtos.EmployeePayrollResponse, int64, error) {
	var results []dtos.EmployeePayrollResponse
	var total int64
	db.DB.Model(&models.EmployeePayroll{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ep.*, 
			epdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM employee_payrolls ep
		LEFT JOIN employee_pay_date_ranges epdr ON ep.pay_date_range_id = epdr.id
		LEFT JOIN users u_post ON ep.posted_by = u_post.id
		LEFT JOIN users u_conf ON ep.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON ep.approved_by = u_appr.id
		WHERE ep.deleted_at IS NULL
		ORDER BY ep.id DESC 
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeePayrollService) GetEmployeePayroll(id string) (*dtos.EmployeePayrollResponse, error) { // Renamed from GetEmployeePayroll to GetPayroll
	var result dtos.EmployeePayrollResponse
	query := `
		SELECT 
			ep.*, 
			epdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM employee_payrolls ep
		LEFT JOIN employee_pay_date_ranges epdr ON ep.pay_date_range_id = epdr.id
		LEFT JOIN users u_post ON ep.posted_by = u_post.id
		LEFT JOIN users u_conf ON ep.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON ep.approved_by = u_appr.id
		WHERE ep.id = ? AND ep.deleted_at IS NULL
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

func (s *EmployeePayrollService) ConfirmPayroll(payrollID uint64, userID uint64) (*models.EmployeePayroll, error) {
	var payroll models.EmployeePayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		if payroll.Status != "draft" && payroll.Status != "incomplete" {
			return fmt.Errorf("payroll must be in 'draft' or 'incomplete' status to be confirmed, current status: %s", payroll.Status)
		}

		// Check for any generation errors
		var errCount int64
		tx.Model(&models.EmployeePayrollGenerationError{}).Where("error LIKE ?", "%"+payroll.PayrollMonth+" "+payroll.PayrollYear+"%").Count(&errCount)
		if errCount > 0 {
			return fmt.Errorf("cannot confirm payroll with %d generation errors. Please resolve them first", errCount)
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

		if err := tx.Model(&models.EmployeePayslip{}).Where("payroll_id = ?", payrollID).Update("status", "confirmed").Error; err != nil {
			return err
		}
		log.Printf("[EmployeePayrollService.ConfirmPayroll] Payroll %d confirmed by user %d", payrollID, userID)
		return nil
	})
	return &payroll, err
}

func (s *EmployeePayrollService) ApprovePayroll(payrollID uint64, userID uint64) (*models.EmployeePayroll, error) {
	var payroll models.EmployeePayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		if payroll.Status != "confirmed" {
			return fmt.Errorf("payroll must be in 'confirmed' status to be approved, current status: %s", payroll.Status)
		}

		var errCount int64
		tx.Model(&models.EmployeePayrollGenerationError{}).Where("error LIKE ?", "%"+payroll.PayrollMonth+" "+payroll.PayrollYear+"%").Count(&errCount)
		if errCount > 0 {
			return fmt.Errorf("cannot approve payroll with %d generation errors", errCount)
		}

		var incompleteCount int64
		tx.Model(&models.EmployeePayslip{}).Where("payroll_id = ? AND status = 'incomplete'", payrollID).Count(&incompleteCount)
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

func (s *EmployeePayrollService) approvePayrollInBackground(payrollID uint64, userID uint64) {
	var payroll models.EmployeePayroll
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
		log.Printf("[EmployeePayrollService] Background approval failed to find payroll %d: %v", payrollID, err)
		return
	}

	now := time.Now()

	// 2. Fetch Account Posting Rules (Example rules, adjust as per actual accounting setup)
	var grossRule, deductionRule, benefitRule, taxRule, reliefRule models.TransactionPostingRule
	db.DB.Where("transaction_type = ?", "EMPLOYEE_SALARY_GROSS").First(&grossRule)
	db.DB.Where("transaction_type = ?", "EMPLOYEE_DEDUCTION_PAYABLE").First(&deductionRule)
	db.DB.Where("transaction_type = ?", "EMPLOYEE_BENEFIT_EXPENSE").First(&benefitRule)
	db.DB.Where("transaction_type = ?", "EMPLOYEE_TAX_PAYABLE").First(&taxRule)
	db.DB.Where("transaction_type = ?", "EMPLOYEE_RELIEF_RECEIVABLE").First(&reliefRule)

	// Basic validation for posting rules
	if grossRule.ID == 0 || deductionRule.ID == 0 || benefitRule.ID == 0 || taxRule.ID == 0 || reliefRule.ID == 0 {
		log.Printf("[EmployeePayrollService] Missing one or more posting rules for payroll %d. Aborting approval.", payrollID)
		db.DB.Model(&payroll).Update("status", "confirmed") // Revert to confirmed
		return
	}

	var payslips []models.EmployeePayslip
	if err := db.DB.Where("payroll_id = ? AND status != 'approved'", payrollID).Find(&payslips).Error; err != nil {
		log.Printf("[EmployeePayrollService] Failed to fetch payslips for payroll %d: %v", payrollID, err)
		db.DB.Model(&payroll).Update("status", "confirmed") // Revert to confirmed
		return
	}

	numWorkers := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	jobs := make(chan models.EmployeePayslip, len(payslips))
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
							log.Printf("[EmployeePayrollService] Worker panicked during payroll approval for payslip %d: %v", ps.ID, r)
							errorMessage := fmt.Sprintf("[%s %s] Panic during approval for payslip %d: %v", payroll.PayrollMonth, payroll.PayrollYear, ps.ID, r)
							db.DB.Create(&models.EmployeePayrollGenerationError{
								BaseModel:  models.BaseModel{CreatedBy: userID},
								EmployeeID: ps.EmployeeID,
								Error:      errorMessage,
							})
							db.DB.Model(&models.EmployeePayslip{}).Where("id = ?", ps.ID).Update("status", "incomplete")
							mu.Lock()
							failedCount++
							mu.Unlock()
						}
					}()
					var err error
					for attempt := 1; attempt <= 3; attempt++ { // Retry logic
						err = db.DB.Transaction(func(tx *gorm.DB) error {
							var currentPS models.EmployeePayslip
							if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&currentPS, ps.ID).Error; err != nil {
								return err
							}
							if currentPS.Status == "approved" {
								return nil // Already processed
							}

							// Create Transaction Header for THIS payslip
							transaction := models.Transaction{
								BaseModel:       models.BaseModel{CreatedBy: userID},
								Reference:       fmt.Sprintf("EPSL-%d", ps.ID),
								TransactionName: fmt.Sprintf("Employee Payslip Approval - %s %s (Employee %d)", payroll.PayrollMonth, payroll.PayrollYear, ps.EmployeeID),
								TransactionType: "EMPLOYEE_PAYSLIP",
								TransactionDate: now,
								Description:     fmt.Sprintf("Payslip for Employee %d in Payroll %s %s", ps.EmployeeID, payroll.PayrollMonth, payroll.PayrollYear),
								Status:          "approved",
							}
							if err := tx.Create(&transaction).Error; err != nil {
								return err
							}

							// GL: DR Salary Expense / CR Employee Payables (Gross Pay)
							if err := s.createGLEntries(tx, transaction.ID, userID, now, grossRule.DebitAccountID, grossRule.CreditAccountID, ps.GrossPay, fmt.Sprintf("Gross Pay: Employee %d", ps.EmployeeID)); err != nil {
								return err
							}

							// GL: DR Employee Payables / CR Deduction Payables (for each deduction)
							var deductions []models.EmployeePayrollDeduction
							if err := tx.Where("payslip_id = ?", ps.ID).Find(&deductions).Error; err != nil {
								return err
							}
							for _, d := range deductions {
								var empDeduction models.EmployeeDeduction
								if err := tx.First(&empDeduction, d.DeductionID).Error; err != nil {
									return err
								}

								var dType models.EmployeeDeductionType
								if err := tx.First(&dType, empDeduction.DeductionTypeID).Error; err != nil {
									return err
								}
								// This is a simplified mapping. In a real system, each deduction type would have its own posting rule.
								// For now, using a generic deductionRule.
								if err := s.createGLEntries(tx, transaction.ID, userID, now, deductionRule.DebitAccountID, deductionRule.CreditAccountID, d.Amount, fmt.Sprintf("Deduction: Employee %d - %s", ps.EmployeeID, dType.Name)); err != nil {
									return err
								}
							}

							// GL: DR Benefit Expense / CR Employee Payables (for each benefit)
							var benefits []models.EmployeePayrollBenefit
							if err := tx.Where("payslip_id = ?", ps.ID).Find(&benefits).Error; err != nil {
								return err
							}
							for _, b := range benefits {
								// Assuming a generic benefit rule
								if err := s.createGLEntries(tx, transaction.ID, userID, now, benefitRule.DebitAccountID, benefitRule.CreditAccountID, b.Amount, fmt.Sprintf("Benefit: Employee %d - Benefit %d", ps.EmployeeID, b.BenefitID)); err != nil {
									return err
								}
							}

							// GL: DR Employee Payables / CR Tax Payables (for tax)
							if ps.TotalTax > 0 {
								if err := s.createGLEntries(tx, transaction.ID, userID, now, taxRule.DebitAccountID, taxRule.CreditAccountID, ps.TotalTax, fmt.Sprintf("Tax: Employee %d", ps.EmployeeID)); err != nil {
									return err
								}
							}

							// GL: DR Employee Payables / CR Relief Receivable (for reliefs)
							if ps.TotalRelief > 0 {
								if err := s.createGLEntries(tx, transaction.ID, userID, now, reliefRule.DebitAccountID, reliefRule.CreditAccountID, ps.TotalRelief, fmt.Sprintf("Relief: Employee %d", ps.EmployeeID)); err != nil {
									return err
								}
							}

							// Update payslip status
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
							break // Success, break retry loop
						}
						if errors.Is(err, gorm.ErrRecordNotFound) { // No need to retry if record not found
							break
						}
						time.Sleep(time.Duration(attempt) * 50 * time.Millisecond) // Exponential backoff
					}

					if err != nil {
						mu.Lock()
						failedCount++
						mu.Unlock()
						log.Printf("[EmployeePayrollService] Approval failed for payslip %d: %v", ps.ID, err)
						errorMessage := fmt.Sprintf("[%s %s] Approval Error for payslip %d: %s", payroll.PayrollMonth, payroll.PayrollYear, ps.ID, err.Error())
						db.DB.Create(&models.EmployeePayrollGenerationError{
							BaseModel:  models.BaseModel{CreatedBy: userID},
							EmployeeID: ps.EmployeeID,
							Error:      errorMessage,
						})
						db.DB.Model(&models.EmployeePayslip{}).Where("id = ?", ps.ID).Update("status", "incomplete")
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
		"updated_at":  time.Now(),
	})
}

func (s *EmployeePayrollService) createGLEntries(tx *gorm.DB, transactionID, userID uint64, transactionDate time.Time, debitAccountID, creditAccountID uint64, amount float64, description string) error {
	glEntries := []models.GeneralLedgerEntry{
		{
			BaseModel:       models.BaseModel{CreatedBy: userID},
			TransactionID:   transactionID,
			AccountID:       debitAccountID,
			Debit:           amount,
			TransactionDate: transactionDate,
			Description:     description,
		},
		{
			BaseModel:       models.BaseModel{CreatedBy: userID},
			TransactionID:   transactionID,
			AccountID:       creditAccountID,
			Credit:          amount,
			TransactionDate: transactionDate,
			Description:     description,
		},
	}
	return tx.Create(&glEntries).Error
}

func (s *EmployeePayrollService) DeleteEmployeePayroll(id string, userID uint64) error {
	// Implement soft delete for payroll header and associated payslips, deductions, benefits, reliefs
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var payroll models.EmployeePayroll
		if err := tx.First(&payroll, id).Error; err != nil {
			return err
		}

		if payroll.Status != "draft" && payroll.Status != "incomplete" {
			return fmt.Errorf("only 'draft' or 'incomplete' payrolls can be deleted")
		}

		if err := tx.Where("payroll_id = ?", id).Delete(&models.EmployeePayslip{}).Error; err != nil {
			return err
		}
		if err := tx.Where("payroll_id = ?", id).Delete(&models.EmployeePayrollDeduction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("payroll_id = ?", id).Delete(&models.EmployeePayrollBenefit{}).Error; err != nil {
			return err
		}
		if err := tx.Where("payroll_id = ?", id).Delete(&models.EmployeePayrollRelief{}).Error; err != nil {
			return err
		}

		return tx.Model(&payroll).Update("updated_by", userID).Delete(&payroll).Error
	})
}
