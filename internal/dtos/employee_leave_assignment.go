package dtos

import "time"

type CreateEmployeeLeaveAssignmentRequest struct {
	EmployeeID         uint64 `json:"employee_id" validate:"required"`
	LeaveApplicationID uint64 `json:"leave_application_id" validate:"required"`
	RelieverID         uint64 `json:"reliever_id" validate:"required"`
}

type UpdateEmployeeLeaveAssignmentRequest struct {
	EmployeeID         uint64 `json:"employee_id" validate:"required"`
	LeaveApplicationID uint64 `json:"leave_application_id" validate:"required"`
	RelieverID         uint64 `json:"reliever_id" validate:"required"`
}

type EmployeeLeaveAssignmentResponse struct {
	ID                 uint64    `json:"id"`
	EmployeeID         uint64    `json:"employee_id"`
	EmployeeName       string    `json:"employee_name"`
	LeaveApplicationID uint64    `json:"leave_application_id"`
	ApplicationNo      string    `json:"application_no"`
	RelieverID         uint64    `json:"reliever_id"`
	RelieverName       string    `json:"reliever_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
