package dtos

type CreateAssetDepreciationRequest struct {
	AssetID                 uint64  `json:"asset_id" validate:"required"`
	DepreciationDate        string  `json:"depreciation_date" validate:"required"`
	DepreciationAmount      float64 `json:"depreciation_amount" validate:"required"`
	AccumulatedDepreciation float64 `json:"accumulated_depreciation"`
	BookValue               float64 `json:"book_value" validate:"required"`
	TransactionID           uint64  `json:"transaction_id"`
	Notes                   string  `json:"notes"`
}

type UpdateAssetDepreciationRequest struct {
	AccumulatedDepreciation float64 `json:"accumulated_depreciation"`
	BookValue               float64 `json:"book_value"`
	DepreciationDate        string  `json:"depreciation_date"`
	DepreciationAmount      float64 `json:"depreciation_amount" validate:"gt=0"`
	Notes                   string  `json:"notes"`
}
