package dtos

import "time"

type CreateEmployeePayslipRequest struct {
	EmployeeID      uint64  `json:"employee_id" validate:"required"`
	PayrollMonth    string  `json:"payroll_month" validate:"required"`
	PayrollYear     string  `json:"payroll_year" validate:"required"`
	GrossPay        float64 `json:"gross_pay"`
	NetPay          float64 `json:"net_pay"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalBenefits   float64 `json:"total_benefits"`
	BasicSalary     float64 `json:"basic_salary"`
	PayrollID       uint64  `json:"payroll_id" validate:"required"`
	TotalTax        float64 `json:"total_tax"`
	TotalRelief     float64 `json:"total_relief"`
}

type UpdateEmployeePayslipRequest struct {
	GrossPay        float64 `json:"gross_pay"`
	NetPay          float64 `json:"net_pay"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalBenefits   float64 `json:"total_benefits"`
	BasicSalary     float64 `json:"basic_salary"`
	TotalTax        float64 `json:"total_tax"`
	TotalRelief     float64 `json:"total_relief"`
}

type EmployeePayslipResponse struct {
	ID              uint64    `json:"id"`
	EmployeeID      uint64    `json:"employee_id"`
	EmployeeName    string    `json:"employee_name"`
	PayrollMonth    string    `json:"payroll_month"`
	PayrollYear     string    `json:"payroll_year"`
	GrossPay        float64   `json:"gross_pay"`
	NetPay          float64   `json:"net_pay"`
	TotalDeductions float64   `json:"total_deductions"`
	TotalBenefits   float64   `json:"total_benefits"`
	BasicSalary     float64   `json:"basic_salary"`
	PayrollID       uint64    `json:"payroll_id"`
	TotalTax        float64   `json:"total_tax"`
	TotalRelief     float64   `json:"total_relief"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
