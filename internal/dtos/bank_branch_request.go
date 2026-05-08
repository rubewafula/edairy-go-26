package dtos

type CreateBankBranchRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	BankID   uint64 `json:"bank_id" validate:"required"`
	Location string `json:"location" validate:"max=255"`
}

type UpdateBankBranchRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	BankID   uint64 `json:"bank_id" validate:"required"`
	Location string `json:"location" validate:"max=255"`
}
