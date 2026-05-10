package dtos

import "time"

type TransporterBankAccountResponse struct {
	ID            uint64    `json:"ID"`
	TransporterID uint64    `json:"TransporterID"`
	TransporterNo string    `json:"TransporterNo"`
	BankID        uint64    `json:"BankID"`
	BankName      string    `json:"BankName"`
	AccountNumber string    `json:"AccountNumber"`
	AccountName   string    `json:"AccountName"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
