package dtos

import "time"

type CreateStoreItemRequest struct {
	ItemName                  string  `json:"ItemName" validate:"required,max=255"`
	Description               string  `json:"Description" validate:"max=128"`
	ReorderPoint              int     `json:"ReorderPoint"`
	DefaultBuyingPrice        float64 `json:"DefaultBuyingPrice"`
	DefaultSellingPrice       float64 `json:"DefaultSellingPrice"`
	DefaultSellingPriceCredit float64 `json:"DefaultSellingPriceCredit"`
	Status                    string  `json:"Status"`
	Thumbnail                 string  `json:"Thumbnail"`
	SKU                       string  `json:"SKU" validate:"max=45"`
	Barcode                   string  `json:"Barcode" validate:"max=45"`
	UnitID                    int64   `json:"UnitID"`
	StoreInventoryID          uint64  `json:"StoreInventoryID"`
}

type UpdateStoreItemRequest struct {
	ItemName                  string  `json:"ItemName" validate:"required,max=255"`
	Description               string  `json:"Description" validate:"max=128"`
	ReorderPoint              int     `json:"ReorderPoint"`
	DefaultBuyingPrice        float64 `json:"DefaultBuyingPrice"`
	DefaultSellingPrice       float64 `json:"DefaultSellingPrice"`
	DefaultSellingPriceCredit float64 `json:"DefaultSellingPriceCredit"`
	Status                    string  `json:"Status"`
	Thumbnail                 string  `json:"Thumbnail"`
	SKU                       string  `json:"SKU" validate:"max=45"`
	Barcode                   string  `json:"Barcode" validate:"max=45"`
	UnitID                    int64   `json:"UnitID"`
	StoreInventoryID          uint64  `json:"StoreInventoryID"`
}

type StoreItemResponse struct {
	ID                        uint64    `json:"ID"`
	ItemName                  string    `json:"ItemName"`
	Description               string    `json:"Description"`
	ReorderPoint              int       `json:"ReorderPoint"`
	DefaultBuyingPrice        float64   `json:"DefaultBuyingPrice"`
	DefaultSellingPrice       float64   `json:"DefaultSellingPrice"`
	DefaultSellingPriceCredit float64   `json:"DefaultSellingPriceCredit"`
	Status                    string    `json:"Status"`
	Thumbnail                 string    `json:"Thumbnail"`
	SKU                       string    `json:"SKU"`
	Barcode                   string    `json:"Barcode"`
	UnitID                    int64     `json:"UnitID"`
	StoreInventoryID          uint64    `json:"StoreInventoryID"`
	InventoryName             string    `json:"InventoryName"`
	CreatedAt                 time.Time `json:"CreatedAt"`
	UpdatedAt                 time.Time `json:"UpdatedAt"`
}
