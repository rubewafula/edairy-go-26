package dtos

type CreateDividendDeclarationRequest struct {
	FiscalYear      int     `json:"fiscal_year" validate:"required"`
	Period          int     `json:"period"`
	TotalPool       float64 `json:"total_pool" validate:"required,min=0"`
	RatePerShare    float64 `json:"rate_per_share"`
	CalculationType string  `json:"calculation_type" validate:"required,oneof=PER_SHARE PERCENTAGE FIXED"`
	Status          string  `json:"status" validate:"omitempty,oneof=DRAFT CALCULATED APPROVED PAID CANCELLED"`
	ApprovedBy      uint64  `json:"approved_by"`
	ApprovedAt      string  `json:"approved_at"`
}

type UpdateDividendDeclarationRequest struct {
	FiscalYear      int     `json:"fiscal_year" validate:"required"`
	Period          int     `json:"period"`
	TotalPool       float64 `json:"total_pool" validate:"required,min=0"`
	RatePerShare    float64 `json:"rate_per_share"`
	CalculationType string  `json:"calculation_type" validate:"required,oneof=PER_SHARE PERCENTAGE FIXED"`
	Status          string  `json:"status" validate:"required,oneof=DRAFT CALCULATED APPROVED PAID CANCELLED"`
	ApprovedBy      uint64  `json:"approved_by"`
	ApprovedAt      string  `json:"approved_at"`
}
