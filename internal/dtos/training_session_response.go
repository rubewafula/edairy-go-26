package dtos

import "time"

type TrainingSessionResponse struct {
	ID               uint64    `json:"ID"`
	TrainingID       uint64    `json:"TrainingID"`
	Topic            string    `json:"Topic"`
	Partner          string    `json:"Partner"`
	SessionStartTime time.Time `json:"SessionStartTime"`
	SessionEndTime   time.Time `json:"SessionEndTime"`
	Trainsers        string    `json:"Trainers"`
	Status           string    `json:"Status"`
	Description      string    `json:"Description"`
	CreatedAt        time.Time `json:"CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt"`
}
