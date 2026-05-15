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
	ID           uint64    `json:"ID"`
	CategoryCode string    `json:"CategoryCode"`
	CategoryName string    `json:"CategoryName"`
	Description  string    `json:"Description"`
	Status       string    `json:"Status"`
	CreatedAt    time.Time `json:"CreatedAt"`
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
	ID             uint64    `json:"ID"`
	SupplierCode   string    `json:"SupplierCode"`
	SupplierType   string    `json:"SupplierType"`
	CompanyName    string    `json:"CompanyName"`
	FullName       string    `json:"FullName"`
	CategoryName   string    `json:"CategoryName"`
	EmailAddress   string    `json:"EmailAddress"`
	PhoneNo        string    `json:"PhoneNo"`
	CurrentBalance float64   `json:"CurrentBalance"`
	Status         string    `json:"Status"`
	CreatedAt      time.Time `json:"CreatedAt"`
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
	ID                 uint64    `json:"ID"`
	SupplierID         uint64    `json:"SupplierID"`
	SupplierName       string    `json:"SupplierName"`
	ContactType        string    `json:"ContactType"`
	FullName           string    `json:"FullName"`
	Designation        string    `json:"Designation"`
	PhoneNo            string    `json:"PhoneNo"`
	AlternativePhoneNo string    `json:"AlternativePhoneNo"`
	EmailAddress       string    `json:"EmailAddress"`
	IsDefault          string    `json:"IsDefault"`
	Notes              string    `json:"Notes"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
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
	ID             uint64    `json:"ID"`
	SupplierID     uint64    `json:"SupplierID"`
	BankName       string    `json:"BankName"`
	BankBranchName string    `json:"BankBranchName"`
	AccountName    string    `json:"AccountName"`
	AccountNumber  string    `json:"AccountNumber"`
	AccountType    string    `json:"AccountType"`
	CurrencyCode   string    `json:"CurrencyCode"`
	SwiftCode      string    `json:"SwiftCode"`
	MobileMoneyNo  string    `json:"MobileMoneyNo"`
	IsDefault      string    `json:"IsDefault"`
	Status         string    `json:"Status"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
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
	ID             uint64     `json:"ID"`
	SupplierID     uint64     `json:"SupplierID"`
	SupplierName   string     `json:"SupplierName"`
	DocumentType   string     `json:"DocumentType"`
	DocumentNumber string     `json:"DocumentNumber"`
	FileName       string     `json:"FileName"`
	FilePath       string     `json:"FilePath"`
	IssueDate      *time.Time `json:"IssueDate"`
	ExpiryDate     *time.Time `json:"ExpiryDate"`
	Verified       string     `json:"Verified"`
	VerifiedBy     *uint64    `json:"VerifiedBy"`
	VerifiedAt     *time.Time `json:"VerifiedAt"`
	Notes          string     `json:"Notes"`
	CreatedAt      time.Time  `json:"CreatedAt"`
	UpdatedAt      time.Time  `json:"UpdatedAt"`
}

// Supplier Quote
type CreateSupplierQuoteRequest struct {
	VendorID         uint64 `json:"vendor_id" validate:"required"`
	Description      string `json:"description" validate:"required"`
	RfqNo            string `json:"rfq_no"`
	SupplierQuoteRef string `json:"supplier_quote_ref"`
}

type SupplierQuoteResponse struct {
	ID               uint64    `json:"ID"`
	VendorID         uint64    `json:"VendorID"`
	VendorName       string    `json:"VendorName"`
	Description      string    `json:"Description"`
	Status           string    `json:"Status"`
	RfqNo            string    `json:"RfqNo"`
	SupplierQuoteRef string    `json:"SupplierQuoteRef"`
	CreatedAt        time.Time `json:"CreatedAt"`
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
	ID                uint64  `json:"ID"`
	SupplierQuoteID   uint64  `json:"SupplierQuoteID"`
	ItemName          string  `json:"ItemName"`
	QuantityRequested float64 `json:"QuantityRequested"`
	QuantitySupplied  float64 `json:"QuantitySupplied"`
	UnitPrice         float64 `json:"UnitPrice"`
	LineTotal         float64 `json:"LineTotal"`
}
