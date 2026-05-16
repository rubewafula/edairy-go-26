package dtos

import "time"

type TrainingSessionResponse struct {
	ID               uint64    `json:"id"`
	TrainingID       uint64    `json:"training_id"`
	Topic            string    `json:"topic"`
	Partner          string    `json:"partner"`
	SessionStartTime time.Time `json:"session_start_time"`
	SessionEndTime   time.Time `json:"session_end_time"`
	Trainsers        string    `json:"trainers"`
	Status           string    `json:"status"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
