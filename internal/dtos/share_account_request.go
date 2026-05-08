package dtos

type CreateShareAccountRequest struct {
	MemberID    uint64 `json:"member_id" validate:"required"`
	ShareTypeID uint64 `json:"share_type_id"`
	Status      string `json:"status" validate:"omitempty,oneof=ACTIVE SUSPENDED CLOSED"`
	OpenedAt    string `json:"opened_at" validate:"omitempty"`
}

type UpdateShareAccountRequest struct {
	MemberID    uint64 `json:"member_id" validate:"required"`
	ShareTypeID uint64 `json:"share_type_id"`
	Status      string `json:"status" validate:"required,oneof=ACTIVE SUSPENDED CLOSED"`
	OpenedAt    string `json:"opened_at" validate:"omitempty"`
}
