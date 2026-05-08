package dtos

import "time"

type ShareDividendResponse struct {
	ID            uint64    `json:"ID"`
	DeclarationID int64     `json:"DeclarationID"`
	MemberID      uint64    `json:"MemberID"`
	MemberNo      string    `json:"MemberNo"`
	FirstName     string    `json:"FirstName"`
	LastName      string    `json:"LastName"`
	FiscalYear    int       `json:"FiscalYear"`
	Period        int       `json:"Period"`
	ShareUnits    float64   `json:"ShareUnits"`
	Status        string    `json:"Status"`
	RatePerShare  float64   `json:"RatePerShare"`
	TaxAmount     float64   `json:"TaxAmount"`
	NetAmount     float64   `json:"NetAmount"`
	TransactionID int64     `json:"TransactionID"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
