package dtos

import "time"

type CreateStoreInventoryRequest struct {
	InventoryName string `json:"InventoryName" validate:"required,max=255"`
	CategoryID    uint64 `json:"CategoryID"`
	IsActive      bool   `json:"IsActive"`
	Description   string `json:"Description" validate:"max=255"`
}

type UpdateStoreInventoryRequest struct {
	InventoryName string `json:"InventoryName" validate:"required,max=255"`
	CategoryID    uint64 `json:"CategoryID"`
	IsActive      bool   `json:"IsActive"`
	Description   string `json:"Description" validate:"max=255"`
}

type StoreInventoryResponse struct {
	ID            uint64    `json:"ID"`
	InventoryName string    `json:"InventoryName"`
	CategoryID    uint64    `json:"CategoryID"`
	CategoryName  string    `json:"CategoryName"`
	IsActive      bool      `json:"IsActive"`
	Description   string    `json:"Description"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
