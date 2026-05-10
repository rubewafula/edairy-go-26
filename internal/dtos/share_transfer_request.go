package dtos

type CreateShareTransferRequest struct {
	TransactionID   uint64  `json:"TransactionID" validate:"required"`
	FromMemberID    uint64  `json:"FromMemberID" validate:"required"`
	ToMemberID      uint64  `json:"ToMemberID" validate:"required"`
	ShareUnits      float64 `json:"ShareUnits" validate:"required,min=0"`
	TransferAmount  float64 `json:"TransferAmount" validate:"required,min=0"`
	Status          string  `json:"Status" validate:"omitempty,oneof=PENDING APPROVED REJECTED COMPLETED"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"ApprovedBy"`
	DateApproved    string  `json:"DateApproved"`
}

type UpdateShareTransferRequest struct {
	TransactionID   uint64  `json:"TransactionID" validate:"required"`
	ShareUnits      float64 `json:"ShareUnits" validate:"required,min=0"`
	TransferAmount  float64 `json:"TransferAmount" validate:"required,min=0"`
	Status          string  `json:"Status" validate:"required,oneof=PENDING APPROVED REJECTED COMPLETED"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"ApprovedBy"`
	DateApproved    string  `json:"DateApproved"`
}
