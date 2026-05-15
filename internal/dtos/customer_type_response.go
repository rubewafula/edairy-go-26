package dtos

import "time"

type CustomerTypeResponse struct {
	ID          uint64    `json:"ID"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	CreatedBy   uint64    `json:"CreatedBy"`
	UpdatedBy   uint64    `json:"UpdatedBy"`
}
