package dtos

import "time"

type TransporterBenefitResponse struct {
	ID          uint64    `json:"ID"`
	Name        string    `json:"Name"`
	MinQuantity string    `json:"MinQuantity"`
	Rate        string    `json:"Rate"`
	RouteID     uint64    `json:"RouteID"`
	RouteName   string    `json:"RouteName"`
	Status      string    `json:"Status"`
	StartDate   time.Time `json:"StartDate"`
	EndDate     time.Time `json:"EndDate"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
