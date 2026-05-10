package dtos

import "time"

type CreateProductGradeRequest struct {
	Name        string `json:"Name" validate:"required,max=45"`
	Description string `json:"Description" validate:"max=255"`
}

type UpdateProductGradeRequest struct {
	Name        string `json:"Name" validate:"required,max=45"`
	Description string `json:"Description" validate:"max=255"`
}

type ProductGradeResponse struct {
	ID          uint64    `json:"ID"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
