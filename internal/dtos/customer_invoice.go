package dtos

import "time"

type CreateCustomerInvoiceRequest struct {
	CustomerID  uint64  `json:"customer_id" validate:"required"`
	BillingID   uint64  `json:"billing_id" validate:"required"`
	InvoiceDate string  `json:"invoice_date" validate:"required,datetime"`
	DueDate     string  `json:"due_date" validate:"required,datetime"`
	GrossAmount float64 `json:"gross_amount" validate:"required"`
	TaxAmount   float64 `json:"tax_amount"`
	TotalAmount float64 `json:"total_amount" validate:"required"`
}

type CustomerInvoiceResponse struct {
	ID           uint64    `json:"ID"`
	InvoiceNo    string    `json:"InvoiceNo"`
	CustomerID   uint64    `json:"CustomerID"`
	CustomerName string    `json:"CustomerName"`
	BillingID    uint64    `json:"BillingID"`
	InvoiceDate  time.Time `json:"InvoiceDate"`
	DueDate      time.Time `json:"DueDate"`
	GrossAmount  float64   `json:"GrossAmount"`
	TaxAmount    float64   `json:"TaxAmount"`
	TotalAmount  float64   `json:"TotalAmount"`
	Balance      float64   `json:"Balance"`
	Status       string    `json:"Status"`
	CreatedBy    uint64    `json:"CreatedBy"`
	UpdatedBy    uint64    `json:"UpdatedBy"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}
