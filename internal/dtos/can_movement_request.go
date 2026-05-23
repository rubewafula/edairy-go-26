package dtos

import "time"

type CreateCanMovementRequest struct {
	CanID             uint64 `json:"can_id" validate:"required"`
	MovementType      string `json:"movement_type" validate:"required"`
	Remarks           string `json:"remarks"`
	ShiftID           uint64 `json:"shift_id" validate:"required"`
	TransporterID     uint64 `json:"transporter_id" validate:"required"`
	RouteID           uint64 `json:"route_id" validate:"required"`
	MovementDate      string `json:"movement_date" validate:"required"`
	ConditionOnReturn string `json:"condition_on_return"`
}

type UpdateCanMovementRequest struct {
	CanID             uint64 `json:"can_id" validate:"required"`
	MovementType      string `json:"movement_type" validate:"required"`
	Remarks           string `json:"remarks"`
	ShiftID           uint64 `json:"shift_id" validate:"required"`
	TransporterID     uint64 `json:"transporter_id" validate:"required"`
	RouteID           uint64 `json:"route_id" validate:"required"`
	MovementDate      string `json:"movement_date" validate:"required"`
	ConditionOnReturn string `json:"condition_on_return"`
}

type CanMovementResponse struct {
	ID                uint64    `json:"id"`
	CanID             uint64    `json:"can_id"`
	CanCode           string    `json:"can_code"`
	MovementType      string    `json:"movement_type"`
	Quantity          float64   `json:"quantity"`
	Remarks           string    `json:"remarks"`
	ShiftID           uint64    `json:"shift_id"`
	ShiftName         string    `json:"shift_name"`
	TransporterID     uint64    `json:"transporter_id"`
	TransporterNo     string    `json:"transporter_no"`
	RouteID           uint64    `json:"route_id"`
	RouteName         string    `json:"route_name"`
	MovementDate      time.Time `json:"movement_date"`
	ConditionOnReturn string    `json:"condition_on_return"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
