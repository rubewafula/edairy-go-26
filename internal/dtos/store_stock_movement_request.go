package dtos

import "time"

type CreateStoreStockMovementRequest struct {
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	StoreID         uint64  `json:"StoreID" validate:"required"`
	ItemID          uint64  `json:"ItemID" validate:"required"`
	MovementType    string  `json:"MovementType" validate:"required,max=50"`
	ReferenceTable  string  `json:"ReferenceTable" validate:"max=100"`
	ReferenceID     uint64  `json:"ReferenceID"`
	QtyIn           float64 `json:"QtyIn"`
	QtyOut          float64 `json:"QtyOut"`
	BalanceAfter    float64 `json:"BalanceAfter" validate:"required"`
	UnitCost        float64 `json:"UnitCost"`
	SellingPrice    float64 `json:"SellingPrice"`
	Remarks         string  `json:"Remarks" validate:"max=255"`
}

type UpdateStoreStockMovementRequest struct {
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	StoreID         uint64  `json:"StoreID" validate:"required"`
	ItemID          uint64  `json:"ItemID" validate:"required"`
	MovementType    string  `json:"MovementType" validate:"required,max=50"`
	ReferenceTable  string  `json:"ReferenceTable" validate:"max=100"`
	ReferenceID     uint64  `json:"ReferenceID"`
	QtyIn           float64 `json:"QtyIn"`
	QtyOut          float64 `json:"QtyOut"`
	BalanceAfter    float64 `json:"BalanceAfter" validate:"required"`
	UnitCost        float64 `json:"UnitCost"`
	SellingPrice    float64 `json:"SellingPrice"`
	Remarks         string  `json:"Remarks" validate:"max=255"`
}

type StoreStockMovementResponse struct {
	ID              uint64    `json:"ID"`
	TransactionDate time.Time `json:"TransactionDate"`
	StoreID         uint64    `json:"StoreID"`
	StoreName       string    `json:"StoreName"`
	ItemID          uint64    `json:"ItemID"`
	ItemName        string    `json:"ItemName"`
	MovementType    string    `json:"MovementType"`
	ReferenceTable  string    `json:"ReferenceTable"`
	ReferenceID     uint64    `json:"ReferenceID"`
	QtyIn           float64   `json:"QtyIn"`
	QtyOut          float64   `json:"QtyOut"`
	BalanceAfter    float64   `json:"BalanceAfter"`
	UnitCost        float64   `json:"UnitCost"`
	SellingPrice    float64   `json:"SellingPrice"`
	Remarks         string    `json:"Remarks"`
	CreatedAt       time.Time `json:"CreatedAt"`
}
