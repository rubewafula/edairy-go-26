package models

import (
	"time"

	"gorm.io/gorm"
)

/*
|--------------------------------------------------------------------------
| Base Model (shared fields)
|--------------------------------------------------------------------------
*/

type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
	CreatedBy uint64         `gorm:"column:created_by"`
	UpdatedBy uint64         `gorm:"column:updated_by"`
}

type Permission struct {
	BaseModel
	Name      string `gorm:"column:name"`
	GuardName string `gorm:"column:guard_name"`
}

type Role struct {
	BaseModel
	Name        string       `gorm:"column:name"`
	GuardName   string       `gorm:"column:guard_name"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;"`
}

type UserRole struct {
	BaseModel
	UserID uint64 `gorm:"primaryKey"`
	RoleID uint64 `gorm:"primaryKey"`
}

type UserPermission struct {
	BaseModel
	UserID       uint64 `gorm:"primaryKey"`
	PermissionID uint64 `gorm:"primaryKey"`
}

type RolePermission struct {
	BaseModel
	RoleID       uint64 `gorm:"primaryKey"`
	PermissionID uint64 `gorm:"primaryKey"`
}

type User struct {
	BaseModel
	Name              string     `gorm:"column:name"`
	Email             string     `gorm:"column:email;uniqueIndex"`
	EmailVerifiedAt   *string    `gorm:"column:email_verified_at"`
	Password          string     `gorm:"column:password"`
	RememberToken     string     `gorm:"column:remember_token"`
	IsVerified        bool       `gorm:"column:is_verified;default:0"`
	VerificationToken string     `gorm:"column:verification_token"`
	ResetToken        string     `gorm:"column:reset_token"`
	ResetTokenExpiry  *time.Time `gorm:"column:reset_token_expiry"`

	Roles       []Role       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Permissions []Permission `gorm:"many2many:user_permissions;constraint:OnDelete:CASCADE;"`
}

// System Models
type Installation struct {
	BaseModel
	InstallationDate time.Time `gorm:"column:installation_date"`
	ExpiryDate       time.Time `gorm:"column:expiry_date"`
}

type ActivityLog struct {
	BaseModel
	LogName     string                 `gorm:"column:log_name"`
	Description string                 `gorm:"column:description"`
	SubjectType string                 `gorm:"index;column:subject_type"`
	BatchUUID   string                 `gorm:"column:batch_uuid"`
	SubjectID   uint64                 `gorm:"column:subject_id"`
	CauserType  string                 `gorm:"column:causer_type"`
	CauserID    uint64                 `gorm:"column:causer_id"`
	Properties  map[string]interface{} `gorm:"column:properties;serializer:json"`
	Event       string                 `gorm:"column:event"`
}

type License struct {
	BaseModel
	Key           string    `gorm:"column:key"`
	Status        string    `gorm:"column:status"`
	LastCheckedAt time.Time `gorm:"column:last_checked_at"`
}

type Notification struct {
	BaseModel
	Type           string    `gorm:"column:type"`
	NotifiableType string    `gorm:"index;column:notifiable_type"`
	NotifiableID   uint64    `gorm:"index;column:notifiable_id"`
	Data           string    `gorm:"type:text;column:data"`
	ReadAt         time.Time `gorm:"column:read_at"`
}

// Accounting
type AccountCategory struct {
	BaseModel
	Name          string `gorm:"column:name"`
	Description   string `gorm:"column:description"`
	AccountTypeID uint64 `gorm:"column:account_type_id"`
}

type AccountType struct {
	BaseModel
	Name string `gorm:"column:name"`
}

type Account struct {
	BaseModel
	AccountCode       string  `gorm:"uniqueIndex;column:account_code"`
	Name              string  `gorm:"column:name"`
	Description       string  `gorm:"column:description"`
	AccountCategoryID uint64  `gorm:"column:account_category_id"`
	ParentAccountID   uint64  `gorm:"column:parent_account_id"`
	IsPostable        bool    `gorm:"column:is_postable"`
	IsActive          bool    `gorm:"default:true;column:is_active"`
	Balance           float64 `gorm:"column:balance"`
}

type AccountSubAccount struct {
	BaseModel
	SubAccountCode string `gorm:"uniqueIndex;column:sub_account_code"`
	Name           string `gorm:"column:name"`
	Description    string `gorm:"column:description"`
	AccountID      uint64 `gorm:"index;column:account_id"`
}

type Transaction struct {
	BaseModel
	AccountID   uint64    `gorm:"index;column:account_id"`
	Amount      float64   `gorm:"column:amount"`
	Type        string    `gorm:"column:type"` // debit/credit
	Reference   string    `gorm:"uniqueIndex;column:reference"`
	Description string    `gorm:"column:description"`
	Date        time.Time `gorm:"index;column:date"`
}

type LedgerEntry struct {
	BaseModel
	TransactionID uint64  `gorm:"index;column:transaction_id"`
	AccountID     uint64  `gorm:"index;column:account_id"`
	SubAccountID  uint64  `gorm:"index;column:sub_account_id"`
	Debit         float64 `gorm:"column:debit"`
	Credit        float64 `gorm:"column:credit"`
}

type CashTransaction struct {
	BaseModel
	ReferenceNumber        string    `gorm:"uniqueIndex;column:reference_number"`
	TransactionDescription string    `gorm:"column:transaction_description"`
	TransactionType        string    `gorm:"column:transaction_type"`
	TransactionDate        time.Time `gorm:"index;column:transaction_date"`
	PaidBy                 string    `gorm:"column:paid_by"`
	TransactionAmount      float64   `gorm:"column:transaction_amount"`
	CustomerType           string    `gorm:"column:customer_type"`
	CustomerID             uint64    `gorm:"index;column:customer_id"`
	PaymentModeID          uint64    `gorm:"index;column:payment_mode_id"`
	PaymentType            string    `gorm:"column:payment_type"`
}

type PaymentType struct {
	BaseModel
	Code        string `gorm:"uniqueIndex;column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Status      string `gorm:"column:status"`
	IsDefault   bool   `gorm:"column:is_default"`
	IsActive    bool   `gorm:"default:true;column:is_active"`
	IsSystem    bool   `gorm:"column:is_system"`
}

// Dairy Module
type SubRoute struct {
	BaseModel
	RouteID     uint64 `gorm:"index;column:route_id"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

type MilkJournal struct {
	BaseModel
	Journal             string    `gorm:"column:journal"`
	JournalDate         time.Time `gorm:"index;column:journal_date"`
	MilkDeliveryShiftID uint64    `gorm:"index;column:milk_delivery_shift_id"`
	RouteID             uint64    `gorm:"index;column:route_id"`
	UserID              uint64    `gorm:"column:user_id"`
	TransporterID       uint64    `gorm:"column:transporter_id"`
	Confirmed           bool      `gorm:"column:confirmed"`
}

type MilkJournalEntry struct {
	BaseModel
	MemberID            uint64    `gorm:"index;column:member_id"`
	MilkJournalID       uint64    `gorm:"index;column:milk_journal_id"`
	MilkJournalBatchID  uint64    `gorm:"index;column:milk_journal_batch_id"`
	RouteID             uint64    `gorm:"index;column:route_id"`
	MilkDeliveryShiftID uint64    `gorm:"index;column:milk_delivery_shift_id"`
	Status              string    `gorm:"column:status"`
	JournalDate         time.Time `gorm:"index;column:journal_date"`
	Quantity            float64   `gorm:"type:decimal(18,2);column:quantity"`
	TransporterID       uint64    `gorm:"index;column:transporter_id"`
	RouteCenterID       uint64    `gorm:"index;column:route_center_id"`
	CanID               uint64    `gorm:"index;column:can_id"`
}

func (MilkJournalEntry) TableName() string {
	return "milk_journal_entries"
}

type Route struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Code        string `gorm:"column:code"`
	LocationID  uint64 `gorm:"column:location_id"`
}

type Location struct {
	BaseModel
	Name string `gorm:"column:name"`
}

type MilkDeliveryShift struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (MilkDeliveryShift) TableName() string {
	return "milk_delivery_shifts"
}

type MilkCollection struct {
	BaseModel
	MemberID           uint64    `gorm:"index;column:member_id"`
	MilkJournalTableID uint64    `gorm:"index;column:milk_journal_table_id"`
	RouteID            uint64    `gorm:"index;column:route_id"`
	MilkJournalBatchID uint64    `gorm:"index;column:milk_journal_batch_id"`
	ShiftID            uint64    `gorm:"index;column:milk_delivery_shift_id"`
	Status             string    `gorm:"column:status"`
	JournalDate        time.Time `gorm:"index;column:journal_date"`
	Quantity           float64   `gorm:"column:quantity"`
	TransporterID      uint64    `gorm:"index;column:transporter_id"`
	RouteCenterID      uint64    `gorm:"index;column:route_center_id"`
	CanID              uint64    `gorm:"index;column:can_id"`
}

