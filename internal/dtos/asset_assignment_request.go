package dtos

type CreateAssetAssignmentRequest struct {
	AssetID        uint64 `json:"asset_id" validate:"required"`
	AssignedToID   uint64 `json:"assigned_to_id" validate:"required"`
	AssignedAt     string `json:"assigned_at" validate:"required"`
	ReturnedAt     string `json:"returned_at" validate:"omitempty"`
	ConditionNotes string `json:"condition_notes"`
	Status         string `json:"status" validate:"omitempty,oneof=ASSIGNED RETURNED"`
}

type UpdateAssetAssignmentRequest struct {
	AssetID        uint64 `json:"asset_id" validate:"required"`
	AssignedToID   uint64 `json:"assigned_to_id" validate:"required"`
	AssignedAt     string `json:"assigned_at" validate:"required,datetime"`
	ReturnedAt     string `json:"returned_at" validate:"omitempty"`
	ConditionNotes string `json:"condition_notes"`
	Status         string `json:"status" validate:"required,oneof=ASSIGNED RETURNED"`
}
