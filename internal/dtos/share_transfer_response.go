package dtos

import "time"

type ShareTransferResponse struct {
	ID                  uint64    `json:"id"`
	TransactionID       uint64    `json:"transaction_id"`
	FromMemberID        uint64    `json:"from_member_id"`
	FromMemberNo        string    `json:"from_member_no"`
	FromMemberFirstName string    `json:"from_member_first_name"`
	FromMemberLastName  string    `json:"from_member_last_name"`
	ToMemberID          uint64    `json:"to_member_id"`
	ToMemberNo          string    `json:"to_member_no"`
	ToMemberFirstName   string    `json:"to_member_first_name"`
	ToMemberLastName    string    `json:"to_member_last_name"`
	ShareUnits          float64   `json:"share_units"`
	TransferAmount      float64   `json:"transfer_amount"`
	ReferenceNo         string    `json:"reference_no"` // This will now come from the transactions table
	Status              string    `json:"status"`
	TransactionDate     time.Time `json:"transaction_date"`
	ApprovedBy          uint64    `json:"approved_by"`
	ApprovedByUserName  string    `json:"approved_by_user_name"`
	DateApproved        time.Time `json:"date_approved"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
