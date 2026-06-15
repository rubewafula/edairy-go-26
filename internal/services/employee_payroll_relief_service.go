package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeePayrollReliefService struct{}

func NewEmployeePayrollReliefService() *EmployeePayrollReliefService {
	return &EmployeePayrollReliefService{}
}

func (s *EmployeePayrollReliefService) GetPayrollReliefs(employeeID string, payrollID string, page, limit int) ([]dtos.EmployeePayrollReliefResponse, int64, error) {
	var results []dtos.EmployeePayrollReliefResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeePayrollRelief{})
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
			epr.id, epr.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			epr.relief_id, tr.relief as relief_name,
			epr.amount, epr.payroll_id,
			epr.created_at, epr.updated_at, epr.created_by, epr.updated_by
		FROM employee_payroll_reliefs epr
		LEFT JOIN employees e ON epr.employee_id = e.id
		LEFT JOIN payroll_reliefs tr ON epr.relief_id = tr.id
		WHERE epr.deleted_at IS NULL 
		AND (? = '' OR epr.employee_id = ?)
		AND (? = '' OR epr.payroll_id = ?)
		ORDER BY epr.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, payrollID, payrollID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeePayrollReliefService) GetPayrollRelief(id string) (*dtos.EmployeePayrollReliefResponse, error) {
	var result dtos.EmployeePayrollReliefResponse
	query := `
		SELECT 
			epr.id, epr.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			epr.relief_id, tr.relief as relief_name,
			epr.amount, epr.payroll_id,
			epr.created_at, epr.updated_at, epr.created_by, epr.updated_by
		FROM employee_payroll_reliefs epr
		LEFT JOIN employees e ON epr.employee_id = e.id
		LEFT JOIN payroll_reliefs tr ON epr.relief_id = tr.id
		WHERE epr.id = ? AND epr.deleted_at IS NULL
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

func (s *EmployeePayrollReliefService) DeletePayrollRelief(id string, userID uint64) error {
	return db.DB.Model(&models.EmployeePayrollRelief{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.EmployeePayrollRelief{}).Error
}
