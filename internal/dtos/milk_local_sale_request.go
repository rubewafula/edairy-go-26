package dtos

import "time"

type CreateMilkLocalSaleRequest struct {
	Quantity        float64 `json:"Quantity" validate:"required"`
	Rate            float64 `json:"Rate" validate:"required"`
	GradeID         uint64  `json:"GradeID" validate:"required"`
	RefNumber       string  `json:"RefNumber" validate:"required"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	TransporterID   uint64  `json:"TransporterID"`
	Amount          float64 `json:"Amount" validate:"required"`
}

type UpdateMilkLocalSaleRequest struct {
	Quantity        float64 `json:"Quantity" validate:"required"`
	Rate            float64 `json:"Rate" validate:"required"`
	GradeID         uint64  `json:"GradeID" validate:"required"`
	RefNumber       string  `json:"RefNumber" validate:"required"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	TransporterID   uint64  `json:"TransporterID"`
	Amount          float64 `json:"Amount" validate:"required"`
}

type MilkLocalSaleResponse struct {
	ID              uint64    `json:"ID"`
	Quantity        float64   `json:"Quantity"`
	Rate            float64   `json:"Rate"`
	GradeID         uint64    `json:"GradeID"`
	GradeName       string    `json:"GradeName"`
	RefNumber       string    `json:"RefNumber"`
	TransactionDate time.Time `json:"TransactionDate"`
	TransporterID   uint64    `json:"TransporterID"`
	TransporterName string    `json:"TransporterName"`
	Amount          float64   `json:"Amount"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
