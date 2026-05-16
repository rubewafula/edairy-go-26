package dtos

import "time"

type CreateCustomerPaymentRequest struct {
	InvoiceID     *uint64 `json:"invoice_id"`
	CustomerID    uint64  `json:"customer_id" validate:"required"`
	PaymentDate   string  `json:"payment_date" validate:"required,datetime"`
	Amount        float64 `json:"amount" validate:"required,min=1"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	ReferenceNo   string  `json:"reference_no" validate:"required"`
	Notes         string  `json:"notes"`
}

type CustomerPaymentResponse struct {
	ID            uint64    `json:"id"`
	ReceiptNumber string    `json:"receipt_number"`
	CustomerID    uint64    `json:"customer_id"`
	CustomerName  string    `json:"customer_name"`
	InvoiceID     *uint64   `json:"invoice_id"`
	InvoiceNo     string    `json:"invoice_no"`
	PaymentDate   time.Time `json:"payment_date"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	ReferenceNo   string    `json:"reference_no"`
	Notes         string    `json:"notes"`
	CreatedBy     uint64    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type CustomerPaymentAllocationResponse struct {
	ID                uint64    `json:"id"`
	InvoiceID         uint64    `json:"invoice_id"`
	InvoiceNo         string    `json:"invoice_no"`
	CustomerPaymentID uint64    `json:"customer_payment_id"`
	ReceiptNumber     string    `json:"receipt_number"`
	AllocatedAmount   float64   `json:"allocated_amount"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         uint64    `json:"created_by"`
}
