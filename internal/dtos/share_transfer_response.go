package dtos

import "time"

type ShareTransferResponse struct {
	ID                  uint64    `json:"ID"`
	TransactionID       uint64    `json:"TransactionID"`
	FromMemberID        uint64    `json:"FromMemberID"`
	FromMemberNo        string    `json:"FromMemberNo"`
	FromMemberFirstName string    `json:"FromMemberFirstName"`
	FromMemberLastName  string    `json:"FromMemberLastName"`
	ToMemberID          uint64    `json:"ToMemberID"`
	ToMemberNo          string    `json:"ToMemberNo"`
	ToMemberFirstName   string    `json:"ToMemberFirstName"`
	ToMemberLastName    string    `json:"ToMemberLastName"`
	ShareUnits          float64   `json:"ShareUnits"`
	TransferAmount      float64   `json:"TransferAmount"`
	ReferenceNo         string    `json:"ReferenceNo"` // This will now come from the transactions table
	Status              string    `json:"Status"`
	TransactionDate     time.Time `json:"TransactionDate"`
	ApprovedBy          uint64    `json:"ApprovedBy"`
	ApprovedByUserName  string    `json:"ApprovedByUserName"`
	DateApproved        time.Time `json:"DateApproved"`
	CreatedAt           time.Time `json:"CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt"`
}
