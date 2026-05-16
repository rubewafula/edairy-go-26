package dtos

import "time"

type TransporterBenefitResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	MinQuantity string    `json:"min_quantity"`
	Rate        string    `json:"rate"`
	RouteID     uint64    `json:"route_id"`
	RouteName   string    `json:"route_name"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
