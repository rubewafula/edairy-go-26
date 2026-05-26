package models

import (
	"time"
)

type Customer struct {
	BaseModel
	CustomerTypeID uint64  `gorm:"column:customer_type_id"`
	FullNames      string  `gorm:"column:full_names"`
	Phone          string  `gorm:"column:phone"`
	EmailAddress   string  `gorm:"column:email_address"`
	CustomerNo     string  `gorm:"uniqueIndex;column:customer_no"`
	KraPin         string  `gorm:"column:kra_pin"`
	Status         string  `gorm:"column:status"`
	CreditLimit    float64 `gorm:"column:credit_limit"`
	PostalAddress  string  `gorm:"column:postal_address"`
	PostalCode     string  `gorm:"column:postal_code"`
	PostalTown     string  `gorm:"column:postal_town"`
	Terms          uint64  `gorm:"column:payment_terms_days"`
	Rate           float64 `gorm:"column:rate"`
}

type CustomerDocument struct {
	BaseModel
	CustomerID     uint64    `gorm:"index;column:customer_id"`
	DocDescription string    `gorm:"column:doc_description"`
	DocBalance     float64   `gorm:"column:doc_balance"`
	DueDate        time.Time `gorm:"column:due_date"`
}

func (CustomerDocument) TableName() string {
	return "customer_documents"
}

type CustomerType struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (CustomerType) TableName() string {
	return "customer_types"
}

type CustomerOpeningBalance struct {
	BaseModel
	CustomerID uint64  `gorm:"index;column:customer_id"`
	Balance    float64 `gorm:"column:balance"`
	Status     string  `gorm:"column:status"`
}

type CustomerClass struct {
	BaseModel
	ClassCode   string `gorm:"uniqueIndex;column:class_code"`
	Description string `gorm:"column:description"`
}

type CustomerMilkRate struct {
	BaseModel
	CustomerID   uint64  `gorm:"index;column:customer_id"`
	Rate         float64 `gorm:"column:rate"`
	GradeID      uint64  `gorm:"column:grade_id"`
	PayDateRange uint64  `gorm:"column:customer_pay_date_range_id"`
}

type CustomerPayDateRange struct {
	BaseModel
	Name      string    `gorm:"column:name"`
	StartDate time.Time `gorm:"column:start_date"`
	EndDate   time.Time `gorm:"column:end_date"`
}

type CustomerCollection struct {
	BaseModel
	PayDateRangeID  uint64  `gorm:"index;column:pay_date_range_id"`
	PayrollMonth    int     `gorm:"column:payroll_month"`
	PayrollYear     int     `gorm:"column:payroll_year"`
	TotalDeliveries float64 `gorm:"column:total_deliveries"`
	TotalAmount     float64 `gorm:"column:total_amount"`
}

type CustomerBilling struct {
	BaseModel
	PayDateRangeID  uint64  `gorm:"column:pay_date_range_id"`
	TotalDeliveries float64 `gorm:"column:total_deliveries"`
	TotalAmount     float64 `gorm:"column:total_amount"`
	Status          string  `gorm:"type:enum('pending','invoiced');default:'pending';column:status"`
	InvoiceID       *uint64 `gorm:"column:invoice_id"`
}

type CustomerBillingItem struct {
	BaseModel
	CustomerBillingID uint64  `gorm:"column:customer_billing_id"`
	ProductGradeID    uint64  `gorm:"column:product_grade_id"`
	TotalQuantity     float64 `gorm:"column:total_quantity"`
	TotalAmount       float64 `gorm:"column:total_amount"`
	UnitPrice         float64 `gorm:"column:unit_price"`
}

type CustomerInvoice struct {
	BaseModel
	TotalAmount float64   `gorm:"column:total_amount"`
	Status      string    `gorm:"column:status"`
	CustomerID  uint64    `gorm:"column:customer_id"`
	InvoiceNo   string    `gorm:"column:invoice_no"`
	BillingID   uint64    `gorm:"column:billing_id"`
	InvoiceDate time.Time `gorm:"column:invoice_date"`
	DueDate     time.Time `gorm:"column:due_date"`
	GrossAmount float64   `gorm:"column:gross_amount"`
	TaxAmount   float64   `gorm:"column:tax_amount"`
	Balance     float64   `gorm:"column:balance"`
}

type CustomerPayment struct {
	BaseModel
	InvoiceID     *uint64   `gorm:"column:invoice_id"`
	ReceiptNumber string    `gorm:"column:receipt_number"`
	PaymentDate   time.Time `gorm:"column:payment_date"`
	Amount        float64   `gorm:"column:amount"`
	PaymentMethod string    `gorm:"column:payment_method"`
	ReferenceNo   string    `gorm:"column:reference_no"`
	Notes         string    `gorm:"column:notes"`
	CustomerID    uint64    `gorm:"column:customer_id"`
}

type CustomerPaymentAllocation struct {
	BaseModel
	InvoiceID         uint64  `gorm:"column:invoice_id"`
	CustomerPaymentID uint64  `gorm:"column:customer_payment_id"`
	AllocatedAmount   float64 `gorm:"column:allocated_amount"`
}