type MilkCan struct {
	BaseModel
	CanID      string  `gorm:"uniqueIndex;column:can_id"`
	CanType    string  `gorm:"column:can_type"`
	CanSize    float64 `gorm:"column:can_size"`
	Units      string  `gorm:"column:units"`
	TareWeight float64 `gorm:"column:tare_weight"`
	RouteID    uint64  `gorm:"index;column:route_id"`
}

func (MilkCan) TableName() string {
	return "milk_cans"
}

type CanMovement struct {
	BaseModel
	CanID             uint64    `gorm:"column:can_id"`
	MovementType      string    `gorm:"column:movement_type"`
	Quantity          float64   `gorm:"column:quantity"`
	Remarks           string    `gorm:"column:remarks"`
	ShiftID           uint64    `gorm:"column:shift_id"`
	TransporterID     uint64    `gorm:"column:transporter_id"`
	RouteID           uint64    `gorm:"column:route_id"`
	MovementDate      time.Time `gorm:"column:movement_date"`
	ConditionOnReturn string    `gorm:"column:condition_on_return"`
}

func (CanMovement) TableName() string {
	return "can_movements"
}

type CoolerMilkCollection struct {
	BaseModel
	TransactionDate     time.Time `gorm:"index;column:transaction_date"`
	Quantity            float64   `gorm:"column:quantity"`
	TransportVehicleID  uint64    `gorm:"index;column:transport_vehicle_id"`
	MilkDeliveryShiftID uint64    `gorm:"index;column:milk_delivery_shift_id"`
	Confirmed           int       `gorm:"column:confirmed"`
	SiteID              uint64    `gorm:"index;column:site_id"`
	TransporterID       uint64    `gorm:"index;column:transporter_id"`
	RouteID             uint64    `gorm:"index;column:route_id"`
}

type Store struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

type MilkReject struct {
	BaseModel
	RouteID             uint64    `gorm:"column:route_id"`
	Quantity            float64   `gorm:"column:quantity"`
	TransactionDate     time.Time `gorm:"index;column:transaction_date"`
	Reason              string    `gorm:"column:reason"`
	Description         string    `gorm:"column:description"`
	Confirmed           int       `gorm:"column:confirmed"`
	MemberID            uint64    `gorm:"index;column:member_id"`
	TransporterID       uint64    `gorm:"index;column:transporter_id"`
	CanID               uint64    `gorm:"index;column:can_id"`
	MilkDeliveryShiftID uint64    `gorm:"index;column:milk_delivery_shift_id"`
}

type MilkSpecialRate struct {
	BaseModel
	MemberID              uint64  `gorm:"index;column:member_id"`
	Rate                  float64 `gorm:"column:rate"`
	MonthlyPayDateRangeID uint64  `gorm:"index;column:monthly_pay_date_range_id"`
	Confirmed             bool    `gorm:"column:confirmed"`
}

type MilkCooler struct {
	BaseModel
	JournalID          string    `gorm:"column:journal_id"`
	TransactionDate    time.Time `gorm:"index;column:transaction_date"`
	Quantity           float64   `gorm:"column:quantity"`
	RegistrationNumber string    `gorm:"column:registration_number"`
	ScaleNumber        string    `gorm:"column:scale_number"`
	Shift              string    `gorm:"column:shift"`
	Confirmed          int       `gorm:"column:confirmed"`
	UserID             uint64    `gorm:"index;column:user_id"`
	MilkBar            float64   `gorm:"column:milk_bar"`
	SiteID             uint64    `gorm:"index;column:site_id"`
}

type MilkDelivery struct {
	BaseModel
	DeliveryNoteNumber string    `gorm:"index;column:delivery_note_number"`
	CustomerID         uint64    `gorm:"index;column:customer_id"`
	QuantityAccepted   float64   `gorm:"column:quantity_accepted"`
	Cooler             string    `gorm:"column:cooler"`
	Invoiced           int       `gorm:"column:invoiced"`
	TransactionDate    time.Time `gorm:"index;column:transaction_date"`
	Amount             float64   `gorm:"column:amount"`
	AmountPaid         float64   `gorm:"column:amount_paid"`
	RouteID            uint64    `gorm:"index;column:route_id"`
	Confirmed          int       `gorm:"column:confirmed"`
	Processed          string    `gorm:"column:processed"`
	TransporterID      uint64    `gorm:"index;column:transporter_id"`
}

func (MilkDelivery) TableName() string {
	return "milk_deliveries"
}

type MilkDeliveryAcceptance struct {
	BaseModel
	DeliveryNoteNumber string    `gorm:"uniqueIndex;column:delivery_note_number"`
	CustomerID         uint64    `gorm:"index;column:customer_id"`
	QuantityAccepted   float64   `gorm:"column:quantity_accepted"`
	Cooler             string    `gorm:"column:cooler"`
	Invoiced           int       `gorm:"column:invoiced"`
	TransactionDate    time.Time `gorm:"index;column:transaction_date"`
	Amount             float64   `gorm:"column:amount"`
	AmountPaid         float64   `gorm:"column:amount_paid"`
	RouteID            uint64    `gorm:"index;column:route_id"`
	Confirmed          int       `gorm:"column:confirmed"`
	Processed          string    `gorm:"column:processed"`
	TransporterID      uint64    `gorm:"index;column:transporter_id"`
}

