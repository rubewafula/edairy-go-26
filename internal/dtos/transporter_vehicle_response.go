package dtos

import "time"

type TransporterVehicleResponse struct {
	ID             uint64    `json:"id"`
	TransporterID  uint64    `json:"transporter_id"`
	RouteID        uint64    `json:"route_id"`
	RegistrationNo string    `json:"registration_no"`
	VehicleType    string    `json:"vehicle_type"`
	CapacityLitres float64   `json:"capacity_litres"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
