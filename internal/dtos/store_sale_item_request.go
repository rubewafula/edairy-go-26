package dtos

import "time"

type CreateStoreSaleItemRequest struct {
	ItemID      uint64  `json:"item_id" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
	Total       float64 `json:"total" validate:"required"`
	StoreSaleID uint64  `json:"store_sale_id" validate:"required"`
}

type UpdateStoreSaleItemRequest struct {
	ItemID      uint64  `json:"item_id" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
	Total       float64 `json:"total" validate:"required"`
	StoreSaleID uint64  `json:"store_sale_id" validate:"required"`
}
type StoreSaleItemResponse struct {
	ID            uint64    `json:"id"`
	ItemID        uint64    `json:"item_id"`
	ItemName      string    `json:"item_name"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	Total         float64   `json:"total"`
	StoreSaleID   uint64    `json:"store_sale_id"`
	SaleReference string    `json:"sale_reference"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
