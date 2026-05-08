package dtos

type CreateMemberBankAccountRequest struct {
	MemberID      uint64 `json:"member_id" validate:"required"`
	BankID        uint64 `json:"bank_id" validate:"required"`
	BankBranchId  uint64 `json:"bank_branch_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	Status        string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateMemberBankAccountRequest struct {
	MemberID      uint64 `json:"member_id" validate:"required"`
	BankID        uint64 `json:"bank_id" validate:"required"`
	BankBranchId  uint64 `json:"bank_branch_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	Status        string `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}
