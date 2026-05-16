package dtos

import "time"

type CreateStoreInventoryRequest struct {
	InventoryName string `json:"inventory_name" validate:"required,max=255"`
	CategoryID    uint64 `json:"category_id"`
	IsActive      bool   `json:"is_active"`
	Description   string `json:"description" validate:"max=255"`
}

type UpdateStoreInventoryRequest struct {
	InventoryName string `json:"inventory_name" validate:"required,max=255"`
	CategoryID    uint64 `json:"category_id"`
	IsActive      bool   `json:"is_active"`
	Description   string `json:"description" validate:"max=255"`
}

type StoreInventoryResponse struct {
	ID            uint64    `json:"id"`
	InventoryName string    `json:"inventory_name"`
	CategoryID    uint64    `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	IsActive      bool      `json:"is_active"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
