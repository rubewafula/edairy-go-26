package models

import (
	"time"
)

type Livestock struct {
	BaseModel
	MemberID            *uint64 `gorm:"column:member_id;index"`
	TagNo               *string `gorm:"column:tag_no;uniqueIndex;size:100"`
	LivestockCategoryID uint64  `gorm:"column:livestock_category_id;index;not null"`
	LivestockBreedID    *uint64 `gorm:"column:livestock_breed_id;index"`
	LivestockName       *string `gorm:"column:livestock_name;size:150"`

	Gender string `gorm:"column:gender;type:enum('male','female');not null"`
	Color  string `gorm:"column:color;size:100"`

	BirthDate    time.Time `gorm:"column:birth_date;type:date"`
	PurchaseDate time.Time `gorm:"column:purchase_date;type:date"`

	Source string `gorm:"column:source;type:enum('born','purchased','donated','transferred');default:'born'"`
	Status string `gorm:"column:status;type:enum('active','sold','dead','slaughtered','missing');default:'active';index"`

	Weight          *float64 `gorm:"column:weight;type:decimal(10,2)"`
	InsuranceNumber *string  `gorm:"column:insurance_number;size:100"`
	Notes           *string  `gorm:"column:notes;type:text"`
}

func (Livestock) TableName() string {
	return "livestocks"
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

type LivestockImportError struct {
	BaseModel
	RowData string `gorm:"column:row_data;type:text"`
	Error   string `gorm:"column:error;type:text"`
}

func (LivestockImportError) TableName() string {
	return "livestock_import_errors"
}
