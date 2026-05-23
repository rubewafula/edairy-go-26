package dtos

import "time"

type CreateCoolerMilkCollectionRequest struct {
	TransactionDate     string  `json:"transaction_date" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required"`
	TransportVehicleID  uint64  `json:"transport_vehicle_id"`
	MilkDeliveryShiftID uint64  `json:"milk_delivery_shift_id"`
	Confirmed           int     `json:"confirmed"`
	SiteID              uint64  `json:"site_id"`
	TransporterID       uint64  `json:"transporter_id"`
	RouteID             uint64  `json:"route_id"`
}

type UpdateCoolerMilkCollectionRequest struct {
	TransactionDate     string  `json:"transaction_date" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required"`
	TransportVehicleID  uint64  `json:"transport_vehicle_id"`
	MilkDeliveryShiftID uint64  `json:"milk_delivery_shift_id"`
	Confirmed           int     `json:"confirmed"`
	SiteID              uint64  `json:"site_id"`
	TransporterID       uint64  `json:"transporter_id"`
	RouteID             uint64  `json:"route_id"`
}

type CoolerMilkCollectionResponse struct {
	ID                  uint64    `json:"id"`
	TransactionDate     time.Time `json:"transaction_date"`
	Quantity            float64   `json:"quantity"`
	TransportVehicleID  uint64    `json:"transport_vehicle_id"`
	VehicleRegNo        string    `json:"vehicle_reg_no"`
	MilkDeliveryShiftID uint64    `json:"milk_delivery_shift_id"`
	MilkDeliveryShift   string    `json:"milk_delivery_shift"`
	Confirmed           int       `json:"confirmed"`
	SiteID              uint64    `json:"site_id"`
	TransporterID       uint64    `json:"transporter_id"`
	TransporterNo       string    `json:"transporter_no"`
	RouteID             uint64    `json:"route_id"`
	RouteName           string    `json:"route_name"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
