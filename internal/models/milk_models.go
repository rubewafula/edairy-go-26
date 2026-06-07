package models

import "time"

type MilkReject struct {
	BaseModel
	RouteID             uint64    `gorm:"column:route_id"`
	Quantity            float64   `gorm:"column:quantity"`
	TransactionDate     time.Time `gorm:"index;column:transaction_date"`
	Reason              string    `gorm:"column:reason"`
	Description         string    `gorm:"column:description"`
	Confirmed           int       `gorm:"column:confirmed"`
	MemberID            uint64    `gorm:"index;column:member_id"`
	TransporterID       uint64    `gorm:"index;column:transporter_id"`
	CanID               uint64    `gorm:"index;column:can_id"`
	MilkDeliveryShiftID uint64    `gorm:"index;column:milk_delivery_shift_id"`
}

type MilkSpecialRate struct {
	BaseModel
	MemberID              *uint64 `gorm:"index;column:member_id"`
	RouteID               *uint64 `gorm:"index;column:route_id"`
	Rate                  float64 `gorm:"column:rate"`
	MonthlyPayDateRangeID uint64  `gorm:"index;column:pay_date_range_id"`
	Confirmed             int     `gorm:"column:confirmed"`
}

func (MilkSpecialRate) TableName() string {
	return "milk_special_rates"
}

type DefaultMilkRate struct {
	BaseModel
	RouteID  *uint64 `gorm:"column:route_id"`
	MemberID *uint64 `gorm:"column:member_id"`
	Rate     float64 `gorm:"column:rate"`
}

func (DefaultMilkRate) TableName() string {
	return "default_milk_rates"
}

type MilkCooler struct {
	BaseModel
	JournalID          string    `gorm:"column:journal_id"`
	TransactionDate    time.Time `gorm:"index;column:transaction_date"`
	Quantity           float64   `gorm:"column:quantity"`
	RegistrationNumber string    `gorm:"column:registration_number"`
	ScaleNumber        string    `gorm:"column:scale_number"`
	Shift              string    `gorm:"column:shift"`
	Confirmed          int       `gorm:"column:confirmed"`
	UserID             uint64    `gorm:"index;column:user_id"`
	MilkBar            float64   `gorm:"column:milk_bar"`
	SiteID             uint64    `gorm:"index;column:site_id"`
}

type MilkDelivery struct {
	BaseModel
	DeliveryNoteNumber string    `gorm:"index;column:delivery_note_number"`
	CustomerID         uint64    `gorm:"index;column:customer_id"`
	ProductGradeID     uint64    `gorm:"index;column:product_grade_id"`
	QuantityAccepted   float64   `gorm:"column:quantity_accepted"`
	Cooler             string    `gorm:"column:cooler"`
	Invoiced           int       `gorm:"column:invoiced"`
	TransactionDate    time.Time `gorm:"index;column:transaction_date"`
	Amount             float64   `gorm:"column:amount"`
	AmountPaid         float64   `gorm:"column:amount_paid"`
	RouteID            uint64    `gorm:"index;column:route_id"`
	Confirmed          int       `gorm:"column:confirmed"`
	Processed          string    `gorm:"column:processed"`
	TransporterID      uint64    `gorm:"index;column:transporter_id"`
}

func (MilkDelivery) TableName() string {
	return "milk_deliveries"
}

type MilkDeliveryAcceptance struct {
	BaseModel
	DeliveryNoteNumber string    `gorm:"uniqueIndex;column:delivery_note_number"`
	CustomerID         uint64    `gorm:"index;column:customer_id"`
	QuantityAccepted   float64   `gorm:"column:quantity_accepted"`
	Cooler             string    `gorm:"column:cooler"`
	Invoiced           int       `gorm:"column:invoiced"`
	TransactionDate    time.Time `gorm:"index;column:transaction_date"`
	Amount             float64   `gorm:"column:amount"`
	AmountPaid         float64   `gorm:"column:amount_paid"`
	RouteID            uint64    `gorm:"index;column:route_id"`
	Confirmed          int       `gorm:"column:confirmed"`
	Processed          string    `gorm:"column:processed"`
	TransporterID      uint64    `gorm:"index;column:transporter_id"`
}

