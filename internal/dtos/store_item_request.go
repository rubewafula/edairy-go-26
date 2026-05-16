package dtos

import "time"

type CreateStoreItemRequest struct {
	ItemName                  string  `json:"item_name" validate:"required,max=255"`
	Description               string  `json:"description" validate:"max=128"`
	ReorderPoint              int     `json:"reorder_point"`
	DefaultBuyingPrice        float64 `json:"default_buying_price"`
	DefaultSellingPrice       float64 `json:"default_selling_price"`
	DefaultSellingPriceCredit float64 `json:"default_selling_price_credit"`
	Status                    string  `json:"status"`
	Thumbnail                 string  `json:"thumbnail"`
	SKU                       string  `json:"sku" validate:"max=45"`
	Barcode                   string  `json:"barcode" validate:"max=45"`
	UnitID                    int64   `json:"unit_id"`
	StoreInventoryID          uint64  `json:"store_inventory_id"`
}

type UpdateStoreItemRequest struct {
	ItemName                  string  `json:"item_name" validate:"required,max=255"`
	Description               string  `json:"description" validate:"max=128"`
	ReorderPoint              int     `json:"reorder_point"`
	DefaultBuyingPrice        float64 `json:"default_buying_price"`
	DefaultSellingPrice       float64 `json:"default_selling_price"`
	DefaultSellingPriceCredit float64 `json:"default_selling_price_credit"`
	Status                    string  `json:"status"`
	Thumbnail                 string  `json:"thumbnail"`
	SKU                       string  `json:"sku" validate:"max=45"`
	Barcode                   string  `json:"barcode" validate:"max=45"`
	UnitID                    int64   `json:"unit_id"`
	StoreInventoryID          uint64  `json:"store_inventory_id"`
}

type StoreItemResponse struct {
	ID                        uint64    `json:"id"`
	ItemName                  string    `json:"item_name"`
	Description               string    `json:"description"`
	ReorderPoint              int       `json:"reorder_point"`
	DefaultBuyingPrice        float64   `json:"default_buying_price"`
	DefaultSellingPrice       float64   `json:"default_selling_price"`
	DefaultSellingPriceCredit float64   `json:"default_selling_price_credit"`
	Status                    string    `json:"status"`
	Thumbnail                 string    `json:"thumbnail"`
	SKU                       string    `json:"sku"`
	Barcode                   string    `json:"barcode"`
	UnitID                    int64     `json:"unit_id"`
	StoreInventoryID          uint64    `json:"store_inventory_id"`
	InventoryName             string    `json:"inventory_name"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
