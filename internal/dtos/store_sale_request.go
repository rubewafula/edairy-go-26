package dtos

import "time"

type CreateStoreSaleRequest struct {
	TotalAmount   float64 `json:"TotalAmount" validate:"required"`
	AmountPaid    float64 `json:"AmountPaid"`
	AmountDue     float64 `json:"AmountDue"`
	Reference     string  `json:"Reference"`
	StoreID       uint64  `json:"StoreID" validate:"required"`
	SaleType      string  `json:"SaleType" validate:"required"`
	CustomerID    uint64  `json:"CustomerID"`
	CustomerType  string  `json:"CustomerType"`
	TransactionID int64   `json:"TransactionID"`
}

type UpdateStoreSaleRequest struct {
	TotalAmount   float64 `json:"TotalAmount" validate:"required"`
	AmountPaid    float64 `json:"AmountPaid"`
	AmountDue     float64 `json:"AmountDue"`
	Reference     string  `json:"Reference"`
	StoreID       uint64  `json:"StoreID" validate:"required"`
	SaleType      string  `json:"SaleType" validate:"required"`
	CustomerID    uint64  `json:"CustomerID"`
	CustomerType  string  `json:"CustomerType"`
	TransactionID int64   `json:"TransactionID"`
}

type StoreSaleResponse struct {
	ID            uint64    `json:"ID"`
	TotalAmount   float64   `json:"TotalAmount"`
	AmountPaid    float64   `json:"AmountPaid"`
	AmountDue     float64   `json:"AmountDue"`
	Reference     string    `json:"Reference"`
	StoreID       uint64    `json:"StoreID"`
	StoreName     string    `json:"StoreName"`
	SaleType      string    `json:"SaleType"`
	CustomerID    uint64    `json:"CustomerID"`
	CustomerType  string    `json:"CustomerType"`
	TransactionID int64     `json:"TransactionID"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
