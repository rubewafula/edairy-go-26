package dtos

import "time"

type CreateEmployeeQualificationRequest struct {
	EmployeeID    uint64 `json:"employee_id" validate:"required"`
	Qualification string `json:"qualification" validate:"required"`
	Institution   string `json:"institution" validate:"required"`
	StartDate     string `json:"start_date" validate:"required,datetime"`
	EndDate       string `json:"end_date" validate:"required,datetime"`
	Score         string `json:"score"`
}

type UpdateEmployeeQualificationRequest struct {
	Qualification string `json:"qualification" validate:"required"`
	Institution   string `json:"institution" validate:"required"`
	StartDate     string `json:"start_date" validate:"required,datetime"`
	EndDate       string `json:"end_date" validate:"required,datetime"`
	Score         string `json:"score"`
}

type EmployeeQualificationResponse struct {
	ID            uint64    `json:"id"`
	EmployeeID    uint64    `json:"employee_id"`
	Qualification string    `json:"qualification"`
	Institution   string    `json:"institution"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Score         string    `json:"score"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
