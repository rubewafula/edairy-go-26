package dtos

import "time"

type CreateCustomerBillingRequest struct {
	PayDateRangeID  uint64  `json:"pay_date_range_id" validate:"required"`
	TotalDeliveries float64 `json:"total_deliveries" validate:"required"`
	TotalAmount     float64 `json:"total_amount" validate:"required"`
}

type CustomerBillingResponse struct {
	ID               uint64    `json:"id"`
	PayDateRangeID   uint64    `json:"pay_date_range_id"`
	PayDateRangeName string    `json:"pay_date_range_name"`
	TotalDeliveries  float64   `json:"total_deliveries"`
	TotalAmount      float64   `json:"total_amount"`
	Status           string    `json:"status"`
	InvoiceID        *uint64   `json:"invoice_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CustomerBillingItemResponse struct {
	ID                uint64    `json:"id"`
	CustomerBillingID uint64    `json:"customer_billing_id"`
	ProductGradeID    uint64    `json:"product_grade_id"`
	GradeName         string    `json:"grade_name"`
	TotalQuantity     float64   `json:"total_quantity"`
	UnitPrice         float64   `json:"unit_price"`
	TotalAmount       float64   `json:"total_amount"`
	CreatedAt         time.Time `json:"created_at"`
}
