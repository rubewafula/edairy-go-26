package dtos

import "time"

type CreateStoreSaleItemRequest struct {
	ItemID      uint64  `json:"ItemID" validate:"required"`
	Quantity    int     `json:"Quantity" validate:"required"`
	UnitPrice   float64 `json:"UnitPrice" validate:"required"`
	Total       float64 `json:"Total" validate:"required"`
	StoreSaleID uint64  `json:"StoreSaleID" validate:"required"`
}

type UpdateStoreSaleItemRequest struct {
	ItemID      uint64  `json:"ItemID" validate:"required"`
	Quantity    int     `json:"Quantity" validate:"required"`
	UnitPrice   float64 `json:"UnitPrice" validate:"required"`
	Total       float64 `json:"Total" validate:"required"`
	StoreSaleID uint64  `json:"StoreSaleID" validate:"required"`
}

type StoreSaleItemResponse struct {
	ID            uint64    `json:"ID"`
	ItemID        uint64    `json:"ItemID"`
	ItemName      string    `json:"ItemName"`
	Quantity      int       `json:"Quantity"`
	UnitPrice     float64   `json:"UnitPrice"`
	Total         float64   `json:"Total"`
	StoreSaleID   uint64    `json:"StoreSaleID"`
	SaleReference string    `json:"SaleReference"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
