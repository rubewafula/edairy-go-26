package dtos

import "time"

type TrainingResponse struct {
	ID           uint64    `json:"id"`
	Topic        string    `json:"topic"`
	Description  string    `json:"description"`
	Venue        string    `json:"venue"`
	Facilitator  string    `json:"facilitator"`
	TrainingDate time.Time `json:"training_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
