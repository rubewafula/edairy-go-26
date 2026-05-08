package dtos

import "time"

type ShareTypeResponse struct {
	ID                uint64    `json:"ID"`
	ShareCode         string    `json:"ShareCode"`
	ShareType         string    `json:"ShareType"`
	Description       string    `json:"Description"`
	Rate              float64   `json:"Rate"`
	Mandatory         int       `json:"Mandatory"`
	HasShareValue     string    `json:"HasShareValue"`
	RepayMethod       string    `json:"RepayMethod"`
	CalculatingMethod string    `json:"CalculatingMethod"`
	ShareValue        float64   `json:"ShareValue"`
	DeductionTypeID   uint64    `json:"DeductionTypeID"`
	DeductionTypeName string    `json:"DeductionTypeName"`
	Priority          int       `json:"Priority"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
}
