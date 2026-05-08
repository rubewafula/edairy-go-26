package dtos

import "time"

type RouteResponse struct {
	ID           uint64    `json:"ID"`
	RouteName    string    `json:"Name"`
	Description  string    `json:"Description"`
	RouteCode    string    `json:"Code"`
	LocationID   uint64    `json:"LocationID"`
	LocationName string    `json:"LocationName"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}
