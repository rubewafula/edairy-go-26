package dtos

import "time"

type CreateMilkJournalRequest struct {
	Journal             string `json:"journal" validate:"required,max=255"`
	JournalDate         string `json:"journal_date" validate:"required,datetime"`
	MilkDeliveryShiftID uint64 `json:"milk_delivery_shift_id" validate:"required"`
	RouteID             uint64 `json:"route_id" validate:"required"`
	UserID              uint64 `json:"user_id"`
	TransporterID       uint64 `json:"transporter_id"`
	Confirmed           bool   `json:"confirmed"`
}

type UpdateMilkJournalRequest struct {
	Journal             string `json:"journal" validate:"required,max=255"`
	JournalDate         string `json:"journal_date" validate:"required,datetime"`
	MilkDeliveryShiftID uint64 `json:"milk_delivery_shift_id" validate:"required"`
	RouteID             uint64 `json:"route_id" validate:"required"`
	UserID              uint64 `json:"user_id"`
	TransporterID       uint64 `json:"transporter_id"`
	Confirmed           bool   `json:"confirmed"`
}

type MilkJournalResponse struct {
	ID                  uint64    `json:"id"`
	Journal             string    `json:"journal"`
	JournalDate         time.Time `json:"journal_date"`
	MilkDeliveryShiftID uint64    `json:"milk_delivery_shift_id"`
	MilkDeliveryShift   string    `json:"milk_delivery_shift"`
	RouteID             uint64    `json:"route_id"`
	RouteName           string    `json:"route_name"`
	UserID              uint64    `json:"user_id"`
	TransporterID       uint64    `json:"transporter_id"`
	Confirmed           bool      `json:"confirmed"`
	EntriesCount        int64     `json:"entries_count"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
