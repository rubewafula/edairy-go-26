package dtos

type CreateTransporterVehicleRequest struct {
	TransporterID  uint64  `json:"transporter_id" validate:"required"`
	RouteID        uint64  `json:"route_id"`
	RegistrationNo string  `json:"registration_no" validate:"required"`
	VehicleType    string  `json:"vehicle_type" validate:"required,oneof=MOTORBIKE PICKUP VAN TANKER TRUCK LORRY BICYCLE TUKTUK OTHER"`
	CapacityLitres float64 `json:"capacity_litres"`
	Active         bool    `json:"active"`
}

type UpdateTransporterVehicleRequest struct {
	RouteID        uint64  `json:"route_id"`
	VehicleType    string  `json:"vehicle_type" validate:"required,oneof=MOTORBIKE PICKUP VAN TANKER TRUCK LORRY BICYCLE TUKTUK OTHER"`
	CapacityLitres float64 `json:"capacity_litres"`
	Active         bool    `json:"active"`
}
