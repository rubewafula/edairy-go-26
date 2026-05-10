package dtos

import "time"

type CreateMilkDeliveryShiftRequest struct {
	Name        string `json:"Name" validate:"required,max=255"`
	Description string `json:"Description" validate:"required,max=255"`
}

type UpdateMilkDeliveryShiftRequest struct {
	Name        string `json:"Name" validate:"required,max=255"`
	Description string `json:"Description" validate:"required,max=255"`
}

type MilkDeliveryShiftResponse struct {
	ID          uint64    `json:"ID"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
