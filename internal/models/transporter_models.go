package models

import "time"

type TransporterPayDateRange struct {
	BaseModel
	Name      string    `gorm:"column:name"`
	StartDate time.Time `gorm:"column:start_date"`
	EndDate   time.Time `gorm:"column:end_date"`
	PayMonth  string    `gorm:"column:pay_month"`
	PayYear   string    `gorm:"column:pay_year"`
	Processed bool      `gorm:"column:processed;default:0"`
	Confirmed bool      `gorm:"column:confirmed"`
}

func (TransporterPayDateRange) TableName() string {
	return "transporter_pay_date_ranges"
}

type TransporterPayroll struct {
	BaseModel
	PayDateRangeID  *uint64    `gorm:"column:pay_date_range_id"`
	PayrollMonth    string     `gorm:"column:payroll_month"`
	PayrollYear     string     `gorm:"column:payroll_year"`
	DateOpened      *time.Time `gorm:"column:date_opened"`
	Description     string     `gorm:"column:description"`
	PhysicalPeriod  string     `gorm:"column:physical_period"`
	TotalKilos      float64    `gorm:"column:total_kilos;default:0.00"`
	TotalDeductions float64    `gorm:"column:total_deductions;default:0.00"`
	GrossPay        float64    `gorm:"column:gross_pay;default:0.00"`
	NetPay          float64    `gorm:"column:net_pay;default:0.00"`
	TotalRejects    float64    `gorm:"column:total_rejects;default:0.00"`
	TotalBenefits   float64    `gorm:"column:total_benefits;default:0.00"`
	Status          string     `gorm:"type:enum('processing','draft','confirmed','approved','closed','cancelled','incomplete');default:'draft';column:status"`
	ConfirmedAt     *time.Time `gorm:"column:confirmed_at"`
	ConfirmedBy     *uint64    `gorm:"column:confirmed_by"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	ApprovedBy      *uint64    `gorm:"column:approved_by"`
	PostedAt        *time.Time `gorm:"column:posted_at"`
	PostedBy        *uint64    `gorm:"column:posted_by"`
	PayrollID       *uint64    `gorm:"column:payroll_id"`
}

func (TransporterPayroll) TableName() string {
	return "transporter_payrolls"
}

type TransporterPayslip struct {
	BaseModel
	TransporterID   uint64     `gorm:"column:transporter_id"`
	PayrollID       uint64     `gorm:"column:payroll_id"`
	PayDateRangeID  *uint64    `gorm:"column:pay_date_range_id"`
	PayrollMonth    string     `gorm:"column:payroll_month"`
	PayrollYear     string     `gorm:"column:payroll_year"`
	TotalKilos      float64    `gorm:"column:total_kilos;default:0.00"`
	PhysicalPeriod  string     `gorm:"column:physical_period"`
	GrossPay        float64    `gorm:"column:gross_pay;default:0.00"`
	TotalDeductions float64    `gorm:"column:total_deductions;default:0.00"`
	TotalBenefits   float64    `gorm:"column:total_benefits;default:0.00"`
	NetPay          float64    `gorm:"column:net_pay;default:0.00"`
	TotalRejects    float64    `gorm:"column:total_rejects;default:0.00"`
	Status          string     `gorm:"type:enum('processing','draft','confirmed','approved','closed','cancelled','incomplete');default:'draft';column:status"`
	ConfirmedAt     *time.Time `gorm:"column:confirmed_at"`
	ConfirmedBy     *uint64    `gorm:"column:confirmed_by"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	ApprovedBy      *uint64    `gorm:"column:approved_by"`
	PostedAt        *time.Time `gorm:"column:posted_at"`
	PostedBy        *uint64    `gorm:"column:posted_by"`
}

func (TransporterPayslip) TableName() string {
	return "transporter_payslips"
}

type TransporterBenefit struct {
	BaseModel
	Name        string     `gorm:"column:name"`
	MinQuantity float64    `gorm:"column:min_quantity"`
	Rate        float64    `gorm:"column:rate"`
	RouteID     *uint64    `gorm:"column:route_id"`
	Status      string     `gorm:"column:status;default:'active'"`
	StartDate   *time.Time `gorm:"column:start_date"`
	EndDate     *time.Time `gorm:"column:end_date"`
}

func (TransporterBenefit) TableName() string {
	return "transporter_benefits"
}

type TransporterPayrollBenefit struct {
	BaseModel
	TransporterID        uint64  `gorm:"column:transporter_id"`
	TransporterBenefitID uint64  `gorm:"column:transporter_benefit_id"`
	Amount               float64 `gorm:"column:amount"`
	BenefitYear          string  `gorm:"column:benefit_year"`
	BenefitMonth         string  `gorm:"column:benefit_month"`
	PayrollID            uint64  `gorm:"column:payroll_id"`
}