type MilkDeliveryItem struct {
	BaseModel
	DeliveryID uint64    `gorm:"index;column:delivery_id"`
	Quantity   float64   `gorm:"column:quantity"`
	Rate       float64   `gorm:"column:rate"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	GradeID    uint64    `gorm:"index;column:grade_id"`
}

type MilkLocalSale struct {
	BaseModel
	Quantity        float64   `gorm:"column:quantity"`
	Rate            float64   `gorm:"column:rate"`
	GradeID         uint64    `gorm:"index;column:grade_id"`
	RefNumber       string    `gorm:"uniqueIndex;column:ref_number"`
	TransactionDate time.Time `gorm:"index;column:transaction_date"`
	TransporterID   uint64    `gorm:"index;column:transporter_id"`
	Amount          float64   `gorm:"column:amount"`
}

func (MilkLocalSale) TableName() string {
	return "milk_local_sales"
}

type MilkSale struct {
	BaseModel
	Quantity float64   `gorm:"column:quantity"`
	Price    float64   `gorm:"column:price"`
	Amount   float64   `gorm:"column:amount"`
	Buyer    string    `gorm:"column:buyer"`
	Date     time.Time `gorm:"index;column:transaction_date"`
}

type MilkTransporterCost struct {
	BaseModel
	MemberID           uint64  `gorm:"index;column:member_id"`
	TransporterID      uint64  `gorm:"index;column:transporter_id"`
	MilkJournalBatchID uint64  `gorm:"index;column:milk_journal_batch_id"`
	PayrollMonth       string  `gorm:"column:payroll_month"`
	PayrollYear        string  `gorm:"column:payroll_year"`
	PayDateRangeID     uint64  `gorm:"index;column:pay_date_range_id"`
	PayrollID          uint64  `gorm:"index;column:payroll_id"`
	Quantity           float64 `gorm:"column:quantity"`
	Rejects            float64 `gorm:"column:rejects"`
}

type DailyMilkVariance struct {
	ID               uint64    `gorm:"column:id"`
	Transporter      string    `gorm:"column:transporter"`
	Day              time.Time `gorm:"column:day"`
	Month            string    `gorm:"column:month"`
	FieldCollections float64   `gorm:"column:field_collections"`
	MCC              float64   `gorm:"column:mcc"`
	CashSales        float64   `gorm:"column:cash_sales"`
	CreditSales      float64   `gorm:"column:credit_sales"`
	Rejects          float64   `gorm:"column:rejects"`
	Balance          float64   `gorm:"column:balance"`
}

func (DailyMilkVariance) TableName() string {
	return "daily_milk_variance"
}

// Member+ Lending
type Member struct {
	BaseModel
	MemberNo          string    `gorm:"uniqueIndex;column:member_no"`
	FirstName         string    `gorm:"column:first_name"`
	LastName          string    `gorm:"column:last_name"`
	OtherNames        string    `gorm:"column:other_names"`
	IDNumber          string    `gorm:"uniqueIndex;column:id_no"`
	Gender            string    `gorm:"column:gender"`
	DateOfBirth       string    `gorm:"column:dob"`
	PrimaryPhone      string    `gorm:"column:primary_phone"`
	SecondaryPhone    string    `gorm:"column:secondary_phone"`
	Email             string    `gorm:"column:email"`
	TaxNumber         string    `gorm:"column:tax_number"`
	MaritalStatus     string    `gorm:"column:marital_status"`
	Status            string    `gorm:"column:status"`
	RouteID           uint64    `gorm:"column:route_id"`
	MemberTypeID      uint64    `gorm:"column:member_type_id"`
	NumberOfCows      int       `gorm:"column:number_of_cows"`
	NextOfKinFullName string    `gorm:"column:next_of_kin_full_name"`
	NextOfKinPhone    string    `gorm:"column:next_of_kin_phone"`
	DateRegistered    time.Time `gorm:"column:date_registered"`
	PassportPhoto     string    `gorm:"column:passport_photo"`
	IdFrontPhoto      string    `gorm:"column:id_front_photo"`
	IdBackPhoto       string    `gorm:"column:id_back_photo"`
	UpdatedBy         string    `gorm:"column:updated_by"`
	Downloaded        string    `gorm:"column:downloaded"`
	BirthCity         string    `gorm:"column:birth_city"`
	IdDateOfIssue     time.Time `gorm:"column:id_date_of_issue"`
	Title             string    `gorm:"column:title"`
	CashoutEnrolled   bool      `gorm:"column:cashout_enrolled"`
}

func (Member) TableName() string {
	return "member_registrations"
}

type MemberProductRejectHistory struct {
	BaseModel
	MemberID uint64  `gorm:"index;column:member_id"`
	Period   int     `gorm:"column:period"`
	Year     int     `gorm:"column:year"`
	Quantity float64 `gorm:"column:quantity"`
	Route    string  `gorm:"column:route"`
}

type Customer struct {
	BaseModel
	ClassID       uint64  `gorm:"column:class_id"`
	FullNames     string  `gorm:"column:full_names"`
	Phone         string  `gorm:"column:phone"`
	EmailAddress  string  `gorm:"column:email_address"`
	CustomerNo    string  `gorm:"uniqueIndex;column:customer_no"`
	KraPin        string  `gorm:"column:kra_pin"`
	Status        string  `gorm:"column:status"`
	CreditLimit   float64 `gorm:"column:credit_limit"`
	PostalAddress string  `gorm:"column:postal_address"`
	PostalCode    string  `gorm:"column:postal_code"`
	PostalTown    string  `gorm:"column:postal_town"`
	SiteID        uint64  `gorm:"column:site_id"`
	Terms         string  `gorm:"column:terms"`
	Rate          float64 `gorm:"column:rate"`
}

type CustomerOpeningBalance struct {
	BaseModel
	CustomerID uint64  `gorm:"index;column:customer_id"`
	Balance    float64 `gorm:"column:balance"`
	Status     string  `gorm:"column:status"`
}

type CreditLimitChangeLog struct {
	BaseModel
	CustomerID   uint64  `gorm:"index;column:customer_id"`
	CustomerType string  `gorm:"column:customer_type"`
	OldLimit     float64 `gorm:"column:old_limit"`
	CreditLimit  float64 `gorm:"column:credit_limit"`
	Action       string  `gorm:"column:action"`
}

type MemberBankAccount struct {
	BaseModel
	MemberID      uint64 `gorm:"index;column:member_id"`
	BankID        uint64 `gorm:"column:bank_id"`
	BankBranchId  uint64 `gorm:"column:bank_branch_id"`
	AccountNumber string `gorm:"column:account_number"`
	AccountName   string `gorm:"column:account_name"`
	Status        string `gorm:"column:status"`
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
	PayMonth  string    `gorm:"column:pay_month"`
	PayYear   string    `gorm:"column:pay_year"`
}

type CustomerCollection struct {
	BaseModel
	PayDateRangeID  uint64  `gorm:"index;column:pay_date_range_id"`
	PayrollMonth    int     `gorm:"column:payroll_month"`
	PayrollYear     int     `gorm:"column:payroll_year"`
	TotalDeliveries float64 `gorm:"column:total_deliveries"`
	TotalAmount     float64 `gorm:"column:total_amount"`
}

type CustomerInvoice struct {
	BaseModel
	CustomerID      uint64  `gorm:"index;column:customer_id"`
	InvoiceNo       string  `gorm:"uniqueIndex;column:invoice_no"`
	TotalDeliveries float64 `gorm:"column:total_deliveries"`
	Rate            float64 `gorm:"column:rate"`
	TotalAmount     float64 `gorm:"column:total_amount"`
	Status          string  `gorm:"column:status"`
}

type CustomerDocument struct {
	BaseModel
	CustomerID     uint64    `gorm:"index;column:customer_id"`
	DocDescription string    `gorm:"column:doc_description"`
	DocBalance     float64   `gorm:"column:doc_balance"`
	DueDate        time.Time `gorm:"column:due_date"`
}

type SalesSummary struct {
	BaseModel
	CustomerID      uint64    `gorm:"index;column:customer_id"`
	CustomerType    string    `gorm:"column:customer_type"`
	TotalSales      float64   `gorm:"column:total_sales"`
	TotalPayments   float64   `gorm:"column:total_payments"`
	Balance         float64   `gorm:"column:balance"`
	LastSaleDate    time.Time `gorm:"column:last_sale_date"`
	LastPaymentDate time.Time `gorm:"column:last_payment_date"`
}

type ProductSale struct {
	BaseModel
	ProductID    uint64    `gorm:"index;column:product_id"`
	CustomerID   uint64    `gorm:"index;column:customer_id"`
	Quantity     float64   `gorm:"column:quantity"`
	UnitPrice    float64   `gorm:"column:unit_price"`
	TotalAmount  float64   `gorm:"column:total_amount"`
	SaleDate     time.Time `gorm:"index;column:sale_date"`
	Reference    string    `gorm:"column:reference"`
	CustomerType string    `gorm:"column:customer_type"`
	Status       string    `gorm:"column:status"`
}

type MemberType struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

type MemberClass struct {
	BaseModel
	ClassName   string `gorm:"column:class_name"`
	Description string `gorm:"column:description"`
}

type WalletType struct {
	BaseModel
	Code        string `gorm:"uniqueIndex;column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

type Wallet struct {
	BaseModel
	WalletID      string  `gorm:"uniqueIndex;column:wallet_id"`
	WalletName    string  `gorm:"column:wallet_name"`
	MemberID      uint64  `gorm:"index;column:member_id"`
	AccountNumber string  `gorm:"column:account_number"`
	Balance       float64 `gorm:"column:balance"`
	UUID          string  `gorm:"uniqueIndex;column:uuid"`
	WalletTypeID  string  `gorm:"column:walletTypeId"`
}

type MemberAuthentication struct {
	BaseModel
	MemberID            uint64    `gorm:"index;column:member_id"`
	PhoneNumber         string    `gorm:"uniqueIndex;column:phone_number"`
	Password            string    `gorm:"column:password"`
	LastPasswordChanged time.Time `gorm:"column:last_password_changed"`
	AppPhoneUsed        string    `gorm:"column:app_phone_used"`
	RememberToken       string    `gorm:"column:remember_token"`
	MemberType          string    `gorm:"column:member_type"`
	PasswordReset       bool      `gorm:"column:password_reset"`
}

type MemberAuthenticationGroup struct {
	BaseModel
	AuthID uint64 `gorm:"index;column:auth_id"`
	Group  string `gorm:"column:group"`
}

type MemberKYCComment struct {
	BaseModel
	MemberID  uint64 `gorm:"index;column:member_id"`
	Issue     string `gorm:"column:issue"`
	Comment   string `gorm:"column:comment"`
	Iteration int    `gorm:"column:iteration"`
}

type MemberDependant struct {
	BaseModel
	MemberID     uint64    `gorm:"index;column:member_id"`
	Name         string    `gorm:"column:name"`
	NationalID   string    `gorm:"column:national_id"`
	Relationship string    `gorm:"column:relationship"`
	MobileNo     string    `gorm:"column:mobile_no"`
	Gender       string    `gorm:"column:gender"`
	DateOfBirth  time.Time `gorm:"column:date_of_birth"`
	BirthCertNo  string    `gorm:"column:birth_cert_no"`
	Email        string    `gorm:"column:email"`
	Address      string    `gorm:"column:address_address"`
	PostalCode   string    `gorm:"column:postal_code"`
	Town         string    `gorm:"column:town"`
}

type MemberBalanceBroughtForward struct {
	BaseModel
	MemberID uint64  `gorm:"index;column:member_id"`
	Month    string  `gorm:"column:month"`
	Year     string  `gorm:"column:year"`
	Amount   float64 `gorm:"column:amount"`
}

type MemberCashTransfer struct {
	BaseModel
	WithdrawalID    uint64  `gorm:"index;column:withdrawal_id"`
	Amount          float64 `gorm:"column:amount"`
	TransferType    string  `gorm:"column:transfer_type"`
	TransactionType string  `gorm:"column:transaction_type"`
}

type MemberOverUnderPayment struct {
	BaseModel
	MemberID         uint64  `gorm:"index;column:member_id"`
	Amount           float64 `gorm:"column:amount"`
	Period           string  `gorm:"column:period"`
	OverUnderPayment string  `gorm:"column:over_under_payment"`
	Reason           string  `gorm:"column:reason"`
}

type MemberDebt struct {
	BaseModel
	MemberID         uint64    `gorm:"index;column:member_id"`
	DeductionType    string    `gorm:"column:deduction_type"`
	TotalDebt        float64   `gorm:"column:total_debt"`
	DebtRecovered    float64   `gorm:"column:debt_recovered"`
	Balance          float64   `gorm:"column:balance"`
	TransactionsDate time.Time `gorm:"column:transactions_date"`
	ValueDelivered   float64   `gorm:"column:value_delivered"`
	Priority         int       `gorm:"column:priority"`
	NetAmount        float64   `gorm:"column:net_amount"`
	Period           int       `gorm:"column:period"`
	Year             int       `gorm:"column:year"`
}

type MemberPayroll struct {
	BaseModel
	PayrollMonth    string    `gorm:"column:payroll_month"`
	PayrollYear     string    `gorm:"column:payroll_year"`
	PayRateID       uint64    `gorm:"column:pay_rate_id"`
	PayDateRangeID  uint64    `gorm:"column:pay_date_range_id"`
	DateOpened      time.Time `gorm:"column:date_opened"`
	TotalKilos      float64   `gorm:"column:total_kilos"`
	TotalDeductions float64   `gorm:"column:total_deductions"`
	GrossPay        float64   `gorm:"column:gross_pay"`
	NetPay          float64   `gorm:"column:net_pay"`
	Complete        string    `gorm:"column:complete"`
	Confirmed       string    `gorm:"column:confirmed"`
	Approved        bool      `gorm:"column:approved"`
}

type MemberPayslip struct {
	BaseModel
	MemberID        uint64  `gorm:"index;column:member_id"`
	PayrollID       uint64  `gorm:"index;column:payroll_id"`
	PayRateID       uint64  `gorm:"column:pay_rate_id"`
	TotalKilos      float64 `gorm:"column:total_kilos"`
	GrossPay        float64 `gorm:"column:gross_pay"`
	TotalDeductions float64 `gorm:"column:total_deductions"`
	NetPay          float64 `gorm:"column:net_pay"`
	PayrollMonth    string  `gorm:"column:payroll_month"`
	PayrollYear     string  `gorm:"column:payroll_year"`
	PayDateRangeID  uint64  `gorm:"column:pay_date_range_id"`
	Complete        string  `gorm:"column:complete"`
}

type MoneyTransfer struct {
	BaseModel
	Type     string  `gorm:"column:type"` // mpesa/wallet_transfer
	MemberID uint64  `gorm:"index;column:member_id"`
	Amount   float64 `gorm:"column:amount"`
	Status   string  `gorm:"column:status"` // pending/success/failed
	Remarks  string  `gorm:"column:remarks"`
}

type MemberMpesaWithdrawal struct {
	BaseModel
	WithdrawalID uint64  `gorm:"index;column:withdrawal_id"`
	LoanID       uint64  `gorm:"index;column:loan_id"`
	WalletID     uint64  `gorm:"index;column:wallet_id"`
	Amount       float64 `gorm:"column:amount"`
	Status       string  `gorm:"column:status"`
}
type WalletWithdrawal struct {
	BaseModel
	WithdrawalUUID string `gorm:"uniqueIndex;column:withdrawal_uuid"`
	Status         string `gorm:"column:status"`
	LoanID         uint64 `gorm:"index;column:loan_id"`
	MemberID       uint64 `gorm:"index;column:member_id"`
}

type Loan struct {
	BaseModel
	MemberID              uint64    `gorm:"index;column:member_id"`
	Amount                float64   `gorm:"column:amount"`
	Interest              float64   `gorm:"column:interest"`
	TotalPayable          float64   `gorm:"column:total_payable"`
	Status                string    `gorm:"column:status"`
	ApprovedAmt           float64   `gorm:"column:approved_amount"`
	ProcessedBy           uint64    `gorm:"column:processed_by"`
	LoanLimitBy           uint64    `gorm:"column:loan_limit_by"`
	CreditLimit           uint64    `gorm:"column:credit_limit"`
	ReviewAccepted        bool      `gorm:"column:review_accepted"`
	UUID                  string    `gorm:"uniqueIndex;column:uuid"`
	DisbursedAt           time.Time `gorm:"column:disbursed_at"`
	ProcessedAt           time.Time `gorm:"column:processed_at"`
	RequestID             string    `gorm:"column:request_id"`
	TotalDue              float64   `gorm:"column:total_due"`
	RepaymentAmount       float64   `gorm:"column:repayment_amount"`
	WithdrawalRequestUUID string    `gorm:"column:withdrawal_request_uuid"`
}

type LoanRepayment struct {
	BaseModel
	LoanID uint64    `gorm:"index;column:loan_id"`
	Amount float64   `gorm:"column:amount"`
	Date   time.Time `gorm:"index;column:date"`
}

type LoanCallback struct {
	BaseModel
	Detail string `gorm:"column:detail"`
	LoanID uint64 `gorm:"index;column:loan_id"`
	Type   string `gorm:"column:type"`
}

type LoanOrganizationProfile struct {
	BaseModel
	NextLevel       string `gorm:"column:next_level"`
	AstraID         string `gorm:"column:astra_id"`
	LinkStatus      string `gorm:"column:link_status"`
	UUID            string `gorm:"uniqueIndex;column:uuid"`
	Version         string `gorm:"column:version"`
	ProductID       string `gorm:"column:product_id"`
	CompanyDetailID uint64 `gorm:"column:company_detail_id"`
	ManuallyRatify  bool   `gorm:"column:manually_ratify"`
}

type LoanOriginationCallbackLog struct {
	BaseModel
	AstraDetail string `gorm:"column:astra_detail"`
	SyncAttempt uint64 `gorm:"column:sync_attempt"`
}

// AI & Insemination
type AIService struct {
	BaseModel
	CattleBreedID uint64    `gorm:"index;column:cattle_breed_id"`
	ServiceDate   time.Time `gorm:"index;column:service_date"`
	MemberID      uint64    `gorm:"index;column:member_id"`
	Owner         string    `gorm:"column:owner"`
	AnimalName    string    `gorm:"column:animal_name"`
	Status        string    `gorm:"column:status"`
	Notes         string    `gorm:"column:notes"`
}

type Insemination struct {
	BaseModel
	AIServiceID         uint64    `gorm:"index;column:ai_service_id"`
	InseminationSemenID uint64    `gorm:"index;column:insemination_semen_id"`
	InseminationDate    time.Time `gorm:"index;column:insemination_date"`
	ExpectedCalvingDate time.Time `gorm:"column:expected_calving_date"`
	InseminationPerson  string    `gorm:"column:insemination_person"`
}

type InseminationSemen struct {
	BaseModel
	CattleBreedID uint64 `gorm:"index;column:cattle_breed_id"`
	Code          string `gorm:"uniqueIndex;column:code"`
	BullName      string `gorm:"column:bull_name"`
	Local         string `gorm:"column:local"`
}

type InseminationItem struct {
	BaseModel
	Code             string `gorm:"uniqueIndex;column:code"`
	Name             string `gorm:"column:name"`
	SupplierRequired bool   `gorm:"column:supplier_required"`
}

type InseminationCost struct {
	BaseModel
	InseminationID uint64  `gorm:"index;column:insemination_id"`
	ItemID         uint64  `gorm:"index;column:insemination_item_id"`
	UnitCost       float64 `gorm:"column:unit_cost"`
	Quantity       float64 `gorm:"column:quantity"`
}

type InseminationItemSupplier struct {
	BaseModel
	InseminationItemID uint64  `gorm:"index;column:insemination_item_id"`
	SupplierID         uint64  `gorm:"index;column:supplier_id"`
	Cost               float64 `gorm:"column:cost"`
}

type CattleBreed struct {
	BaseModel
	Code string `gorm:"uniqueIndex;column:code"`
	Name string `gorm:"column:name"`
}

// HR Payroll
type Employee struct {
	BaseModel
	UserID            uint64    `gorm:"column:user_id"`
	Surname           string    `gorm:"column:surname"`
	FirstName         string    `gorm:"column:first_name"`
	MiddleName        string    `gorm:"column:middle_name"`
	EmployeeNo        string    `gorm:"uniqueIndex;column:employee_no"`
	IDNo              string    `gorm:"uniqueIndex;column:id_no"`
	KraPin            string    `gorm:"column:kra_pin"`
	NssfNo            string    `gorm:"column:nssf_no"`
	NhifNo            string    `gorm:"column:nhif_no"`
	Gender            string    `gorm:"column:gender"`
	DateOfBirth       time.Time `gorm:"column:date_of_birth"`
	Phone             string    `gorm:"column:phone_number"`
	Email             string    `gorm:"column:email_address"`
	JobPositionID     uint64    `gorm:"column:job_position_id"`
	Status            int       `gorm:"column:status"`
	Title             string    `gorm:"column:title"`
	PassportNo        string    `gorm:"column:passport_no"`
	Town              string    `gorm:"column:town"`
	SiteID            uint64    `gorm:"column:site_id"`
	SalesSummary      string    `gorm:"column:sales_summary"`
	MaritalStatus     string    `gorm:"column:marital_status"`
	Religion          string    `gorm:"column:religion"`
	Disabled          bool      `gorm:"column:disabled"`
	StoreID           uint64    `gorm:"column:store_id"`
	PostalAddress     string    `gorm:"column:postal_address"`
	PostalCode        string    `gorm:"column:postal_code"`
	BirthCity         string    `gorm:"column:birth_city"`
	NextOfKinFullName string    `gorm:"column:next_of_kin_full_name"`
	NextOfKinPhone    string    `gorm:"column:next_of_kin_phone"`
	PassportPhoto     string    `gorm:"column:passport_photo"`
	IdFrontPhoto      string    `gorm:"column:id_front_photo"`
}

type EmployeeDetail struct {
	BaseModel
	EmployeeID    uint64    `gorm:"index;column:employee_id"`
	Gender        string    `gorm:"column:gender"`
	MaritalStatus string    `gorm:"column:marital_status"`
	Religion      string    `gorm:"column:religion"`
	Ethnicity     string    `gorm:"column:ethnicity"`
	Disabled      bool      `gorm:"column:employee_disabled"`
	DateOfBirth   time.Time `gorm:"column:date_of_birth"`
	JobCategoryID uint64    `gorm:"column:job_categroy_id"`
	Seconded      bool      `gorm:"column:seconded"`
}

type EmployeeBenefit struct {
	BaseModel
	EmployeeID uint64  `gorm:"index;column:employee_id"`
	BenefitID  uint64  `gorm:"index;column:benefit_id"`
	Amount     float64 `gorm:"column:amount"`
	Status     string  `gorm:"column:status"`
}

type EmployeeDeductionType struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	IsStatutory bool   `gorm:"column:is_statutory"`
}

