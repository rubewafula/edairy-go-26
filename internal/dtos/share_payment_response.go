package dtos

import "time"

type SharePaymentResponse struct {
	ID                 uint64    `json:"ID"`
	TransactionID      uint64    `json:"TransactionID"`
	MemberID           uint64    `json:"MemberID"`
	MemberNo           string    `json:"MemberNo"`
	MemberFirstName    string    `json:"MemberFirstName"`
	MemberLastName     string    `json:"MemberLastName"`
	ShareAccountID     uint64    `json:"ShareAccountID"`
	AmountPaid         float64   `json:"AmountPaid"`
	ShareUnits         float64   `json:"ShareUnits"`
	PaymentModeID      uint64    `json:"PaymentModeID"`
	PaymentModeName    string    `json:"PaymentModeName"`
	ReferenceNo        string    `json:"ReferenceNo"`
	Description        string    `json:"Description"`
	Status             string    `json:"Status"`
	TransactionDate    time.Time `json:"TransactionDate"`
	ApprovedBy         uint64    `json:"ApprovedBy"`
	ApprovedByUserName string    `json:"ApprovedByUserName"`
	DateApproved       time.Time `json:"DateApproved"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
}
