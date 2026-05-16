package dtos

type CreateShareTransactionRequest struct {
	TransactionID   uint64  `json:"transaction_id" validate:"required"`
	ShareAccountID  uint64  `json:"share_account_id" validate:"required"`
	MemberID        uint64  `json:"member_id" validate:"required"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=PURCHASE TRANSFER_IN TRANSFER_OUT ADJUSTMENT DIVIDEND_REINVEST REFUND WITHDRAWAL"`
	ShareUnits      float64 `json:"share_units"`
	UnitPrice       float64 `json:"unit_price"`
	Debit           float64 `json:"debit"`
	Credit          float64 `json:"credit"`
	BalanceAfter    float64 `json:"balance_after"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
}

type UpdateShareTransactionRequest struct {
	TransactionID   uint64  `json:"transaction_id" validate:"required"`
	ShareAccountID  uint64  `json:"share_account_id" validate:"required"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=PURCHASE TRANSFER_IN TRANSFER_OUT ADJUSTMENT DIVIDEND_REINVEST REFUND WITHDRAWAL"`
	ShareUnits      float64 `json:"share_units"`
	UnitPrice       float64 `json:"unit_price"`
	Debit           float64 `json:"debit"`
	Credit          float64 `json:"credit"`
	BalanceAfter    float64 `json:"balance_after"`
	TransactionDate string  `json:"transaction_date" validate:"required,datetime"`
}
