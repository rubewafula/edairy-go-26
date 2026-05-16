package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type EmployeePayrollService struct{}

func NewEmployeePayrollService() *EmployeePayrollService {
	return &EmployeePayrollService{}
}

func (s *EmployeePayrollService) CreateEmployeePayroll(req dtos.CreateEmployeePayrollRequest, userID uint64) (*dtos.EmployeePayrollResponse, error) {
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
	return s.GetEmployeePayroll(utils.Uint64ToString(payroll.ID))
}

func (s *EmployeePayrollService) GetEmployeePayrolls(page, limit int) ([]dtos.EmployeePayrollResponse, int64, error) {
	var results []dtos.EmployeePayrollResponse
	var total int64
	db.DB.Model(&models.EmployeePayroll{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Model(&models.EmployeePayroll{}).Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *EmployeePayrollService) GetEmployeePayroll(id string) (*dtos.EmployeePayrollResponse, error) {
	var result dtos.EmployeePayrollResponse
	err := db.DB.Model(&models.EmployeePayroll{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *EmployeePayrollService) UpdateEmployeePayroll(id string, req dtos.UpdateEmployeePayrollRequest, userID uint64) error {
	var payroll models.EmployeePayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"payroll_month":    req.PayrollMonth,
		"payroll_year":     req.PayrollYear,
		"total_deductions": req.TotalDeductions,
		"gross_pay":        req.GrossPay,
		"net_pay":          req.NetPay,
		"complete":         req.Complete,
		"confirmed":        req.Confirmed,
		"approved":         req.Approved,
		"total_benefits":   req.TotalBenefits,
		"total_tax":        req.TotalTax,
		"total_relief":     req.TotalRelief,
		"period":           req.Period,
		"updated_by":       userID,
	}

	if req.DateOpened != "" {
		updates["date_opened"] = utils.ParseDate(req.DateOpened)
	}
	if req.PaidAt != "" {
		updates["paid_at"] = utils.ParseDate(req.PaidAt)
	}

	return db.DB.Model(&payroll).Updates(updates).Error
}

func (s *EmployeePayrollService) DeleteEmployeePayroll(id string, userID uint64) error {
	var payroll models.EmployeePayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&payroll).Update("updated_by", userID).Delete(&payroll).Error
}
