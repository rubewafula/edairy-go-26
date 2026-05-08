package dtos

import "time"

type PermissionResponse struct {
	ID        uint64    `json:"ID"`
	Name      string    `json:"Name"`
	GuardName string    `json:"GuardName"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}
