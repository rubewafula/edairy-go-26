package dtos

import "time"

type CreateCustomerPayDateRangeRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type UpdateCustomerPayDateRangeRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type CustomerPayDateRangeResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint64    `json:"created_by"`
	UpdatedBy uint64    `json:"updated_by"`
}
