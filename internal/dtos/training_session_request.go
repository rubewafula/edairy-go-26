package dtos

type CreateTrainingSessionRequest struct {
	TrainingID uint64 `json:"training_id" validate:"required"`
	MemberID   uint64 `json:"member_id" validate:"required"`
	Status     string `json:"status" validate:"omitempty,oneof=INVITED ATTENDED ABSENT"`
	Remarks    string `json:"remarks"`
}

type UpdateTrainingSessionRequest struct {
	TrainingID uint64 `json:"training_id" validate:"required"`
	MemberID   uint64 `json:"member_id" validate:"required"`
	Status     string `json:"status" validate:"required,oneof=INVITED ATTENDED ABSENT"`
	Remarks    string `json:"remarks"`
}
