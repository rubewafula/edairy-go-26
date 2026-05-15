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
	DateOfBirth       string `json:"date_of_birth" validate:"required,datetime"`
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
	DateOfBirth       string `json:"date_of_birth" validate:"required,datetime"`
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
	ID          uint64    `json:"ID"`
	UserID      uint64    `json:"UserID"`
	Surname     string    `json:"Surname"`
	FirstName   string    `json:"FirstName"`
	MiddleName  string    `json:"MiddleName"`
	EmployeeNo  string    `json:"EmployeeNo"`
	IDNo        string    `json:"IDNo"`
	Gender      string    `json:"Gender"`
	DateOfBirth time.Time `json:"DateOfBirth"`
	Phone       string    `json:"Phone"`
	Email       string    `json:"Email"`
	Status      int       `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type CreateEmployeeSalaryRequest struct {
	EmployeeID  uint64  `json:"employee_id" validate:"required"`
	BasicSalary float64 `json:"basic_salary" validate:"required,min=0"`
	Status      string  `json:"status"`
}

type CreateEmployeeBankAccountRequest struct {
	EmployeeID    uint64 `json:"employee_id" validate:"required"`
	BankID        uint64 `json:"bank_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
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
	ID          uint64    `json:"ID"`
	EmployeeID  uint64    `json:"EmployeeID"`
	BenefitID   uint64    `json:"BenefitID"`
	BenefitName string    `json:"BenefitName"`
	Amount      float64   `json:"Amount"`
	Status      string    `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
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
	ID              uint64    `json:"ID"`
	EmployeeID      uint64    `json:"EmployeeID"`
	DocumentTypeID  uint64    `json:"DocumentTypeID"`
	DocumentType    string    `json:"DocumentType"` // Assuming a DocumentType model exists
	FileName        string    `json:"FileName"`
	FileDescription string    `json:"FileDescription"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
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
	ID          uint64    `json:"ID"`
	Code        string    `json:"Code"`
	Description string    `json:"Description"`
	Days        float64   `json:"Days"`
	Gender      string    `json:"Gender"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
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
	PayrollMonth string `json:"payroll_month"`
	PayrollYear  string `json:"payroll_year"`
	// ... other fields that can be updated
}

type EmployeePayrollResponse struct {
	ID              uint64    `json:"ID"`
	PayrollMonth    string    `json:"PayrollMonth"`
	PayrollYear     string    `json:"PayrollYear"`
	DateOpened      time.Time `json:"DateOpened"`
	TotalDeductions float64   `json:"TotalDeductions"`
	GrossPay        float64   `json:"GrossPay"`
	NetPay          float64   `json:"NetPay"`
	Complete        string    `json:"Complete"`
	Confirmed       string    `json:"Confirmed"`
	Approved        string    `json:"Approved"`
	TotalBenefits   float64   `json:"TotalBenefits"`
	TotalTax        float64   `json:"TotalTax"`
	TotalRelief     float64   `json:"TotalRelief"`
	Period          string    `json:"Period"`
	PaidAt          time.Time `json:"PaidAt"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
