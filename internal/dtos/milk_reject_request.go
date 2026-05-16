package dtos

import "time"

type CreateMilkRejectRequest struct {
	RouteID             uint64  `json:"route_id" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required"`
	TransactionDate     string  `json:"transaction_date" validate:"required,datetime"`
	Confirmed           int     `json:"confirmed"`
	Reason              string  `json:"reason" validate:"required,max=255"`
	Description         string  `json:"description" validate:"max=255"`
	TransporterID       uint64  `json:"transporter_id"`
	CanID               uint64  `json:"can_id"`
	MemberID            uint64  `json:"member_id"`
	MilkDeliveryShiftID uint64  `json:"milk_delivery_shift"`
}

type UpdateMilkRejectRequest struct {
	RouteID             uint64  `json:"route_id" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required"`
	TransactionDate     string  `json:"transaction_date" validate:"required,datetime"`
	Confirmed           int     `json:"confirmed"`
	Reason              string  `json:"reason" validate:"required,max=255"`
	Description         string  `json:"description" validate:"max=255"`
	TransporterID       uint64  `json:"transporter_id"`
	CanID               uint64  `json:"can_id"`
	MemberID            uint64  `json:"member_id"`
	MilkDeliveryShiftID uint64  `json:"milk_delivery_shift"`
}

type MilkRejectResponse struct {
	ID                uint64    `json:"id"`
	RouteName         string    `json:"route_name"`
	Quantity          float64   `json:"quantity"`
	TransactionDate   time.Time `json:"transaction_date"`
	Confirmed         int       `json:"confirmed"`
	Reason            string    `json:"reason"`
	Description       string    `json:"description"`
	TransporterName   string    `json:"transporter_name"`
	MemberName        string    `json:"member_name"`
	MilkDeliveryShift string    `json:"milk_delivery_shift"`
	CreatedAt         time.Time `json:"created_at"`
}
