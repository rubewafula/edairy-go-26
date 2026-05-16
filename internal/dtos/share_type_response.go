package dtos

import "time"

type ShareTypeResponse struct {
	ID                uint64    `json:"id"`
	ShareCode         string    `json:"share_code"`
	ShareType         string    `json:"share_type"`
	Description       string    `json:"description"`
	Rate              float64   `json:"rate"`
	Mandatory         int       `json:"mandatory"`
	HasShareValue     string    `json:"has_share_value"`
	RepayMethod       string    `json:"repay_method"`
	CalculatingMethod string    `json:"calculating_method"`
	ShareValue        float64   `json:"share_value"`
	DeductionTypeID   uint64    `json:"deduction_type_id"`
	DeductionTypeName string    `json:"deduction_type_name"`
	Priority          int       `json:"priority"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
