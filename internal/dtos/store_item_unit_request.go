package dtos

import "time"

type CreateStoreItemUnitRequest struct {
	Name        string `json:"Name" validate:"required,max=100"`
	Symbol      string `json:"Symbol" validate:"required,max=20"`
	Description string `json:"Description" validate:"max=255"`
}

type UpdateStoreItemUnitRequest struct {
	Name        string `json:"Name" validate:"required,max=100"`
	Symbol      string `json:"Symbol" validate:"required,max=20"`
	Description string `json:"Description" validate:"max=255"`
}

type StoreItemUnitResponse struct {
	ID          uint64    `json:"ID"`
	Name        string    `json:"Name"`
	Symbol      string    `json:"Symbol"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
