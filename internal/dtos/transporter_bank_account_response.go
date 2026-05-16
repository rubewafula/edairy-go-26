package dtos

import "time"

type TransporterBankAccountResponse struct {
	ID            uint64    `json:"id"`
	TransporterID uint64    `json:"transporter_id"`
	TransporterNo string    `json:"transporter_no"`
	BankID        uint64    `json:"bank_id"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
