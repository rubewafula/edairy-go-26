package dtos

import "time"

type CreateStoreStockTakingItemRequest struct {
	ItemID           uint64  `json:"item_id" validate:"required"`
	SystemQuantity   float64 `json:"system_quantity"` // Fetched from DB, provided value ignored
	PhysicalQuantity float64 `json:"physical_quantity" validate:"required"`
}

type CreateStoreStockTakingRequest struct {
	StockTakeNo   string                              `json:"stock_take_no" validate:"required,max=100"`
	StoreID       uint64                              `json:"store_id" validate:"required"`
	Remarks       string                              `json:"remarks" validate:"max=255"`
	StockTakeDate string                              `json:"stock_take_date" validate:"required"`
	Items         []CreateStoreStockTakingItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateStoreStockTakingRequest struct {
	StockTakeNo      string  `json:"stock_take_no" validate:"required,max=100"`
	StoreID          uint64  `json:"store_id" validate:"required"`
	ItemID           uint64  `json:"item_id" validate:"required"`
	SystemQuantity   float64 `json:"system_quantity" validate:"required"`
	PhysicalQuantity float64 `json:"physical_quantity" validate:"required"`
	VarianceQuantity float64 `json:"variance_quantity" validate:"required"`
	Remarks          string  `json:"remarks" validate:"max=255"`
	StockTakeDate    string  `json:"stock_take_date" validate:"required,datetime"`
}

type StoreStockTakingResponse struct {
	ID               uint64    `json:"id"`
	StockTakeNo      string    `json:"stock_take_no"`
	StoreID          uint64    `json:"store_id"`
	StoreName        string    `json:"store_name"`
	ItemID           uint64    `json:"item_id"`
	ItemName         string    `json:"item_name"`
	SystemQuantity   float64   `json:"system_quantity"`
	PhysicalQuantity float64   `json:"physical_quantity"`
	VarianceQuantity float64   `json:"variance_quantity"`
	Remarks          string    `json:"remarks"`
	StockTakeDate    time.Time `json:"stock_take_date"`
	CreatedAt        time.Time `json:"created_at"`
}
