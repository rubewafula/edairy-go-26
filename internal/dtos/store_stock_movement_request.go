package dtos

import "time"

type CreateStoreStockMovementItemRequest struct {
	ItemID       uint64  `json:"item_id" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	UnitCost     float64 `json:"unit_cost"`
	SellingPrice float64 `json:"selling_price"`
	Remarks      string  `json:"remarks" validate:"max=255"`
}

type CreateStoreStockMovementRequest struct {
	TransactionDate string                                `json:"transaction_date" validate:"required"`
	StoreID         uint64                                `json:"store_id" validate:"required"`
	MovementTypeID  uint64                                `json:"movement_type_id" validate:"required,max=50"`
	Remarks         string                                `json:"remarks" validate:"max=255"`
	Items           []CreateStoreStockMovementItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateStoreStockMovementRequest struct {
	TransactionDate string                                `json:"transaction_date" validate:"required,datetime"`
	StoreID         uint64                                `json:"store_id" validate:"required"`
	MovementType    string                                `json:"movement_type_id" validate:"required,max=50"`
	ReferenceTable  string                                `json:"reference_table" validate:"max=100"`
	ReferenceID     uint64                                `json:"reference_id"`
	Remarks         string                                `json:"remarks" validate:"max=255"`
	Items           []CreateStoreStockMovementItemRequest `json:"items" validate:"required,min=1"`
}

type StoreStockMovementItemResponse struct {
	ID           uint64    `json:"id"`
	MovementID   uint64    `json:"movement_id"`
	ItemID       uint64    `json:"item_id"`
	ItemName     string    `json:"item_name"`
	Quantity     float64   `json:"quantity"`
	UnitCost     float64   `json:"unit_cost"`
	SellingPrice float64   `json:"selling_price"`
	BalanceAfter float64   `json:"balance_after"`
	Remarks      string    `json:"remarks"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StoreStockMovementResponse struct {
	ID              uint64                           `json:"id"`
	TransactionDate time.Time                        `json:"transaction_date"`
	StoreID         uint64                           `json:"store_id"`
	StoreName       string                           `json:"store_name"`
	MovementType    string                           `json:"movement_type"`
	Remarks         string                           `json:"remarks"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       time.Time                        `json:"updated_at"`
	Items           []StoreStockMovementItemResponse `json:"items" gorm:"-"`
}
