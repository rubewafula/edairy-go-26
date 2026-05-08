package dtos

import "time"

type RouteCenterResponse struct {
	ID        uint64    `json:"ID"`
	RouteID   uint64    `json:"RouteID"`
	RouteName string    `json:"RouteName"`
	Center    string    `json:"Center"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}
