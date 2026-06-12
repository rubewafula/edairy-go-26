package dtos

import "time"

type CreateEmployeeLeaveDetailRequest struct {
	EmployeeID    uint64 `json:"employee_id" validate:"required"`
	BalanceBF     string `json:"balance_bf"`
	AllocatedDays int    `json:"allocated_days"`
	LeaveTypeID   uint64 `json:"leave_type_id" validate:"required"`
}

type UpdateEmployeeLeaveDetailRequest struct {
	BalanceBF     string `json:"balance_bf"`
	AllocatedDays int    `json:"allocated_days"`
	LeaveTypeID   uint64 `json:"leave_type_id" validate:"required"`
}

type EmployeeLeaveDetailResponse struct {
	ID                  uint64    `json:"id"`
	EmployeeID          uint64    `json:"employee_id"`
	EmployeeName        string    `json:"employee_name"`
	BalanceBF           string    `json:"balance_bf"`
	AllocatedDays       int       `json:"allocated_days"`
	EmployeeLeaveTypeID uint64    `json:"employee_leave_type_id"`
	LeaveTypeCode       string    `json:"leave_type_code"`
	LeaveTypeName       string    `json:"leave_type_name"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
