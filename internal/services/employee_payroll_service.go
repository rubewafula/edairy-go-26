package services

import (
	"bytes"
	"encoding/csv"
	"errors"
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

type EmployeePayrollService struct {
	notificationService *UINotificationService
}

func NewEmployeePayrollService() *EmployeePayrollService {
	return &EmployeePayrollService{
		notificationService: NewUINotificationService(),
	}
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
			tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.EmployeePayrollGenerationError{})
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
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Critical Error fetching employees: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Employee Payroll Generation Failed",
			Message:          fmt.Sprintf("Could not fetch active employees: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

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

	var empReliefs []models.EmployeeRelief
	db.DB.Where("employee_id IN ? AND status = 'ACTIVE' AND deleted_at IS NULL", employeeIDs).Find(&empReliefs)
	reliefMap := make(map[uint64][]models.EmployeeRelief)
	for _, r := range empReliefs {
		reliefMap[r.EmployeeID] = append(reliefMap[r.EmployeeID], r)
	}

	var reliefConfigs []models.EmployeePayrollRelief
	db.DB.Where("deleted_at IS NULL").Find(&reliefConfigs)
	reliefConfigMap := make(map[uint64]models.EmployeePayrollRelief)
	for _, rc := range reliefConfigs {
		reliefConfigMap[rc.ID] = rc
	}

	// 3. Create Payroll Header
	payrollHeader := models.EmployeePayroll{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		PayrollMonth: req.PayrollMonth,
		PayrollYear:  req.PayrollYear,
		DateOpened:   time.Now(),
		Status:       "processing",
		PaidAt:       nil,
	}
	if err := db.DB.Create(&payrollHeader).Error; err != nil {
		s.handleProcessingError(0, userID, "Employee Payroll Header Error", "Failed to create payroll header record", err, "incomplete")
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
						reliefMap[emp.ID],
						reliefConfigMap,
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
	}

	// Send UI notification
	totalProcessed := len(employees)
	successCount := totalProcessed - failedEmployeeCount
	message := fmt.Sprintf("Employee payroll generation completed. Success: %d, Failed: %d out of %d records.", successCount, failedEmployeeCount, totalProcessed)
	notificationType := "SUCCESS"
	errorLink := ""

	if failedEmployeeCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/employee-payrolls/generation-errors/%d", payrollHeader.ID)
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Completed with %d failures for payroll %d", failedEmployeeCount, payrollHeader.ID)
	} else if totalProcessed == 0 {
		message = "Employee payroll generation completed. No active employees were found to process."
		notificationType = "INFO"
	} else {
		log.Printf("[EmployeePayrollService.generatePayrollInBackground] Completed successfully for payroll %d", payrollHeader.ID)
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Employee Payroll Generation Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
		ReferenceID:      &payrollHeader.ID,
		ReferenceType:    utils.StringPtr("EMPLOYEE_PAYROLL"),
	})

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
	employeeReliefs []models.EmployeeRelief,
	reliefConfigMap map[uint64]models.EmployeePayrollRelief,
) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create EmployeePayslip first to get its ID
		payslip := models.EmployeePayslip{
			BaseModel:    models.BaseModel{CreatedBy: userID},
			EmployeeID:   employee.ID,
			PayrollID:    payrollID,
			PayrollMonth: payMonth,
			PayrollYear:  payYear,
			Status:       "draft", // Initial status
		}
		if err := tx.Create(&payslip).Error; err != nil {
			return err
		}

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
				PayslipID:  payslip.ID, // Assign payslip ID directly
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
				PayslipID:   payslip.ID, // Assign payslip ID directly
				Amount:      deductionAmount,
				PayrollID:   payrollID,
				// Month and Year can be derived from payDateRangeID
			}
			if err := tx.Create(&epd).Error; err != nil {
				return err
			}
		}

		// 4. Process Reliefs (if any)
		for _, er := range employeeReliefs {
			if config, ok := reliefConfigMap[er.ReliefID]; ok {
				reliefAmount := config.Amount
				totalRelief += reliefAmount

				epr := models.EmployeePayrollRelief{
					BaseModel:  models.BaseModel{CreatedBy: userID},
					EmployeeID: employee.ID,
					ReliefID:   er.ReliefID,
					Amount:     reliefAmount,
					PayrollID:  payrollID,
				}
				if err := tx.Create(&epr).Error; err != nil {
					return err
				}
			}
		}

		// 5. Calculate Net Pay
		netPay = grossPay - totalDeductions - totalTax + totalRelief // totalTax and totalRelief would be calculated from grossPay and statutory rules

		// Update the payslip with calculated values
		if err := tx.Model(&payslip).Updates(map[string]interface{}{
			"basic_salary":     basicSalary,
			"gross_pay":        grossPay,
			"total_benefits":   totalBenefits,
			"total_deductions": totalDeductions,
			"total_tax":        totalTax,
			"total_relief":     totalRelief,
			"net_pay":          netPay,
		}).Error; err != nil {
			return err
		}

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
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM employee_payrolls ep
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
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM employee_payrolls ep
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

