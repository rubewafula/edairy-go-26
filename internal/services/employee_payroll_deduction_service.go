package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeePayrollDeductionService struct{}

func NewEmployeePayrollDeductionService() *EmployeePayrollDeductionService {
	return &EmployeePayrollDeductionService{}
}

func (s *EmployeePayrollDeductionService) CreatePayrollDeduction(req dtos.CreateEmployeePayrollDeductionRequest, userID uint64) (*models.EmployeePayrollDeduction, error) {
	deduction := &models.EmployeePayrollDeduction{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		EmployeeID:  req.EmployeeID,
		DeductionID: req.DeductionID,
		Amount:      req.Amount,
		Year:        req.Year,
		Month:       req.Month,
		PayrollID:   req.PayrollID,
	}

	if err := db.DB.Create(deduction).Error; err != nil {
		return nil, err
	}
	return deduction, nil
}

func (s *EmployeePayrollDeductionService) GetPayrollDeductions(employeeID string, payrollID string, page, limit int) ([]dtos.EmployeePayrollDeductionResponse, int64, error) {
	var results []dtos.EmployeePayrollDeductionResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeePayrollDeduction{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}
	if payrollID != "" {
		queryBuilder = queryBuilder.Where("payroll_id = ?", payrollID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			epd.id, epd.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			epd.employee_deduction_id as deduction_id, edt.name as deduction_name,
			epd.amount, epd.deduction_year as year, epd.deduction_month as month, epd.payroll_id,
			epd.created_at, epd.updated_at
		FROM employee_payroll_deductions epd
		LEFT JOIN employees e ON epd.employee_id = e.id
		LEFT JOIN employee_deduction_types edt ON epd.employee_deduction_id = edt.id
		WHERE epd.deleted_at IS NULL 
		AND (? = '' OR epd.employee_id = ?)
		AND (? = '' OR epd.payroll_id = ?)
		ORDER BY epd.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, payrollID, payrollID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeePayrollDeductionService) GetPayrollDeduction(id string) (*dtos.EmployeePayrollDeductionResponse, error) {
	var result dtos.EmployeePayrollDeductionResponse
	query := `
		SELECT 
			epd.id, epd.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			epd.employee_deduction_id as deduction_id, edt.name as deduction_name,
			epd.amount, epd.deduction_year as year, epd.deduction_month as month, epd.payroll_id,
			epd.created_at, epd.updated_at
		FROM employee_payroll_deductions epd
		LEFT JOIN employees e ON epd.employee_id = e.id
		LEFT JOIN employee_deduction_types edt ON epd.employee_deduction_id = edt.id
		WHERE epd.id = ? AND epd.deleted_at IS NULL
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

func (s *EmployeePayrollDeductionService) UpdatePayrollDeduction(id string, req dtos.UpdateEmployeePayrollDeductionRequest, userID uint64) error {
	var deduction models.EmployeePayrollDeduction
	if err := db.DB.First(&deduction, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"employee_deduction_id": req.DeductionID,
		"amount":                req.Amount,
		"deduction_year":        req.Year,
		"deduction_month":       req.Month,
		"payroll_id":            req.PayrollID,
		"updated_by":            userID,
	}

	return db.DB.Model(&deduction).Updates(updates).Error
}

func (s *EmployeePayrollDeductionService) DeletePayrollDeduction(id string, userID uint64) error {
	return db.DB.Model(&models.EmployeePayrollDeduction{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.EmployeePayrollDeduction{}).Error
}
