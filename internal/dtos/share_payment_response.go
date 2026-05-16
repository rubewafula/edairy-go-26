package dtos

import "time"

type SharePaymentResponse struct {
	ID                 uint64    `json:"id"`
	TransactionID      uint64    `json:"transaction_id"`
	MemberID           uint64    `json:"member_id"`
	MemberNo           string    `json:"member_no"`
	MemberFirstName    string    `json:"member_first_name"`
	MemberLastName     string    `json:"member_last_name"`
	ShareAccountID     uint64    `json:"share_account_id"`
	AmountPaid         float64   `json:"amount_paid"`
	ShareUnits         float64   `json:"share_units"`
	PaymentModeID      uint64    `json:"payment_mode_id"`
	PaymentModeName    string    `json:"payment_mode_name"`
	ReferenceNo        string    `json:"reference_no"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	TransactionDate    time.Time `json:"transaction_date"`
	ApprovedBy         uint64    `json:"approved_by"`
	ApprovedByUserName string    `json:"approved_by_user_name"`
	DateApproved       time.Time `json:"date_approved"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
