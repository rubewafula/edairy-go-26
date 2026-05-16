package dtos

import "time"

type CreateStoreItemUnitRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Symbol      string `json:"symbol" validate:"required,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateStoreItemUnitRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Symbol      string `json:"symbol" validate:"required,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type StoreItemUnitResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