type MilkDeliveryItem struct {
	BaseModel
	DeliveryID uint64    `gorm:"index;column:delivery_id"`
	Quantity   float64   `gorm:"column:quantity"`
	Rate       float64   `gorm:"column:rate"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	GradeID    uint64    `gorm:"index;column:grade_id"`
}

type MilkLocalSale struct {
	BaseModel
	Quantity        float64   `gorm:"column:quantity"`
	Rate            float64   `gorm:"column:rate"`
	GradeID         uint64    `gorm:"index;column:grade_id"`
	RefNumber       string    `gorm:"uniqueIndex;column:ref_number"`
	TransactionDate time.Time `gorm:"index;column:transaction_date"`
	TransporterID   uint64    `gorm:"index;column:transporter_id"`
	Amount          float64   `gorm:"column:amount"`
}

func (MilkLocalSale) TableName() string {
	return "milk_local_sales"
}

type MilkSale struct {
	BaseModel
	Quantity float64   `gorm:"column:quantity"`
	Price    float64   `gorm:"column:price"`
	Amount   float64   `gorm:"column:amount"`
	Buyer    string    `gorm:"column:buyer"`
	Date     time.Time `gorm:"index;column:transaction_date"`
}

type DailyMilkVariance struct {
	ID               uint64    `gorm:"column:id"`
	Transporter      string    `gorm:"column:transporter"`
	Day              time.Time `gorm:"column:day"`
	Month            string    `gorm:"column:month"`
	FieldCollections float64   `gorm:"column:field_collections"`
	MCC              float64   `gorm:"column:mcc"`
	CashSales        float64   `gorm:"column:cash_sales"`
	CreditSales      float64   `gorm:"column:credit_sales"`
	Rejects          float64   `gorm:"column:rejects"`
	Balance          float64   `gorm:"column:balance"`
}

func (DailyMilkVariance) TableName() string {
	return "daily_milk_variances"
}

type MilkJournal struct {
	BaseModel
	Journal             string     `gorm:"column:journal"`
	JournalDate         *time.Time `gorm:"index;column:journal_date"`
	MilkDeliveryShiftID uint64     `gorm:"index;column:milk_delivery_shift_id"`
	RouteID             uint64     `gorm:"index;column:route_id"`
	UserID              uint64     `gorm:"column:user_id"`
	TransporterID       uint64     `gorm:"column:transporter_id"`
	Confirmed           bool       `gorm:"column:confirmed"`
}

type MilkJournalBatch struct {
	BaseModel
	MilkJournalID uint64 `gorm:"column:milk_journal_id"`
	BatchNo       string `gorm:"column:batch_no"`
}

func (MilkJournalBatch) TableName() string {
	return "milk_journal_batches"
}

type MilkJournalEntry struct {
	BaseModel
	MemberID           uint64  `gorm:"index;column:member_id"`
	MilkJournalID      uint64  `gorm:"index;column:milk_journal_id"`
	MilkJournalBatchID uint64  `gorm:"index;column:milk_journal_batch_id"`
	Status             string  `gorm:"column:status"`
	Quantity           float64 `gorm:"type:decimal(18,2);column:quantity"`
	RouteCenterID      uint64  `gorm:"index;column:route_center_id"`
	CanID              uint64  `gorm:"index;column:can_id"`
}

func (MilkJournalEntry) TableName() string {
	return "milk_journal_entries"
}

type MilkDeliveryShift struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (MilkDeliveryShift) TableName() string {
	return "milk_delivery_shifts"
}

type MilkCollection struct {
	BaseModel
	MemberID           uint64    `gorm:"index;column:member_id"`
	MilkJournalTableID uint64    `gorm:"index;column:milk_journal_table_id"`
	RouteID            uint64    `gorm:"index;column:route_id"`
	MilkJournalBatchID uint64    `gorm:"index;column:milk_journal_batch_id"`
	ShiftID            uint64    `gorm:"index;column:milk_delivery_shift_id"`
	Status             string    `gorm:"column:status"`
	JournalDate        time.Time `gorm:"index;column:journal_date"`
	Quantity           float64   `gorm:"column:quantity"`
	TransporterID      uint64    `gorm:"index;column:transporter_id"`
	RouteCenterID      uint64    `gorm:"index;column:route_center_id"`
	CanID              uint64    `gorm:"index;column:can_id"`
}

type MilkCan struct {
	BaseModel
	CanID      string  `gorm:"uniqueIndex;column:can_id"`
	CanType    string  `gorm:"column:can_type"`
	CanSize    float64 `gorm:"column:can_size"`
	Units      string  `gorm:"column:units"`
	TareWeight float64 `gorm:"column:tare_weight"`
	RouteID    uint64  `gorm:"index;column:route_id"`
}

type CoolerMilkCollection struct {
	BaseModel
	TransactionDate     time.Time `gorm:"index;column:transaction_date"`
	Quantity            float64   `gorm:"column:quantity"`
	TransportVehicleID  uint64    `gorm:"index;column:transport_vehicle_id"`
	MilkDeliveryShiftID uint64    `gorm:"index;column:milk_delivery_shift_id"`
	Confirmed           int       `gorm:"column:confirmed"`
	SiteID              uint64    `gorm:"index;column:site_id"`
	TransporterID       uint64    `gorm:"index;column:transporter_id"`
	RouteID             uint64    `gorm:"index;column:route_id"`
}

func (MilkCan) TableName() string {
	return "milk_cans"
}
