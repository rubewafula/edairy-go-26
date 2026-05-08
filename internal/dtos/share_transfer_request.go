package dtos

type CreateShareTransferRequest struct {
	TransactionID   uint64  `json:"transaction_id" validate:"required"`
	FromMemberID    uint64  `json:"from_member_id" validate:"required"`
	ToMemberID      uint64  `json:"to_member_id" validate:"required"`
	ShareUnits      float64 `json:"share_units" validate:"required,min=0"`
	TransferAmount  float64 `json:"transfer_amount" validate:"required,min=0"`
	Status          string  `json:"status" validate:"omitempty,oneof=PENDING APPROVED REJECTED COMPLETED"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"approved_by"`
	DateApproved    string  `json:"date_approved"`
}

type UpdateShareTransferRequest struct {
	TransactionID   uint64  `json:"transaction_id" validate:"required"`
	ShareUnits      float64 `json:"share_units" validate:"required,min=0"`
	TransferAmount  float64 `json:"transfer_amount" validate:"required,min=0"`
	Status          string  `json:"status" validate:"required,oneof=PENDING APPROVED REJECTED COMPLETED"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"approved_by"`
	DateApproved    string  `json:"date_approved"`
}
