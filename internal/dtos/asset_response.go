package dtos

import "time"

type AssetResponse struct {
	ID                      uint64    `json:"id"`
	AssetCode               string    `json:"asset_code"`
	AssetName               string    `json:"asset_name"`
	CategoryID              uint64    `json:"category_id"`
	CategoryName            string    `json:"category_name"`
	SerialNo                string    `json:"serial_no"`
	Barcode                 string    `json:"barcode"`
	Manufacturer            string    `json:"manufacturer"`
	VendorID                uint64    `json:"vendor_id"`
	PurchaseCost            float64   `json:"purchase_cost"`
	SalvageValue            float64   `json:"salvage_value"`
	AcquisitionDate         time.Time `json:"acquisition_date"`
	UsefulLifeYears         int       `json:"useful_life_years"`
	DepreciationMethod      string    `json:"depreciation_method"`
	DepreciationRate        float64   `json:"depreciation_rate"`
	AccumulatedDepreciation float64   `json:"accumulated_depreciation"`
	BookValue               float64   `json:"book_value"`
	WarrantyEndDate         time.Time `json:"warranty_end_date"`
	CurrentLocation         string    `json:"current_location"`
	Status                  string    `json:"status"`
	Loanable                bool      `json:"loanable"`
	Comments                string    `json:"comments"`
	CreatedBy               uint64    `json:"created_by"`
	UpdatedBy               uint64    `json:"updated_by"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
