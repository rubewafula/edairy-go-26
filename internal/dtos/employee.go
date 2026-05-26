package dtos

import "time"

type CreateEmployeeRequest struct {
	UserID            uint64 `json:"user_id"`
	Surname           string `json:"surname" validate:"required"`
	FirstName         string `json:"first_name" validate:"required"`
	MiddleName        string `json:"middle_name"`
	EmployeeNo        string `json:"employee_no"`
	IDNo              string `json:"id_no" validate:"required"`
	KraPin            string `json:"kra_pin"`
	NssfNo            string `json:"nssf_no"`
	NhifNo            string `json:"nhif_no"`
	Gender            string `json:"gender"`
	DateOfBirth       string `json:"date_of_birth" validate:"required"`
	Phone             string `json:"phone_number"`
	Email             string `json:"email_address" validate:"omitempty,email"`
	JobPositionID     uint64 `json:"job_position_id"`
	Status            int    `json:"status"`
	Title             string `json:"title"`
	PassportNo        string `json:"passport_no"`
	Town              string `json:"town"`
	SiteID            uint64 `json:"site_id"`
	MaritalStatus     string `json:"marital_status"`
	Religion          string `json:"religion"`
	Disabled          bool   `json:"disabled"`
	StoreID           uint64 `json:"store_id"`
	PostalAddress     string `json:"postal_address"`
	PostalCode        string `json:"postal_code"`
	BirthCity         string `json:"birth_city"`
	NextOfKinFullName string `json:"next_of_kin_full_name"`
	NextOfKinPhone    string `json:"next_of_kin_phone"`
}

type UpdateEmployeeRequest struct {
	UserID            uint64 `json:"user_id"`
	Surname           string `json:"surname" validate:"required"`
	FirstName         string `json:"first_name" validate:"required"`
	MiddleName        string `json:"middle_name"`
	EmployeeNo        string `json:"employee_no"`
	IDNo              string `json:"id_no" validate:"required"`
	KraPin            string `json:"kra_pin"`
	NssfNo            string `json:"nssf_no"`
	NhifNo            string `json:"nhif_no"`
	Gender            string `json:"gender"`
	DateOfBirth       string `json:"date_of_birth" validate:"required"`
	Phone             string `json:"phone_number"`
	Email             string `json:"email_address" validate:"omitempty,email"`
	JobPositionID     uint64 `json:"job_position_id"`
	Status            int    `json:"status"`
	Title             string `json:"title"`
	PassportNo        string `json:"passport_no"`
	Town              string `json:"town"`
	SiteID            uint64 `json:"site_id"`
	MaritalStatus     string `json:"marital_status"`
	Religion          string `json:"religion"`
	Disabled          bool   `json:"disabled"`
	StoreID           uint64 `json:"store_id"`
	PostalAddress     string `json:"postal_address"`
	PostalCode        string `json:"postal_code"`
	BirthCity         string `json:"birth_city"`
	NextOfKinFullName string `json:"next_of_kin_full_name"`
	NextOfKinPhone    string `json:"next_of_kin_phone"`
}

