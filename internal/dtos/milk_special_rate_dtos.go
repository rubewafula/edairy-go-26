package dtos

import "time"

type CreateMilkSpecialRateRequest struct {
	MonthlyPayDateRangeID uint64  `json:"monthly_pay_date_range_id" binding:"required"`
	MemberID              uint64  `json:"member_id"`
	RouteID               uint64  `json:"route_id"`
	Rate                  float64 `json:"rate" binding:"required"`
}

type UpdateMilkSpecialRateRequest struct {
	PayDateRangeID uint64  `json:"monthly_pay_date_range_id"`
	MemberID       uint64  `json:"member_id"`
	RouteID        uint64  `json:"route_id"`
	Rate           float64 `json:"rate"`
	Confirmed      int     `json:"confirmed"`
}

type MilkSpecialRateResponse struct {
	ID               uint64    `json:"id"`
	PayDateRangeID   uint64    `json:"pay_date_range_id"`
	PayDateRangeName string    `json:"pay_date_range_name"`
	MemberID         uint64    `json:"member_id"`
	MemberNo         string    `json:"member_no"`
	MemberName       string    `json:"member_name"`
	RouteID          uint64    `json:"route_id"`
	RouteName        string    `json:"route_name"`
	Rate             float64   `json:"rate"`
	Confirmed        int       `json:"confirmed"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
