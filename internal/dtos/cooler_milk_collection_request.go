package dtos

import "time"

type CreateCoolerMilkCollectionRequest struct {
	TransactionDate     string  `json:"TransactionDate" validate:"required,datetime"`
	Quantity            float64 `json:"Quantity" validate:"required"`
	TransportVehicleID  uint64  `json:"TransportVehicleID"`
	MilkDeliveryShiftID uint64  `json:"MilkDeliveryShiftID"`
	Confirmed           int     `json:"Confirmed"`
	SiteID              uint64  `json:"SiteID"`
	TransporterID       uint64  `json:"TransporterID"`
	RouteID             uint64  `json:"RouteID"`
}

type UpdateCoolerMilkCollectionRequest struct {
	TransactionDate     string  `json:"TransactionDate" validate:"required,datetime"`
	Quantity            float64 `json:"Quantity" validate:"required"`
	TransportVehicleID  uint64  `json:"TransportVehicleID"`
	MilkDeliveryShiftID uint64  `json:"MilkDeliveryShiftID"`
	Confirmed           int     `json:"Confirmed"`
	SiteID              uint64  `json:"SiteID"`
	TransporterID       uint64  `json:"TransporterID"`
	RouteID             uint64  `json:"RouteID"`
}

type CoolerMilkCollectionResponse struct {
	ID                  uint64    `json:"ID"`
	TransactionDate     time.Time `json:"TransactionDate"`
	Quantity            float64   `json:"Quantity"`
	TransportVehicleID  uint64    `json:"TransportVehicleID"`
	VehicleRegNo        string    `json:"VehicleRegNo"`
	MilkDeliveryShiftID uint64    `json:"MilkDeliveryShiftID"`
	MilkDeliveryShift   string    `json:"MilkDeliveryShift"`
	Confirmed           int       `json:"Confirmed"`
	SiteID              uint64    `json:"SiteID"`
	TransporterID       uint64    `json:"TransporterID"`
	TransporterNo       string    `json:"TransporterNo"`
	RouteID             uint64    `json:"RouteID"`
	RouteName           string    `json:"RouteName"`
	CreatedAt           time.Time `json:"CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt"`
}
