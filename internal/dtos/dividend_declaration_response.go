package dtos

import "time"

type DividendDeclarationResponse struct {
	ID                 uint64    `json:"ID"`
	FiscalYear         int       `json:"FiscalYear"`
	Period             int       `json:"Period"`
	TotalPool          float64   `json:"TotalPool"`
	RatePerShare       float64   `json:"RatePerShare"`
	CalculationType    string    `json:"CalculationType"`
	Status             string    `json:"Status"`
	ApprovedBy         uint64    `json:"ApprovedBy"`
	ApprovedByUserName string    `json:"ApprovedByUserName"` // Joined from users table
	ApprovedAt         time.Time `json:"ApprovedAt"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
}
