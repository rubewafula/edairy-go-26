package dtos

import "time"

type CreateMilkJournalEntryRequest struct {
	MemberID            uint64  `json:"member_id" validate:"required"`
	MilkJournalID       uint64  `json:"milk_journal_id" validate:"required"`
	MilkJournalBatchID  uint64  `json:"milk_journal_batch_id"`
	RouteID             uint64  `json:"route_id" validate:"required"`
	MilkDeliveryShiftID uint64  `json:"milk_delivery_shift_id" validate:"required"`
	Status              string  `json:"status"`
	JournalDate         string  `json:"journal_date" validate:"required,datetime"`
	Quantity            float64 `json:"quantity" validate:"required"`
	TransporterID       uint64  `json:"transporter_id"`
	RouteCenterID       uint64  `json:"route_center_id"`
	CanID               uint64  `json:"can_id"`
}

type MilkJournalEntryResponse struct {
	ID                  uint64    `json:"id"`
	MemberID            uint64    `json:"member_id"`
	Journal             string    `json:"journal"`
	BatchNo             string    `json:"batch_no"`
	MemberNo            string    `json:"member_no"`
	MemberName          string    `json:"member_name"`
	MilkJournalID       uint64    `json:"milk_journal_id"`
	MilkJournalBatchID  uint64    `json:"milk_journal_batch_id"`
	RouteID             uint64    `json:"route_id"`
	RouteName           string    `json:"route_name"`
	MilkDeliveryShiftID uint64    `json:"milk_delivery_shift_id"`
	MilkDeliveryShift   string    `json:"milk_delivery_shift"`
	Status              string    `json:"status"`
	JournalDate         string    `json:"journal_date"`
	Quantity            float64   `json:"quantity"`
	TransporterID       uint64    `json:"transporter_id"`
	RouteCenterID       uint64    `json:"route_center_id"`
	CanID               uint64    `json:"can_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type UpdateMilkJournalEntryRequest struct {
	MemberID            uint64  `json:"member_id" validate:"required"`
	MilkJournalID       uint64  `json:"milk_journal_id" validate:"required"`
	MilkJournalBatchID  uint64  `json:"milk_journal_batch_id"`
	RouteID             uint64  `json:"route_id" validate:"required"`
	MilkDeliveryShiftID uint64  `json:"milk_delivery_shift_id" validate:"required"`
	Status              string  `json:"status"`
	JournalDate         string  `json:"journal_date" validate:"required,datetime"`
	Quantity            float64 `json:"quantity" validate:"required"`
	TransporterID       uint64  `json:"transporter_id"`
	RouteCenterID       uint64  `json:"route_center_id"`
	CanID               uint64  `json:"can_id"`
}

type StrayMilkCollectionResponse struct {
	ID                uint64    `json:"id"`
	MemberID          uint64    `json:"member_id"`
	MemberNo          string    `json:"member_no"`
	MemberName        string    `json:"member_name"`
	MemberRouteID     uint64    `json:"member_route_id"`
	MemberRoute       string    `json:"member_route"`
	JournalRouteID    uint64    `json:"journal_route_id"`
	StrayRoute        string    `json:"stray_route"`
	Quantity          float64   `json:"quantity"`
	JournalDate       time.Time `json:"journal_date"`
	MilkDeliveryShift string    `json:"milk_delivery_shift"`
	CreatedAt         time.Time `json:"created_at"`
}
