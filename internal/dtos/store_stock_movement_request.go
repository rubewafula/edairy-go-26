package dtos

import "time"

type CreateStoreStockMovementRequest struct {
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	StoreID         uint64  `json:"store_id" validate:"required"`
	ItemID          uint64  `json:"item_id" validate:"required"`
	MovementType    string  `json:"movement_type" validate:"required,max=50"`
	ReferenceTable  string  `json:"reference_table" validate:"max=100"`
	ReferenceID     uint64  `json:"reference_id"`
	QtyIn           float64 `json:"qty_in"`
	QtyOut          float64 `json:"qty_out"`
	BalanceAfter    float64 `json:"balance_after" validate:"required"`
	UnitCost        float64 `json:"unit_cost"`
	SellingPrice    float64 `json:"selling_price"`
	Remarks         string  `json:"remarks" validate:"max=255"`
}

type UpdateStoreStockMovementRequest struct {
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	StoreID         uint64  `json:"store_id" validate:"required"`
	ItemID          uint64  `json:"item_id" validate:"required"`
	MovementType    string  `json:"movement_type" validate:"required,max=50"`
	ReferenceTable  string  `json:"reference_table" validate:"max=100"`
	ReferenceID     uint64  `json:"reference_id"`
	QtyIn           float64 `json:"qty_in"`
	QtyOut          float64 `json:"qty_out"`
	BalanceAfter    float64 `json:"balance_after" validate:"required"`
	UnitCost        float64 `json:"unit_cost"`
	SellingPrice    float64 `json:"selling_price"`
	Remarks         string  `json:"remarks" validate:"max=255"`
}

type StoreStockMovementResponse struct {
	ID              uint64    `json:"id"`
	TransactionDate time.Time `json:"transaction_date"`
	StoreID         uint64    `json:"store_id"`
	StoreName       string    `json:"store_name"`
	ItemID          uint64    `json:"item_id"`
	ItemName        string    `json:"item_name"`
	MovementType    string    `json:"movement_type"`
	ReferenceTable  string    `json:"reference_table"`
	ReferenceID     uint64    `json:"reference_id"`
	QtyIn           float64   `json:"qty_in"`
	QtyOut          float64   `json:"qty_out"`
	BalanceAfter    float64   `json:"balance_after"`
	UnitCost        float64   `json:"unit_cost"`
	SellingPrice    float64   `json:"selling_price"`
	Remarks         string    `json:"remarks"`
	CreatedAt       time.Time `json:"created_at"`
}
