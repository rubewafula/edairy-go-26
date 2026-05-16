package dtos

import "time"

type CreateEmployeeExitDetailRequest struct {
	EmployeeID      uint64 `json:"employee_id" validate:"required"`
	ContractType    string `json:"contract_type" validate:"required"`
	ContractEndDate string `json:"contract_end_date" validate:"required,datetime"`
	DateOfLeaving   string `json:"date_of_leaving" validate:"required,datetime"`
	ExitCategory    string `json:"exit_category" validate:"required"`
	Reasons         string `json:"reasons" validate:"required"`
}

type UpdateEmployeeExitDetailRequest struct {
	ContractType    string `json:"contract_type" validate:"required"`
	ContractEndDate string `json:"contract_end_date" validate:"required,datetime"`
	DateOfLeaving   string `json:"date_of_leaving" validate:"required,datetime"`
	ExitCategory    string `json:"exit_category" validate:"required"`
	Reasons         string `json:"reasons" validate:"required"`
}

type EmployeeExitDetailResponse struct {
	ID              uint64    `json:"id"`
	EmployeeID      uint64    `json:"employee_id"`
	ContractType    string    `json:"contract_type"`
	ContractEndDate time.Time `json:"contract_end_date"`
	DateOfLeaving   time.Time `json:"date_of_leaving"`
	ExitCategory    string    `json:"exit_category"`
	Reasons         string    `json:"reasons"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
