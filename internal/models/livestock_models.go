package models

import (
	"time"
)

type Livestock struct {
	BaseModel
	TagNo       string    `gorm:"column:tag_no;uniqueIndex;not null" json:"tag_no"`
	BreedID     uint64    `gorm:"column:breed_id" json:"breed_id"`
	Gender      string    `gorm:"column:gender" json:"gender"`
	DateOfBirth time.Time `gorm:"column:date_of_birth" json:"date_of_birth"`
	Status      string    `gorm:"column:status;default:ACTIVE" json:"status"`
	Description string    `gorm:"column:description" json:"description"`
}

func (Livestock) TableName() string {
	return "livestock"
}

type LivestockCategory struct {
	BaseModel
	CategoryName string `gorm:"column:category_name;not null"`
	Description  string `gorm:"column:description"`
}

func (LivestockCategory) TableName() string {
	return "livestock_categories"
}

type LivestockBreed struct {
	BaseModel
	LivestockCategoryID uint64 `gorm:"column:livestock_category_id;not null"`
	BreedName           string `gorm:"column:breed_name;not null"`
	Description         string `gorm:"column:description"`
}

func (LivestockBreed) TableName() string {
	return "livestock_breeds"
}

type LivestockDeath struct {
	BaseModel
	LivestockID    uint64    `gorm:"column:livestock_id;not null"`
	DeathDate      time.Time `gorm:"column:death_date;not null"`
	CauseOfDeath   string    `gorm:"column:cause_of_death"`
	DisposalMethod string    `gorm:"column:disposal_method"`
	Remarks        string    `gorm:"column:remarks"`
}

type LivestockFeeding struct {
	BaseModel
	LivestockID uint64    `gorm:"column:livestock_id;not null"`
	FeedName    string    `gorm:"column:feed_name;not null"`
	Quantity    float64   `gorm:"column:quantity;type:decimal(12,2);not null"`
	Unit        string    `gorm:"column:unit;not null"`
	FeedingDate time.Time `gorm:"column:feeding_date;not null"`
	Cost        float64   `gorm:"column:cost;type:decimal(12,2);default:0.00"`
	Notes       string    `gorm:"column:notes"`
}

type LivestockHealthRecord struct {
	BaseModel
	LivestockID   uint64     `gorm:"column:livestock_id;not null"`
	RecordType    string     `gorm:"column:record_type;type:enum('vaccination','treatment','disease','checkup','deworming');not null"`
	Diagnosis     string     `gorm:"column:diagnosis"`
	Medication    string     `gorm:"column:medication"`
	Dosage        string     `gorm:"column:dosage"`
	Veterinarian  string     `gorm:"column:veterinarian"`
	TreatmentDate time.Time  `gorm:"column:treatment_date;not null"`
	NextVisitDate *time.Time `gorm:"column:next_visit_date"`
	Notes         string     `gorm:"column:notes"`
}

type LivestockMovement struct {
	BaseModel
	LivestockID  uint64    `gorm:"column:livestock_id;not null"`
	FromLocation string    `gorm:"column:from_location"`
	ToLocation   string    `gorm:"column:to_location;not null"`
	MovementDate time.Time `gorm:"column:movement_date;not null"`
	MovementType string    `gorm:"column:movement_type;type:enum('transfer','sale','grazing','medical','show');default:'transfer'"`
	Transporter  string    `gorm:"column:transporter"`
	Remarks      string    `gorm:"column:remarks"`
}

type LivestockPhoto struct {
	BaseModel
	LivestockID uint64 `gorm:"column:livestock_id;not null"`
	PhotoURL    string `gorm:"column:photo_url;not null"`
	Description string `gorm:"column:description"`
}

type LivestockProductionRecord struct {
	BaseModel
	LivestockID    uint64    `gorm:"column:livestock_id;not null"`
	ProductionType string    `gorm:"column:production_type;type:enum('milk','eggs','wool','meat');not null"`
	ProductionDate time.Time `gorm:"column:production_date;not null"`
	Quantity       float64   `gorm:"column:quantity;type:decimal(12,2);not null"`
	Unit           string    `gorm:"column:unit;not null"`
	Remarks        string    `gorm:"column:remarks"`
}

func (LivestockProductionRecord) TableName() string {
	return "livestock_production_records"
}

type LivestockSale struct {
	BaseModel
	LivestockID   uint64    `gorm:"column:livestock_id;not null"`
	CustomerID    uint64    `gorm:"column:customer_id"`
	SaleDate      time.Time `gorm:"column:sale_date;not null"`
	Quantity      int       `gorm:"column:quantity;default:1"`
	SalePrice     float64   `gorm:"column:sale_price;type:decimal(12,2);not null"`
	PaymentStatus string    `gorm:"column:payment_status;type:enum('pending','paid','partial');default:'pending'"`
	Notes         string    `gorm:"column:notes"`
}

type LivestockWeightRecord struct {
	BaseModel
	LivestockID uint64    `gorm:"column:livestock_id;not null"`
	Weight      float64   `gorm:"column:weight;type:decimal(10,2);not null"`
	Remarks     string    `gorm:"column:remarks"`
	RecordedAt  time.Time `gorm:"column:recorded_at;not null"`
}
