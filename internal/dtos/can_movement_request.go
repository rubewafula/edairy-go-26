package dtos

import "time"

type CreateCanMovementRequest struct {
	CanID             uint64  `json:"CanID" validate:"required"`
	MovementType      string  `json:"MovementType" validate:"required"`
	Quantity          float64 `json:"Quantity" validate:"required"`
	Remarks           string  `json:"Remarks"`
	ShiftID           uint64  `json:"ShiftID" validate:"required"`
	TransporterID     uint64  `json:"TransporterID" validate:"required"`
	RouteID           uint64  `json:"RouteID" validate:"required"`
	MovementDate      string  `json:"MovementDate" validate:"required,datetime"`
	ConditionOnReturn string  `json:"ConditionOnReturn"`
}

type UpdateCanMovementRequest struct {
	CanID             uint64  `json:"CanID" validate:"required"`
	MovementType      string  `json:"MovementType" validate:"required"`
	Quantity          float64 `json:"Quantity" validate:"required"`
	Remarks           string  `json:"Remarks"`
	ShiftID           uint64  `json:"ShiftID" validate:"required"`
	TransporterID     uint64  `json:"TransporterID" validate:"required"`
	RouteID           uint64  `json:"RouteID" validate:"required"`
	MovementDate      string  `json:"MovementDate" validate:"required,datetime"`
	ConditionOnReturn string  `json:"ConditionOnReturn"`
}

type CanMovementResponse struct {
	ID                uint64    `json:"ID"`
	CanID             uint64    `json:"CanID"`
	CanCode           string    `json:"CanCode"`
	MovementType      string    `json:"MovementType"`
	Quantity          float64   `json:"Quantity"`
	Remarks           string    `json:"Remarks"`
	ShiftID           uint64    `json:"ShiftID"`
	ShiftName         string    `json:"ShiftName"`
	TransporterID     uint64    `json:"TransporterID"`
	TransporterNo     string    `json:"TransporterNo"`
	RouteID           uint64    `json:"RouteID"`
	RouteName         string    `json:"RouteName"`
	MovementDate      time.Time `json:"MovementDate"`
	ConditionOnReturn string    `json:"ConditionOnReturn"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
}
