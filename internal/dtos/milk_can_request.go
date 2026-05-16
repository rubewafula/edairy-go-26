package dtos

import "time"

type CreateMilkCanRequest struct {
	CanID      string  `json:"can_id" validate:"required"`
	CanType    string  `json:"can_type" validate:"required"`
	CanSize    float64 `json:"can_size" validate:"required"`
	Units      string  `json:"units"`
	TareWeight float64 `json:"tare_weight"`
	RouteID    uint64  `json:"route_id" validate:"required"`
}

type UpdateMilkCanRequest struct {
	CanID      string  `json:"can_id" validate:"required"`
	CanType    string  `json:"can_type" validate:"required"`
	CanSize    float64 `json:"can_size" validate:"required"`
	Units      string  `json:"units"`
	TareWeight float64 `json:"tare_weight"`
	RouteID    uint64  `json:"route_id" validate:"required"`
}

type MilkCanResponse struct {
	ID         uint64    `json:"id"`
	CanID      string    `json:"can_id"`
	CanType    string    `json:"can_type"`
	CanSize    float64   `json:"can_size"`
	Units      string    `json:"units"`
	TareWeight float64   `json:"tare_weight"`
	RouteID    uint64    `json:"route_id"`
	RouteName  string    `json:"route_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
