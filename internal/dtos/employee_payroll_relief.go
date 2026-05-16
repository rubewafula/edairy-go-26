package dtos

import "time"

type CreateEmployeePayrollReliefRequest struct {
	EmployeeID uint64 `json:"employee_id" validate:"required"`
	ReliefID   uint64 `json:"relief_id" validate:"required"`
	Amount     string `json:"amount" validate:"required"`
	PayrollID  uint64 `json:"payroll_id" validate:"required"`
}

type UpdateEmployeePayrollReliefRequest struct {
	ReliefID  uint64 `json:"relief_id"`
	Amount    string `json:"amount"`
	PayrollID uint64 `json:"payroll_id"`
}

type EmployeePayrollReliefResponse struct {
	ID           uint64    `json:"id"`
	EmployeeID   uint64    `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	ReliefID     uint64    `json:"relief_id"`
	ReliefName   string    `json:"relief_name"`
	Amount       string    `json:"amount"`
	PayrollID    uint64    `json:"payroll_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uint64    `json:"created_by"`
	UpdatedBy    uint64    `json:"updated_by"`
}
