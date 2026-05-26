package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*
|--------------------------------------------------------------------------
| Base Model (shared fields)
|--------------------------------------------------------------------------
*/

type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
	CreatedBy uint64         `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy uint64         `gorm:"column:updated_by" json:"updated_by,omitempty"`
}

// System Models
type Installation struct {
	BaseModel
	InstallationDate time.Time `gorm:"column:installation_date"`
	ExpiryDate       time.Time `gorm:"column:expiry_date"`
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

type CanMovement struct {
	BaseModel
	CanID             uint64    `gorm:"column:can_id"`
	MovementType      string    `gorm:"column:movement_type"`
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

type StoreStockMovement struct {
	BaseModel
	TransactionID   uint64    `gorm:"index;column:transaction_id"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	StoreID         uint64    `gorm:"index;column:store_id"`
	MovementTypeID  uint64    `gorm:"column:movement_type_id"` // e.g., IN, OUT, ADJUSTMENT
	Remarks         string    `gorm:"column:remarks"`
}

type StoreStockMovementItem struct {
	BaseModel
	StoreStockMovementID uint64  `gorm:"index;column:store_stock_movement_id"`
	ItemID               uint64  `gorm:"index;column:item_id"`
	Quantity             float64 `gorm:"column:quantity"`
	UnitCost             float64 `gorm:"column:unit_cost"`
	SellingPrice         float64 `gorm:"column:selling_price"`
	BalanceAfter         float64 `gorm:"column:balance_after"`
}

type LivestockBreedingRecord struct {
	BaseModel
	LivestockID         uint64     `gorm:"index;column:livestock_id"`
	BreedingDate        time.Time  `gorm:"column:breeding_date"`
	BreedingType        string     `gorm:"column:breeding_type"` // AI or Natural
	SireID              *uint64    `gorm:"column:sire_id"`
	TechnicianName      string     `gorm:"column:technician_name"`
	PregnancyCheckDate  *time.Time `gorm:"column:pregnancy_check_date"`
	PregnancyStatus     string     `gorm:"column:pregnancy_status"` // e.g., PENDING, POSITIVE, NEGATIVE
	ExpectedCalvingDate *time.Time `gorm:"column:expected_calving_date"`
	ActualCalvingDate   *time.Time `gorm:"column:actual_calving_date"`
	Remarks             string     `gorm:"column:remarks"`
}

func (LivestockBreedingRecord) TableName() string {
	return "livestock_breeding_records"
}

