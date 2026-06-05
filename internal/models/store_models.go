package models

import "time"

type Store struct {
	BaseModel
	Name        string `gorm:"column:store"`
	Description string `gorm:"column:description"`
}

type StoreInventory struct {
	BaseModel
	InventoryName string `gorm:"column:inventory_name"`
	CategoryID    uint64 `gorm:"column:category_id"`
	IsActive      bool   `gorm:"column:is_active;default:1"`
	Description   string `gorm:"column:description"`
}

func (StoreInventory) TableName() string {
	return "store_inventories"
}

type StoreItem struct {
	BaseModel
	Description               string  `gorm:"column:description"`
	ReorderPoint              int     `gorm:"column:reorder_point"`
	DefaultBuyingPrice        float64 `gorm:"column:default_buying_price"`
	DefaultSellingPrice       float64 `gorm:"column:default_selling_price"`
	Status                    string  `gorm:"column:status;default:0"`
	Thumbnail                 string  `gorm:"column:thumbnail"`
	ItemName                  string  `gorm:"column:item_name"`
	SKU                       string  `gorm:"column:sku"`
	Barcode                   string  `gorm:"column:barcode"`
	UnitID                    int64   `gorm:"column:unit_id"`
	DefaultSellingPriceCredit float64 `gorm:"column:default_selling_price_credit"`
	StoreInventoryID          uint64  `gorm:"column:store_inventory_id"`
}

func (StoreItem) TableName() string {
	return "store_items"
}

type StoreStock struct {
	BaseModel
	ItemID             uint64  `gorm:"column:item_id"`
	StoreID            uint64  `gorm:"column:store_id"`
	Quantity           float64 `gorm:"column:quantity"`
	Unit               string  `gorm:"column:unit;default:KG"`
	BuyingPrice        float64 `gorm:"column:buying_price"`
	SellingPrice       float64 `gorm:"column:selling_price"`
	CreditSellingPrice float64 `gorm:"column:credit_selling_price"`
}

func (StoreStock) TableName() string {
	return "store_stocks"
}

type StoreStockTaking struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	StockTakeNo      string    `gorm:"column:stock_take_no"`
	StoreID          uint64    `gorm:"column:store_id"`
	ItemID           uint64    `gorm:"column:item_id"`
	SystemQuantity   float64   `gorm:"column:system_quantity"`
	PhysicalQuantity float64   `gorm:"column:physical_quantity"`
	VarianceQuantity float64   `gorm:"column:variance_quantity"`
	Remarks          string    `gorm:"column:remarks"`
	StockTakeDate    time.Time `gorm:"column:stock_take_date"`
	CreatedBy        uint64    `gorm:"column:created_by"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (StoreStockTaking) TableName() string {
	return "store_stock_takings"
}

type StoreStockMovementType struct {
	BaseModel
	MovementCode string `gorm:"uniqueIndex;column:movement_code"`
	MovementName string `gorm:"column:movement_name"`
	Direction    string `gorm:"column:direction"` // enum('IN','OUT')
	AffectsStock bool   `gorm:"column:affects_stock;default:1"`
	Description  string `gorm:"column:description"`
	IsSystem     bool   `gorm:"column:is_system;default:1"`
}

func (StoreStockMovementType) TableName() string {
	return "store_stock_movement_types"
}

type StoreSale struct {
	BaseModel
	TotalAmount   float64 `gorm:"column:total_amount"`
	AmountPaid    float64 `gorm:"column:amount_paid;default:0.00"`
	AmountDue     float64 `gorm:"column:amount_due;default:0.00"`
	Reference     string  `gorm:"column:reference"`
	StoreID       uint64  `gorm:"column:store_id"`
	SaleType      string  `gorm:"column:sale_type;default:cash"`
	CustomerID    uint64  `gorm:"column:customer_id"`
	CustomerType  string  `gorm:"column:customer_type"`
	TransactionID uint64  `gorm:"column:transaction_id"`
}

func (StoreSale) TableName() string {
	return "store_sales"
}

type StoreSaleItem struct {
	BaseModel
	ItemID      uint64  `gorm:"column:item_id"`
	Quantity    int     `gorm:"column:quantity"`
	UnitPrice   float64 `gorm:"column:unit_price"`
	Total       float64 `gorm:"column:total"`
	StoreSaleID uint64  `gorm:"column:store_sale_id"`
}

func (StoreSaleItem) TableName() string {
	return "store_sale_items"
}

type StoreItemUnit struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;column:name"`
	Symbol      string `gorm:"uniqueIndex;column:symbol"`
	Description string `gorm:"column:description"`
}

