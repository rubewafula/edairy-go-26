package dtos

import "time"

type CreateMilkJournalEntryRequest struct {
	MemberID            uint64  `json:"MemberID" validate:"required"`
	MilkJournalID       uint64  `json:"MilkJournalID" validate:"required"`
	MilkJournalBatchID  uint64  `json:"MilkJournalBatchID"`
	RouteID             uint64  `json:"RouteID" validate:"required"`
	MilkDeliveryShiftID uint64  `json:"MilkDeliveryShiftID" validate:"required"`
	Status              string  `json:"Status"`
	JournalDate         string  `json:"JournalDate" validate:"required,datetime"`
	Quantity            float64 `json:"Quantity" validate:"required"`
	TransporterID       uint64  `json:"TransporterID"`
	RouteCenterID       uint64  `json:"RouteCenterID"`
	CanID               uint64  `json:"CanID"`
}

type MilkJournalEntryResponse struct {
	ID                  uint64    `json:"ID"`
	MemberID            uint64    `json:"MemberID"`
	MemberNo            string    `json:"MemberNo"`
	MemberName          string    `json:"MemberName"`
	MilkJournalID       uint64    `json:"MilkJournalID"`
	MilkJournalBatchID  uint64    `json:"MilkJournalBatchID"`
	RouteID             uint64    `json:"RouteID"`
	RouteName           string    `json:"RouteName"`
	MilkDeliveryShiftID uint64    `json:"MilkDeliveryShiftID"`
	MilkDeliveryShift   string    `json:"MilkDeliveryShift"`
	Status              string    `json:"Status"`
	JournalDate         time.Time `json:"JournalDate"`
	Quantity            float64   `json:"Quantity"`
	TransporterID       uint64    `json:"TransporterID"`
	RouteCenterID       uint64    `json:"RouteCenterID"`
	CanID               uint64    `json:"CanID"`
	CreatedAt           time.Time `json:"CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt"`
}

type UpdateMilkJournalEntryRequest struct {
	MemberID            uint64  `json:"MemberID" validate:"required"`
	MilkJournalID       uint64  `json:"MilkJournalID" validate:"required"`
	MilkJournalBatchID  uint64  `json:"MilkJournalBatchID"`
	RouteID             uint64  `json:"RouteID" validate:"required"`
	MilkDeliveryShiftID uint64  `json:"MilkDeliveryShiftID" validate:"required"`
	Status              string  `json:"Status"`
	JournalDate         string  `json:"JournalDate" validate:"required,datetime"`
	Quantity            float64 `json:"Quantity" validate:"required"`
	TransporterID       uint64  `json:"TransporterID"`
	RouteCenterID       uint64  `json:"RouteCenterID"`
	CanID               uint64  `json:"CanID"`
}

type StrayMilkCollectionResponse struct {
	ID                uint64    `json:"ID"`
	MemberID          uint64    `json:"MemberID"`
	MemberNo          string    `json:"MemberNo"`
	MemberName        string    `json:"MemberName"`
	MemberRouteID     uint64    `json:"MemberRouteID"`
	MemberRoute       string    `json:"MemberRoute"`
	JournalRouteID    uint64    `json:"JournalRouteID"`
	StrayRoute        string    `json:"StrayRoute"`
	Quantity          float64   `json:"Quantity"`
	JournalDate       time.Time `json:"JournalDate"`
	MilkDeliveryShift string    `json:"MilkDeliveryShift"`
	CreatedAt         time.Time `json:"CreatedAt"`
}
