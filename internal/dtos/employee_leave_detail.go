package dtos

import "time"

type CreateEmployeeLeaveDetailRequest struct {
	EmployeeID    uint64 `json:"employee_id" validate:"required"`
	BalanceBF     string `json:"balance_bf"`
	AllocatedDays int    `json:"allocated_days"`
}

type UpdateEmployeeLeaveDetailRequest struct {
	BalanceBF     string `json:"balance_bf"`
	AllocatedDays int    `json:"allocated_days"`
}

type EmployeeLeaveDetailResponse struct {
	ID            uint64    `json:"id"`
	EmployeeID    uint64    `json:"employee_id"`
	BalanceBF     string    `json:"balance_bf"`
	AllocatedDays int       `json:"allocated_days"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
