package dtos

type CreateSharePaymentRequest struct {
	TransactionID   uint64  `json:"TransactionID"`
	MemberID        uint64  `json:"MemberID" validate:"required"`
	ShareAccountID  uint64  `json:"ShareAccountID"`
	AmountPaid      float64 `json:"AmountPaid" validate:"required,min=0"`
	ShareUnits      float64 `json:"ShareUnits"`
	PaymentModeID   uint64  `json:"PaymentModeID"`
	ReferenceNo     string  `json:"ReferenceNo"`
	Description     string  `json:"Description"`
	Status          string  `json:"Status" validate:"omitempty,oneof=PENDING CONFIRMED FAILED REVERSED"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"ApprovedBy"`
	DateApproved    string  `json:"DateApproved"`
}

type UpdateSharePaymentRequest struct {
	TransactionID   uint64  `json:"TransactionID"`
	ShareAccountID  uint64  `json:"ShareAccountID"`
	AmountPaid      float64 `json:"AmountPaid" validate:"required,min=0"`
	ShareUnits      float64 `json:"ShareUnits"`
	PaymentModeID   uint64  `json:"PaymentModeID"`
	ReferenceNo     string  `json:"ReferenceNo"`
	Description     string  `json:"Description"`
	Status          string  `json:"Status" validate:"required,oneof=PENDING CONFIRMED FAILED REVERSED"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
	ApprovedBy      uint64  `json:"ApprovedBy"`
	DateApproved    string  `json:"DateApproved"`
}
