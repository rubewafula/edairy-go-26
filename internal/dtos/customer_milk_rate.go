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
	ID               uint64    `json:"ID"`
	CustomerID       uint64    `json:"CustomerID"`
	CustomerName     string    `json:"CustomerName"`
	Rate             float64   `json:"Rate"`
	GradeID          uint64    `json:"GradeID"`
	GradeName        string    `json:"GradeName"`
	PayDateRange     uint64    `json:"PayDateRangeID" gorm:"column:customer_pay_date_range_id"`
	PayDateRangeName string    `json:"PayDateRangeName"`
	CreatedBy        uint64    `json:"CreatedBy"`
	UpdatedBy        uint64    `json:"UpdatedBy"`
	CreatedAt        time.Time `json:"CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt"`
}
