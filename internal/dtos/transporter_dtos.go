package dtos

import "time"

// TransporterPayDateRange DTOs
type CreateTransporterPayDateRangeRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	PayMonth  string `json:"pay_month" binding:"required"`
	PayYear   string `json:"pay_year" binding:"required"`
}

type UpdateTransporterPayDateRangeRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	PayMonth  string `json:"pay_month"`
	PayYear   string `json:"pay_year"`
}

type TransporterPayDateRangeResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	PayMonth  string    `json:"pay_month"`
	PayYear   string    `json:"pay_year"`
	Processed bool      `json:"processed"`
	Confirmed bool      `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransporterPayroll DTOs
type CreateTransporterPayrollRequest struct {
	PayDateRangeID uint64 `json:"pay_date_range_id" validate:"required"`
	PayrollMonth   string `json:"payroll_month" validate:"required"`
	PayrollYear    string `json:"payroll_year" validate:"required"`
	DateOpened     string `json:"date_opened" validate:"required"`
	Description    string `json:"description"`
	PhysicalPeriod string `json:"physical_period"`
}

type UpdateTransporterPayrollRequest struct {
	TotalKilos      float64 `json:"total_kilos"`
	TotalDeductions float64 `json:"total_deductions"`
	GrossPay        float64 `json:"gross_pay"`
	NetPay          float64 `json:"net_pay"`
	TotalRejects    float64 `json:"total_rejects"`
	TotalBenefits   float64 `json:"total_benefits"`
	Status          string  `json:"status" validate:"omitempty,oneof=draft confirmed approved closed cancelled incomplete"`
}

type TransporterPayrollResponse struct {
	ID               uint64     `json:"id"`
	PayDateRangeID   *uint64    `json:"pay_date_range_id"`
	PayDateRangeName string     `json:"pay_date_range_name"`
	PayrollMonth     string     `json:"payroll_month"`
	PayrollYear      string     `json:"payroll_year"`
	DateOpened       *time.Time `json:"date_opened"`
	Description      string     `json:"description"`
	PhysicalPeriod   string     `json:"physical_period"`
	TotalKilos       float64    `json:"total_kilos"`
	TotalDeductions  float64    `json:"total_deductions"`
	GrossPay         float64    `json:"gross_pay"`
	NetPay           float64    `json:"net_pay"`
	TotalRejects     float64    `json:"total_rejects"`
	TotalBenefits    float64    `json:"total_benefits"`
	Status           string     `json:"status"`
	ConfirmedAt      *time.Time `json:"confirmed_at"`
	ConfirmedByName  string     `json:"confirmed_by_name"`
	ApprovedAt       *time.Time `json:"approved_at"`
	ApprovedByName   string     `json:"approved_by_name"`
	PostedAt         *time.Time `json:"posted_at"`
	PostedByName     string     `json:"posted_by_name"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	PayrollID        *uint64    `json:"payroll_id"`
}

// TransporterPayslip DTOs
type TransporterPayslipResponse struct {
	ID              uint64     `json:"id"`
	TransporterID   uint64     `json:"transporter_id"`
	TransporterName string     `json:"transporter_name"`
	PayrollID       uint64     `json:"payroll_id"`
	PayDateRangeID  *uint64    `json:"pay_date_range_id"`
	PayrollMonth    string     `json:"payroll_month"`
	PayrollYear     string     `json:"payroll_year"`
	PhysicalPeriod  string     `json:"physical_period"`
	TotalKilos      float64    `json:"total_kilos"`
	GrossPay        float64    `json:"gross_pay"`
	TotalDeductions float64    `json:"total_deductions"`
	TotalBenefits   float64    `json:"total_benefits"`
	NetPay          float64    `json:"net_pay"`
	TotalRejects    float64    `json:"total_rejects"`
	Status          string     `json:"status"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
	ConfirmedByName string     `json:"confirmed_by_name"`
	ApprovedAt      *time.Time `json:"approved_at"`
	ApprovedByName  string     `json:"approved_by_name"`
	PostedAt        *time.Time `json:"posted_at"`
	PostedByName    string     `json:"posted_by_name"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TransporterBenefit DTOs
type CreateTransporterBenefitRequest struct {
	Name        string  `json:"name" validate:"required"`
	MinQuantity float64 `json:"min_quantity" validate:"required,gte=0"`
	Rate        float64 `json:"rate" validate:"required,gte=0"`
	RouteID     uint64  `json:"route_id"`
	Status      string  `json:"status"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}

type UpdateTransporterBenefitRequest struct {
	Name        string  `json:"name"`
	MinQuantity float64 `json:"min_quantity"`
	Rate        float64 `json:"rate"`
	RouteID     *uint64 `json:"route_id"`
	Status      string  `json:"status"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}

// TransporterPayrollBenefit DTOs
type CreateTransporterPayrollBenefitRequest struct {
	TransporterID        uint64  `json:"transporter_id" validate:"required"`
	TransporterBenefitID uint64  `json:"transporter_benefit_id" validate:"required"`
	Amount               float64 `json:"amount" validate:"required,gte=0"`
	BenefitYear          string  `json:"benefit_year" validate:"required"`
	BenefitMonth         string  `json:"benefit_month" validate:"required"`
	PayrollID            uint64  `json:"payroll_id" validate:"required"`
}

type UpdateTransporterPayrollBenefitRequest struct {
	TransporterID        uint64  `json:"transporter_id"`
	TransporterBenefitID uint64  `json:"transporter_benefit_id"`
	Amount               float64 `json:"amount"`
	BenefitYear          string  `json:"benefit_year"`
	BenefitMonth         string  `json:"benefit_month"`
	PayrollID            uint64  `json:"payroll_id"`
}

type TransporterPayrollBenefitResponse struct {
	ID                   uint64    `json:"id"`
	TransporterID        uint64    `json:"transporter_id"`
	TransporterName      string    `json:"transporter_name"`
	TransporterBenefitID uint64    `json:"transporter_benefit_id"`
	BenefitName          string    `json:"benefit_name"`
	Amount               float64   `json:"amount"`
	BenefitYear          string    `json:"benefit_year"`
	BenefitMonth         string    `json:"benefit_month"`
	PayrollID            uint64    `json:"payroll_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TransporterPayrollDeduction DTOs (assuming it will be similar to MemberPayrollDeduction)
type TransporterPayrollDeductionResponse struct {
	ID              uint64    `json:"id"`
	TransporterID   uint64    `json:"transporter_id"`
	PayrollID       uint64    `json:"payroll_id"`
	DeductionTypeID uint64    `json:"deduction_type_id"`
	DeductionName   string    `json:"deduction_name"`
	Amount          float64   `json:"amount"`
	Reference       string    `json:"reference"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
