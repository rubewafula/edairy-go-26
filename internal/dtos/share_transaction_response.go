package dtos

import "time"

type ShareTransactionResponse struct {
	ID              uint64    `json:"id"`
	TransactionID   uint64    `json:"transaction_id"`
	ShareAccountID  uint64    `json:"share_account_id"`
	MemberID        uint64    `json:"member_id"`
	MemberNo        string    `json:"member_no"`
	MemberFirstName string    `json:"member_first_name"`
	MemberLastName  string    `json:"member_last_name"`
	TransactionType string    `json:"transaction_type"`
	ShareUnits      float64   `json:"share_units"`
	UnitPrice       float64   `json:"unit_price"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	BalanceAfter    float64   `json:"balance_after"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
