package dtos

import "time"

type AssetResponse struct {
	ID                      uint64    `json:"ID"`
	AssetCode               string    `json:"AssetCode"`
	AssetName               string    `json:"AssetName"`
	CategoryID              uint64    `json:"CategoryID"`
	CategoryName            string    `json:"CategoryName"`
	SerialNo                string    `json:"SerialNo"`
	Barcode                 string    `json:"Barcode"`
	Manufacturer            string    `json:"Manufacturer"`
	VendorID                uint64    `json:"VendorID"`
	PurchaseCost            float64   `json:"PurchaseCost"`
	SalvageValue            float64   `json:"SalvageValue"`
	AcquisitionDate         time.Time `json:"AcquisitionDate"`
	UsefulLifeYears         int       `json:"UsefulLifeYears"`
	DepreciationMethod      string    `json:"DepreciationMethod"`
	DepreciationRate        float64   `json:"DepreciationRate"`
	AccumulatedDepreciation float64   `json:"AccumulatedDepreciation"`
	BookValue               float64   `json:"BookValue"`
	WarrantyEndDate         time.Time `json:"WarrantyEndDate"`
	CurrentLocation         string    `json:"CurrentLocation"`
	Status                  string    `json:"Status"`
	Loanable                bool      `json:"Loanable"`
	Comments                string    `json:"Comments"`
	CreatedBy               uint64    `json:"CreatedBy"`
	UpdatedBy               uint64    `json:"UpdatedBy"`
	CreatedAt               time.Time `json:"CreatedAt"`
	UpdatedAt               time.Time `json:"UpdatedAt"`
}
