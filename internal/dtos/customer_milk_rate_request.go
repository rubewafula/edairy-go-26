package dtos

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
