package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeePayrollBenefitService struct{}

func NewEmployeePayrollBenefitService() *EmployeePayrollBenefitService {
	return &EmployeePayrollBenefitService{}
}

func (s *EmployeePayrollBenefitService) CreatePayrollBenefit(req dtos.CreateEmployeePayrollBenefitRequest, userID uint64) (*models.EmployeePayrollBenefit, error) {
	benefit := &models.EmployeePayrollBenefit{
		BaseModel:  models.BaseModel{CreatedBy: userID},
		EmployeeID: req.EmployeeID,
		BenefitID:  req.BenefitID,
		Amount:     req.Amount,
		Year:       req.Year,
		Month:      req.Month,
		PayrollID:  req.PayrollID,
	}

	if err := db.DB.Create(benefit).Error; err != nil {
		return nil, err
	}
	return benefit, nil
}

func (s *EmployeePayrollBenefitService) GetPayrollBenefits(employeeID string, payrollID string, page, limit int) ([]dtos.EmployeePayrollBenefitResponse, int64, error) {
	var results []dtos.EmployeePayrollBenefitResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeePayrollBenefit{})
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
			epb.id, epb.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			epb.employee_benefit_id as benefit_id, b.name as benefit_name,
			epb.amount, epb.benefit_year as year, epb.benefit_month as month, epb.payroll_id,
			epb.created_at, epb.updated_at
		FROM employee_payroll_benefits epb
		LEFT JOIN employees e ON epb.employee_id = e.id
		LEFT JOIN benefits b ON epb.employee_benefit_id = b.id
		WHERE epb.deleted_at IS NULL 
		AND (? = '' OR epb.employee_id = ?)
		AND (? = '' OR epb.payroll_id = ?)
		ORDER BY epb.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, payrollID, payrollID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeePayrollBenefitService) GetPayrollBenefit(id string) (*dtos.EmployeePayrollBenefitResponse, error) {
	var result dtos.EmployeePayrollBenefitResponse
	query := `
		SELECT 
			epb.id, epb.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			epb.employee_benefit_id as benefit_id, b.name as benefit_name,
			epb.amount, epb.benefit_year as year, epb.benefit_month as month, epb.payroll_id,
			epb.created_at, epb.updated_at
		FROM employee_payroll_benefits epb
		LEFT JOIN employees e ON epb.employee_id = e.id
		LEFT JOIN benefits b ON epb.employee_benefit_id = b.id
		WHERE epb.id = ? AND epb.deleted_at IS NULL
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

func (s *EmployeePayrollBenefitService) UpdatePayrollBenefit(id string, req dtos.UpdateEmployeePayrollBenefitRequest, userID uint64) error {
	var benefit models.EmployeePayrollBenefit
	if err := db.DB.First(&benefit, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"employee_benefit_id": req.BenefitID,
		"amount":              req.Amount,
		"benefit_year":        req.Year,
		"benefit_month":       req.Month,
		"payroll_id":          req.PayrollID,
		"updated_by":          userID,
	}

	return db.DB.Model(&benefit).Updates(updates).Error
}

func (s *EmployeePayrollBenefitService) DeletePayrollBenefit(id string, userID uint64) error {
	return db.DB.Model(&models.EmployeePayrollBenefit{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.EmployeePayrollBenefit{}).Error
}