type EmployeePayrollBenefit struct {
	BaseModel
	EmployeeID uint64  `gorm:"index;column:employee_id"`
	BenefitID  uint64  `gorm:"index;column:employee_benefit_id"`
	Amount     float64 `gorm:"column:amount"`
	Year       string  `gorm:"column:benefit_year"`
	Month      string  `gorm:"column:benefit_month"`
	PayrollID  uint64  `gorm:"index;column:payroll_id"`
}

type EmployeeDeduction struct {
	BaseModel
	EmployeeID      uint64  `gorm:"index;column:employee_id"`
	DeductionTypeID uint64  `gorm:"index;column:deduction_type_id"`
	Amount          float64 `gorm:"column:amount"`
	Status          bool    `gorm:"column:status"`
}

type EmployeeLeaveDetail struct {
	BaseModel
	EmployeeID    uint64 `gorm:"index;column:employee_id"`
	BalanceBF     string `gorm:"column:balance_bf"`
	AllocatedDays int    `gorm:"column:allocated_days"`
}

type EmployeePayrollDeduction struct {
	BaseModel
	EmployeeID  uint64  `gorm:"index;column:employee_id"`
	DeductionID uint64  `gorm:"index;column:employee_deduction_id"`
	Amount      float64 `gorm:"column:amount"`
	Month       string  `gorm:"column:deduction_month"`
	Year        string  `gorm:"column:deduction_year"`
	PayrollID   uint64  `gorm:"index;column:payroll_id"`
}

