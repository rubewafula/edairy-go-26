package dtos

import "time"

type ShareDividendResponse struct {
	ID            uint64    `json:"id"`
	DeclarationID int64     `json:"declaration_id"`
	MemberID      uint64    `json:"member_id"`
	MemberNo      string    `json:"member_no"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	FiscalYear    int       `json:"fiscal_year"`
	Period        int       `json:"period"`
	ShareUnits    float64   `json:"share_units"`
	Status        string    `json:"status"`
	RatePerShare  float64   `json:"rate_per_share"`
	TaxAmount     float64   `json:"tax_amount"`
	NetAmount     float64   `json:"net_amount"`
	TransactionID int64     `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
