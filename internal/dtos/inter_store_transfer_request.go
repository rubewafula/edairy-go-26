package dtos

import "time"

type CreateInterStoreTransferRequest struct {
	FromStoreID  uint64 `json:"from_store_id" validate:"required"`
	ToStoreID    uint64 `json:"to_store_id" validate:"required"`
	Reference    string `json:"reference" validate:"required,max=255"`
	TransferDate string `json:"transfer_date" validate:"required"`
	Status       string `json:"status" validate:"required,max=255"`
}

type UpdateInterStoreTransferRequest struct {
	FromStoreID  uint64 `json:"FromStoreID" validate:"required"`
	ToStoreID    uint64 `json:"ToStoreID" validate:"required"`
	Reference    string `json:"Reference" validate:"required,max=255"`
	TransferDate string `json:"TransferDate" validate:"required,datetime"`
	Status       string `json:"Status" validate:"required,max=255"`
}

type InterStoreTransferResponse struct {
	ID            uint64    `json:"id"`
	FromStoreID   uint64    `json:"from_store_id"`
	FromStoreName string    `json:"from_store_name"`
	ToStoreID     uint64    `json:"to_store_id"`
	ToStoreName   string    `json:"to_store_name"`
	Reference     string    `json:"reference"`
	TransferDate  time.Time `json:"transfer_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"`
}

type CreateInterStoreTransferItemRequest struct {
	TransferID uint64  `json:"transfer_id" validate:"required"`
	ItemID     uint64  `json:"item_id" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required"`
	StockID    uint64  `json:"stock_id"`
}

type UpdateInterStoreTransferItemRequest struct {
	TransferID uint64  `json:"transfer_id" validate:"required"`
	ItemID     uint64  `json:"item_id" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required"`
	StockID    uint64  `json:"stock_id"`
}

type InterStoreTransferItemResponse struct {
	ID         uint64    `json:"id"`
	TransferID uint64    `json:"transfer_id"`
	Reference  string    `json:"sale_reference"`
	ItemID     uint64    `json:"item_id"`
	ItemName   string    `json:"item_name"`
	Quantity   float64   `json:"quantity"`
	StockID    uint64    `json:"stock_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
