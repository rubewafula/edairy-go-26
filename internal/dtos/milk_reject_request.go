package dtos

import "time"

type CreateMilkRejectRequest struct {
	RouteID             uint64  `json:"RouteID" validate:"required"`
	Quantity            float64 `json:"Quantity" validate:"required"`
	TransactionDate     string  `json:"TransactionDate" validate:"required,datetime"`
	Confirmed           int     `json:"Confirmed"`
	Reason              string  `json:"Reason" validate:"required,max=255"`
	Description         string  `json:"Description" validate:"max=255"`
	TransporterID       uint64  `json:"TransporterID"`
	CanID               uint64  `json:"CanID"`
	MemberID            uint64  `json:"MemberID"`
	MilkDeliveryShiftID uint64  `json:"MilkDeliveryShiftID"`
}

type UpdateMilkRejectRequest struct {
	RouteID             uint64  `json:"RouteID" validate:"required"`
	Quantity            float64 `json:"Quantity" validate:"required"`
	TransactionDate     string  `json:"TransactionDate" validate:"required,datetime"`
	Confirmed           int     `json:"Confirmed"`
	Reason              string  `json:"Reason" validate:"required,max=255"`
	Description         string  `json:"Description" validate:"max=255"`
	TransporterID       uint64  `json:"TransporterID"`
	CanID               uint64  `json:"CanID"`
	MemberID            uint64  `json:"MemberID"`
	MilkDeliveryShiftID uint64  `json:"MilkDeliveryShiftID"`
}

type MilkRejectResponse struct {
	ID                uint64    `json:"ID"`
	RouteName         string    `json:"RouteName"`
	Quantity          float64   `json:"Quantity"`
	TransactionDate   time.Time `json:"TransactionDate"`
	Confirmed         int       `json:"Confirmed"`
	Reason            string    `json:"Reason"`
	Description       string    `json:"Description"`
	TransporterName   string    `json:"TransporterName"`
	MemberName        string    `json:"MemberName"`
	MilkDeliveryShift string    `json:"MilkDeliveryShift"`
	CreatedAt         time.Time `json:"CreatedAt"`
}
