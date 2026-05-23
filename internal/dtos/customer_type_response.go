package dtos

import "time"

type CustomerTypeResponse struct {
	ID          uint64    `json:"id"`
	TypeCode    string    `json:"type_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
