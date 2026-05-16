package dtos

import "time"

type CreateDeductionTypeRequest struct {
	Code        string `json:"code" validate:"max=255"`
	Description string `json:"description" validate:"required,max=255"`
	Status      string `json:"status" validate:"max=255"`
	IsStatutory string `json:"is_statutory" validate:"max=255"`
}

type UpdateDeductionTypeRequest struct {
	Code        string `json:"code" validate:"max=255"`
	Description string `json:"description" validate:"required,max=255"`
	Status      string `json:"status" validate:"max=255"`
	IsStatutory string `json:"is_statutory" validate:"max=255"`
}

type DeductionTypeResponse struct {
	ID          uint64    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	IsStatutory string    `json:"is_statutory"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
