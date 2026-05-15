package models

import (
	"time"
)

/*
|--------------------------------------------------------------------------
| Employee / HR Models
|--------------------------------------------------------------------------
*/

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

type EmployeePayrollRelief struct {
	BaseModel
	ReliefID   uint64 `gorm:"index;column:relief_id"`
	EmployeeID uint64 `gorm:"index;column:employee_id"`
	Amount     string `gorm:"column:amount"`
	PayrollID  uint64 `gorm:"index;column:payroll_id"`
}