type EmployeeResponse struct {
	ID              uint64    `json:"id"`
	UserID          uint64    `json:"user_id"`
	Surname         string    `json:"surname"`
	FirstName       string    `json:"first_name"`
	MiddleName      string    `json:"middle_name"`
	EmployeeNo      string    `json:"employee_no"`
	IDNo            string    `json:"id_no"`
	Gender          string    `json:"gender"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Status          int       `json:"status"`
	JobPositionName string    `json:"job_position_name"`
	DepartmentName  string    `json:"department_name"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateEmployeeSalaryRequest struct {
	EmployeeID  uint64  `json:"employee_id" validate:"required"`
	BasicSalary float64 `json:"basic_salary" validate:"required,min=0"`
	Status      string  `json:"status"`
}

type UpdateEmployeeSalaryRequest struct {
	BasicSalary float64 `json:"basic_salary"`
	Status      string  `json:"status"`
}

type EmployeeSalaryResponse struct {
	ID          uint64    `json:"id"`
	EmployeeID  uint64    `json:"employee_id"`
	BasicSalary float64   `json:"basic_salary"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateEmployeeBankAccountRequest struct {
	EmployeeID    uint64 `json:"employee_id" validate:"required"`
	BankID        uint64 `json:"bank_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
}

type UpdateEmployeeBankAccountRequest struct {
	BankID        uint64 `json:"bank_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

type EmployeeBankAccountResponse struct {
	ID            uint64    `json:"id"`
	EmployeeID    uint64    `json:"employee_id"`
	EmployeeName  string    `json:"employee_name"`
	EmployeeNo    string    `json:"employee_no"`
	BankID        uint64    `json:"bank_id"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// EmployeeBenefit DTOs
type CreateEmployeeBenefitRequest struct {
	EmployeeID uint64  `json:"employee_id" validate:"required"`
	BenefitID  uint64  `json:"benefit_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,min=0"`
	Status     string  `json:"status"`
}

type UpdateEmployeeBenefitRequest struct {
	BenefitID uint64  `json:"benefit_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
}

type EmployeeBenefitResponse struct {
	ID          uint64    `json:"id"`
	EmployeeID  uint64    `json:"employee_id"`
	BenefitID   uint64    `json:"benefit_id"`
	BenefitName string    `json:"benefit_name"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EmployeeDocument DTOs
type CreateEmployeeDocumentRequest struct {
	EmployeeID      uint64 `json:"employee_id" validate:"required"`
	DocumentTypeID  uint64 `json:"document_type_id" validate:"required"`
	FileName        string `json:"file_name" validate:"required"`
	FileDescription string `json:"file_description"`
}

type UpdateEmployeeDocumentRequest struct {
	DocumentTypeID  uint64 `json:"document_type_id"`
	FileName        string `json:"file_name"`
	FileDescription string `json:"file_description"`
}

type EmployeeDocumentResponse struct {
	ID              uint64    `json:"id"`
	EmployeeID      uint64    `json:"employee_id"`
	DocumentTypeID  uint64    `json:"document_type_id"`
	DocumentType    string    `json:"document_type"` // Assuming a DocumentType model exists
	FileName        string    `json:"file_name"`
	FileDescription string    `json:"file_description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// EmployeeDeductionType DTOs
type CreateEmployeeDeductionTypeRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsStatutory bool   `json:"is_statutory"`
}

type UpdateEmployeeDeductionTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsStatutory bool   `json:"is_statutory"`
}

type EmployeeDeductionTypeResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsStatutory bool      `json:"is_statutory"`
	CreatedAt   time.Time `json:"created_at"`
}

// EmployeeLeaveApplication DTOs
type CreateEmployeeLeaveApplicationRequest struct {
	EmployeeID  uint64  `json:"employee_id" validate:"required"`
	LeaveTypeID uint64  `json:"leave_type_id" validate:"required"`
	DaysApplied float64 `json:"days_applied" validate:"required,gt=0"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     string  `json:"end_date" validate:"required"`
	ReturnDate  string  `json:"return_date" validate:"required"`
}

type UpdateEmployeeLeaveApplicationRequest struct {
	ApproverID   uint64  `json:"approver_id"`
	DaysApproved float64 `json:"days_approved"`
	Status       string  `json:"status"`
	Approved     bool    `json:"approved"`
}

type EmployeeLeaveApplicationResponse struct {
	ID            uint64    `json:"id"`
	ApplicationNo string    `json:"application_no"`
	EmployeeID    uint64    `json:"employee_id"`
	EmployeeName  string    `json:"employee_name"`
	LeaveTypeID   uint64    `json:"leave_type_id"`
	LeaveType     string    `json:"leave_type"`
	DaysApplied   float64   `json:"days_applied"`
	DaysApproved  float64   `json:"days_approved"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	ReturnDate    time.Time `json:"return_date"`
	ApproverID    uint64    `json:"approver_id"`
	ApproverName  string    `json:"approver_name"`
	Status        string    `json:"status"`
	Approved      bool      `json:"approved"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// EmployeeDependant DTOs
type CreateEmployeeDependantRequest struct {
	EmployeeID   uint64 `json:"employee_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Relationship string `json:"relationship" validate:"required"`
}

type UpdateEmployeeDependantRequest struct {
	Name         string `json:"name" validate:"required"`
	Relationship string `json:"relationship" validate:"required"`
}

type EmployeeDependantResponse struct {
	ID           uint64 `json:"id"`
	EmployeeID   uint64 `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	EmployeeNo   string `json:"employee_no"`
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
}

// EmployeeLeaveType DTOs
type CreateEmployeeLeaveTypeRequest struct {
	Code        string  `json:"code" validate:"required,max=255"`
	Description string  `json:"description"`
	Days        float64 `json:"days" validate:"required,min=0"`
	Gender      string  `json:"gender"`
}

type UpdateEmployeeLeaveTypeRequest struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Days        float64 `json:"days"`
	Gender      string  `json:"gender"`
}

type EmployeeLeaveTypeResponse struct {
	ID          uint64    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Days        float64   `json:"days"`
	Gender      string    `json:"gender"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EmployeePayroll DTOs
type CreateEmployeePayrollRequest struct {
	PayrollMonth    string  `json:"payroll_month" validate:"required"`
	PayrollYear     string  `json:"payroll_year" validate:"required"`
	DateOpened      string  `json:"date_opened" validate:"required,datetime"`
	TotalDeductions float64 `json:"total_deductions"`
	GrossPay        float64 `json:"gross_pay"`
	NetPay          float64 `json:"net_pay"`
	Complete        string  `json:"complete"`
	Confirmed       string  `json:"confirmed"`
	Approved        string  `json:"approved"`
	TotalBenefits   float64 `json:"total_benefits"`
	TotalTax        float64 `json:"total_tax"`
	TotalRelief     float64 `json:"total_relief"`
	Period          string  `json:"period"`
	PaidAt          string  `json:"paid_at"`
}

type UpdateEmployeePayrollRequest struct {
	PayrollMonth    string  `json:"payroll_month"`
	PayrollYear     string  `json:"payroll_year"`
	DateOpened      string  `json:"date_opened" validate:"omitempty,datetime"`
	TotalDeductions float64 `json:"total_deductions"`
	GrossPay        float64 `json:"gross_pay"`
	NetPay          float64 `json:"net_pay"`
	Complete        string  `json:"complete"`
	Confirmed       string  `json:"confirmed"`
	Approved        string  `json:"approved"`
	TotalBenefits   float64 `json:"total_benefits"`
	TotalTax        float64 `json:"total_tax"`
	TotalRelief     float64 `json:"total_relief"`
	Period          string  `json:"period"`
	PaidAt          string  `json:"paid_at" validate:"omitempty,datetime"`
}

type EmployeePayrollResponse struct {
	ID              uint64    `json:"id"`
	PayrollMonth    string    `json:"payroll_month"`
	PayrollYear     string    `json:"payroll_year"`
	DateOpened      time.Time `json:"date_opened"`
	TotalDeductions float64   `json:"total_deductions"`
	GrossPay        float64   `json:"gross_pay"`
	NetPay          float64   `json:"net_pay"`
	Complete        string    `json:"complete"`
	Confirmed       string    `json:"confirmed"`
	Approved        string    `json:"approved"`
	TotalBenefits   float64   `json:"total_benefits"`
	TotalTax        float64   `json:"total_tax"`
	TotalRelief     float64   `json:"total_relief"`
	Period          string    `json:"period"`
	PaidAt          time.Time `json:"paid_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uint64    `json:"created_by"`
	UpdatedBy       uint64    `json:"updated_by"`
}
