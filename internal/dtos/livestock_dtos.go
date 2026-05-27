package dtos

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/models"
)

// Livestock
type CreateLivestockRequest struct {
	MemberID            *uint64 `json:"member_id"`
	TagNo               *string `json:"tag_no"`
	LivestockCategoryID uint64  `json:"livestock_category_id" binding:"required"`
	LivestockBreedID    *uint64 `json:"livestock_breed_id"`
	LivestockName       *string `json:"livestock_name"`

	Gender string `json:"gender" binding:"required"` // male | female
	Color  string `json:"color"`

	BirthDate    string `json:"birth_date"`
	PurchaseDate string `json:"purchase_date"`

	Source *string  `json:"source"` // born | purchased | donated | transferred
	Weight *float64 `json:"weight"`

	InsuranceNumber *string `json:"insurance_number"`
	Notes           *string `json:"notes"`
}

type UpdateLivestockRequest struct {
	TagNo       string `json:"tag_no"`
	BreedID     uint64 `json:"breed_id"`
	Gender      string `json:"gender" validate:"omitempty,oneof=MALE FEMALE"`
	DateOfBirth string `json:"date_of_birth"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type LivestockResponse struct {
	models.BaseModel
	TagNo        string    `json:"tag_no"`
	BreedName    string    `json:"breed_name"`
	CategoryName string    `json:"category_name"`
	Gender       string    `json:"gender"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
}

// Category
type CreateLivestockCategoryRequest struct {
	CategoryName string `json:"category_name" validate:"required"`
	Description  string `json:"description"`
}

type LivestockCategoryResponse struct {
	ID           uint64    `json:"id"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateLivestockCategoryRequest struct {
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
	// No validation here, as it's an update
}

// Breed
type CreateLivestockBreedRequest struct {
	LivestockCategoryID uint64 `json:"livestock_category_id" validate:"required"`
	BreedName           string `json:"breed_name" validate:"required"`
	Description         string `json:"description"`
}

type LivestockBreedResponse struct {
	ID           uint64    `json:"id"`
	CategoryName string    `json:"category_name"`
	BreedName    string    `json:"breed_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateLivestockBreedRequest struct {
	LivestockCategoryID uint64 `json:"livestock_category_id"`
	BreedName           string `json:"breed_name"`
	Description         string `json:"description"`
}

// Death
type CreateLivestockDeathRequest struct {
	LivestockID    uint64 `json:"livestock_id" validate:"required"`
	DeathDate      string `json:"death_date" validate:"required"`
	CauseOfDeath   string `json:"cause_of_death"`
	DisposalMethod string `json:"disposal_method"`
	Remarks        string `json:"remarks"`
}

type UpdateLivestockDeathRequest struct {
	DeathDate      string `json:"death_date"`
	CauseOfDeath   string `json:"cause_of_death"`
	DisposalMethod string `json:"disposal_method"`
	Remarks        string `json:"remarks"`
}

type LivestockDeathResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	DeathDate      time.Time `json:"death_date"`
	CauseOfDeath   string    `json:"cause_of_death"`
	DisposalMethod string    `json:"disposal_method"`
	Remarks        string    `json:"remarks"`
}

// Feeding
type CreateLivestockFeedingRequest struct {
	LivestockID uint64  `json:"livestock_id" validate:"required"`
	FeedName    string  `json:"feed_name" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required"`
	Unit        string  `json:"unit" validate:"required"`
	FeedingDate string  `json:"feeding_date" validate:"required"`
	Cost        float64 `json:"cost"`
	Notes       string  `json:"notes"`
}

type UpdateLivestockFeedingRequest struct {
	FeedName    string  `json:"feed_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	FeedingDate string  `json:"feeding_date"`
	Cost        float64 `json:"cost"`
	Notes       string  `json:"notes"`
}

type LivestockFeedingResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	FeedName       string    `json:"feed_name"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	FeedingDate    time.Time `json:"feeding_date"`
	Cost           float64   `json:"cost"`
	Notes          string    `json:"notes"`
}

// Health Record
type CreateLivestockHealthRequest struct {
	LivestockID   uint64 `json:"livestock_id" validate:"required"`
	RecordType    string `json:"record_type" validate:"required,oneof=vaccination treatment disease checkup deworming"`
	Diagnosis     string `json:"diagnosis"`
	Medication    string `json:"medication"`
	Dosage        string `json:"dosage"`
	Veterinarian  string `json:"veterinarian"`
	TreatmentDate string `json:"treatment_date" validate:"required"`
	NextVisitDate string `json:"next_visit_date"`
	Notes         string `json:"notes"`
}

type UpdateLivestockHealthRequest struct {
	RecordType    string `json:"record_type" validate:"oneof=vaccination treatment disease checkup deworming"`
	Diagnosis     string `json:"diagnosis"`
	Medication    string `json:"medication"`
	Dosage        string `json:"dosage"`
	Veterinarian  string `json:"veterinarian"`
	TreatmentDate string `json:"treatment_date"`
	NextVisitDate string `json:"next_visit_date"`
	Notes         string `json:"notes"`
}

