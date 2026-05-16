package dtos

import "time"

type TransporterRouteAssignmentResponse struct {
	ID            uint64    `json:"id"`
	TransporterID uint64    `json:"transporter_id"`
	TransporterNo string    `json:"transporter_no"`
	RouteID       uint64    `json:"route_id"`
	RouteName     string    `json:"route_name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
