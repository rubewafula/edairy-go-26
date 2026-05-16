package dtos

import "time"

type CreateEmployeeProfessionalTitleRequest struct {
	Code  string `json:"code" validate:"required"`
	Title string `json:"title" validate:"required"`
}

type UpdateEmployeeProfessionalTitleRequest struct {
	Title string `json:"title" validate:"required"`
}

type EmployeeProfessionalTitleResponse struct {
	ID         uint64    `json:"id"`
	EmployeeID uint64    `json:"employee_id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  uint64    `json:"created_by"`
	UpdatedBy  uint64    `json:"updated_by"`
}
