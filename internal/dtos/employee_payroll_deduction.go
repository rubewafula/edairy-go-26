package dtos

import "time"

type CreateEmployeePayrollDeductionRequest struct {
	EmployeeID  uint64  `json:"employee_id" validate:"required"`
	DeductionID uint64  `json:"deduction_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Year        string  `json:"year" validate:"required"`
	Month       string  `json:"month" validate:"required"`
	PayrollID   uint64  `json:"payroll_id" validate:"required"`
}

type UpdateEmployeePayrollDeductionRequest struct {
	DeductionID uint64  `json:"deduction_id"`
	Amount      float64 `json:"amount"`
	Year        string  `json:"year"`
	Month       string  `json:"month"`
	PayrollID   uint64  `json:"payroll_id"`
}

type EmployeePayrollDeductionResponse struct {
	ID            uint64    `json:"id"`
	EmployeeID    uint64    `json:"employee_id"`
	EmployeeName  string    `json:"employee_name"`
	DeductionID   uint64    `json:"deduction_id"`
	DeductionName string    `json:"deduction_name"`
	Amount        float64   `json:"amount"`
	Year          string    `json:"year"`
	Month         string    `json:"month"`
	PayrollID     uint64    `json:"payroll_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
