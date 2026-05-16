package dtos

import "time"

type CreateMilkLocalSaleRequest struct {
	Quantity        float64 `json:"quantity" validate:"required"`
	Rate            float64 `json:"rate" validate:"required"`
	GradeID         uint64  `json:"grade_id" validate:"required"`
	RefNumber       string  `json:"ref_number" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	TransporterID   uint64  `json:"transporter_id"`
	Amount          float64 `json:"amount" validate:"required"`
}

type UpdateMilkLocalSaleRequest struct {
	Quantity        float64 `json:"Quantity" validate:"required"`
	Rate            float64 `json:"rate" validate:"required"`
	GradeID         uint64  `json:"grade_id" validate:"required"`
	RefNumber       string  `json:"ref_number" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	TransporterID   uint64  `json:"transporter_id"`
	Amount          float64 `json:"amount" validate:"required"`
}

type MilkLocalSaleResponse struct {
	ID              uint64    `json:"id"`
	Quantity        float64   `json:"quantity"`
	Rate            float64   `json:"rate"`
	GradeID         uint64    `json:"grade_id"`
	GradeName       string    `json:"grade_name"`
	RefNumber       string    `json:"ref_number"`
	TransactionDate time.Time `json:"transaction_date"`
	TransporterID   uint64    `json:"transporter_id"`
	TransporterName string    `json:"transporter_name"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