type EmployeeDependant struct {
	BaseModel
	EmployeeID   uint64 `gorm:"index;column:employee_id"`
	Name         string `gorm:"column:name"`
	Relationship string `gorm:"column:relationship"`
}

type EmployeeContractDetail struct {
	BaseModel
	ContractType    string    `gorm:"column:contract_type"`
	ContractEndDate time.Time `gorm:"column:contract_end_date"`
	NoticePeriod    string    `gorm:"column:notice_period"`
	RetirementDate  time.Time `gorm:"column:retirement_date"`
}

type EmployeeExitDetail struct {
	BaseModel
	ContractType    string    `gorm:"column:contract_type"`
	ContractEndDate time.Time `gorm:"column:contract_end_date"`
	DateOfLeaving   time.Time `gorm:"column:dateof_living"`
	ExitCategory    string    `gorm:"column:exit_category"`
	Reasons         string    `gorm:"column:reasonsfoexit"`
}

type EmployeeTerminationCategory struct {
	BaseModel
	TerminationCategory string `gorm:"column:termination_category"`
}

type EmployeeDocument struct {
	BaseModel
	EmployeeID      uint64 `gorm:"index;column:employee_id"`
	DocumentTypeID  uint64 `gorm:"column:document_type_id"`
	FileName        string `gorm:"column:file_name"`
	FileDescription string `gorm:"column:file_description"`
}

type EmployeeProfessionalTitle struct {
	BaseModel
	Code string `gorm:"uniqueIndex;column:code"`
	Name string `gorm:"column:name"`
}

type EmployeeQualification struct {
	BaseModel
	EmployeeID    uint64    `gorm:"index;column:employee_id"`
	Qualification string    `gorm:"column:qualification"`
	Institution   string    `gorm:"column:institution"`
	StartDate     time.Time `gorm:"column:start_date"`
	EndDate       time.Time `gorm:"column:end_date"`
}

type InsuranceDetail struct {
	BaseModel
	EmployeeID       uint64    `gorm:"index;column:employee_id"`
	CompanyName      string    `gorm:"column:company_name"`
	PolicyNo         string    `gorm:"column:policy_no"`
	CommencementDate time.Time `gorm:"column:commencement_date"`
	MaturityDate     time.Time `gorm:"column:maturity_date"`
	SumAssured       float64   `gorm:"column:sum_assured"`
}

type EmployeeSalary struct {
	BaseModel
	EmployeeID  uint64  `gorm:"index;column:employee_id"`
	BasicSalary float64 `gorm:"column:basic_salary"`
	Status      string  `gorm:"column:status"`
}

type EmployeeBankAccount struct {
	BaseModel
	EmployeeID    uint64 `gorm:"index;column:employee_id"`
	BankID        uint64 `gorm:"index;column:bank_id"`
	AccountNumber string `gorm:"column:account_number"`
	AccountName   string `gorm:"column:account_name"`
}

type EmployeeLeaveApplication struct {
	BaseModel
	ApplicationNo string    `gorm:"uniqueIndex;column:application_no"`
	EmployeeID    uint64    `gorm:"index;column:employee_id"`
	LeaveTypeID   uint64    `gorm:"index;column:leave_type_id"`
	DaysApplied   float64   `gorm:"column:days_applied"`
	DaysApproved  float64   `gorm:"column:days_approved"`
	StartDate     time.Time `gorm:"column:start_date"`
	EndDate       time.Time `gorm:"column:end_date"`
	ReturnDate    time.Time `gorm:"column:return_date"`
	ApproverID    uint64    `gorm:"column:approver_id"`
	Status        string    `gorm:"column:status"`
	Approved      bool      `gorm:"column:approved"`
}

