package dtos

import "time"

type CreateMemberPayDateRangeRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	PayMonth  string `json:"pay_month" binding:"required"`
	PayYear   string `json:"pay_year" binding:"required"`
}

type UpdateMemberPayDateRangeRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	PayMonth  string `json:"pay_month"`
	PayYear   string `json:"pay_year"`
}

type MemberPayDateRangeResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	PayMonth  string    `json:"pay_month"`
	PayYear   string    `json:"pay_year"`
	Processed bool      `json:"processed"`
	Confirmed bool      `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
