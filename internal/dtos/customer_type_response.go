package dtos

import "time"

type CustomerTypeResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint64    `json:"created_by"`
	UpdatedBy   uint64    `json:"updated_by"`
}
