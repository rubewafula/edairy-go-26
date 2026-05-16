package models

import (
	"time"

	"gorm.io/gorm"
)

type Asset struct {
	BaseModel
	AssetCode               string    `gorm:"uniqueIndex;column:asset_code"`
	AssetName               string    `gorm:"column:asset_name"`
	CategoryID              uint64    `gorm:"column:asset_category_id"`
	SerialNo                string    `gorm:"column:serial_no"`
	Barcode                 string    `gorm:"column:barcode"`
	Manufacturer            string    `gorm:"column:manufacturer"`
	VendorID                uint64    `gorm:"column:vendor_id"`
	PurchaseCost            float64   `gorm:"column:purchase_cost"`
	SalvageValue            float64   `gorm:"column:salvage_value;default:0.00"`
	AcquisitionDate         time.Time `gorm:"column:acquisition_date"`
	UsefulLifeYears         int       `gorm:"column:useful_life_years"`
	DepreciationMethod      string    `gorm:"column:depreciation_method"`
	DepreciationRate        float64   `gorm:"column:depreciation_rate"`
	AccumulatedDepreciation float64   `gorm:"column:accumulated_depreciation;default:0.00"`
	BookValue               float64   `gorm:"column:book_value"`
	WarrantyEndDate         time.Time `gorm:"column:warranty_end_date"`
	CurrentLocation         string    `gorm:"column:current_location"`
	Status                  string    `gorm:"column:status;default:ACTIVE"`
	Loanable                bool      `gorm:"column:loanable;default:0"`
	Comments                string    `gorm:"column:comments"`
}

func (Asset) TableName() string {
	return "fixed_assets"
}

type AssetCategory struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;column:ud" json:"ID"`
	Name        string         `gorm:"column:name" json:"Name"`
	Description string         `gorm:"column:description" json:"Description"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"UpdatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

func (AssetCategory) TableName() string {
	return "asset_categories"
}

type AssetAssignment struct {
	BaseModel
	AssetID        uint64    `gorm:"column:asset_id"`
	AssignedToID   uint64    `gorm:"column:assigned_to_id"`
	AssignedAt     time.Time `gorm:"column:assigned_at"`
	ReturnedAt     time.Time `gorm:"column:returned_at"`
	ConditionNotes string    `gorm:"column:condition_notes"`
	Status         string    `gorm:"column:status;default:ASSIGNED"`
}

func (AssetAssignment) TableName() string {
	return "asset_assignments"
}

type AssetDepreciationEntry struct {
	ID                      uint64         `gorm:"primaryKey;autoIncrement;column:id"`
	AssetID                 uint64         `gorm:"column:asset_id"`
	DepreciationDate        time.Time      `gorm:"column:depreciation_date"`
	DepreciationAmount      float64        `gorm:"column:depreciation_amount"`
	AccumulatedDepreciation float64        `gorm:"column:accumulated_depreciation"`
	BookValue               float64        `gorm:"column:book_value"`
	TransactionID           uint64         `gorm:"column:transaction_id"`
	CreatedAt               time.Time      `gorm:"column:created_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
