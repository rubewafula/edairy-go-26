package dtos

import "time"

type TrainingResponse struct {
	ID           uint64    `json:"ID"`
	Topic        string    `json:"Topic"`
	Description  string    `json:"Description"`
	Venue        string    `json:"Venue"`
	Facilitator  string    `json:"Facilitator"`
	TrainingDate time.Time `json:"TrainingDate"`
	Status       string    `json:"Status"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}
