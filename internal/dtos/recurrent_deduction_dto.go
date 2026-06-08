package dtos

import (
	"time"
)

// RecurrentDeductionResponse defines the structure for a recurrent deduction response with enriched details.
type RecurrentDeductionResponse struct {
	ID                uint64     `json:"id"`
	CustomerID        uint64     `json:"customer_id"`
	CustomerType      string     `json:"customer_type"`
	MemberNo          string     `json:"member_no,omitempty"`
	Names             string     `json:"names,omitempty"`
	TotalAmount       float64    `json:"total_amount"`
	PaidAmount        float64    `json:"paid_amount"`
	RecurrentAmount   float64    `json:"recurrent_amount"`
	PrincipalAmount   float64    `json:"principal_amount"`
	DeductionTypeID   uint64     `json:"deduction_type_id"`
	DeductionTypeName string     `json:"deduction_type_name"`
	Reference         string     `json:"reference"`
	Settled           int        `json:"settled"`
	TransactionDate   *time.Time `json:"transaction_date"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
