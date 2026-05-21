package dtos

import "time"

// TransportRateResponse represents the transport rate (likely transporter related).
type TransportRateResponse struct {
	ID              uint64    `json:"id"`
	MemberNo        string    `json:"member_no"`
	MemberFirstName string    `json:"member_first_name"`
	MemberLastName  string    `json:"member_last_name"`
	TransporterNo   string    `json:"transporter_no"`
	RouteName       string    `json:"route_name"`
	Rate            float64   `json:"rate"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
