package dtos

import "time"

// CreateMemberTypeRequest defines the structure for creating a new member type.
type CreateMemberTypeRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=500"`
}

// UpdateMemberTypeRequest defines the structure for updating an existing member type.
type UpdateMemberTypeRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=500"`
}

// MemberTypeResponse defines the structure for returning member type data.
type MemberTypeResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
