package dtos

type CreateTransporterBankAccountRequest struct {
	TransporterID uint64 `json:"transporter_id" validate:"required"`
	BankID        uint64 `json:"bank_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
}

type UpdateTransporterBankAccountRequest struct {
	BankID        uint64 `json:"bank_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
}