type RecurrentDeduction struct {
	BaseModel
	CustomerID      uint64    `gorm:"column:customer_id"`
	TotalAmount     float64   `gorm:"column:total_amount"`
	PaidAmount      float64   `gorm:"column:paid_amount"` // Using float64 for logic, DB is decimal
	RecurrentAmount float64   `gorm:"column:recurrent_amount"`
	DeductionTypeID uint64    `gorm:"column:deduction_type_id"`
	Reference       string    `gorm:"column:reference"`
	CustomerType    string    `gorm:"column:customer_type"`
	Settled         int       `gorm:"column:settled"`
	PrincipalAmount float64   `gorm:"column:principal_amount"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
}

func (RecurrentDeduction) TableName() string {
	return "recurrent_deductions"
}

type Department struct {
	BaseModel
	DepartmentCode string `gorm:"column:department_code"`
	DepartmentName string `gorm:"column:department_name"`
	Description    string `gorm:"column:description"`
}

func (Department) TableName() string {
	return "departments"
}

type JobPosition struct {
	BaseModel
	Code              string `gorm:"column:code"`
	Name              string `gorm:"column:name"`
	JobDescription    string `gorm:"column:job_description"`
	DepartmentID      uint64 `gorm:"column:department_id"`
	GradeID           string `gorm:"column:grade_id"` // Assuming grade_id is a string as per schema
	NoOfPosts         int    `gorm:"column:no_of_posts"`
	OccupiedPositions int    `gorm:"column:occupied_positions"`
	VaccantPositions  int    `gorm:"column:vaccant_positions"`
}

func (JobPosition) TableName() string {
	return "job_positions"
}

type CashTransaction struct {
	BaseModel
	ReferenceNumber        string    `gorm:"column:reference_number"`
	TransactionDescription string    `gorm:"column:transaction_description"`
	TransactionType        string    `gorm:"column:transaction_type"`
	TransactionDate        time.Time `gorm:"column:transaction_date"`
	PaidBy                 string    `gorm:"column:paid_by"`
	TransactionAmount      float64   `gorm:"column:transaction_amount"`
	CustomerType           string    `gorm:"column:customer_type"`
	CustomerID             uint64    `gorm:"column:customer_id"`
	PaymentModeID          uint64    `gorm:"column:payment_mode_id"`
	PaymentType            string    `gorm:"column:payment_type"`
	TransactionID          uint64    `gorm:"column:transaction_id"`
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

type CreditLimitChangeLog struct {
	BaseModel
	CustomerID   uint64  `gorm:"index;column:customer_id"`
	CustomerType string  `gorm:"column:customer_type"`
	OldLimit     float64 `gorm:"column:old_limit"`
	CreditLimit  float64 `gorm:"column:credit_limit"`
	Action       string  `gorm:"column:action"`
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

type InsuranceDetail struct {
	BaseModel
	EmployeeID       uint64    `gorm:"index;column:employee_id"`
	CompanyName      string    `gorm:"column:company_name"`
	PolicyNo         string    `gorm:"column:policy_no"`
	CommencementDate time.Time `gorm:"column:commencement_date"`
	MaturityDate     time.Time `gorm:"column:maturity_date"`
	SumAssured       float64   `gorm:"column:sum_assured"`
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
	AstraID        uint64 `gorm:"index;column:astra_id"`
	DocumentTypeID uint64 `gorm:"column:document_type_id"`   // Corrected column name to match schema
	DocumentName   string `gorm:"column:document_name"`      // New field for original file name
	Document       string `gorm:"type:text;column:document"` // This will store the URL/path
	Submitted      bool   `gorm:"column:submitted"`
}

type Bank struct {
	BaseModel
	BankName string `gorm:"column:bank_name"`
	BankCode string `gorm:"column:bank_code"`
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

type Guest struct {
	BaseModel
	Name        string `gorm:"column:name"`
	PhoneNumber string `gorm:"uniqueIndex;column:phone_number"`
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
	TrainingID uint64 `gorm:"column:training_id" json:"training_id"`
	Partner    string `gorm:"column:partner;not null" json:"partner"`

	SessionStartTime time.Time `gorm:"column:session_start_time;not null" json:"session_start_time"`
	SessionEndTime   time.Time `gorm:"column:session_end_time;not null" json:"session_end_time"`

	Topic       string `gorm:"column:topic;not null" json:"topic"`
	Description string `gorm:"column:description;not null" json:"description"`
	Trainers    string `gorm:"column:trainers;not null" json:"trainers"`

	Status string `gorm:"column:status" json:"status"`
}

func (TrainingSession) TableName() string {
	return "training_sessions"
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

type ActivityLog struct {
	ID          uint64         `gorm:"column:id;primaryKey"`
	LogName     *string        `gorm:"column:log_name"`
	Description string         `gorm:"column:description"`
	SubjectType *string        `gorm:"column:subject_type"`
	BatchUUID   *string        `gorm:"column:batch_uuid"`
	SubjectID   *uint64        `gorm:"column:subject_id"`
	CauserType  *string        `gorm:"column:causer_type"`
	CauserID    *uint64        `gorm:"column:causer_id"`
	Properties  datatypes.JSON `gorm:"column:properties"`
	CreatedAt   *time.Time     `gorm:"column:created_at"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at"`
	Event       *string        `gorm:"column:event"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ActivityLog) TableName() string {
	return "activity_log"
}
