package dtos

import "time"

type CreateProductGradeRequest struct {
	Name        string `json:"name" validate:"required,max=45"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateProductGradeRequest struct {
	Name        string `json:"name" validate:"required,max=45"`
	Description string `json:"description" validate:"max=255"`
}

type ProductGradeResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