func (TransporterPayrollBenefit) TableName() string {
	return "transporter_payroll_benefits"
}

type MilkTransporterCost struct {
	BaseModel
	MemberID           uint64  `gorm:"column:member_id"`
	TransporterID      uint64  `gorm:"column:transporter_id"`
	MilkJournalBatchID uint64  `gorm:"column:milk_journal_batch_id"`
	PayrollMonth       string  `gorm:"column:payroll_month"`
	PayrollYear        string  `gorm:"column:payroll_year"`
	PayDateRangeID     uint64  `gorm:"column:pay_date_range_id"`
	Quantity           float64 `gorm:"column:quantity"`
	Rejects            float64 `gorm:"column:rejects"`
	PayrollID          uint64  `gorm:"column:payroll_id"`
}

func (MilkTransporterCost) TableName() string {
	return "milk_transporter_cost"
}

type TransportRate struct {
	BaseModel
	TransportRate float64 `gorm:"type:decimal(10,2);column:rate"`
	RouteID       uint64  `gorm:"column:route_id"`
	MemberID      uint64  `gorm:"column:member_id"`
	TransporterID uint64  `gorm:"column:transporter_id"`
	Status        string  `gorm:"type:varchar(45);default:'active';column:status"`
}

func (TransportRate) TableName() string {
	return "transport_rates"
}

type TransporterPayrollGenerationError struct {
	BaseModel
	TransporterID uint64 `gorm:"column:transporter_id"`
	PayrollID     uint64 `gorm:"column:payroll_id"`
	Error         string `gorm:"column:error;type:text"`
}

func (TransporterPayrollGenerationError) TableName() string {
	return "transporter_payroll_generation_errors"
}

type TransporterBankAccount struct {
	BaseModel
	TransporterID uint64 `gorm:"column:transporter_id"`
	BankID        uint64 `gorm:"column:bank_id"`
	AccountNumber string `gorm:"column:account_number"`
	AccountName   string `gorm:"column:account_name"`
}

func (TransporterBankAccount) TableName() string {
	return "transporter_bank_accounts"
}

type TransporterDriverAssignment struct {
	BaseModel
	TransporterDriverID  uint64     `gorm:"column:transporter_driver_id"`
	TransporterVehicleID uint64     `gorm:"column:transporter_vehicle_id"`
	AssignedFrom         *time.Time `gorm:"column:assigned_from"`
	AssignedTo           *time.Time `gorm:"column:assigned_to"`
	AssignmentType       string     `gorm:"column:assignment_type"`
	Active               bool       `gorm:"column:active"`
	Notes                string     `gorm:"column:notes"`
}

func (TransporterDriverAssignment) TableName() string {
	return "transporter_driver_assignments"
}

type Transporter struct {
	BaseModel
	TransporterNo string                 `gorm:"uniqueIndex;column:transporter_no"`
	Category      string                 `gorm:"type:enum('INDIVIDUAL','COMPANY');column:category"`
	PrimaryPhone  string                 `gorm:"column:primary_phone"`
	EmailAddress  string                 `gorm:"column:email_address"`
	Status        string                 `gorm:"column:status;default:'ACTIVE'"`
	Restricted    bool                   `gorm:"column:restricted;default:0"`
	RouteID       uint64                 `gorm:"column:route_id"`
	Individual    *IndividualTransporter `gorm:"foreignKey:TransporterID"`
	Company       *CompanyTransporter    `gorm:"foreignKey:TransporterID"`
}

func (Transporter) TableName() string {
	return "transporters"
}

type IndividualTransporter struct {
	BaseModel
	TransporterID     uint64     `gorm:"column:transporter_id"`
	FirstName         string     `gorm:"column:first_name"`
	LastName          string     `gorm:"column:last_name"`
	OtherNames        string     `gorm:"column:other_names"`
	Gender            string     `gorm:"column:gender"`
	DateOfBirth       *time.Time `gorm:"column:date_of_birth"`
	NationalIDNo      string     `gorm:"column:national_id_no"`
	KraPin            string     `gorm:"column:kra_pin"`
	MaritalStatus     string     `gorm:"column:marital_status"`
	NextOfKinFullName string     `gorm:"column:next_of_kin_full_name"`
	NextOfKinPhone    string     `gorm:"column:next_of_kin_phone"`
	PassportPhoto     string     `gorm:"column:passport_photo"`
	IDFrontPhoto      string     `gorm:"column:id_front_photo"`
	IDBackPhoto       string     `gorm:"column:id_back_photo"`
}