func (StoreItemUnit) TableName() string {
	return "store_item_units"
}

type Inventory struct {
	BaseModel
	InventoryName      string    `gorm:"column:inventory_name"`
	MovementType       string    `gorm:"column:movement_type"`
	Direction          string    `gorm:"column:direction"`
	DateCaptured       time.Time `gorm:"index;column:date_captured"`
	Quantity           float64   `gorm:"column:quantity"`
	BuyingPrice        float64   `gorm:"column:buying_price"`
	SellingPriceCash   float64   `gorm:"column:selling_price_cash"`
	SellingPriceCredit float64   `gorm:"column:selling_price_credit"`
	ReorderLevel       float64   `gorm:"column:reorder_level"`
	InventoryCategory  string    `gorm:"column:inventory_category"`
	InvoiceNumber      string    `gorm:"column:invoice_number"`
	ValuationMethod    string    `gorm:"column:valuation_method"`
	TransactionID      uint64    `gorm:"index;column:transaction_id"`
	VendorID           uint64    `gorm:"index;column:vendor_id"`
	Status             string    `gorm:"column:status"`
}

type InventoryStockMovement struct {
	BaseModel
	InventoryID uint64  `gorm:"index;column:inventory_id"`
	OpeningBal  float64 `gorm:"column:openning_bal"`
	Receipts    float64 `gorm:"column:receipts"`
	Sales       float64 `gorm:"column:sales"`
	Transfers   float64 `gorm:"column:transfers"`
	Adjustments float64 `gorm:"column:adjustments"`
	Closing     float64 `gorm:"column:closing"`
}

type InventoryOpeningBalance struct {
	BaseModel
	MemberID        uint64    `gorm:"index;column:member_id"`
	TransactionDate time.Time `gorm:"index;column:transaction_date"`
	Type            string    `gorm:"column:type"`
	Amount          float64   `gorm:"column:amount"`
	SiteID          uint64    `gorm:"column:site_id"`
}

type InterStoreTransfer struct {
	BaseModel
	FromStoreID  uint64    `gorm:"index;column:from_store_id"`
	ToStoreID    uint64    `gorm:"index;column:to_store_id"`
	Reference    string    `gorm:"uniqueIndex;column:reference"`
	TransferDate time.Time `gorm:"index;column:transfer_date"`
	Status       string    `gorm:"column:status"`
}

func (InterStoreTransfer) TableName() string {
	return "inter_store_transfers"
}

type InterStoreTransferItem struct {
	BaseModel
	TransferID uint64 `gorm:"index;column:inter_store_transfer_id"`
	ItemID     uint64 `gorm:"index;column:item_id"`
	Quantity   string `gorm:"column:quantity"`
	StockID    uint64 `gorm:"column:stock_id"`
}

func (InterStoreTransferItem) TableName() string {
	return "inter_store_transfer_items"
}

type StockAdjustment struct {
	BaseModel
	InventoryID    uint64    `gorm:"index;column:inventory_id"`
	AdjustmentType string    `gorm:"column:adjustment_type"`
	Quantity       float64   `gorm:"column:quantity"`
	Reason         string    `gorm:"column:reason"`
	AdjustmentDate time.Time `gorm:"index;column:adjustment_date"`
}

type StockTransfer struct {
	InterStoreTransfer // Embed for common fields, or define separately if distinct
}
