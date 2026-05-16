package dtos

import "time"

type CreateDefaultMilkRateRequest struct {
	Rate    float64 `json:"rate" validate:"required"`
	RouteID uint64  `json:"route_id" validate:"required"`
}

type UpdateDefaultMilkRateRequest struct {
	Rate    float64 `json:"rate" validate:"required"`
	RouteID uint64  `json:"route_id" validate:"required"`
}

type DefaultMilkRateResponse struct {
	ID        uint64    `json:"id"`
	Rate      float64   `json:"rate"`
	RouteID   uint64    `json:"route_id"`
	RouteName string    `json:"route_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
