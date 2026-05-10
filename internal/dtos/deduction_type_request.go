package dtos

import "time"

type CreateDeductionTypeRequest struct {
	Code        string `json:"Code" validate:"max=255"`
	Description string `json:"Description" validate:"required,max=255"`
	Status      string `json:"Status" validate:"max=255"`
	IsStatutory string `json:"IsStatutory" validate:"max=255"`
}

type UpdateDeductionTypeRequest struct {
	Code        string `json:"Code" validate:"max=255"`
	Description string `json:"Description" validate:"required,max=255"`
	Status      string `json:"Status" validate:"max=255"`
	IsStatutory string `json:"IsStatutory" validate:"max=255"`
}

type DeductionTypeResponse struct {
	ID          uint64    `json:"ID"`
	Code        string    `json:"Code"`
	Description string    `json:"Description"`
	Status      string    `json:"Status"`
	IsStatutory string    `json:"IsStatutory"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
