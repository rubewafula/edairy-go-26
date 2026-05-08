package dtos

import "time"

type MemberBankAccountResponse struct {
	ID            uint64    `json:"ID"`
	MemberID      uint64    `json:"MemberID"`
	MemberNo      string    `json:"MemberNo"`
	FirstName     string    `json:"FirstName"`
	LastName      string    `json:"LastName"`
	BankID        uint64    `json:"BankID"`
	BankName      string    `json:"BankName"`
	BankBranchId  uint64    `json:"BankBranchId"`
	AccountNumber string    `json:"AccountNumber"`
	AccountName   string    `json:"AccountName"`
	Status        string    `json:"Status"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
