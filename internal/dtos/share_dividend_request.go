package dtos

type CreateShareDividendRequest struct {
	DeclarationID int64   `json:"declaration_id"`
	MemberID      uint64  `json:"member_id" validate:"required"`
	FiscalYear    int     `json:"fiscal_year" validate:"required"`
	Period        int     `json:"period"`
	ShareUnits    float64 `json:"share_units"`
	Status        string  `json:"status" validate:"omitempty,oneof=CALCULATED APPROVED PAID"`
	RatePerShare  float64 `json:"rate_per_share"`
	TaxAmount     float64 `json:"tax_amount"`
	NetAmount     float64 `json:"net_amount" validate:"required"`
	TransactionID int64   `json:"transaction_id"`
}

type UpdateShareDividendRequest struct {
	DeclarationID int64   `json:"declaration_id"`
	MemberID      uint64  `json:"member_id" validate:"required"`
	FiscalYear    int     `json:"fiscal_year" validate:"required"`
	Period        int     `json:"period"`
	ShareUnits    float64 `json:"share_units"`
	Status        string  `json:"status" validate:"required,oneof=CALCULATED APPROVED PAID"`
	RatePerShare  float64 `json:"rate_per_share"`
	TaxAmount     float64 `json:"tax_amount"`
	NetAmount     float64 `json:"net_amount" validate:"required"`
	TransactionID int64   `json:"transaction_id"`
}
