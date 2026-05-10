package dtos

import "time"

type CreateStoreStockMovementTypeRequest struct {
	MovementCode string `json:"MovementCode" validate:"required,max=100"`
	MovementName string `json:"MovementName" validate:"required,max=255"`
	Direction    string `json:"Direction" validate:"required,oneof=IN OUT"`
	AffectsStock bool   `json:"AffectsStock"`
	Description  string `json:"Description" validate:"max=255"`
	IsSystem     bool   `json:"IsSystem"`
}

type UpdateStoreStockMovementTypeRequest struct {
	MovementCode string `json:"MovementCode" validate:"required,max=100"`
	MovementName string `json:"MovementName" validate:"required,max=255"`
	Direction    string `json:"Direction" validate:"required,oneof=IN OUT"`
	AffectsStock bool   `json:"AffectsStock"`
	Description  string `json:"Description" validate:"max=255"`
	IsSystem     bool   `json:"IsSystem"`
}

type StoreStockMovementTypeResponse struct {
	ID           uint64    `json:"ID"`
	MovementCode string    `json:"MovementCode"`
	MovementName string    `json:"MovementName"`
	Direction    string    `json:"Direction"`
	AffectsStock bool      `json:"AffectsStock"`
	Description  string    `json:"Description"`
	IsSystem     bool      `json:"IsSystem"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}
