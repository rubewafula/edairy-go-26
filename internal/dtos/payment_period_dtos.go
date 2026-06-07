package dtos

import "time"

type CreatePaymentPeriodRequest struct {
	Name          string `json:"name" validate:"required,oneof=WEEKLY BI-WEEKLY MONTHLY"`
	Description   string `json:"description" validate:"required"`
	DefaultPeriod int    `json:"default_period"`
}

type UpdatePaymentPeriodRequest struct {
	Name          string `json:"name" validate:"omitempty,oneof=WEEKLY BI-WEEKLY MONTHLY"`
	Description   string `json:"description"`
	DefaultPeriod int    `json:"default_period"`
}

type PaymentPeriodResponse struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	DefaultPeriod int       `json:"default_period"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     uint64    `json:"created_by"`
	UpdatedBy     uint64    `json:"updated_by"`
}
