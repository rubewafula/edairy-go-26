package models

import "time"

// Member+ Lending
type Member struct {
	BaseModel
	MemberNo          string     `gorm:"uniqueIndex;column:member_no"`
	IDNo              string     `gorm:"uniqueIndex;column:id_no"`
	FirstName         string     `gorm:"column:first_name"`
	LastName          string     `gorm:"column:last_name"`
	OtherNames        string     `gorm:"column:other_names"`
	Gender            string     `gorm:"column:gender"`
	DateOfBirth       time.Time  `gorm:"column:dob"`
	PrimaryPhone      string     `gorm:"column:primary_phone"`
	SecondaryPhone    string     `gorm:"column:secondary_phone"`
	Email             string     `gorm:"column:email"`
	TaxNumber         string     `gorm:"column:tax_number"`
	MaritalStatus     string     `gorm:"column:marital_status"`
	Status            string     `gorm:"column:status"`
	RouteID           uint64     `gorm:"column:route_id"`
	MemberTypeID      uint64     `gorm:"column:member_type_id"`
	NumberOfCows      int        `gorm:"column:number_of_cows"`
	NextOfKinFullName string     `gorm:"column:next_of_kin_full_name"`
	NextOfKinPhone    string     `gorm:"column:next_of_kin_phone"`
	DateRegistered    time.Time  `gorm:"column:date_registered"`
	PassportPhoto     string     `gorm:"column:passport_photo"`
	IdFrontPhoto      string     `gorm:"column:id_front_photo"`
	IdBackPhoto       string     `gorm:"column:id_back_photo"`
	UpdatedBy         string     `gorm:"column:updated_by"`
	Downloaded        string     `gorm:"column:downloaded"`
	BirthCity         string     `gorm:"column:birth_city"`
	IdDateOfIssue     *time.Time `gorm:"column:id_date_of_issue"`
	Title             string     `gorm:"column:title"`
	CashoutEnrolled   bool       `gorm:"column:cashout_enrolled"`
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

type MemberBankAccount struct {
	BaseModel
	MemberID      uint64 `gorm:"index;column:member_id"`
	BankID        uint64 `gorm:"column:bank_id"`
	BankBranchId  uint64 `gorm:"column:bank_branch_id"`
	AccountNumber string `gorm:"column:account_number"`
	AccountName   string `gorm:"column:account_name"`
	Status        string `gorm:"column:status"`
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

type MemberNextOfKin struct {
	BaseModel
	MemberID               uint64  `gorm:"column:member_id;not null;index"`
	FullName               string  `gorm:"column:full_name;not null;size:255"`
	Relationship           *string `gorm:"column:relationship;size:100"`
	PhoneNumber            *string `gorm:"column:phone_number;size:30"`
	AlternativePhoneNumber *string `gorm:"column:alternative_phone_number;size:30"`
	EmailAddress           *string `gorm:"column:email_address;size:255"`
	NationalIDNo           *string `gorm:"column:national_id_no;size:100"`
	PostalAddress          *string `gorm:"column:postal_address;size:255"`
	PhysicalAddress        *string `gorm:"column:physical_address;size:255"`
	Occupation             *string `gorm:"column:occupation;size:255"`
	IsPrimary              bool    `gorm:"column:is_primary;not null;default:0"`
	Status                 bool    `gorm:"column:status;not null;default:1"` // TINYINT(1) maps to bool
	Remarks                *string `gorm:"column:remarks;type:text"`
}

func (MemberNextOfKin) TableName() string {
	return "member_next_of_kins"
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

type MemberMpesaWithdrawal struct {
	BaseModel
	WithdrawalID uint64  `gorm:"index;column:withdrawal_id"`
	LoanID       uint64  `gorm:"index;column:loan_id"`
	WalletID     uint64  `gorm:"index;column:wallet_id"`
	Amount       float64 `gorm:"column:amount"`
	Status       string  `gorm:"column:status"`
}

type MemberPayroll struct {
	BaseModel
	DateOpened      *time.Time `gorm:"column:date_opened"`
	Description     string     `gorm:"column:description"`
	Status          string     `gorm:"type:enum('draft','confirmed','approved','closed','cancelled');default:'draft';column:status"`
	PostedAt        *time.Time `gorm:"column:posted_at"`
	PostedBy        *uint64    `gorm:"column:posted_by"`
	ConfirmedAt     *time.Time `gorm:"column:confirmed_at"`
	ConfirmedBy     *uint64    `gorm:"column:confirmed_by"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	ApprovedBy      *uint64    `gorm:"column:approved_by"`
	PayDateRangeID  *uint64    `gorm:"column:pay_date_range_id"`
	DateClosed      *time.Time `gorm:"column:date_closed"`
	Closed          string     `gorm:"column:closed"`
	PhysicalPeriod  string     `gorm:"column:physical_period"`
	GrossKilos      float64    `gorm:"column:gross_kilos;default:0.00"`
	RejectKilos     float64    `gorm:"column:reject_kilos;default:0.00"`
	NetKilos        float64    `gorm:"column:net_kilos;default:0.00"`
	TotalKilos      string     `gorm:"column:total_kilos"` // Varchar in schema
	TotalDeductions float64    `gorm:"column:total_deductions;default:0.00"`
	GrossPay        float64    `gorm:"column:gross_pay;default:0.00"`
	NetPay          float64    `gorm:"column:net_pay;default:0.00"`
	TotalRejects    float64    `gorm:"column:total_rejects;default:0.00"`
	TransportCost   float64    `gorm:"column:transport_cost;default:0.00"`
}

func (MemberPayroll) TableName() string {
	return "member_payrolls"
}

type MemberPayslip struct {
	BaseModel
	MemberID        uint64     `gorm:"column:member_id"`
	PayRateID       uint64     `gorm:"column:pay_rate_id"`
	PayrollID       uint64     `gorm:"column:payroll_id"`
	DateOpened      *time.Time `gorm:"column:date_opened"`
	Description     string     `gorm:"column:description"`
	Status          string     `gorm:"type:enum('processing','draft','confirmed','approved','closed','cancelled','incomplete');default:'draft';column:status"`
	PostedAt        *time.Time `gorm:"column:posted_at"`
	PostedBy        *uint64    `gorm:"column:posted_by"`
	ConfirmedAt     *time.Time `gorm:"column:confirmed_at"`
	ConfirmedBy     *uint64    `gorm:"column:confirmed_by"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	ApprovedBy      *uint64    `gorm:"column:approved_by"`
	PayDateRangeID  *uint64    `gorm:"column:pay_date_range_id"`
	DateClosed      *time.Time `gorm:"column:date_closed"`
	Closed          bool       `gorm:"column:closed"`
	PhysicalPeriod  string     `gorm:"column:physical_period"`
	GrossKilos      float64    `gorm:"column:gross_kilos"`
	RejectKilos     float64    `gorm:"column:reject_kilos"`
	NetKilos        float64    `gorm:"column:net_kilos"`
	TotalKilos      float64    `gorm:"column:total_kilos"`
	TotalDeductions float64    `gorm:"column:total_deductions"`
	GrossPay        float64    `gorm:"column:gross_pay"`
	NetPay          float64    `gorm:"column:net_pay"`
	TotalRejects    float64    `gorm:"column:total_rejects"`
	TransportCost   float64    `gorm:"column:transport_cost"`
}

func (MemberPayslip) TableName() string {
	return "member_payslips"
}

type MemberPayrollDeduction struct {
	BaseModel
	MemberID        int64     `gorm:"column:member_id"`
	DeductionMonth  string    `gorm:"column:deduction_month"`
	FiscalYear      int       `gorm:"column:fiscal_year"`
	DeductionTypeID uint64    `gorm:"column:deduction_type_id"`
	Amount          float64   `gorm:"column:amount"` // Varchar in schema
	Priority        int       `gorm:"column:priority"`
	Settled         string    `gorm:"column:settled"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	DateCaptured    time.Time `gorm:"column:date_captured"`
	Confirmed       string    `gorm:"column:confirmed;default:'0'"`
	PayrollID       uint64    `gorm:"column:payroll_id"`
	Reference       string    `gorm:"column:reference"`
	SettlementType  string    `gorm:"column:settlement_type"`
}

func (MemberPayrollDeduction) TableName() string {
	return "member_payroll_deductions"
}

type MemberPayrollGenerationError struct {
	BaseModel
	MemberID  uint64 `gorm:"column:member_id"`
	PayrollID uint64 `gorm:"column:payroll_id"`
	Error     string `gorm:"column:error;type:text"`
}

func (MemberPayrollGenerationError) TableName() string {
	return "member_payroll_generation_errors"
}

type MemberPayDateRange struct {
	BaseModel
	Name      string    `gorm:"column:name"`
	StartDate time.Time `gorm:"column:start_date"`
	EndDate   time.Time `gorm:"column:end_date"`
	Processed bool      `gorm:"column:processed;default:0"`
	Confirmed int       `gorm:"column:confirmed"`
}

func (MemberPayDateRange) TableName() string {
	return "member_pay_date_ranges"
}
