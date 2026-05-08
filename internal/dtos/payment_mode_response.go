package dtos

import "time"

type PaymentModeResponse struct {
	ID        uint64    `json:"ID"`
	Code      string    `json:"Code"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}
