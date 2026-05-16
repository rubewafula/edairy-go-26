package dtos

import "time"

type CreateStoreStockMovementTypeRequest struct {
	MovementCode string `json:"movement_code" validate:"required,max=100"`
	MovementName string `json:"movement_name" validate:"required,max=255"`
	Direction    string `json:"direction" validate:"required,oneof=IN OUT"`
	AffectsStock bool   `json:"affects_stock"`
	Description  string `json:"description" validate:"max=255"`
	IsSystem     bool   `json:"is_system"`
}

type UpdateStoreStockMovementTypeRequest struct {
	MovementCode string `json:"movement_code" validate:"required,max=100"`
	MovementName string `json:"movement_name" validate:"required,max=255"`
	Direction    string `json:"direction" validate:"required,oneof=IN OUT"`
	AffectsStock bool   `json:"affects_stock"`
	Description  string `json:"description" validate:"max=255"`
	IsSystem     bool   `json:"is_system"`
}

type StoreStockMovementTypeResponse struct {
	ID           uint64    `json:"id"`
	MovementCode string    `json:"movement_code"`
	MovementName string    `json:"movement_name"`
	Direction    string    `json:"direction"`
	AffectsStock bool      `json:"affects_stock"`
	Description  string    `json:"description"`
	IsSystem     bool      `json:"is_system"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