func (s *EmployeePayrollService) ApprovePayroll(payrollID uint64, userID uint64, isApproved bool) (*models.EmployeePayroll, error) {
	var payroll models.EmployeePayroll
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", payrollID).First(&payroll).Error; err != nil {
			return err
		}

		if !isApproved {
			if payroll.Status != "confirmed" {
				return fmt.Errorf("payroll must be in 'confirmed' status to be rejected, current status: %s", payroll.Status)
			}

			// Reject logic
			now := time.Now()
			if err := tx.Model(&payroll).Updates(map[string]interface{}{
				"status":      "rejected",
				"rejected_at": &now,
				"rejected_by": &userID,
				"updated_by":  userID,
				"updated_at":  now,
			}).Error; err != nil {
				return err
			}

			// Mark all associated payslips as rejected
			if err := tx.Model(&models.EmployeePayslip{}).Where("payroll_id = ?", payrollID).Update("status", "rejected").Error; err != nil {
				return err
			}

			// Delete associated payroll deductions, benefits, and reliefs
			if err := tx.Where("payroll_id = ?", payrollID).Delete(&models.EmployeePayrollDeduction{}).Error; err != nil {
				return err
			}
			if err := tx.Where("payroll_id = ?", payrollID).Delete(&models.EmployeePayrollBenefit{}).Error; err != nil {
				return err
			}
			if err := tx.Where("payroll_id = ?", payrollID).Delete(&models.EmployeePayrollRelief{}).Error; err != nil {
				return err
			}

			log.Printf("[EmployeePayrollService.ApprovePayroll] Payroll %d rejected by user %d", payrollID, userID)
			return nil // Transaction successful for rejection
		} else {
			// Approve logic (existing)
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
		}
	})

	if err != nil {
		return nil, err
	}

	if isApproved {
		go s.approvePayrollInBackground(payrollID, userID)
	}

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
		s.handleProcessingError(payrollID, userID, "Employee Payroll Approval Failed", "Missing one or more accounting posting rules. Please check configurations.", nil, "confirmed")
		return
	}

	var payslips []models.EmployeePayslip
	if err := db.DB.Where("payroll_id = ? AND status != 'approved'", payrollID).Find(&payslips).Error; err != nil {
		s.handleProcessingError(payrollID, userID, "Employee Payroll Approval Failed", "Failed to fetch payslips from database", err, "confirmed")
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

	// Send UI notification
	totalPayslips := len(payslips)
	message := fmt.Sprintf("Employee payroll approval completed. Success: %d, Failed: %d out of %d payslips.", totalPayslips-int(failedCount), failedCount, totalPayslips)
	notificationType := "SUCCESS"
	errorLink := ""

	if failedCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/employee-payrolls/approval-errors/%d", payrollID)
		log.Printf("[EmployeePayrollService.approvePayrollInBackground] Approval completed with %d failures for payroll %d.", failedCount, payrollID)
	} else if totalPayslips == 0 {
		message = "Employee payroll approval completed. No pending payslips were found."
		notificationType = "INFO"
	} else {
		log.Printf("[EmployeePayrollService.approvePayrollInBackground] Approval completed successfully for payroll %d.", payrollID)
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Employee Payroll Approval Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("EMPLOYEE_PAYROLL"),
	})

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

