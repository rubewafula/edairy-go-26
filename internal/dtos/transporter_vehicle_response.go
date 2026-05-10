package dtos

import "time"

type TransporterVehicleResponse struct {
	ID             uint64    `json:"ID"`
	TransporterID  uint64    `json:"TransporterID"`
	RouteID        uint64    `json:"RouteID"`
	RegistrationNo string    `json:"RegistrationNo"`
	VehicleType    string    `json:"VehicleType"`
	CapacityLitres float64   `json:"CapacityLitres"`
	Active         bool      `json:"Active"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}