type EmployeeLeaveAssignment struct {
	BaseModel
	EmployeeID         uint64 `gorm:"index;column:employee_id"`
	LeaveApplicationID uint64 `gorm:"index;column:leave_application_id"`
	RelieverID         uint64 `gorm:"column:reliever_id"`
}

type EmployeeLeaveType struct {
	BaseModel
	Code        string  `gorm:"uniqueIndex;column:code"`
	Description string  `gorm:"column:description"`
	Days        float64 `gorm:"column:days"`
	Gender      string  `gorm:"column:gender"`
}

type EmployeeRelief struct {
	BaseModel
	EmployeeID uint64 `gorm:"index;column:employee_id"`
	ReliefID   uint64 `gorm:"index;column:relief_id"`
	Status     string `gorm:"column:status"`
}

type Benefit struct {
	BaseModel
	Name         string `gorm:"column:name"`
	IsTaxable    bool   `gorm:"column:is_taxable"`
	DefaultValue string `gorm:"column:default_value"`
}

type DeductionType struct {
	BaseModel
	Code        string `gorm:"column:code"`
	Description string `gorm:"column:description"`
	Status      string `gorm:"column:status"`
	IsStatutory string `gorm:"column:is_statutory"`
}

func (DeductionType) TableName() string {
	return "deduction_types"
}

type DeductionPricingRule struct {
	BaseModel
	DeductionTypeID uint64  `gorm:"column:deduction_type_id"`
	MinCreditLimit  float64 `gorm:"column:min_credit_limit"`
	MaxLimit        float64 `gorm:"column:max_limit"`
	BoardingFee     float64 `gorm:"column:boarding_fee"`
	ProcessingFee   float64 `gorm:"column:processing_fee"`
	InsuranceFee    float64 `gorm:"column:insurance_fee"`
	LegalFee        float64 `gorm:"column:legal_fee"`
	InterestRate    float64 `gorm:"column:interest_rate"`
	Status          string  `gorm:"column:status;default:ACTIVE"`
}

func (DeductionPricingRule) TableName() string {
	return "deduction_pricing_rules"
}

type Deduction struct {
	BaseModel
	CustomerID      uint64  `gorm:"index;column:customer_id"`
	DeductionTypeID uint64  `gorm:"index;column:deduction_type_id"`
	Amount          float64 `gorm:"column:amount"`
	CustomerType    string  `gorm:"column:customer_type"` // member, employee, etc.
	Confirmed       bool    `gorm:"column:confirmed"`
}

type DeductionTypeRateCard struct {
	BaseModel
	DeductionTypeID uint64  `gorm:"index;column:deduction_type_id"`
	MaxLimit        float64 `gorm:"column:max_limit"`
	InterestRate    float64 `gorm:"column:interest_rate"`
	ProcessingFee   float64 `gorm:"column:processing_fee"`
}

type JobCategory struct {
	BaseModel
	Code string `gorm:"uniqueIndex;column:code"`
	Name string `gorm:"column:name"`
}

type JobDetail struct {
	BaseModel
	JobPositionID uint64 `gorm:"index;column:job_position_id"`
	JobTitle      string `gorm:"column:job_title"`
	Department    string `gorm:"column:department"`
}

type JobHODRequisition struct {
	BaseModel
	JobPositionID uint64 `gorm:"index;column:job_position_id"`
	HOD           string `gorm:"column:hod"`
}

type JobGrade struct {
	BaseModel
	Code            string  `gorm:"uniqueIndex;column:code"`
	Name            string  `gorm:"column:name"`
	MinSalary       float64 `gorm:"column:min_salary"`
	MaxSalary       float64 `gorm:"column:max_salary"`
	YearlyIncrement float64 `gorm:"column:yearly_increment"`
}

type JobReasonToFillVacancy struct {
	BaseModel
	Reason string `gorm:"column:reason"`
}

type JobRequisition struct {
	BaseModel
	JobPositionID     uint64 `gorm:"index;column:job_position_id"`
	RequiredPositions int    `gorm:"column:required_positions"`
	Status            string `gorm:"column:status"`
}

type JobPosition struct {
	BaseModel
	Code           string `gorm:"uniqueIndex;column:code"`
	Name           string `gorm:"column:name"`
	JobDescription string `gorm:"column:job_description"`
	DepartmentID   uint64 `gorm:"index;column:department_id"`
}

type JobQualificationType struct {
	BaseModel
	QualificationType string `gorm:"column:qualification_type"`
	Description       string `gorm:"column:description"`
}

type JobQualification struct {
	BaseModel
	QualificationTypeID uint64 `gorm:"index;column:qualification_type_id"`
	Code                string `gorm:"uniqueIndex;column:qualification_code"`
	Qualification       string `gorm:"column:qualification"`
}

type EmployeePayroll struct { // Renamed from Payroll to avoid conflict and represent the header
	BaseModel
	PayrollMonth    string    `gorm:"column:payroll_month"`
	PayrollYear     string    `gorm:"column:payroll_year"`
	DateOpened      time.Time `gorm:"column:date_opened"`
	TotalDeductions float64   `gorm:"column:total_deductions"`
	GrossPay        float64   `gorm:"column:gross_pay"`
	NetPay          float64   `gorm:"column:net_pay"`
	Complete        string    `gorm:"column:complete"`
	Confirmed       string    `gorm:"column:confirmed"`
	Approved        string    `gorm:"column:approved"`
	TotalBenefits   float64   `gorm:"column:total_benefits"`
	TotalTax        float64   `gorm:"column:total_tax"`
	TotalRelief     float64   `gorm:"column:total_relief"`
	Period          string    `gorm:"column:period"`
	PaidAt          time.Time `gorm:"column:paid_at"`
}

type EmployeePayslip struct {
	BaseModel
	EmployeeID      uint64  `gorm:"index;column:employee_id"`
	PayrollMonth    string  `gorm:"column:payroll_month"`
	PayrollYear     string  `gorm:"column:payroll_year"`
	GrossPay        float64 `gorm:"column:gross_pay"`
	NetPay          float64 `gorm:"column:net_pay"`
	TotalDeductions float64 `gorm:"column:total_deductions"`
	TotalBenefits   float64 `gorm:"column:total_benefits"`
	BasicSalary     float64 `gorm:"column:basic_salary"`
	PayrollID       uint64  `gorm:"index;column:payroll_id"`
	TotalTax        float64 `gorm:"column:total_tax"`
	TotalRelief     float64 `gorm:"column:total_relief"`
}

type EmployeePayrollEntry struct { // Original Payroll model, renamed to reflect its role as an entry
	BaseModel
	EmployeeID uint64    `gorm:"index;column:employee_id"`
	Amount     float64   `gorm:"column:amount"`
	Period     string    `gorm:"column:period"`
	PaidAt     time.Time `gorm:"column:paid_at"`
}

type EmployerPayrollDeduction struct {
	BaseModel
	EmployeeID  uint64  `gorm:"index;column:employee_id"`
	DeductionID uint64  `gorm:"index;column:employee_deduction_id"`
	Amount      float64 `gorm:"column:amount"`
	Month       string  `gorm:"column:month"`
	Year        string  `gorm:"column:year"`
	PayrollID   uint64  `gorm:"index;column:payroll_id"`
}

type EmployeePayrollRelief struct {
	BaseModel
	ReliefID   uint64 `gorm:"index;column:relief_id"`
	EmployeeID uint64 `gorm:"index;column:employee_id"`
	Amount     string `gorm:"column:amount"`
	PayrollID  uint64 `gorm:"index;column:payroll_id"`
}

// Inventory
type ItemCategory struct {
	BaseModel
	Name             string `gorm:"column:name"`
	Description      string `gorm:"column:description"`
	ParentCategoryID uint64 `gorm:"column:parent_category_id"`
}

