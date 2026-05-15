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
	ID            uint64    `json:"ID"`
	ReceiptNumber string    `json:"ReceiptNumber"`
	CustomerID    uint64    `json:"CustomerID"`
	CustomerName  string    `json:"CustomerName"`
	InvoiceID     *uint64   `json:"InvoiceID"`
	InvoiceNo     string    `json:"InvoiceNo"`
	PaymentDate   time.Time `json:"PaymentDate"`
	Amount        float64   `json:"Amount"`
	PaymentMethod string    `json:"PaymentMethod"`
	ReferenceNo   string    `json:"ReferenceNo"`
	Notes         string    `json:"Notes"`
	CreatedBy     uint64    `json:"CreatedBy"`
	CreatedAt     time.Time `json:"CreatedAt"`
}

type CustomerPaymentAllocationResponse struct {
	ID                uint64    `json:"ID"`
	InvoiceID         uint64    `json:"InvoiceID"`
	InvoiceNo         string    `json:"InvoiceNo"`
	CustomerPaymentID uint64    `json:"CustomerPaymentID"`
	ReceiptNumber     string    `json:"ReceiptNumber"`
	AllocatedAmount   float64   `json:"AllocatedAmount"`
	CreatedAt         time.Time `json:"CreatedAt"`
	CreatedBy         uint64    `json:"CreatedBy"`
}
