package dtos

import "time"

type AssetAssignmentResponse struct {
	ID             uint64    `json:"id"`
	AssetID        uint64    `json:"asset_id"`
	AssetName      string    `json:"asset_name"`
	AssetCode      string    `json:"asset_code"`
	AssignedToID   uint64    `json:"assigned_to_id"`
	AssignedToName string    `json:"assigned_to_name"`
	AssignedAt     time.Time `json:"assigned_at"`
	ReturnedAt     time.Time `json:"returned_at"`
	ConditionNotes string    `json:"condition_notes"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
