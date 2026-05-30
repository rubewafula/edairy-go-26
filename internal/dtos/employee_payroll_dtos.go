package dtos

import "time"

// EmployeePayDateRange DTOs
type CreateEmployeePayDateRangeRequest struct {
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
	PayMonth  string `json:"pay_month"`
	PayYear   string `json:"pay_year"`
}

type UpdateEmployeePayDateRangeRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	PayMonth  string `json:"pay_month"`
	PayYear   string `json:"pay_year"`
}

type EmployeePayDateRangeResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	PayMonth  string    `json:"pay_month"`
	PayYear   string    `json:"pay_year"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StatutoryDeductionConfiguration DTOs
type CreateStatutoryDeductionConfigurationRequest struct {
	DeductionID           uint64  `json:"deduction_id" validate:"required"`
	EmployeeDeductionRate float64 `json:"employee_deduction_rate"`
	EmployerDeductionRate float64 `json:"employer_deduction_rate"`
	MinAmount             float64 `json:"min_amount"`
	FixedAmount           float64 `json:"fixed_amount"`
	BandLowerLimitAmount  float64 `json:"band_lower_limit_amount"`
	BandUpperLimitAmount  float64 `json:"band_upper_limit_amount"`
	MinApplicableAmount   float64 `json:"min_applicable_amount"`
}

type UpdateStatutoryDeductionConfigurationRequest struct {
	DeductionID           uint64  `json:"deduction_id"`
	EmployeeDeductionRate float64 `json:"employee_deduction_rate"`
	EmployerDeductionRate float64 `json:"employer_deduction_rate"`
	MinAmount             float64 `json:"min_amount"`
	FixedAmount           float64 `json:"fixed_amount"`
	BandLowerLimitAmount  float64 `json:"band_lower_limit_amount"`
	BandUpperLimitAmount  float64 `json:"band_upper_limit_amount"`
	MinApplicableAmount   float64 `json:"min_applicable_amount"`
}

type StatutoryDeductionConfigurationResponse struct {
	ID                    uint64    `json:"id"`
	DeductionID           uint64    `json:"deduction_id"`
	DeductionTypeName     string    `json:"deduction_type_name"`
	EmployeeDeductionRate float64   `json:"employee_deduction_rate"`
	EmployerDeductionRate float64   `json:"employer_deduction_rate"`
	MinAmount             float64   `json:"min_amount"`
	FixedAmount           float64   `json:"fixed_amount"`
	BandLowerLimitAmount  float64   `json:"band_lower_limit_amount"`
	BandUpperLimitAmount  float64   `json:"band_upper_limit_amount"`
	MinApplicableAmount   float64   `json:"min_applicable_amount"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
