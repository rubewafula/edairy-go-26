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
	ID           uint64    `json:"id"`
	InvoiceNo    string    `json:"invoice_no"`
	CustomerID   uint64    `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	BillingID    uint64    `json:"billing_id"`
	InvoiceDate  time.Time `json:"invoice_date"`
	DueDate      time.Time `json:"due_date"`
	GrossAmount  float64   `json:"gross_amount"`
	TaxAmount    float64   `json:"tax_amount"`
	TotalAmount  float64   `json:"total_amount"`
	Balance      float64   `json:"balance"`
	Status       string    `json:"status"`
	CreatedBy    uint64    `json:"created_by"`
	UpdatedBy    uint64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
