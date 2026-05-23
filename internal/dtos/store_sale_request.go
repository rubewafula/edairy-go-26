package dtos

import "time"

type CreateStoreSaleRequest struct {
	TotalAmount     float64                      `json:"total_amount" validate:"required"`
	AmountPaid      float64                      `json:"amount_paid"`
	AmountDue       float64                      `json:"amount_due"`
	Reference       string                       `json:"reference"`
	StoreID         uint64                       `json:"store_id" validate:"required"`
	SaleType        string                       `json:"sale_type" validate:"required"`
	CustomerID      uint64                       `json:"customer_id"`
	CustomerType    string                       `json:"customer_type"`
	TransactionID   int64                        `json:"transaction_id"`
	TransactionDate string                       `json:"transaction_date" validate:"required"`
	Items           []CreateStoreSaleItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateStoreSaleRequest struct {
	TotalAmount   float64 `json:"total_amount" validate:"required"`
	AmountPaid    float64 `json:"amount_paid"`
	AmountDue     float64 `json:"amount_due"`
	Reference     string  `json:"reference"`
	StoreID       uint64  `json:"store_id" validate:"required"`
	SaleType      string  `json:"sale_type" validate:"required"`
	CustomerID    uint64  `json:"customer_id"`
	CustomerType  string  `json:"customer_type"`
	TransactionID int64   `json:"transaction_id"`
}
type StoreSaleResponse struct {
	ID            uint64    `json:"id"`
	TotalAmount   float64   `json:"total_amount"`
	AmountPaid    float64   `json:"amount_paid"`
	AmountDue     float64   `json:"amount_due"`
	Reference     string    `json:"reference"`
	StoreID       uint64    `json:"store_id"`
	StoreName     string    `json:"store_name"`
	SaleType      string    `json:"sale_type"`
	CustomerID    uint64    `json:"customer_id"`
	CustomerType  string    `json:"customer_type"`
	TransactionID int64     `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
