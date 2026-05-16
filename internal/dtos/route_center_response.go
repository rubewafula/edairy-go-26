package dtos

import "time"

type RouteCenterResponse struct {
	ID        uint64    `json:"id"`
	RouteID   uint64    `json:"route_id"`
	RouteName string    `json:"route_name"`
	Center    string    `json:"center"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
