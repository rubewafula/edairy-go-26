package dtos

import "time"

type TrainingAttendeeResponse struct {
	ID                uint64    `json:"id"`
	TrainingSessionID uint64    `json:"training_session_id"`
	Topic             string    `json:"topic"`
	Names             string    `json:"names"`
	IDNumber          string    `json:"id_number"`
	PhoneNumber       string    `json:"phone_number"`
	MembershipNumber  string    `json:"membership_number"`
	Comments          string    `json:"comments"`
	MemberID          uint64    `json:"member_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
