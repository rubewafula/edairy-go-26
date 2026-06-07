package dtos

import "time"

type CreateDefaultMilkRateRequest struct {
	Rate     float64 `json:"rate" validate:"required"`
	RouteID  uint64  `json:"route_id"`
	MemberID uint64  `json:"member_id"`
}

type UpdateDefaultMilkRateRequest struct {
	Rate     float64 `json:"rate" validate:"required"`
	RouteID  uint64  `json:"route_id" `
	MemberID uint64  `json:"member_id"`
}

type DefaultMilkRateResponse struct {
	ID         uint64    `json:"id"`
	Rate       float64   `json:"rate"`
	RouteID    uint64    `json:"route_id"`
	RouteName  string    `json:"route_name"`
	MemberName string    `json:"member_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
