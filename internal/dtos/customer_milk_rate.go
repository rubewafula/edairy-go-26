package dtos

import "time"

type CreateCustomerMilkRateRequest struct {
	CustomerID   uint64  `json:"customer_id" validate:"required"`
	Rate         float64 `json:"rate" validate:"required,min=0"`
	GradeID      uint64  `json:"grade_id" validate:"required"`
	PayDateRange uint64  `json:"customer_pay_date_range_id" validate:"required"`
}

type UpdateCustomerMilkRateRequest struct {
	CustomerID   uint64  `json:"customer_id" validate:"required"`
	Rate         float64 `json:"rate" validate:"required,min=0"`
	GradeID      uint64  `json:"grade_id" validate:"required"`
	PayDateRange uint64  `json:"customer_pay_date_range_id" validate:"required"`
}

type CustomerMilkRateResponse struct {
	ID               uint64    `json:"id"`
	CustomerID       uint64    `json:"customer_id"`
	CustomerName     string    `json:"customer_name"`
	Rate             float64   `json:"rate"`
	GradeID          uint64    `json:"grade_id"`
	GradeName        string    `json:"grade_name"`
	PayDateRange     uint64    `json:"pay_date_range_id"`
	PayDateRangeName string    `json:"pay_date_range_name"`
	CreatedBy        uint64    `json:"created_by"`
	UpdatedBy        uint64    `json:"updated_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
