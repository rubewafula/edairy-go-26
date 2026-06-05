package dtos

type CreateAssetRequest struct {
	AssetCode               string  `json:"asset_code" validate:"required,max=100"`
	AssetName               string  `json:"asset_name" validate:"required,max=255"`
	CategoryID              uint64  `json:"asset_category_id" validate:"required"`
	SerialNo                string  `json:"serial_no" validate:"max=255"`
	Barcode                 string  `json:"barcode" validate:"max=255"`
	Manufacturer            string  `json:"manufacturer" validate:"max=255"`
	VendorID                uint64  `json:"vendor_id"`
	PurchaseCost            float64 `json:"purchase_cost" validate:"required,min=0"`
	SalvageValue            float64 `json:"salvage_value"`
	AcquisitionDate         string  `json:"acquisition_date" validate:"required"`
	UsefulLifeYears         int     `json:"useful_life_years"`
	DepreciationMethod      string  `json:"depreciation_method" validate:"omitempty,oneof=STRAIGHT_LINE DECLINING_BALANCE UNITS_OF_PRODUCTION"`
	DepreciationRate        float64 `json:"depreciation_rate"`
	AccumulatedDepreciation float64 `json:"accumulated_depreciation"`
	BookValue               float64 `json:"book_value"`
	WarrantyEndDate         string  `json:"warranty_end_date" validate:"required"`
	LocationID              uint64  `json:"location_id"`
	Status                  string  `json:"status" validate:"omitempty,oneof=ACTIVE MAINTENANCE DISPOSED WRITTEN_OFF"`
	Loanable                bool    `json:"loanable"`
	Comments                string  `json:"comments"`
}

type UpdateAssetRequest struct {
	AssetCode               string  `json:"asset_code" validate:"required,max=100"`
	AssetName               string  `json:"asset_name" validate:"required,max=255"`
	CategoryID              uint64  `json:"asset_category_id" validate:"required"`
	SerialNo                string  `json:"serial_no" validate:"max=255"`
	Barcode                 string  `json:"barcode" validate:"max=255"`
	Manufacturer            string  `json:"manufacturer" validate:"max=255"`
	VendorID                uint64  `json:"vendor_id"`
	PurchaseCost            float64 `json:"purchase_cost" validate:"required,min=0"`
	SalvageValue            float64 `json:"salvage_value"`
	AcquisitionDate         string  `json:"acquisition_date" validate:"required,datetime"`
	UsefulLifeYears         int     `json:"useful_life_years"`
	DepreciationMethod      string  `json:"depreciation_method" validate:"omitempty,oneof=STRAIGHT_LINE DECLINING_BALANCE UNITS_OF_PRODUCTION"`
	DepreciationRate        float64 `json:"depreciation_rate"`
	AccumulatedDepreciation float64 `json:"accumulated_depreciation"`
	BookValue               float64 `json:"book_value"`
	WarrantyEndDate         string  `json:"warranty_end_date" validate:"required"`
	LocationID              uint64  `json:"location_id"`
	Status                  string  `json:"status" validate:"omitempty,oneof=ACTIVE MAINTENANCE DISPOSED WRITTEN_OFF"`
	Loanable                bool    `json:"loanable"`
	Comments                string  `json:"comments"`
}
