package dtos

import "time"

type CreateStoreStockRequest struct {
	ItemID             uint64  `json:"ItemID" validate:"required"`
	StoreID            uint64  `json:"StoreID" validate:"required"`
	Quantity           float64 `json:"Quantity"`
	Unit               string  `json:"Unit"`
	BuyingPrice        float64 `json:"BuyingPrice" validate:"required"`
	SellingPrice       float64 `json:"SellingPrice"`
	CreditSellingPrice float64 `json:"CreditSellingPrice"`
}

type UpdateStoreStockRequest struct {
	ItemID             uint64  `json:"ItemID" validate:"required"`
	StoreID            uint64  `json:"StoreID" validate:"required"`
	Quantity           float64 `json:"Quantity"`
	Unit               string  `json:"Unit"`
	BuyingPrice        float64 `json:"BuyingPrice" validate:"required"`
	SellingPrice       float64 `json:"SellingPrice"`
	CreditSellingPrice float64 `json:"CreditSellingPrice"`
}

type StoreStockResponse struct {
	ID                 uint64    `json:"ID"`
	ItemID             uint64    `json:"ItemID"`
	ItemName           string    `json:"ItemName"`
	StoreID            uint64    `json:"StoreID"`
	StoreName          string    `json:"StoreName"`
	Quantity           float64   `json:"Quantity"`
	Unit               string    `json:"Unit"`
	BuyingPrice        float64   `json:"BuyingPrice"`
	SellingPrice       float64   `json:"SellingPrice"`
	CreditSellingPrice float64   `json:"CreditSellingPrice"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
}
