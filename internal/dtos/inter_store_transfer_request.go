package dtos

import "time"

type CreateInterStoreTransferRequest struct {
	FromStoreID  uint64 `json:"FromStoreID" validate:"required"`
	ToStoreID    uint64 `json:"ToStoreID" validate:"required"`
	Reference    string `json:"Reference" validate:"required,max=255"`
	TransferDate string `json:"TransferDate" validate:"required,datetime"`
	Status       string `json:"Status" validate:"required,max=255"`
}

type UpdateInterStoreTransferRequest struct {
	FromStoreID  uint64 `json:"FromStoreID" validate:"required"`
	ToStoreID    uint64 `json:"ToStoreID" validate:"required"`
	Reference    string `json:"Reference" validate:"required,max=255"`
	TransferDate string `json:"TransferDate" validate:"required,datetime"`
	Status       string `json:"Status" validate:"required,max=255"`
}

type InterStoreTransferResponse struct {
	ID            uint64    `json:"ID"`
	FromStoreID   uint64    `json:"FromStoreID"`
	FromStoreName string    `json:"FromStoreName"`
	ToStoreID     uint64    `json:"ToStoreID"`
	ToStoreName   string    `json:"ToStoreName"`
	Reference     string    `json:"Reference"`
	TransferDate  time.Time `json:"TransferDate"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
	Status        string    `json:"Status"`
}

type CreateInterStoreTransferItemRequest struct {
	TransferID uint64  `json:"TransferID" validate:"required"`
	ItemID     uint64  `json:"ItemID" validate:"required"`
	Quantity   float64 `json:"Quantity" validate:"required"`
	StockID    uint64  `json:"StockID"`
}

type UpdateInterStoreTransferItemRequest struct {
	TransferID uint64  `json:"TransferID" validate:"required"`
	ItemID     uint64  `json:"ItemID" validate:"required"`
	Quantity   float64 `json:"Quantity" validate:"required"`
	StockID    uint64  `json:"StockID"`
}

type InterStoreTransferItemResponse struct {
	ID         uint64    `json:"ID"`
	TransferID uint64    `json:"TransferID"`
	Reference  string    `json:"SaleReference"`
	ItemID     uint64    `json:"ItemID"`
	ItemName   string    `json:"ItemName"`
	Quantity   float64   `json:"Quantity"`
	StockID    uint64    `json:"StockID"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
}
