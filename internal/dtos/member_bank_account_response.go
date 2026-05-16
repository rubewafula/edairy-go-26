package dtos

import "time"

type MemberBankAccountResponse struct {
	ID            uint64    `json:"id"`
	MemberID      uint64    `json:"member_id"`
	MemberNo      string    `json:"member_no"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BankID        uint64    `json:"bank_id"`
	BankName      string    `json:"bank_name"`
	BankBranchId  uint64    `json:"bank_branch_id"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
