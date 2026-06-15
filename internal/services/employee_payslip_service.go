package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeePayslipService struct{}

func NewEmployeePayslipService() *EmployeePayslipService {
	return &EmployeePayslipService{}
}

func (s *EmployeePayslipService) GetPayslips(employeeID string, payrollID string, page, limit int) ([]dtos.EmployeePayslipResponse, int64, error) {
	var results []dtos.EmployeePayslipResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeePayslip{})
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
			ep.*, CONCAT(e.first_name, ' ', e.surname) as employee_name
		FROM employee_payslips ep
		LEFT JOIN employees e ON ep.employee_id = e.id
		WHERE ep.deleted_at IS NULL 
		AND (? = '' OR ep.employee_id = ?)
		AND (? = '' OR ep.payroll_id = ?)
		ORDER BY ep.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, payrollID, payrollID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeePayslipService) GetPayslip(id string) (*dtos.EmployeePayslipResponse, error) {
	var result dtos.EmployeePayslipResponse
	query := `
		SELECT 
			ep.*, CONCAT(e.first_name, ' ', e.surname) as employee_name
		FROM employee_payslips ep
		LEFT JOIN employees e ON ep.employee_id = e.id
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

func (s *EmployeePayslipService) DeletePayslip(id string, userID uint64) error {
	var payslip models.EmployeePayslip
	if err := db.DB.First(&payslip, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&payslip).Update("updated_by", userID).Delete(&payslip).Error
}
