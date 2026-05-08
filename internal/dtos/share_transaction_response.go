package dtos

import "time"

type ShareTransactionResponse struct {
	ID              uint64    `json:"ID"`
	TransactionID   uint64    `json:"TransactionID"`
	ShareAccountID  uint64    `json:"ShareAccountID"`
	MemberID        uint64    `json:"MemberID"`
	MemberNo        string    `json:"MemberNo"`
	MemberFirstName string    `json:"MemberFirstName"`
	MemberLastName  string    `json:"MemberLastName"`
	TransactionType string    `json:"TransactionType"`
	ShareUnits      float64   `json:"ShareUnits"`
	UnitPrice       float64   `json:"UnitPrice"`
	Debit           float64   `json:"Debit"`
	Credit          float64   `json:"Credit"`
	BalanceAfter    float64   `json:"BalanceAfter"`
	TransactionDate time.Time `json:"TransactionDate"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
