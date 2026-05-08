package dtos

import "time"

type AssetAssignmentResponse struct {
	ID             uint64    `json:"ID"`
	AssetID        uint64    `json:"AssetID"`
	AssetName      string    `json:"AssetName"`
	AssetCode      string    `json:"AssetCode"`
	AssignedToID   uint64    `json:"AssignedToID"`
	AssignedToName string    `json:"AssignedToName"`
	AssignedAt     time.Time `json:"AssignedAt"`
	ReturnedAt     time.Time `json:"ReturnedAt"`
	ConditionNotes string    `json:"ConditionNotes"`
	Status         string    `json:"Status"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}
