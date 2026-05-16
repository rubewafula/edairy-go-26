package dtos

import "time"

// Supplier Category
type CreateSupplierCategoryRequest struct {
	CategoryCode string  `json:"category_code" validate:"required"`
	CategoryName string  `json:"category_name" validate:"required"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	SiteID       *uint64 `json:"site_id"`
}

type SupplierCategoryResponse struct {
	ID           uint64    `json:"id"`
	CategoryCode string    `json:"category_code"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// Supplier
type CreateSupplierRequest struct {
	SupplierCategoryID uint64  `json:"supplier_category_id" validate:"required"`
	SupplierCode       string  `json:"supplier_code" validate:"required"`
	SupplierType       string  `json:"supplier_type" validate:"required,oneof=individual company"`
	CompanyName        string  `json:"company_name"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	Gender             string  `json:"gender"`
	PhoneNo            string  `json:"phone_no" validate:"required"`
	EmailAddress       string  `json:"email_address" validate:"required,email"`
	KraPin             string  `json:"kra_pin"`
	OpeningBalance     float64 `json:"opening_balance"`
	CreditLimit        float64 `json:"credit_limit"`
	PaymentTermsDays   int     `json:"payment_terms_days"`
	Status             string  `json:"status"`
	Notes              string  `json:"notes"`
}

type SupplierResponse struct {
	ID             uint64    `json:"id"`
	SupplierCode   string    `json:"supplier_code"`
	SupplierType   string    `json:"supplier_type"`
	CompanyName    string    `json:"company_name"`
	FullName       string    `json:"full_name"`
	CategoryName   string    `json:"category_name"`
	EmailAddress   string    `json:"email_address"`
	PhoneNo        string    `json:"phone_no"`
	CurrentBalance float64   `json:"current_balance"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// Supplier Contact
type CreateSupplierContactRequest struct {
	SupplierID         uint64 `json:"supplier_id" validate:"required"`
	ContactType        string `json:"contact_type" validate:"required,oneof=primary finance procurement technical other"`
	FullName           string `json:"full_name" validate:"required"`
	Designation        string `json:"designation"`
	PhoneNo            string `json:"phone_no" validate:"required"`
	AlternativePhoneNo string `json:"alternative_phone_no"`
	EmailAddress       string `json:"email_address" validate:"required,email"`
	IsDefault          string `json:"is_default" validate:"required,oneof=yes no"`
	Notes              string `json:"notes"`
}

type UpdateSupplierContactRequest struct {
	ContactType        string `json:"contact_type" validate:"required,oneof=primary finance procurement technical other"`
	FullName           string `json:"full_name" validate:"required"`
	Designation        string `json:"designation"`
	PhoneNo            string `json:"phone_no" validate:"required"`
	AlternativePhoneNo string `json:"alternative_phone_no"`
	EmailAddress       string `json:"email_address" validate:"required,email"`
	IsDefault          string `json:"is_default" validate:"required,oneof=yes no"`
	Notes              string `json:"notes"`
}

type SupplierContactResponse struct {
	ID                 uint64    `json:"id"`
	SupplierID         uint64    `json:"supplier_id"`
	SupplierName       string    `json:"supplier_name"`
	ContactType        string    `json:"contact_type"`
	FullName           string    `json:"full_name"`
	Designation        string    `json:"designation"`
	PhoneNo            string    `json:"phone_no"`
	AlternativePhoneNo string    `json:"alternative_phone_no"`
	EmailAddress       string    `json:"email_address"`
	IsDefault          string    `json:"is_default"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Supplier Bank Account
type CreateSupplierBankAccountRequest struct {
	SupplierID    uint64 `json:"supplier_id" validate:"required"`
	BankID        uint64 `json:"bank_id"`
	BankBranchID  uint64 `json:"bank_branch_id"`
	AccountName   string `json:"account_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountType   string `json:"account_type" validate:"required"`
	CurrencyCode  string `json:"currency_code"`
	IsDefault     string `json:"is_default"`
}

type UpdateSupplierBankAccountRequest struct {
	BankID        *uint64 `json:"bank_id"`
	BankBranchID  *uint64 `json:"bank_branch_id"`
	AccountName   string  `json:"account_name" validate:"required"`
	AccountNumber string  `json:"account_number" validate:"required"`
	AccountType   string  `json:"account_type" validate:"required,oneof=bank mobile_money"`
	CurrencyCode  string  `json:"currency_code"`
	SwiftCode     string  `json:"swift_code"`
	MobileMoneyNo string  `json:"mobile_money_no"`
	IsDefault     string  `json:"is_default" validate:"oneof=yes no"`
	Status        string  `json:"status" validate:"oneof=active inactive"`
}

type SupplierBankAccountResponse struct {
	ID             uint64    `json:"id"`
	SupplierID     uint64    `json:"supplier_id"`
	BankName       string    `json:"bank_name"`
	BankBranchName string    `json:"bank_branch_name"`
	AccountName    string    `json:"account_name"`
	AccountNumber  string    `json:"account_number"`
	AccountType    string    `json:"account_type"`
	CurrencyCode   string    `json:"currency_code"`
	SwiftCode      string    `json:"swift_code"`
	MobileMoneyNo  string    `json:"mobile_money_no"`
	IsDefault      string    `json:"is_default"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Supplier Document
type CreateSupplierDocumentRequest struct {
	SupplierID     uint64 `json:"supplier_id" validate:"required"`
	DocumentType   string `json:"document_type" validate:"required,oneof=national_id passport kra_pin business_registration license contract bank_letter certificate other"`
	DocumentNumber string `json:"document_number"`
	IssueDate      string `json:"issue_date"`
	ExpiryDate     string `json:"expiry_date"`
	Notes          string `json:"notes"`
}

type UpdateSupplierDocumentRequest struct {
	DocumentType   string `json:"document_type" validate:"required,oneof=national_id passport kra_pin business_registration license contract bank_letter certificate other"`
	DocumentNumber string `json:"document_number"`
	IssueDate      string `json:"issue_date"`
	ExpiryDate     string `json:"expiry_date"`
	Notes          string `json:"notes"`
}

type VerifySupplierDocumentRequest struct {
	Verified string `json:"verified" validate:"required,oneof=yes no"`
}

type SupplierDocumentResponse struct {
	ID             uint64     `json:"id"`
	SupplierID     uint64     `json:"supplier_id"`
	SupplierName   string     `json:"supplier_name"`
	DocumentType   string     `json:"document_type"`
	DocumentNumber string     `json:"document_number"`
	FileName       string     `json:"file_name"`
	FilePath       string     `json:"file_path"`
	IssueDate      *time.Time `json:"issue_date"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	Verified       string     `json:"verified"`
	VerifiedBy     *uint64    `json:"verified_by"`
	VerifiedAt     *time.Time `json:"verified_at"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Supplier Quote
type CreateSupplierQuoteRequest struct {
	VendorID         uint64 `json:"vendor_id" validate:"required"`
	Description      string `json:"description" validate:"required"`
	RfqNo            string `json:"rfq_no"`
	SupplierQuoteRef string `json:"supplier_quote_ref"`
}

type SupplierQuoteResponse struct {
	ID               uint64    `json:"id"`
	VendorID         uint64    `json:"vendor_id"`
	VendorName       string    `json:"vendor_name"`
	Description      string    `json:"description"`
	Status           string    `json:"status"`
	RfqNo            string    `json:"rfq_no"`
	SupplierQuoteRef string    `json:"supplier_quote_ref"`
	CreatedAt        time.Time `json:"created_at"`
}

// Supplier Quote Item
type CreateSupplierQuoteItemRequest struct {
	SupplierQuoteID   uint64  `json:"supplier_quote_id" validate:"required"`
	ItemID            uint64  `json:"item_id" validate:"required"`
	QuantityRequested float64 `json:"quantity_requested" validate:"required"`
	UnitPrice         float64 `json:"unit_price" validate:"required"`
	Notes             string  `json:"notes"`
}

type UpdateSupplierQuoteItemRequest struct {
	QuantitySupplied float64 `json:"quantity_supplied" validate:"required,min=0"`
	UnitPrice        float64 `json:"unit_price" validate:"required,min=0"`
	Notes            string  `json:"notes" validate:"max=255"`
}

type SupplierQuoteItemResponse struct {
	ID                uint64  `json:"id"`
	SupplierQuoteID   uint64  `json:"supplier_quote_id"`
	ItemName          string  `json:"item_name"`
	QuantityRequested float64 `json:"quantity_requested"`
	QuantitySupplied  float64 `json:"quantity_supplied"`
	UnitPrice         float64 `json:"unit_price"`
	LineTotal         float64 `json:"line_total"`
}
