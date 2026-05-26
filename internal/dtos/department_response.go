package dtos

import "time"

type DepartmentResponse struct {
	ID             uint64    `json:"id"`
	DepartmentCode string    `json:"department_code"`
	DepartmentName string    `json:"department_name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
