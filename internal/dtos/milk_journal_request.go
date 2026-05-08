package dtos

type CreateMilkJournalRequest struct {
	Journal             string `json:"journal" validate:"required,max=255"`
	JournalDate         string `json:"journal_date" validate:"required,datetime"`
	MilkDeliveryShiftID uint64 `json:"milk_delivery_shift_id" validate:"required"`
	RouteID             uint64 `json:"route_id" validate:"required"`
	Confirmed           bool   `json:"confirmed"`
}

type UpdateMilkJournalRequest struct {
	Journal             string `json:"journal" validate:"required,max=255"`
	JournalDate         string `json:"journal_date" validate:"required,datetime"`
	MilkDeliveryShiftID uint64 `json:"milk_delivery_shift_id" validate:"required"`
	RouteID             uint64 `json:"route_id" validate:"required"`
	Confirmed           bool   `json:"confirmed"`
}
