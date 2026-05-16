package dtos

import "time"

type CreateEmployeeContractRequest struct {
	EmployeeID      uint64 `json:"employee_id" validate:"required"`
	ContractType    string `json:"contract_type" validate:"required"`
	ContractEndDate string `json:"contract_end_date" validate:"required,datetime"`
	NoticePeriod    string `json:"notice_period"`
	RetirementDate  string `json:"retirement_date" validate:"required,datetime"`
}

type UpdateEmployeeContractRequest struct {
	ContractType    string `json:"contract_type" validate:"required"`
	ContractEndDate string `json:"contract_end_date" validate:"required,datetime"`
	NoticePeriod    string `json:"notice_period"`
	RetirementDate  string `json:"retirement_date" validate:"required,datetime"`
}

type EmployeeContractResponse struct {
	ID              uint64    `json:"id"`
	EmployeeID      uint64    `json:"employee_id"`
	ContractType    string    `json:"contract_type"`
	ContractEndDate time.Time `json:"contract_end_date"`
	NoticePeriod    string    `json:"notice_period"`
	RetirementDate  time.Time `json:"retirement_date"`
	CreatedAt       time.Time `json:"created_at"`
}
