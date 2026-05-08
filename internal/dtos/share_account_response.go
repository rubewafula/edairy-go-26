package dtos

import "time"

type ShareAccountResponse struct {
	ID            uint64    `json:"ID"`
	MemberID      uint64    `json:"MemberID"`
	MemberNo      string    `json:"MemberNo"`
	FirstName     string    `json:"FirstName"`
	LastName      string    `json:"LastName"`
	ShareTypeID   uint64    `json:"ShareTypeID"`
	ShareCode     string    `json:"ShareCode"`
	ShareTypeName string    `json:"ShareTypeName"`
	Description   string    `json:"Description"`
	Status        string    `json:"Status"`
	OpenedAt      time.Time `json:"OpenedAt"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}
