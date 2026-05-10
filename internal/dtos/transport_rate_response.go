package dtos

import "time"

type TransportRateResponse struct {
	ID              uint64
	MemberNo        string
	MemberFirstName string
	MemberLastName  string
	RouteName       string
	TransporterNo   string
	Rate            float64
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
