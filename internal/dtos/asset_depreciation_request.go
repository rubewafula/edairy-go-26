package dtos

type CreateAssetDepreciationRequest struct {
	AssetID                 uint64  `json:"asset_id" validate:"required"`
	DepreciationDate        string  `json:"depreciation_date" validate:"required,datetime"`
	DepreciationAmount      float64 `json:"depreciation_amount" validate:"required"`
	AccumulatedDepreciation float64 `json:"accumulated_depreciation" validate:"required"`
	BookValue               float64 `json:"book_value" validate:"required"`
	TransactionID           *uint64 `json:"transaction_id"`
}
