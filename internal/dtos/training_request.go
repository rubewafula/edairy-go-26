package dtos

type CreateTrainingRequest struct {
	Topic          string `json:"topic" validate:"required,max=255"`
	Description    string `json:"description"`
	Venue          string `json:"venue" validate:"required"`
	TrainingUserID uint64 `json:"training_user_id" validate:"required"`
	TrainingDate   string `json:"training_date" validate:"required,datetime"`
	Status         string `json:"status" validate:"omitempty,oneof=SCHEDULED COMPLETED CANCELLED"`
}

type UpdateTrainingRequest struct {
	Topic          string `json:"topic" validate:"required,max=255"`
	Description    string `json:"description"`
	Venue          string `json:"venue" validate:"required"`
	TrainingUserID uint64 `json:"training_user_id" validate:"required"`
	TrainingDate   string `json:"training_date" validate:"required,datetime"`
	Status         string `json:"status" validate:"required,oneof=SCHEDULED COMPLETED CANCELLED"`
}
