package models

import "time"

type SupplierCategory struct {
	BaseModel
	CategoryCode string  `gorm:"uniqueIndex;column:category_code"`
	CategoryName string  `gorm:"column:category_name"`
	Description  string  `gorm:"column:description"`
	Status       string  `gorm:"type:enum('active','inactive');default:'active';column:status"`
	SiteID       *uint64 `gorm:"column:site_id"`
}

type Supplier struct {
	BaseModel
	SupplierCategoryID uint64     `gorm:"column:supplier_category_id"`
	SupplierCode       string     `gorm:"uniqueIndex;column:supplier_code"`
	SupplierType       string     `gorm:"type:enum('individual','company');default:'individual';column:supplier_type"`
	CompanyName        string     `gorm:"column:company_name"`
	Title              string     `gorm:"column:title"`
	FirstName          string     `gorm:"column:first_name"`
	LastName           string     `gorm:"column:last_name"`
	OtherNames         string     `gorm:"column:other_names"`
	Gender             string     `gorm:"type:enum('male','female','other');column:gender"`
	Dob                *time.Time `gorm:"column:dob"`
	IdNo               string     `gorm:"column:id_no"`
	PassportNo         string     `gorm:"column:passport_no"`
	KraPin             string     `gorm:"column:kra_pin"`
	LicenseNumber      string     `gorm:"column:license_number"`
	PhoneNo            string     `gorm:"column:phone_no"`
	AlternativePhoneNo string     `gorm:"column:alternative_phone_no"`
	EmailAddress       string     `gorm:"column:email_address"`
	PostalAddress      string     `gorm:"column:postal_address"`
	PostalCode         string     `gorm:"column:postal_code"`
	Town               string     `gorm:"column:town"`
	Residence          string     `gorm:"column:residence"`
	ContactPerson      string     `gorm:"column:contact_person"`
	OpeningBalance     float64    `gorm:"column:opening_balance;default:0.00"`
	CurrentBalance     float64    `gorm:"column:current_balance;default:0.00"`
	CreditLimit        float64    `gorm:"column:credit_limit;default:0.00"`
	PaymentTermsDays   int        `gorm:"column:payment_terms_days;default:0"`
	Status             string     `gorm:"type:enum('active','inactive','blacklisted','suspended');default:'active';column:status"`
	Notes              string     `gorm:"type:text;column:notes"`
	SiteID             *uint64    `gorm:"column:site_id"`
}

type SupplierContact struct {
	BaseModel
	SupplierID         uint64 `gorm:"column:supplier_id"`
	ContactType        string `gorm:"type:enum('primary','finance','procurement','technical','other');default:'primary';column:contact_type"`
	FullName           string `gorm:"column:full_name"`
	Designation        string `gorm:"column:designation"`
	PhoneNo            string `gorm:"column:phone_no"`
	AlternativePhoneNo string `gorm:"column:alternative_phone_no"`
	EmailAddress       string `gorm:"column:email_address"`
	IsDefault          string `gorm:"type:enum('yes','no');default:'no';column:is_default"`
	Notes              string `gorm:"column:notes"`
}

type SupplierBankAccount struct {
	BaseModel
	SupplierID    uint64 `gorm:"column:supplier_id"`
	BankID        uint64 `gorm:"column:bank_id"`
	BankBranchID  uint64 `gorm:"column:bank_branch_id"`
	AccountName   string `gorm:"column:account_name"`
	AccountNumber string `gorm:"column:account_number"`
	AccountType   string `gorm:"type:enum('bank','mobile_money');default:'bank';column:account_type"`
	CurrencyCode  string `gorm:"column:currency_code;default:'KES'"`
	SwiftCode     string `gorm:"column:swift_code"`
	MobileMoneyNo string `gorm:"column:mobile_money_no"`
	IsDefault     string `gorm:"type:enum('yes','no');default:'no';column:is_default"`
	Status        string `gorm:"type:enum('active','inactive');default:'active';column:status"`
}

type SupplierDocument struct {
	BaseModel
	SupplierID     uint64     `gorm:"column:supplier_id"`
	DocumentType   string     `gorm:"type:enum('national_id','passport','kra_pin','business_registration','license','contract','bank_letter','certificate','other');column:document_type"`
	DocumentNumber string     `gorm:"column:document_number"`
	FileName       string     `gorm:"column:file_name"`
	FilePath       string     `gorm:"column:file_path"`
	IssueDate      *time.Time `gorm:"column:issue_date"`
	ExpiryDate     *time.Time `gorm:"column:expiry_date"`
	Verified       string     `gorm:"type:enum('yes','no');default:'no';column:verified"`
	VerifiedBy     *uint64    `gorm:"column:verified_by"`
	VerifiedAt     *time.Time `gorm:"column:verified_at"`
	Notes          string     `gorm:"column:notes"`
}

