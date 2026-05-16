package dtos

import "time"

type AssetDepreciationResponse struct {
	ID                      uint64    `json:"id"`
	AssetID                 uint64    `json:"asset_id"`
	AssetName               string    `json:"asset_name"`
	AssetCode               string    `json:"asset_code"`
	DepreciationDate        time.Time `json:"depreciation_date"`
	DepreciationAmount      float64   `json:"depreciation_amount"`
	AccumulatedDepreciation float64   `json:"accumulated_depreciation"`
	BookValue               float64   `json:"book_value"`
	TransactionID           *uint64   `json:"transaction_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
