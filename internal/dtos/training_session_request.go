package dtos

type CreateTrainingSessionRequest struct {
	TrainingID       uint64 `json:"training_id" validate:"required"`
	Partner          string `json:"partner"`
	SessionStartTime string `json:"session_start_time" validate:"required"`
	SessionEndTime   string `json:"session_end_time" validate:"required"`
	Topic            string `json:"topic" validate:"required"`
	Description      string `json:"description"`
	Trainers         string `json:"trainers"`
}

type UpdateTrainingSessionRequest struct {
	TrainingID       uint64 `json:"training_id" validate:"required"`
	Partner          string `json:"partner"`
	SessionStartTime string `json:"session_start_time" validate:"required"`
	SessionEndTime   string `json:"session_end_time" validate:"required"`
	Topic            string `json:"topic" validate:"required"`
	Description      string `json:"description"`
	Trainers         string `json:"trainers"`
}