func (s *EmployeePayrollService) handleProcessingError(payrollID uint64, userID uint64, title, message string, err error, status string) {
	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
	}
	log.Printf("[%s] %s for payroll %d", title, errMsg, payrollID)
	if payrollID > 0 {
		db.DB.Model(&models.EmployeePayroll{}).Where("id = ?", payrollID).Update("status", status)
	}
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            title,
		Message:          errMsg,
		NotificationType: "ERROR",
		ReferenceID:      &payrollID,
		ReferenceType:    utils.StringPtr("EMPLOYEE_PAYROLL"),
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

func (s *EmployeePayrollService) GetPayslipStatement(employeeID, payrollID string) (*dtos.EmployeePayslipStatementResponse, error) {
	var payslip models.EmployeePayslip
	if err := db.DB.Where("employee_id = ? AND payroll_id = ? AND deleted_at IS NULL", employeeID, payrollID).First(&payslip).Error; err != nil {
		return nil, err
	}

	var employee dtos.EmployeeResponse
	// Aggregating employee details for the statement
	db.DB.Model(&models.Employee{}).
		Select("employees.*, jp.name as job_position_name").
		Joins("left join job_positions jp on employees.job_position_id = jp.id").
		Where("employees.id = ?", payslip.EmployeeID).
		Scan(&employee)

	// Fetch Benefits breakdown
	var benefits []dtos.EmployeePayrollBenefitResponse
	db.DB.Raw(`
		SELECT epb.id, eb.benefit_id, b.name as benefit_name, epb.amount
		FROM employee_payroll_benefits epb
		JOIN employee_benefits eb ON epb.employee_benefit_id = eb.id
		JOIN benefits b ON eb.benefit_id = b.id
		WHERE epb.payslip_id = ? AND epb.deleted_at IS NULL
	`, payslip.ID).Scan(&benefits)

	// Fetch Deductions breakdown
	var deductions []dtos.EmployeePayrollDeductionResponse
	db.DB.Raw(`
		SELECT epd.id, ed.deduction_type_id as deduction_id, edt.name as deduction_name, epd.amount
		FROM employee_payroll_deductions epd
		JOIN employee_deductions ed ON epd.employee_deduction_id = ed.id
		JOIN employee_deduction_types edt ON ed.deduction_type_id = edt.id
		WHERE epd.payslip_id = ? AND epd.deleted_at IS NULL
	`, payslip.ID).Scan(&deductions)

	// Fetch Reliefs breakdown
	var reliefs []dtos.EmployeePayrollReliefResponse
	db.DB.Raw(`
		SELECT epr.id, epr.relief_id, pr.relief as relief_name, epr.amount
		FROM employee_payroll_reliefs epr
		JOIN payroll_reliefs pr ON epr.relief_id = pr.id
		WHERE epr.employee_id = ? AND epr.payroll_id = ? AND epr.deleted_at IS NULL
	`, payslip.EmployeeID, payslip.PayrollID).Scan(&reliefs)

	return &dtos.EmployeePayslipStatementResponse{
		Payslip: dtos.EmployeePayslipResponse{
			ID:              payslip.ID,
			EmployeeID:      payslip.EmployeeID,
			PayrollMonth:    payslip.PayrollMonth,
			PayrollYear:     payslip.PayrollYear,
			GrossPay:        payslip.GrossPay,
			NetPay:          payslip.NetPay,
			TotalDeductions: payslip.TotalDeductions,
			TotalBenefits:   payslip.TotalBenefits,
			BasicSalary:     payslip.BasicSalary,
			PayrollID:       payslip.PayrollID,
			TotalTax:        payslip.TotalTax,
			TotalRelief:     payslip.TotalRelief,
			Status:          payslip.Status,
			CreatedAt:       payslip.CreatedAt,
		},
		Employee:   employee,
		Benefits:   benefits,
		Deductions: deductions,
		Reliefs:    reliefs,
	}, nil
}

func (s *EmployeePayrollService) ExportPayslipStatement(userID uint64, employeeID, payrollID, reportType string) error {
	go s.processPayslipStatementExportInBackground(userID, employeeID, payrollID, reportType)
	return nil
}

