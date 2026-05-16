package dtos

import "time"

type CreateStoreStockRequest struct {
	ItemID             uint64  `json:"item_id" validate:"required"`
	StoreID            uint64  `json:"store_id" validate:"required"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	BuyingPrice        float64 `json:"buying_price" validate:"required"`
	SellingPrice       float64 `json:"selling_price"`
	CreditSellingPrice float64 `json:"credit_selling_price"`
}

type UpdateStoreStockRequest struct {
	ItemID             uint64  `json:"item_id" validate:"required"`
	StoreID            uint64  `json:"store_id" validate:"required"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	BuyingPrice        float64 `json:"buying_price" validate:"required"`
	SellingPrice       float64 `json:"selling_price"`
	CreditSellingPrice float64 `json:"credit_selling_price"`
}

type StoreStockResponse struct {
	ID                 uint64    `json:"id"`
	ItemID             uint64    `json:"item_id"`
	ItemName           string    `json:"item_name"`
	StoreID            uint64    `json:"store_id"`
	StoreName          string    `json:"store_name"`
	Quantity           float64   `json:"quantity"`
	Unit               string    `json:"unit"`
	BuyingPrice        float64   `json:"buying_price"`
	SellingPrice       float64   `json:"selling_price"`
	CreditSellingPrice float64   `json:"credit_selling_price"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
