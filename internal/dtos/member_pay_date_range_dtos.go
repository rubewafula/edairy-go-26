package dtos

import "time"

type CreateMemberPayDateRangeRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type UpdateMemberPayDateRangeRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Confirmed int    `json:"confirmed"`
}

type MemberPayDateRangeResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Processed int       `json:"processed"`
	Confirmed int       `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
