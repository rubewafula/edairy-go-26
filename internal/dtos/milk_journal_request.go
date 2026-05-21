package dtos

import "time"

type CreateMilkJournalRequest struct {
	Journal             string                    `json:"journal" validate:"required,max=255"`
	JournalDate         string                    `json:"journal_date" validate:"required"`
	MilkDeliveryShiftID uint64                    `json:"milk_delivery_shift_id" validate:"required"`
	RouteID             uint64                    `json:"route_id" validate:"required"`
	UserID              uint64                    `json:"user_id"`
	TransporterID       uint64                    `json:"transporter_id"`
	Confirmed           bool                      `json:"confirmed"`
	Batches             []MilkJournalBatchRequest `json:"batches"`
}

type MilkJournalBatchRequest struct {
	BatchNo string                    `json:"batch_no"`
	Entries []MilkJournalEntryRequest `json:"entries"`
}

type MilkJournalEntryRequest struct {
	MemberID      uint64  `json:"member_id"`
	Status        string  `json:"status"`
	Quantity      float64 `json:"quantity"`
	RouteCenterID uint64  `json:"route_center_id"`
	CanID         uint64  `json:"can_id"`
}

type UpdateMilkJournalRequest struct {
	Journal             string                    `json:"journal" validate:"required,max=255"`
	JournalDate         string                    `json:"journal_date" validate:"required"`
	MilkDeliveryShiftID uint64                    `json:"milk_delivery_shift_id" validate:"required"`
	RouteID             uint64                    `json:"route_id" validate:"required"`
	UserID              uint64                    `json:"user_id"`
	TransporterID       uint64                    `json:"transporter_id"`
	Confirmed           bool                      `json:"confirmed"`
	Batches             []MilkJournalBatchRequest `json:"batches"`
}

type MilkJournalResponse struct {
	ID                  uint64    `json:"id"`
	Journal             string    `json:"journal"`
	JournalDate         string    `json:"journal_date"`
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
	Collections         float64   `json:"collections"`
	BatchNo             string    `json:"batch_no"`
}

type DailyJournalSummaryResponse struct {
	Journal string   `json:"journal"`
	Route   string   `json:"route"`
	Shift   string   `json:"shift"`
	Batches []string `json:"batches"`
}
