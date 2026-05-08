package dtos

import "time"

type TrainingAttendeeResponse struct {
	ID                uint64    `json:"ID"`
	TrainingSessionID uint64    `json:"TrainingSessionID"`
	Topic             string    `json:"Topic"`
	Names             string    `json:"Names"`
	IDNumber          string    `json:"IDNumber"`
	PhoneNumber       string    `json:"PhoneNumber"`
	MembershipNumber  string    `json:"MembershipNumber"`
	Comments          string    `json:"Comments"`
	MemberID          uint64    `json:"MemberID"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
}
