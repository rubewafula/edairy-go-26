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
