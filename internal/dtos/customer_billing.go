package dtos

import "time"

type CreateCustomerBillingRequest struct {
	PayDateRangeID  uint64  `json:"pay_date_range_id" validate:"required"`
	TotalDeliveries float64 `json:"total_deliveries" validate:"required"`
	TotalAmount     float64 `json:"total_amount" validate:"required"`
}

type CustomerBillingResponse struct {
	ID               uint64    `json:"ID"`
	PayDateRangeID   uint64    `json:"PayDateRangeID"`
	PayDateRangeName string    `json:"PayDateRangeName"`
	TotalDeliveries  float64   `json:"TotalDeliveries"`
	TotalAmount      float64   `json:"TotalAmount"`
	Status           string    `json:"Status"`
	InvoiceID        *uint64   `json:"InvoiceID"`
	CreatedAt        time.Time `json:"CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt"`
}

type CustomerBillingItemResponse struct {
	ID                uint64    `json:"ID"`
	CustomerBillingID uint64    `json:"CustomerBillingID"`
	ProductGradeID    uint64    `json:"ProductGradeID"`
	GradeName         string    `json:"GradeName"`
	TotalQuantity     float64   `json:"TotalQuantity"`
	UnitPrice         float64   `json:"UnitPrice"`
	TotalAmount       float64   `json:"TotalAmount"`
	CreatedAt         time.Time `json:"CreatedAt"`
}
