package dtos

import "time"

type CreateMilkJournalRequest struct {
	Journal             string `json:"Journal" validate:"required,max=255"`
	JournalDate         string `json:"JournalDate" validate:"required,datetime"`
	MilkDeliveryShiftID uint64 `json:"MilkDeliveryShiftID" validate:"required"`
	RouteID             uint64 `json:"RouteID" validate:"required"`
	UserID              uint64 `json:"UserID"`
	TransporterID       uint64 `json:"TransporterID"`
	Confirmed           bool   `json:"Confirmed"`
}

type MilkJournalResponse struct {
	ID                  uint64    `json:"ID"`
	Journal             string    `json:"Journal"`
	JournalDate         time.Time `json:"JournalDate"`
	MilkDeliveryShiftID uint64    `json:"MilkDeliveryShiftID"`
	MilkDeliveryShift   string    `json:"MilkDeliveryShift"`
	RouteID             uint64    `json:"RouteID"`
	RouteName           string    `json:"RouteName"`
	UserID              uint64    `json:"UserID"`
	TransporterID       uint64    `json:"TransporterID"`
	Confirmed           bool      `json:"Confirmed"`
	EntriesCount        int64     `json:"EntriesCount"`
	CreatedAt           time.Time `json:"CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt"`
}

type UpdateMilkJournalRequest struct {
	Journal             string `json:"Journal" validate:"required,max=255"`
	JournalDate         string `json:"JournalDate" validate:"required,datetime"`
	MilkDeliveryShiftID uint64 `json:"MilkDeliveryShiftID" validate:"required"`
	RouteID             uint64 `json:"RouteID" validate:"required"`
	UserID              uint64 `json:"UserID"`
	TransporterID       uint64 `json:"TransporterID"`
	Confirmed           bool   `json:"Confirmed"`
}