func (IndividualTransporter) TableName() string {
	return "individual_transporters"
}

type CompanyTransporter struct {
	BaseModel
	TransporterID              uint64       `gorm:"column:transporter_id"`
	CompanyName                string       `gorm:"column:company_name"`
	RegistrationNo             string       `gorm:"column:registration_no"`
	KraPin                     string       `gorm:"column:kra_pin"`
	ContactPersonName          string       `gorm:"column:contact_person_name"`
	ContactPersonPhone         string       `gorm:"column:contact_person_phone"`
	PostalAddress              string       `gorm:"column:postal_address"`
	PostalCode                 string       `gorm:"column:postal_code"`
	Town                       string       `gorm:"column:town"`
	CertificateOfIncorporation string       `gorm:"column:certificate_of_incorporation"`
	Transporter                *Transporter `gorm:"foreignKey:TransporterID"`
}

func (CompanyTransporter) TableName() string {
	return "company_transporters"
}

type TransporterVehicle struct {
	BaseModel
	TransporterID  uint64  `gorm:"column:transporter_id"`
	RouteID        uint64  `gorm:"column:route_id"`
	RegistrationNo string  `gorm:"column:registration_no"`
	VehicleType    string  `gorm:"column:vehicle_type"`
	CapacityLitres float64 `gorm:"column:capacity_litres"`
	Active         bool    `gorm:"column:active"`
}

func (TransporterVehicle) TableName() string {
	return "transporter_vehicles"
}

type TransporterDriver struct {
	BaseModel
	TransporterID            uint64     `gorm:"column:transporter_id"`
	DriverNo                 string     `gorm:"uniqueIndex;column:driver_no"`
	FirstName                string     `gorm:"column:first_name"`
	LastName                 string     `gorm:"column:last_name"`
	OtherNames               string     `gorm:"column:other_names"`
	Gender                   string     `gorm:"column:gender"`
	DateOfBirth              *time.Time `gorm:"column:date_of_birth"`
	NationalIDNo             string     `gorm:"column:national_id_no"`
	KraPin                   string     `gorm:"column:kra_pin"`
	PrimaryPhone             string     `gorm:"column:primary_phone"`
	SecondaryPhone           string     `gorm:"column:secondary_phone"`
	EmailAddress             string     `gorm:"column:email_address"`
	DrivingLicenseNo         string     `gorm:"column:driving_license_no"`
	DrivingLicenseClass      string     `gorm:"column:driving_license_class"`
	DrivingLicenseExpiry     *time.Time `gorm:"column:driving_license_expiry"`
	EmploymentDate           *time.Time `gorm:"column:employment_date"`
	Status                   string     `gorm:"column:status"`
	NextOfKinFullName        string     `gorm:"column:next_of_kin_full_name"`
	NextOfKinPhone           string     `gorm:"column:next_of_kin_phone"`
	PassportPhoto            string     `gorm:"column:passport_photo"`
	DrivingLicenseFrontPhoto string     `gorm:"column:driving_license_front_photo"`
	DrivingLicenseBackPhoto  string     `gorm:"column:driving_license_back_photo"`
}

func (TransporterDriver) TableName() string {
	return "transporter_drivers"
}

type TransporterRouteAssignment struct {
	BaseModel
	TransporterID uint64    `gorm:"column:transporter_id"`
	RouteID       uint64    `gorm:"column:route_id"`
	StartDate     time.Time `gorm:"column:start_date"`
	EndDate       time.Time `gorm:"column:end_date"`
	Active        bool      `gorm:"column:active"`
}

func (TransporterRouteAssignment) TableName() string {
	return "transporter_route_assignments"
}

type TransporterPayrollDeduction struct {
	BaseModel
	TransporterID   uint64    `gorm:"column:transporter_id"`
	DeductionMonth  string    `gorm:"column:deduction_month"`
	FiscalYear      int       `gorm:"column:fiscal_year"`
	DeductionTypeID uint64    `gorm:"column:deduction_type_id"`
	Amount          float64   `gorm:"column:amount"`
	Priority        int       `gorm:"column:priority"`
	Settled         string    `gorm:"column:settled"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	DateCaptured    time.Time `gorm:"column:date_captured"`
	Confirmed       string    `gorm:"column:confirmed;default:'0'"`
	PayrollID       uint64    `gorm:"column:payroll_id"`
	Reference       string    `gorm:"column:reference"`
	SettlementType  string    `gorm:"column:settlement_type"`
}

func (TransporterPayrollDeduction) TableName() string {
	return "transporter_payroll_deductions"
}