func (ItemCategory) TableName() string {
	return "item_categories"
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

type StoreStockMovement struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	StoreID         uint64    `gorm:"column:store_id"`
	ItemID          uint64    `gorm:"column:item_id"`
	MovementType    string    `gorm:"column:movement_type"`
	ReferenceTable  string    `gorm:"column:reference_table"`
	ReferenceID     uint64    `gorm:"column:reference_id"`
	QtyIn           float64   `gorm:"column:qty_in"`
	QtyOut          float64   `gorm:"column:qty_out"`
	BalanceAfter    float64   `gorm:"column:balance_after"`
	UnitCost        float64   `gorm:"column:unit_cost"`
	SellingPrice    float64   `gorm:"column:selling_price"`
	Remarks         string    `gorm:"column:remarks"`
	CreatedBy       uint64    `gorm:"column:created_by"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (StoreStockMovement) TableName() string {
	return "store_stock_movements"
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
	TransactionID int64   `gorm:"column:transaction_id"`
}

func (StoreSale) TableName() string {
	return "store_sales"
}

type StoreSaleItem struct {
	BaseModel
	ItemID      uint64 `gorm:"column:item_id"`
	Quantity    int    `gorm:"column:quantity"`
	UnitPrice   string `gorm:"column:unit_price"`
	Total       string `gorm:"column:total"`
	StoreSaleID uint64 `gorm:"column:store_sale_id"`
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

type Product struct {
	BaseModel
	ItemCategoryID uint64  `gorm:"index;column:item_category_id"`
	Description    string  `gorm:"column:description"`
	ReorderPoint   int     `gorm:"column:reorder_point"`
	BuyingPrice    float64 `gorm:"column:buying_price"`
	SellingPrice   float64 `gorm:"column:selling_price"`
	Status         string  `gorm:"column:status"`
	InventoryID    int64   `gorm:"index;column:inventory_id"`
	Thumbnail      string  `gorm:"column:thumbnail"`
	InventoryName  string  `gorm:"column:inventory_name"`
}

type ProductGrade struct {
	BaseModel
	Name        string `gorm:"column:name"` // Changed from uniqueIndex to match schema
	Description string `gorm:"column:description"`
}

func (ProductGrade) TableName() string {
	return "product_grades"
}

type DefaultMilkRate struct {
	BaseModel
	Rate    float64 `gorm:"column:rate"`
	RouteID uint64  `gorm:"index;column:route_id"`
}

func (DefaultMilkRate) TableName() string {
	return "default_milk_rates"
}

type ProductPrice struct {
	BaseModel
	ProductID     uint64    `gorm:"index;column:product_id"`
	Price         float64   `gorm:"column:price"`
	EffectiveDate time.Time `gorm:"index;column:effective_date"`
	Currency      string    `gorm:"column:currency"`
}

type ProductOpeningBalance struct {
	BaseModel
	ProductID   uint64    `gorm:"index;column:product_id"`
	Quantity    float64   `gorm:"column:quantity"`
	UnitPrice   float64   `gorm:"column:unit_price"`
	TotalValue  float64   `gorm:"column:total_value"`
	BalanceDate time.Time `gorm:"column:balance_date"`
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

type IncomingGoodsReceipt struct {
	BaseModel
	ReceiptNumber   string    `gorm:"uniqueIndex;column:receipt_number"`
	Description     string    `gorm:"type:text;column:description"`
	TransactionDate time.Time `gorm:"index;column:transaction_date"`
	InventoryID     uint64    `gorm:"index;column:inventory_id"`
	Status          string    `gorm:"column:status"`
}

type IncomingGoodsReceiptItem struct {
	BaseModel
	IncomingGoodsReceiptID uint64    `gorm:"index;column:incoming_goods_receipt_id"`
	TransactionDate        time.Time `gorm:"index;column:transaction_date"`
	InventoryID            uint64    `gorm:"index;column:inventory_id"`
	Quantity               float64   `gorm:"column:quantity"`
	UnitCost               float64   `gorm:"column:unit_cost"`
	Confirmed              bool      `gorm:"column:confirmed"`
}

type SystemSetting struct {
	BaseModel
	Key   string `gorm:"uniqueIndex;column:key"`
	Value string `gorm:"type:text;column:value"`
}

// SMS

type OrganizationAddress struct {
	BaseModel
	AddressType string `gorm:"column:address_type"`
	City        string `gorm:"column:city"`
	Code        string `gorm:"column:code"`
	Country     string `gorm:"column:country"`
	Line1       string `gorm:"column:line1"`
	Line2       string `gorm:"column:line2"`
	Line3       string `gorm:"column:line3"`
	State       string `gorm:"column:state"`
}

type OrganizationDocument struct {
	BaseModel
	AstraID      uint64 `gorm:"index;column:astra_id"`
	DocumentType string `gorm:"column:document_type"`
	Document     string `gorm:"type:text;column:document"`
	Submitted    bool   `gorm:"column:submitted"`
}

type Bank struct {
	BaseModel
	Name        string `gorm:"column:name"`
	SwiftCode   string `gorm:"column:swift_code"`
	Description string `gorm:"column:description"`
}

type BankBranch struct {
	BaseModel
	Name     string `gorm:"column:name"`
	BankId   uint64 `gorm:"column:bank_id"`
	Location string `gorm:"column:location"`
}

type OrganizationBank struct {
	BaseModel
	Name string `gorm:"column:name"`
}

type PaymentMode struct {
	BaseModel
	Code string `gorm:"uniqueIndex;column:code"`
	Name string `gorm:"column:name"`
}

type PaymentTerm struct {
	BaseModel
	Code        string `gorm:"uniqueIndex;column:code"`
	Description string `gorm:"type:text;column:description"`
}

type OrganizationMember struct {
	BaseModel
	ManuallyRatify bool   `gorm:"column:manually_ratify"`
	NextLevel      string `gorm:"column:next_level"`
	Status         string `gorm:"column:status"`
	AstraID        string `gorm:"column:astra_id"`
	CreditLimit    uint64 `gorm:"column:credit_limit"`
	LinkStatus     string `gorm:"column:link_status"`
	LivenessPassed bool   `gorm:"column:liveness_passed"`
	AstraRemarks   string `gorm:"type:text;column:astra_remarks"`
	UUID           string `gorm:"uniqueIndex;column:uuid"`
	AuthCreated    bool   `gorm:"column:auth_created"`
	Locale         string `gorm:"column:locale"`
	CustomerID     uint64 `gorm:"index;column:customer_id"`
	CustomerType   string `gorm:"column:customer_type"`
}

type OrganizationLeadership struct {
	BaseModel
	FirstName      string `gorm:"column:first_name"`
	LastName       string `gorm:"column:last_name"`
	PrimaryPhone   string `gorm:"column:primary_phone"`
	IDDateOfIssue  string `gorm:"column:id_date_of_issue"`
	NextLevel      string `gorm:"column:next_level"`
	SecondaryPhone string `gorm:"column:secondary_phone"`
	IDNo           string `gorm:"uniqueIndex;column:id_no"`
	IDFrontPhoto   string `gorm:"column:id_front_photo"`
	IDBackPhoto    string `gorm:"column:id_back_photo"`
	BirthCity      string `gorm:"column:birth_city"`
	BirthCountry   string `gorm:"column:birth_country"`
	Email          string `gorm:"column:email"`
	Title          string `gorm:"column:title"`
	Position       string `gorm:"column:position"`
	Locale         string `gorm:"column:locale"`
	Status         string `gorm:"column:status"`
	TaxNumber      string `gorm:"column:tax_number"`
	MaritalStatus  string `gorm:"column:marital_status"`
	Gender         string `gorm:"column:gender"`
	AstraID        string `gorm:"column:astra_id"`
	LinkStatus     string `gorm:"column:link_status"`
	LivenessPassed bool   `gorm:"column:liveness_passed"`
	Submitted      bool   `gorm:"column:submitted"`
	KraNumber      string `gorm:"column:kra_number"`
	UUID           string `gorm:"uniqueIndex;column:uuid"`
	AstraRemarks   string `gorm:"type:text;column:astra_remarks"`
}

type OrganizationWallet struct {
	BaseModel
	WalletTypeID uint64 `gorm:"column:walletTypeId"`
	WalletID     string `gorm:"uniqueIndex;column:wallet_id"`
	WalletName   string `gorm:"column:wallet_name"`
}

type OrganizationKybComment struct {
	BaseModel
	Issue     string `gorm:"type:text;column:issue"`
	Comment   string `gorm:"type:text;column:comment"`
	Iteration int    `gorm:"column:iteration"`
}

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

type Guest struct {
	BaseModel
	Name        string `gorm:"column:name"`
	PhoneNumber string `gorm:"uniqueIndex;column:phone_number"`
}

type TaxRelief struct {
	BaseModel
	Code        string `gorm:"uniqueIndex;column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"type:text;column:description"`
}

type BoardMember struct {
	BaseModel
	IDNo     string `gorm:"uniqueIndex;column:id_no"`
	Names    string `gorm:"column:names"`
	Position string `gorm:"column:position"`
	Phone    string `gorm:"column:phone"`
	Status   string `gorm:"column:status"`
}

