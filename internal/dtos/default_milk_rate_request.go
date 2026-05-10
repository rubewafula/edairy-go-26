package dtos

import "time"

type CreateDefaultMilkRateRequest struct {
	Rate    float64 `json:"Rate" validate:"required"`
	RouteID uint64  `json:"RouteID" validate:"required"`
}

type UpdateDefaultMilkRateRequest struct {
	Rate    float64 `json:"Rate" validate:"required"`
	RouteID uint64  `json:"RouteID" validate:"required"`
}

type DefaultMilkRateResponse struct {
	ID        uint64    `json:"ID"`
	Rate      float64   `json:"Rate"`
	RouteID   uint64    `json:"RouteID"`
	RouteName string    `json:"RouteName"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}