type SupplierQuote struct {
	BaseModel
	VendorID              uint64  `gorm:"column:vendor_id"`
	Description           string  `gorm:"column:description"`
	Status                string  `gorm:"type:enum('pending','approved','rejected','cancelled','ordered','suspended','expired');column:status"`
	RfqNo                 string  `gorm:"column:rfq_no"`
	SupplierQuoteRef      string  `gorm:"column:supplier_quote_ref"`
	RequestForQuotationID *uint64 `gorm:"column:request_for_quotation_id"`
}

type SupplierQuoteItem struct {
	BaseModel
	SupplierQuoteID   uint64  `gorm:"column:supplier_quote_id"`
	ItemID            uint64  `gorm:"column:item_id"`
	QuantityRequested float64 `gorm:"column:quantity_requested"`
	QuantitySupplied  float64 `gorm:"column:quantity_supplied"`
	UnitPrice         float64 `gorm:"column:unit_price"`
	LineTotal         float64 `gorm:"->;column:line_total"` // Read-only generated field
	Notes             string  `gorm:"column:notes"`
}

type Supply struct {
	BaseModel
	VendorID        uint64    `gorm:"column:vendor_id"`
	PaymentTypeID   *uint64   `gorm:"column:payment_type_id"`
	PurchaseOrderID *uint64   `gorm:"column:purchase_order_id"`
	ItemCount       uint64    `gorm:"column:item_count"`
	TotalAmount     float64   `gorm:"column:total_amount"`
	Reference       string    `gorm:"column:reference"`
	Activity        string    `gorm:"column:activity"`
	SuppliedDate    time.Time `gorm:"column:supplied_date"`
	Settled         bool      `gorm:"column:settled;default:0"`
	StoreID         *uint64   `gorm:"column:store_id"`
	PaymentTermID   *uint64   `gorm:"column:payment_term_id"`
}

type SuppliedItem struct {
	BaseModel
	SupplyID   uint64  `gorm:"column:supply_id"`
	ItemID     uint64  `gorm:"column:item_id"`
	Quantity   int     `gorm:"column:quantity;default:0"`
	UnitPrice  float64 `gorm:"column:unit_price"`
	TotalPrice float64 `gorm:"column:total_price"`
}

type SupplyReject struct {
	BaseModel
	ItemID   uint64 `gorm:"column:item_id"`
	SupplyID uint64 `gorm:"column:supply_id"`
	Quantity string `gorm:"column:quantity"` // Varchar in schema
	Reason   string `gorm:"column:reason"`
}

// Order Process

type PurchaseOrder struct {
	BaseModel
	SupplierID      *uint64   `gorm:"column:supplier_id"`
	SupplierQuoteID *uint64   `gorm:"column:supplier_quote_id"`
	PoNumber        string    `gorm:"uniqueIndex;column:po_number"`
	PoDate          time.Time `gorm:"column:po_date"`
	Status          string    `gorm:"column:status;default:'draft'"`
	TotalAmount     float64   `gorm:"column:total_amount;default:0.00"`
}

type PurchaseOrderItem struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	PurchaseOrderID uint64    `gorm:"column:purchase_order_id"`
	ItemID          uint64    `gorm:"column:item_id"`
	Description     string    `gorm:"column:description"`
	Quantity        float64   `gorm:"column:quantity"`
	UnitPrice       float64   `gorm:"column:unit_price"`
	TotalPrice      float64   `gorm:"column:total_price"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

type PurchaseRequisition struct {
	BaseModel
	RequisitionNo   string    `gorm:"uniqueIndex;column:requisition_no"`
	RequisitionDate time.Time `gorm:"column:requisition_date"`
	Description     string    `gorm:"column:description"`
	Status          string    `gorm:"column:status;default:'draft'"`
}

type PurchaseRequisitionItem struct {
	BaseModel
	PurchaseRequisitionID *uint64 `gorm:"column:purchase_requisition_id"`
	ItemID                *uint64 `gorm:"column:item_id"`
	Quantity              float64 `gorm:"column:quantity"`
	Status                string  `gorm:"column:status;default:'pending'"`
}

func (PurchaseRequisitionItem) TableName() string {
	return "purchase_requisition_items"
}

type ProductStockMovement struct {
	BaseModel
	ProductID    uint64    `gorm:"index;column:product_id"`
	MovementType string    `gorm:"column:movement_type"`
	Quantity     float64   `gorm:"column:quantity"`
	MovementDate time.Time `gorm:"index;column:movement_date"`
}

type Procurement struct {
	BaseModel
	ProductID uint64    `gorm:"index;column:product_id"`
	Quantity  float64   `gorm:"column:quantity"`
	Cost      float64   `gorm:"column:cost"`
	Date      time.Time `gorm:"index;column:date"`
	Supplier  string    `gorm:"column:supplier"`
}

type ProcurementItem struct {
	BaseModel
	ProcurementID uint64  `gorm:"index;column:procurement_id"`
	ItemID        uint64  `gorm:"index;column:item_id"`
	Quantity      float64 `gorm:"column:quantity"`
	UnitPrice     float64 `gorm:"column:unit_price"`
}
