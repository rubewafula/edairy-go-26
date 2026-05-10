package dtos

import "time"

type CreateStoreStockTakingRequest struct {
	StockTakeNo      string  `json:"StockTakeNo" validate:"required,max=100"`
	StoreID          uint64  `json:"StoreID" validate:"required"`
	ItemID           uint64  `json:"ItemID" validate:"required"`
	SystemQuantity   float64 `json:"SystemQuantity" validate:"required"`
	PhysicalQuantity float64 `json:"PhysicalQuantity" validate:"required"`
	VarianceQuantity float64 `json:"VarianceQuantity" validate:"required"`
	Remarks          string  `json:"Remarks" validate:"max=255"`
	StockTakeDate    string  `json:"StockTakeDate" validate:"required,datetime"`
}

type UpdateStoreStockTakingRequest struct {
	StockTakeNo      string  `json:"StockTakeNo" validate:"required,max=100"`
	StoreID          uint64  `json:"StoreID" validate:"required"`
	ItemID           uint64  `json:"ItemID" validate:"required"`
	SystemQuantity   float64 `json:"SystemQuantity" validate:"required"`
	PhysicalQuantity float64 `json:"PhysicalQuantity" validate:"required"`
	VarianceQuantity float64 `json:"VarianceQuantity" validate:"required"`
	Remarks          string  `json:"Remarks" validate:"max=255"`
	StockTakeDate    string  `json:"StockTakeDate" validate:"required,datetime"`
}

type StoreStockTakingResponse struct {
	ID               uint64    `json:"ID"`
	StockTakeNo      string    `json:"StockTakeNo"`
	StoreID          uint64    `json:"StoreID"`
	StoreName        string    `json:"StoreName"`
	ItemID           uint64    `json:"ItemID"`
	ItemName         string    `json:"ItemName"`
	SystemQuantity   float64   `json:"SystemQuantity"`
	PhysicalQuantity float64   `json:"PhysicalQuantity"`
	VarianceQuantity float64   `json:"VarianceQuantity"`
	Remarks          string    `json:"Remarks"`
	StockTakeDate    time.Time `json:"StockTakeDate"`
	CreatedAt        time.Time `json:"CreatedAt"`
}
