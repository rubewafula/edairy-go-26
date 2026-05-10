package dtos

import "time"

type TransporterRouteAssignmentResponse struct {
	ID            uint64    `json:"ID"`
	TransporterID uint64    `json:"TransporterID"`
	TransporterNo string    `json:"TransporterNo"`
	RouteID       uint64    `json:"RouteID"`
	RouteName     string    `json:"RouteName"`
	StartDate     time.Time `json:"StartDate"`
	EndDate       time.Time `json:"EndDate"`
	Active        bool      `json:"Active"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
