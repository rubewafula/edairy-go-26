package dtos

import "time"

type PermissionResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	GuardName string    `json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
