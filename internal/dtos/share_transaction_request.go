package dtos

type CreateShareTransactionRequest struct {
	TransactionID   uint64  `json:"TransactionID" validate:"required"`
	ShareAccountID  uint64  `json:"ShareAccountID" validate:"required"`
	MemberID        uint64  `json:"MemberID" validate:"required"`
	TransactionType string  `json:"TransactionType" validate:"required,oneof=PURCHASE TRANSFER_IN TRANSFER_OUT ADJUSTMENT DIVIDEND_REINVEST REFUND WITHDRAWAL"`
	ShareUnits      float64 `json:"ShareUnits"`
	UnitPrice       float64 `json:"UnitPrice"`
	Debit           float64 `json:"Debit"`
	Credit          float64 `json:"Credit"`
	BalanceAfter    float64 `json:"BalanceAfter"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
}

type UpdateShareTransactionRequest struct {
	TransactionID   uint64  `json:"TransactionID" validate:"required"`
	ShareAccountID  uint64  `json:"ShareAccountID" validate:"required"`
	TransactionType string  `json:"TransactionType" validate:"required,oneof=PURCHASE TRANSFER_IN TRANSFER_OUT ADJUSTMENT DIVIDEND_REINVEST REFUND WITHDRAWAL"`
	ShareUnits      float64 `json:"ShareUnits"`
	UnitPrice       float64 `json:"UnitPrice"`
	Debit           float64 `json:"Debit"`
	Credit          float64 `json:"Credit"`
	BalanceAfter    float64 `json:"BalanceAfter"`
	TransactionDate string  `json:"TransactionDate" validate:"required,datetime"`
}
