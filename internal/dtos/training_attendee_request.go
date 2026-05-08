package dtos

type CreateTrainingAttendeeRequest struct {
	TrainingSessionID uint64 `json:"training_session_id" validate:"required"`
	Names             string `json:"names" validate:"required"`
	IDNumber          string `json:"id_number" validate:"required"`
	PhoneNumber       string `json:"phone_number" validate:"required"`
	MembershipNumber  string `json:"membership_number"`
	Comments          string `json:"comments"`
	MemberID          uint64 `json:"member_id"`
}

type UpdateTrainingAttendeeRequest struct {
	TrainingSessionID uint64 `json:"training_session_id" validate:"required"`
	Names             string `json:"names" validate:"required"`
	IDNumber          string `json:"id_number" validate:"required"`
	PhoneNumber       string `json:"phone_number" validate:"required"`
	MembershipNumber  string `json:"membership_number"`
	Comments          string `json:"comments"`
	MemberID          uint64 `json:"member_id"`
}
