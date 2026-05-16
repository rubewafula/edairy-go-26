package dtos

import "time"

type DividendDeclarationResponse struct {
	ID                 uint64    `json:"id"`
	FiscalYear         int       `json:"fiscal_year"`
	Period             int       `json:"period"`
	TotalPool          float64   `json:"total_pool"`
	RatePerShare       float64   `json:"rate_per_share"`
	CalculationType    string    `json:"calculation_type"`
	Status             string    `json:"status"`
	ApprovedBy         uint64    `json:"approved_by"`
	ApprovedByUserName string    `json:"approved_by_user_name"` // Joined from users table
	ApprovedAt         time.Time `json:"approved_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
