package dtos

import "time"

type CreateCustomerPayDateRangeRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	StartDate string `json:"start_date" validate:"required,datetime"`
	EndDate   string `json:"end_date" validate:"required,datetime"`
}

type UpdateCustomerPayDateRangeRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	StartDate string `json:"start_date" validate:"required,datetime"`
	EndDate   string `json:"end_date" validate:"required,datetime"`
}

type CustomerPayDateRangeResponse struct {
	ID        uint64    `json:"ID"`
	Name      string    `json:"Name"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	CreatedBy uint64    `json:"CreatedBy"`
	UpdatedBy uint64    `json:"UpdatedBy"`
}
