package dtos

import "time"

type CreateEmployeePayrollBenefitRequest struct {
	EmployeeID uint64  `json:"employee_id" validate:"required"`
	BenefitID  uint64  `json:"benefit_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Year       string  `json:"year" validate:"required"`
	Month      string  `json:"month" validate:"required"`
	PayrollID  uint64  `json:"payroll_id" validate:"required"`
}

type UpdateEmployeePayrollBenefitRequest struct {
	BenefitID uint64  `json:"benefit_id"`
	Amount    float64 `json:"amount"`
	Year      string  `json:"year"`
	Month     string  `json:"month"`
	PayrollID uint64  `json:"payroll_id"`
}

type EmployeePayrollBenefitResponse struct {
	ID           uint64    `json:"id"`
	EmployeeID   uint64    `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	BenefitID    uint64    `json:"benefit_id"`
	BenefitName  string    `json:"benefit_name"`
	Amount       float64   `json:"amount"`
	Year         string    `json:"year"`
	Month        string    `json:"month"`
	PayrollID    uint64    `json:"payroll_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
