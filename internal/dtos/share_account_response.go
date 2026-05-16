package dtos

import "time"

type ShareAccountResponse struct {
	ID            uint64    `json:"id"`
	MemberID      uint64    `json:"member_id"`
	MemberNo      string    `json:"member_no"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ShareTypeID   uint64    `json:"share_type_id"`
	ShareCode     string    `json:"share_code"`
	ShareTypeName string    `json:"share_type_name"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	OpenedAt      time.Time `json:"opened_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
