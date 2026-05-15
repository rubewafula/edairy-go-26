package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type EmployeePayrollService struct{}

func NewEmployeePayrollService() *EmployeePayrollService {
	return &EmployeePayrollService{}
}

func (s *EmployeePayrollService) CreateEmployeePayroll(req dtos.CreateEmployeePayrollRequest, userID uint64) (*models.EmployeePayroll, error) {
	payroll := &models.EmployeePayroll{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		PayrollMonth:    req.PayrollMonth,
		PayrollYear:     req.PayrollYear,
		DateOpened:      utils.ParseDate(req.DateOpened),
		TotalDeductions: req.TotalDeductions,
		GrossPay:        req.GrossPay,
		NetPay:          req.NetPay,
		Complete:        req.Complete,
		Confirmed:       req.Confirmed,
		Approved:        req.Approved,
		TotalBenefits:   req.TotalBenefits,
		TotalTax:        req.TotalTax,
		TotalRelief:     req.TotalRelief,
		Period:          req.Period,
		PaidAt:          utils.ParseDate(req.PaidAt),
	}
	if err := db.DB.Create(payroll).Error; err != nil {
		return nil, err
	}
	return payroll, nil
}

func (s *EmployeePayrollService) GetEmployeePayrolls(page, limit int) ([]models.EmployeePayroll, int64, error) {
	var payrolls []models.EmployeePayroll
	var total int64
	db.DB.Model(&models.EmployeePayroll{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&payrolls).Error
	return payrolls, total, err
}

func (s *EmployeePayrollService) GetEmployeePayroll(id string) (*models.EmployeePayroll, error) {
	var payroll models.EmployeePayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (s *EmployeePayrollService) UpdateEmployeePayroll(id string, req dtos.UpdateEmployeePayrollRequest, userID uint64) error {
	var payroll models.EmployeePayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"payroll_month": req.PayrollMonth,
		"payroll_year":  req.PayrollYear,
		// Add other updatable fields here
		"updated_by": userID,
	}
	// Example of conditional update for date fields if they are part of UpdateEmployeePayrollRequest
	// if req.DateOpened != "" {
	// 	updates["date_opened"] = utils.ParseDate(req.DateOpened)
	// }
	// if req.PaidAt != "" {
	// 	updates["paid_at"] = utils.ParseDate(req.PaidAt)
	// }

	return db.DB.Model(&payroll).Updates(updates).Error
}

func (s *EmployeePayrollService) DeleteEmployeePayroll(id string, userID uint64) error {
	var payroll models.EmployeePayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&payroll).Update("updated_by", userID).Delete(&payroll).Error
}
