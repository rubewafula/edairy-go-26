package dtos

import "time"

// BankResponse represents a bank entity in JSON responses.
type BankResponse struct {
	ID        uint64    `json:"id"`
	BankName  string    `json:"bank_name"`
	BankCode  string    `json:"bank_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BankBranchResponse represents a bank branch in JSON responses.
type BankBranchResponse struct {
	ID         uint64    `json:"id"`
	BankID     uint64    `json:"bank_id"`
	BankName   string    `json:"bank_name"`
	BranchName string    `json:"branch_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// OrganizationBankResponse represents a bank registered to the organization.
type OrganizationBankResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOrganizationBankRequest defines the request body for creating an organization bank.
type CreateOrganizationBankRequest struct {
	Name string `json:"name" validate:"required"`
}

// UpdateOrganizationBankRequest defines the request body for updating an organization bank.
type UpdateOrganizationBankRequest struct {
	Name string `json:"name" validate:"required"`
}

// MemberBankAccountResponse represents a member's bank account details.
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

// CreateMemberBankAccountRequest defines the request body for adding a member's bank account.
type CreateMemberBankAccountRequest struct {
	MemberID      uint64 `json:"member_id" validate:"required"`
	BankID        uint64 `json:"bank_id" validate:"required"`
	BankBranchId  uint64 `json:"bank_branch_id"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	Status        string `json:"status"`
}

// UpdateMemberBankAccountRequest defines the request body for updating a member's bank account.
type UpdateMemberBankAccountRequest struct {
	MemberID      uint64 `json:"member_id"`
	BankID        uint64 `json:"bank_id"`
	BankBranchId  uint64 `json:"bank_branch_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Status        string `json:"status"`
}