func (s *EmployeePayrollService) processPayslipStatementExportInBackground(userID uint64, employeeID, payrollID, reportType string) {
	result, err := s.GetPayslipStatement(employeeID, payrollID)
	if err != nil {
		log.Printf("[EmployeePayrollService.processPayslipStatementExportInBackground] Error fetching data: %v", err)
		return
	}

	var fileData []byte
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateDetailedPayslipPDF(result)
	} else {
		fileData, err = s.generateDetailedPayslipCSV(result)
	}

	if err != nil {
		log.Printf("[EmployeePayrollService] Generation error: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("employee_payslip_%s_%d.%s", employeeID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[EmployeePayrollService] File write error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Employee Payslip Statement Ready",
		Message:          fmt.Sprintf("Your detailed payslip statement for %s %s is ready for download.", result.Payslip.PayrollMonth, result.Payslip.PayrollYear),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/employee-payslips/export/download/%s", filename),
	})
}

func (s *EmployeePayrollService) generateDetailedPayslipCSV(data *dtos.EmployeePayslipStatementResponse) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"Employee No", "Name", "Period", "Category", "Description", "Amount"})

	p := data.Payslip
	e := data.Employee
	period := fmt.Sprintf("%s %s", p.PayrollMonth, p.PayrollYear)
	fullName := fmt.Sprintf("%s %s", e.FirstName, e.Surname)

	writer.Write([]string{e.EmployeeNo, fullName, period, "EARNING", "Basic Salary", fmt.Sprintf("%.2f", p.BasicSalary)})
	for _, b := range data.Benefits {
		writer.Write([]string{e.EmployeeNo, fullName, period, "BENEFIT", b.BenefitName, fmt.Sprintf("%.2f", b.Amount)})
	}
	for _, d := range data.Deductions {
		writer.Write([]string{e.EmployeeNo, fullName, period, "DEDUCTION", d.DeductionName, fmt.Sprintf("%.2f", d.Amount)})
	}
	for _, r := range data.Reliefs {
		writer.Write([]string{e.EmployeeNo, fullName, period, "RELIEF", r.ReliefName, fmt.Sprintf("%.2f", r.Amount)})
	}
	writer.Write([]string{e.EmployeeNo, fullName, period, "SUMMARY", "NET PAY", fmt.Sprintf("%.2f", p.NetPay)})

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *EmployeePayrollService) generateDetailedPayslipPDF(data *dtos.EmployeePayslipStatementResponse) ([]byte, error) {
	var org struct{ RegisteredName, Address string }
	db.DB.Table("organization_details").First(&org)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "EMPLOYEE PAYSLIP STATEMENT", "B", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Employee No:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 7, data.Employee.EmployeeNo, "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Employee Name:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 7, fmt.Sprintf("%s %s", data.Employee.FirstName, data.Employee.Surname), "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 7, "Period:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 7, fmt.Sprintf("%s %s", data.Payslip.PayrollMonth, data.Payslip.PayrollYear), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(140, 8, "Description", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 8, "Amount", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(140, 8, "Basic Salary", "1", 0, "L", false, 0, "")
	pdf.CellFormat(50, 8, fmt.Sprintf("%.2f", data.Payslip.BasicSalary), "1", 1, "R", false, 0, "")
	for _, b := range data.Benefits {
		pdf.CellFormat(140, 8, b.BenefitName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 8, fmt.Sprintf("%.2f", b.Amount), "1", 1, "R", false, 0, "")
	}
	for _, d := range data.Deductions {
		pdf.CellFormat(140, 8, d.DeductionName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 8, fmt.Sprintf("-%.2f", d.Amount), "1", 1, "R", false, 0, "")
	}
	for _, r := range data.Reliefs {
		pdf.CellFormat(140, 8, r.ReliefName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 8, fmt.Sprintf("%.2f", r.Amount), "1", 1, "R", false, 0, "")
	}

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(140, 10, "NET PAYABLE", "1", 0, "R", false, 0, "")
	pdf.CellFormat(50, 10, fmt.Sprintf("%.2f", data.Payslip.NetPay), "1", 1, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

type payslipExportData struct {
	EmployeeNo      string  `gorm:"column:employee_no"`
	EmployeeName    string  `gorm:"column:employee_name"`
	PayrollMonth    string  `gorm:"column:payroll_month"`
	PayrollYear     string  `gorm:"column:payroll_year"`
	BasicSalary     float64 `gorm:"column:basic_salary"`
	GrossPay        float64 `gorm:"column:gross_pay"`
	TotalBenefits   float64 `gorm:"column:total_benefits"`
	TotalDeductions float64 `gorm:"column:total_deductions"`
	TotalTax        float64 `gorm:"column:total_tax"`
	TotalRelief     float64 `gorm:"column:total_relief"`
	NetPay          float64 `gorm:"column:net_pay"`
	Status          string  `gorm:"column:status"`
}

func (s *EmployeePayrollService) ExportPayslips(userID uint64, filters map[string]string, format string) error {
	go s.processPayslipsExportInBackground(userID, filters, format)
	return nil
}

func (s *EmployeePayrollService) processPayslipsExportInBackground(userID uint64, filters map[string]string, format string) {

	var data []payslipExportData
	query := db.DB.Table("employee_payslips ep").
		Select("e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name, ep.payroll_month, ep.payroll_year, ep.basic_salary, ep.gross_pay, ep.total_benefits, ep.total_deductions, ep.total_tax, ep.total_relief, ep.net_pay, ep.status").
		Joins("JOIN employees e ON ep.employee_id = e.id").
		Where("ep.deleted_at IS NULL")

	if val, ok := filters["payroll_id"]; ok {
		query = query.Where("ep.payroll_id = ?", val)
	}
	if val, ok := filters["payroll_month"]; ok {
		query = query.Where("ep.payroll_month = ?", val)
	}
	if val, ok := filters["payroll_year"]; ok {
		query = query.Where("ep.payroll_year = ?", val)
	}

	if err := query.Scan(&data).Error; err != nil {
		log.Printf("[EmployeePayrollService.processPayslipsExportInBackground] Fetch Error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(format) == "pdf" {
		ext = "pdf"
		fileData, err = s.generatePayslipsListPDF(data)
	} else {
		fileData, err = s.generatePayslipsListCSV(data)
	}

	if err != nil {
		log.Printf("[EmployeePayrollService.processPayslipsExportInBackground] Generation error: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("employee_payslips_export_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[EmployeePayrollService.processPayslipsExportInBackground] File write error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Employee Payslips Export Ready",
		Message:          "The bulk export of employee payslips is ready for download.",
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/employee-payslips/export/download/%s", filename),
	})
}

func (s *EmployeePayrollService) generatePayslipsListCSV(data []payslipExportData) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"Employee No", "Employee Name", "Month", "Year", "Basic Salary", "Gross Pay", "Benefits", "Deductions", "Tax", "Relief", "Net Pay", "Status"})

	for _, row := range data {
		writer.Write([]string{
			row.EmployeeNo, row.EmployeeName, row.PayrollMonth, row.PayrollYear,
			fmt.Sprintf("%.2f", row.BasicSalary), fmt.Sprintf("%.2f", row.GrossPay),
			fmt.Sprintf("%.2f", row.TotalBenefits), fmt.Sprintf("%.2f", row.TotalDeductions),
			fmt.Sprintf("%.2f", row.TotalTax), fmt.Sprintf("%.2f", row.TotalRelief),
			fmt.Sprintf("%.2f", row.NetPay), row.Status,
		})
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *EmployeePayrollService) generatePayslipsListPDF(data []payslipExportData) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "EMPLOYEE PAYROLL SUMMARY", "B", 1, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"Emp No", "Name", "Month/Year", "Basic", "Gross", "Benefits", "Deduc.", "Tax", "Relief", "Net"}
	widths := []float64{20, 50, 25, 23, 23, 23, 23, 23, 23, 23}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, row := range data {
		pdf.CellFormat(widths[0], 7, row.EmployeeNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 7, row.EmployeeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 7, fmt.Sprintf("%s %s", row.PayrollMonth, row.PayrollYear), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[3], 7, fmt.Sprintf("%.2f", row.BasicSalary), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[4], 7, fmt.Sprintf("%.2f", row.GrossPay), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[5], 7, fmt.Sprintf("%.2f", row.TotalBenefits), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[6], 7, fmt.Sprintf("%.2f", row.TotalDeductions), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[7], 7, fmt.Sprintf("%.2f", row.TotalTax), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[8], 7, fmt.Sprintf("%.2f", row.TotalRelief), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[9], 7, fmt.Sprintf("%.2f", row.NetPay), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
