package dtos

import "time"

type RouteResponse struct {
	ID           uint64    `json:"id"`
	RouteName    string    `json:"name"`
	Description  string    `json:"description"`
	RouteCode    string    `json:"code"`
	LocationID   uint64    `json:"location_id"`
	LocationName string    `json:"location_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
