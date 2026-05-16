package dtos

import "time"

type CreateDeductionPricingRuleRequest struct {
	DeductionTypeID uint64  `json:"deduction_type_id" validate:"required"`
	MinCreditLimit  float64 `json:"min_credit_limit"`
	MaxLimit        float64 `json:"max_limit"`
	BoardingFee     float64 `json:"boarding_fee"`
	ProcessingFee   float64 `json:"processing_fee"`
	InsuranceFee    float64 `json:"insurance_fee"`
	LegalFee        float64 `json:"legal_fee"`
	InterestRate    float64 `json:"interest_rate"`
	Status          string  `json:"status"`
}

type UpdateDeductionPricingRuleRequest struct {
	DeductionTypeID uint64  `json:"deduction_type_id" validate:"required"`
	MinCreditLimit  float64 `json:"min_credit_limit"`
	MaxLimit        float64 `json:"max_limit"`
	BoardingFee     float64 `json:"boarding_fee"`
	ProcessingFee   float64 `json:"processing_fee"`
	InsuranceFee    float64 `json:"insurance_fee"`
	LegalFee        float64 `json:"legal_fee"`
	InterestRate    float64 `json:"interest_rate"`
	Status          string  `json:"status"`
}

type DeductionPricingRuleResponse struct {
	ID                uint64    `json:"id"`
	DeductionTypeID   uint64    `json:"deduction_type_id"`
	DeductionTypeName string    `json:"deduction_type_name"`
	MinCreditLimit    float64   `json:"min_credit_limit"`
	MaxLimit          float64   `json:"max_limit"`
	BoardingFee       float64   `json:"boarding_fee"`
	ProcessingFee     float64   `json:"processing_fee"`
	InsuranceFee      float64   `json:"insurance_fee"`
	LegalFee          float64   `json:"legal_fee"`
	InterestRate      float64   `json:"interest_rate"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
