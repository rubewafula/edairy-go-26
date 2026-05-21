package dtos

import (
	"time"
)

type CreateEmployeeReliefRequest struct {
	EmployeeID uint64 `json:"employee_id" validate:"required"`
	ReliefID   uint64 `json:"relief_id" validate:"required"`
	Status     string `json:"status"`
}

type UpdateEmployeeReliefRequest struct {
	EmployeeID uint64 `json:"employee_id"`
	ReliefID   uint64 `json:"relief_id"`
	Status     string `json:"status"`
}

type EmployeeReliefResponse struct {
	ID           uint64    `json:"id"`
	EmployeeID   uint64    `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	ReliefID     uint64    `json:"relief_id"`
	ReliefName   string    `json:"relief_name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
