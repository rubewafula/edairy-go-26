package models

import "time"

// StoreTransaction represents a single transaction affecting store stock.
type StoreTransaction struct {
	BaseModel
	TransactionID   uint64    `gorm:"not null;index" json:"transaction_id"`
	StoreID         uint64    `gorm:"not null;index" json:"store_id"`
	StoreItemID     uint64    `gorm:"not null;index" json:"store_item_id"`
	TransactionType string    `gorm:"type:enum('PURCHASE','SALE','TRANSFER_IN','TRANSFER_OUT','ADJUSTMENT_IN','ADJUSTMENT_OUT','RETURN_IN','RETURN_OUT');not null;index" json:"transaction_type"`
	Quantity        float64   `gorm:"type:decimal(18,4);not null" json:"quantity"`
	UnitCost        float64   `gorm:"type:decimal(18,4);default:0.0000" json:"unit_cost"`
	SellingPrice    float64   `gorm:"type:decimal(18,4);default:0.0000" json:"selling_price"`
	ReferenceType   string    `gorm:"type:varchar(50);index" json:"reference_type"`
	ReferenceID     uint64    `json:"reference_id"`
	BalanceAfter    float64   `gorm:"type:decimal(18,4);not null;default:0.0000" json:"balance_after"`
	Notes           string    `gorm:"type:text" json:"notes"`
	TransactedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"transacted_at"`
}

// TableName specifies the table name for StoreTransaction model.
func (StoreTransaction) TableName() string {
	return "store_transactions"
}

// StoreTransactionType defines the enum values for transaction_type.
type StoreTransactionType string

const (
	AdjustmentIn StoreTransactionType = "ADJUSTMENT_IN"
)
