package dtos

type CreateSharePaymentRequest struct {
	TransactionID   uint64  `json:"transaction_id"`
	MemberID        uint64  `json:"member_id" validate:"required"`
	ShareAccountID  uint64  `json:"share_account_id"`
	AmountPaid      float64 `json:"amount_paid" validate:"required,min=0"`
	ShareUnits      float64 `json:"share_units"`
	PaymentModeID   uint64  `json:"payment_mode_id"`
	ReferenceNo     string  `json:"reference_no"`
	Description     string  `json:"description"`
	Status          string  `json:"status" validate:"omitempty,oneof=PENDING CONFIRMED FAILED REVERSED"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	ApprovedBy      uint64  `json:"approved_by"`
	DateApproved    string  `json:"date_approved"`
}

type UpdateSharePaymentRequest struct {
	TransactionID   uint64  `json:"transaction_id"`
	ShareAccountID  uint64  `json:"share_account_id"`
	AmountPaid      float64 `json:"amount_paid" validate:"required,min=0"`
	ShareUnits      float64 `json:"share_units"`
	PaymentModeID   uint64  `json:"payment_mode_id"`
	ReferenceNo     string  `json:"reference_no"`
	Description     string  `json:"description"`
	Status          string  `json:"status" validate:"required,oneof=PENDING CONFIRMED FAILED REVERSED"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"approved_by"`
	DateApproved    string  `json:"date_approved"`
}
