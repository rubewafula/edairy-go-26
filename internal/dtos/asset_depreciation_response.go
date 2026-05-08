package dtos

import "time"

type AssetDepreciationResponse struct {
	ID                      uint64    `json:"ID"`
	AssetID                 uint64    `json:"AssetID"`
	AssetName               string    `json:"AssetName"`
	AssetCode               string    `json:"AssetCode"`
	DepreciationDate        time.Time `json:"DepreciationDate"`
	DepreciationAmount      float64   `json:"DepreciationAmount"`
	AccumulatedDepreciation float64   `json:"AccumulatedDepreciation"`
	BookValue               float64   `json:"BookValue"`
	TransactionID           *uint64   `json:"TransactionID"`
	CreatedAt               time.Time `json:"CreatedAt"`
}