type LivestockHealthResponse struct {
	models.BaseModel
	LivestockID    uint64     `json:"livestock_id"`
	LivestockTagNo string     `json:"livestock_tag_no"`
	RecordType     string     `json:"record_type"`
	Diagnosis      string     `json:"diagnosis"`
	Medication     string     `json:"medication"`
	Dosage         string     `json:"dosage"`
	Veterinarian   string     `json:"veterinarian"`
	TreatmentDate  time.Time  `json:"treatment_date"`
	NextVisitDate  *time.Time `json:"next_visit_date,omitempty"`
	Notes          string     `json:"notes"`
}

// Movement
type CreateLivestockMovementRequest struct {
	LivestockID  uint64 `json:"livestock_id" validate:"required"`
	FromLocation string `json:"from_location"`
	ToLocation   string `json:"to_location" validate:"required"`
	MovementDate string `json:"movement_date" validate:"required"`
	MovementType string `json:"movement_type" validate:"oneof=transfer sale grazing medical show"`
	Transporter  string `json:"transporter"`
	Remarks      string `json:"remarks"`
}

type UpdateLivestockMovementRequest struct {
	FromLocation string `json:"from_location"`
	ToLocation   string `json:"to_location"`
	MovementDate string `json:"movement_date"`
	MovementType string `json:"movement_type" validate:"oneof=transfer sale grazing medical show"`
	Transporter  string `json:"transporter"`
	Remarks      string `json:"remarks"`
}

type LivestockMovementResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	FromLocation   string    `json:"from_location"`
	ToLocation     string    `json:"to_location"`
	MovementDate   time.Time `json:"movement_date"`
	MovementType   string    `json:"movement_type"`
	Transporter    string    `json:"transporter"`
	Remarks        string    `json:"remarks"`
}

// Photo
type CreateLivestockPhotoRequest struct {
	LivestockID uint64 `json:"livestock_id" validate:"required"`
	PhotoURL    string `json:"photo_url" validate:"required"`
	Description string `json:"description"`
}

type UpdateLivestockPhotoRequest struct {
	PhotoURL    string `json:"photo_url"`
	Description string `json:"description"`
}

type LivestockPhotoResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	PhotoURL       string    `json:"photo_url"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Production
type CreateLivestockProductionRequest struct {
	LivestockID    uint64  `json:"livestock_id" validate:"required"`
	ProductionType string  `json:"production_type" validate:"required,oneof=milk eggs wool meat"`
	ProductionDate string  `json:"production_date" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	Remarks        string  `json:"remarks"`
}

type UpdateLivestockProductionRequest struct {
	ProductionType string  `json:"production_type" validate:"oneof=milk eggs wool meat"`
	ProductionDate string  `json:"production_date"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Remarks        string  `json:"remarks"`
}

type LivestockProductionResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	ProductionType string    `json:"production_type"`
	ProductionDate time.Time `json:"production_date"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	Remarks        string    `json:"remarks"`
}

// Sale
type CreateLivestockSaleRequest struct {
	LivestockID   uint64  `json:"livestock_id" validate:"required"`
	CustomerID    uint64  `json:"customer_id"`
	SaleDate      string  `json:"sale_date" validate:"required"`
	Quantity      int     `json:"quantity"`
	SalePrice     float64 `json:"sale_price" validate:"required"`
	PaymentStatus string  `json:"payment_status" validate:"oneof=pending paid partial"`
	Notes         string  `json:"notes"`
}

type UpdateLivestockSaleRequest struct {
	CustomerID    uint64  `json:"customer_id"`
	SaleDate      string  `json:"sale_date"`
	Quantity      int     `json:"quantity"`
	SalePrice     float64 `json:"sale_price"`
	PaymentStatus string  `json:"payment_status" validate:"oneof=pending paid partial"`
	Notes         string  `json:"notes"`
}

type LivestockSaleResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	CustomerID     uint64    `json:"customer_id"`
	SaleDate       time.Time `json:"sale_date"`
	Quantity       int       `json:"quantity"`
	SalePrice      float64   `json:"sale_price"`
	PaymentStatus  string    `json:"payment_status"`
	Notes          string    `json:"notes"`
}

// Weight
type CreateLivestockWeightRequest struct {
	LivestockID uint64  `json:"livestock_id" validate:"required"`
	Weight      float64 `json:"weight" validate:"required"`
	RecordedAt  string  `json:"recorded_at" validate:"required"`
	Remarks     string  `json:"remarks"`
}

type UpdateLivestockWeightRequest struct {
	Weight     float64 `json:"weight"`
	RecordedAt string  `json:"recorded_at"`
	Remarks    string  `json:"remarks"`
}

type LivestockWeightResponse struct {
	models.BaseModel
	LivestockID    uint64    `json:"livestock_id"`
	LivestockTagNo string    `json:"livestock_tag_no"`
	Weight         float64   `json:"weight"`
	RecordedAt     time.Time `json:"recorded_at"`
	Remarks        string    `json:"remarks"`
}

// Generic Response structure used for listings
type LivestockGenericResponse struct {
	ID          uint64      `json:"id"`
	LivestockID uint64      `json:"livestock_id"`
	TagNo       string      `json:"tag_no"` // Joined field
	Date        time.Time   `json:"date"`
	Type        string      `json:"type,omitempty"`
	Value       float64     `json:"value,omitempty"`
	Unit        string      `json:"unit,omitempty"`
	Details     string      `json:"details,omitempty"`
	Status      string      `json:"status,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	Data        interface{} `json:"metadata,omitempty"`
}
