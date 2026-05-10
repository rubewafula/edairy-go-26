package dtos

import "time"

type CreateDeductionPricingRuleRequest struct {
	DeductionTypeID uint64  `json:"DeductionTypeID" validate:"required"`
	MinCreditLimit  float64 `json:"MinCreditLimit"`
	MaxLimit        float64 `json:"MaxLimit"`
	BoardingFee     float64 `json:"BoardingFee"`
	ProcessingFee   float64 `json:"ProcessingFee"`
	InsuranceFee    float64 `json:"InsuranceFee"`
	LegalFee        float64 `json:"LegalFee"`
	InterestRate    float64 `json:"InterestRate"`
	Status          string  `json:"Status"`
}

type UpdateDeductionPricingRuleRequest struct {
	DeductionTypeID uint64  `json:"DeductionTypeID" validate:"required"`
	MinCreditLimit  float64 `json:"MinCreditLimit"`
	MaxLimit        float64 `json:"MaxLimit"`
	BoardingFee     float64 `json:"BoardingFee"`
	ProcessingFee   float64 `json:"ProcessingFee"`
	InsuranceFee    float64 `json:"InsuranceFee"`
	LegalFee        float64 `json:"LegalFee"`
	InterestRate    float64 `json:"InterestRate"`
	Status          string  `json:"Status"`
}

type DeductionPricingRuleResponse struct {
	ID                uint64    `json:"ID"`
	DeductionTypeID   uint64    `json:"DeductionTypeID"`
	DeductionTypeName string    `json:"DeductionTypeName"`
	MinCreditLimit    float64   `json:"MinCreditLimit"`
	MaxLimit          float64   `json:"MaxLimit"`
	BoardingFee       float64   `json:"BoardingFee"`
	ProcessingFee     float64   `json:"ProcessingFee"`
	InsuranceFee      float64   `json:"InsuranceFee"`
	LegalFee          float64   `json:"LegalFee"`
	InterestRate      float64   `json:"InterestRate"`
	Status            string    `json:"Status"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
}
