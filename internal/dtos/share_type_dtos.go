package dtos

import "time"

type CreateShareTypeRequest struct {
	ShareCode          string  `json:"share_code" validate:"required"`
	ShareType          string  `json:"share_type" validate:"required"`
	Description        string  `json:"description"`
	Rate               float64 `json:"rate"`
	Mandatory          int     `json:"mandatory"`
	HasShareValue      int     `json:"has_share_value"`
	CalculatingMethod  string  `json:"calculating_method"`
	ShareValue         float64 `json:"share_value"`
	DeductionTypeID    uint64  `json:"deduction_type_id"`
	Priority           int     `json:"priority"`
	IsPayrollDeduction int     `json:"is_payroll_deduction"`
	EarnsDividend      int     `json:"earns_dividend"`
	IsTransferable     int     `json:"is_transferable"`
	MinimumShares      float64 `json:"minimum_shares"`
	MaxmumShares       float64 `json:"maxmum_shares"`
}

type UpdateShareTypeRequest struct {
	ShareCode          string  `json:"share_code"`
	ShareType          string  `json:"share_type"`
	Description        string  `json:"description"`
	Rate               float64 `json:"rate"`
	Mandatory          int     `json:"mandatory"`
	HasShareValue      int     `json:"has_share_value"`
	CalculatingMethod  string  `json:"calculating_method"`
	ShareValue         float64 `json:"share_value"`
	DeductionTypeID    uint64  `json:"deduction_type_id"`
	Priority           int     `json:"priority"`
	IsPayrollDeduction int     `json:"is_payroll_deduction"`
	EarnsDividend      int     `json:"earns_dividend"`
	IsTransferable     int     `json:"is_transferable"`
	MinimumShares      float64 `json:"minimum_shares"`
	MaxmumShares       float64 `json:"maxmum_shares"`
}

type ShareTypeResponse struct {
	ID                 uint64    `json:"id"`
	ShareCode          string    `json:"share_code"`
	ShareType          string    `json:"share_type"`
	Description        string    `json:"description"`
	Rate               float64   `json:"rate"`
	Mandatory          int       `json:"mandatory"`
	HasShareValue      int       `json:"has_share_value"`
	CalculatingMethod  string    `json:"calculating_method"`
	ShareValue         float64   `json:"share_value"`
	DeductionTypeID    uint64    `json:"deduction_type_id"`
	DeductionTypeName  string    `json:"deduction_type_name"`
	Priority           int       `json:"priority"`
	IsPayrollDeduction bool      `json:"is_payroll_deduction"`
	EarnsDividend      bool      `json:"earns_dividend"`
	IsTransferable     bool      `json:"is_transferable"`
	MinimumShares      float64   `json:"minimum_shares"`
	MaxmumShares       float64   `json:"maxmum_shares"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