type BoardMemberPayment struct {
	BaseModel
	IDNumber          string  `gorm:"index;column:id_number"`
	TransactionNumber string  `gorm:"uniqueIndex;column:transaction_number"`
	Amount            float64 `gorm:"column:amount"`
	Paye              float64 `gorm:"column:paye"`
	Month             int     `gorm:"column:month"`
	Year              int     `gorm:"column:year"`
}

type ExchangeVisit struct {
	BaseModel
	Partner    string    `gorm:"column:exchange_visit_partner"`
	VisitDate  time.Time `gorm:"index;column:exchange_visit_date"`
	Purpose    string    `gorm:"type:text;column:purpose"`
	Venue      string    `gorm:"column:venue"`
	EmployeeID uint64    `gorm:"index;column:exchange_visit_employee_id"`
	VisitNotes string    `gorm:"type:text;column:visit_notes"`
}

type ExchangeVisitAttendee struct {
	BaseModel
	ExchangeVisitID      uint64 `gorm:"column:exchange_visit_id"`
	Attendee             string `gorm:"column:attendee"`
	AttendeeOrganization string `gorm:"column:attendee_organization"`
	AttendeeDesignation  string `gorm:"column:attendee_designation"`
	Attended             string `gorm:"column:attended;default:0"`
	Comments             string `gorm:"column:comments"`
	AttendanceEmployeeID uint64 `gorm:"column:attendance_employee_id"`
}

func (ExchangeVisitAttendee) TableName() string {
	return "exchange_visit_attendees"
}

type SMSLog struct {
	BaseModel
	Phone   string `gorm:"index;column:phone"`
	Message string `gorm:"type:text;column:message"`
	Status  string `gorm:"column:status"`
	Error   string `gorm:"type:text;column:error"`
}

type TransportRate struct {
	BaseModel
	RouteID       uint64      `gorm:"index;column:route_id"`
	Route         Route       `gorm:"foreignKey:RouteID" json:"route"`
	TransporterID uint64      `gorm:"index;column:transporter_id"`
	Transporter   Transporter `gorm:"foreignKey:TransporterID" json:"transporter"`
	Rate          float64     `gorm:"column:transport_rate"`
	MemberID      uint64      `gorm:"index;column:member_id"`
	Member        Member      `gorm:"foreignKey:MemberID" json:"member"`
	Status        string      `gorm:"column:status"`
}

type Training struct {
	BaseModel
	Topic          string    `gorm:"column:topic"`
	Description    string    `gorm:"column:description"`
	Venue          string    `gorm:"column:venue"`
	TrainingUserID uint64    `gorm:"column:training_user_id"`
	TrainingDate   time.Time `gorm:"column:training_date"`
	Status         string    `gorm:"column:status"`
}

type TrainingSession struct {
	BaseModel
	TrainingID uint64 `gorm:"column:training_id"`
	MemberID   uint64 `gorm:"column:member_id"`
	Status     string `gorm:"column:status"`
	Remarks    string `gorm:"column:remarks"`
}

type TrainingAttendee struct {
	BaseModel
	TrainingSessionID uint64 `gorm:"column:training_session_id"`
	Names             string `gorm:"column:names"`
	IDNumber          string `gorm:"column:id_number"`
	PhoneNumber       string `gorm:"column:phone_number"`
	MembershipNumber  string `gorm:"column:membership_number"`
	Comments          string `gorm:"column:comments"`
	MemberID          uint64 `gorm:"column:member_id"`
}

func (TrainingAttendee) TableName() string {
	return "training_session_attendees"
}

type RouteCenter struct {
	BaseModel
	RouteID uint64 `gorm:"column:route_id"`
	Center  string `gorm:"column:center"`
}

func (RouteCenter) TableName() string {
	return "route_centers"
}

type ShareType struct {
	BaseModel
	ShareCode         string  `gorm:"column:share_code"`
	ShareType         string  `gorm:"column:share_type"`
	Description       string  `gorm:"column:description"`
	Rate              float64 `gorm:"column:rate"`
	Mandatory         int     `gorm:"column:madatory"` // Matches 'madatory' in schema
	HasShareValue     string  `gorm:"column:has_share_value"`
	RepayMethod       string  `gorm:"column:repay_method"`
	CalculatingMethod string  `gorm:"column:calculating_method"`
	ShareValue        float64 `gorm:"column:share_value"`
	DeductionTypeID   uint64  `gorm:"column:deduction_type_id"`
	Priority          int     `gorm:"column:priority"`
}

func (ShareType) TableName() string {
	return "share_types"
}

type ShareDividend struct {
	BaseModel
	DeclarationID int64   `gorm:"column:declaration_id"`
	MemberID      uint64  `gorm:"column:member_id"`
	FiscalYear    int     `gorm:"column:fiscal_year"`
	Period        int     `gorm:"column:period"`
	ShareUnits    float64 `gorm:"column:share_units"`
	Status        string  `gorm:"column:status;default:CALCULATED"`
	RatePerShare  float64 `gorm:"column:rate_per_share"`
	TaxAmount     float64 `gorm:"column:tax_amount"`
	NetAmount     float64 `gorm:"column:net_amount"`
	TransactionID int64   `gorm:"column:transaction_id"`
}

func (ShareDividend) TableName() string {
	return "share_dividends"
}

type SharePayment struct {
	BaseModel
	TransactionID   uint64    `gorm:"column:transaction_id"`
	MemberID        uint64    `gorm:"column:member_id"`
	ShareAccountID  uint64    `gorm:"column:share_account_id"`
	AmountPaid      float64   `gorm:"column:amount_paid"`
	ShareUnits      float64   `gorm:"column:share_units"`
	PaymentModeID   uint64    `gorm:"column:payment_mode_id"`
	Description     string    `gorm:"column:description"`
	Status          string    `gorm:"column:status;default:PENDING"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	ApprovedBy      uint64    `gorm:"column:approved_by"`
	DateApproved    time.Time `gorm:"column:date_approved"`
}

func (SharePayment) TableName() string {
	return "share_payments"
}

type ShareTransaction struct {
	BaseModel
	TransactionID   uint64    `gorm:"column:transaction_id"`
	ShareAccountID  uint64    `gorm:"column:share_account_id"`
	MemberID        uint64    `gorm:"column:member_id"`
	TransactionType string    `gorm:"column:transaction_type"`
	ShareUnits      float64   `gorm:"column:share_units;default:0.0000"`
	UnitPrice       float64   `gorm:"column:unit_price;default:0.00"`
	Debit           float64   `gorm:"column:debit;default:0.00"`
	Credit          float64   `gorm:"column:credit;default:0.00"`
	BalanceAfter    float64   `gorm:"column:balance_after;default:0.00"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
}

func (ShareTransaction) TableName() string {
	return "share_transactions"
}

type ShareTransfer struct {
	BaseModel
	TransactionID   uint64    `gorm:"column:transaction_id"`
	FromMemberID    uint64    `gorm:"column:from_member_id"`
	ToMemberID      uint64    `gorm:"column:to_member_id"`
	ShareUnits      float64   `gorm:"column:share_units"`
	TransferAmount  float64   `gorm:"column:transfer_amount"`
	Status          string    `gorm:"column:status;default:PENDING"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	ApprovedBy      uint64    `gorm:"column:approved_by"`
	DateApproved    time.Time `gorm:"column:date_approved"`
}

func (ShareTransfer) TableName() string {
	return "share_transfers"
}

type ShareAccount struct {
	BaseModel
	MemberID    uint64    `gorm:"column:member_id"`
	ShareTypeID uint64    `gorm:"column:share_type_id"`
	Status      string    `gorm:"column:status;default:ACTIVE"`
	OpenedAt    time.Time `gorm:"column:opened_at;default:CURRENT_TIMESTAMP"`
}

func (ShareAccount) TableName() string {
	return "share_accounts"
}

type DividendDeclaration struct {
	BaseModel
	FiscalYear      int       `gorm:"column:fiscal_year"`
	Period          int       `gorm:"column:period"`
	TotalPool       float64   `gorm:"column:total_pool"`
	RatePerShare    float64   `gorm:"column:rate_per_share"`
	CalculationType string    `gorm:"column:calculation_type"`
	Status          string    `gorm:"column:status;default:DRAFT"`
	ApprovedBy      uint64    `gorm:"column:approved_by"`
	ApprovedAt      time.Time `gorm:"column:approved_at"`
}

func (DividendDeclaration) TableName() string {
	return "dividend_declarations"
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
	TransactionID           *uint64        `gorm:"column:transaction_id"`
	CreatedAt               time.Time      `gorm:"column:created_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

func (AssetDepreciationEntry) TableName() string {
	return "asset_depreciation_entries"
}
