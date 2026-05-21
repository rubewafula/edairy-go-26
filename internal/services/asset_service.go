package services

import (
	"log"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type AssetService struct{}

func NewAssetService() *AssetService {
	return &AssetService{}
}

func (s *AssetService) CreateAsset(req dtos.CreateAssetRequest) (*models.Asset, error) {
	asset := &models.Asset{
		AssetCode:               req.AssetCode,
		AssetName:               req.AssetName,
		CategoryID:              req.CategoryID,
		SerialNo:                req.SerialNo,
		Barcode:                 req.Barcode,
		Manufacturer:            req.Manufacturer,
		VendorID:                req.VendorID,
		PurchaseCost:            req.PurchaseCost,
		SalvageValue:            req.SalvageValue,
		AcquisitionDate:         utils.ParseDate(req.AcquisitionDate),
		UsefulLifeYears:         req.UsefulLifeYears,
		DepreciationMethod:      req.DepreciationMethod,
		DepreciationRate:        req.DepreciationRate,
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		BookValue:               req.BookValue,
		WarrantyEndDate:         utils.ParseDate(req.WarrantyEndDate),
		CurrentLocation:         req.CurrentLocation,
		Status:                  req.Status,
		Loanable:                req.Loanable,
		Comments:                req.Comments,
	}

	if err := db.DB.Create(asset).Error; err != nil {
		log.Println("Found error trying to created asset: %s", err.Error())
		return nil, err
	}
	return asset, nil
}

func (s *AssetService) GetAssets() ([]dtos.AssetResponse, int64, error) {
	var results []dtos.AssetResponse
	var total int64
	db.DB.Model(&models.Asset{}).Count(&total)

	query := `
		SELECT 
			a.id, a.asset_code, a.asset_name, a.asset_category_id as category_id, ac.name as category_name,
			a.serial_no, a.barcode, a.manufacturer, a.vendor_id, a.purchase_cost, a.salvage_value,
			a.acquisition_date, a.useful_life_years, a.depreciation_method, a.depreciation_rate,
			a.accumulated_depreciation, a.book_value, a.warranty_end_date, a.current_location,
			a.status, a.loanable, a.comments, a.created_by, a.updated_by, a.created_at, a.updated_at
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.ud
		WHERE a.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *AssetService) GetAsset(id string) (*dtos.AssetResponse, error) {
	var result dtos.AssetResponse
	query := `
		SELECT 
			a.id, a.asset_code, a.asset_name, a.asset_category_id as category_id, ac.name as category_name,
			a.serial_no, a.barcode, a.manufacturer, a.vendor_id, a.purchase_cost, a.salvage_value,
			a.acquisition_date, a.useful_life_years, a.depreciation_method, a.depreciation_rate,
			a.accumulated_depreciation, a.book_value, a.warranty_end_date, a.current_location,
			a.status, a.loanable, a.comments, a.created_by, a.updated_by, a.created_at, a.updated_at
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.ud
		WHERE a.id = ? AND a.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *AssetService) UpdateAsset(id string, req dtos.UpdateAssetRequest) error {
	var asset models.Asset
	if err := db.DB.First(&asset, id).Error; err != nil {
		return err
	}

	asset.AssetCode = req.AssetCode
	asset.AssetName = req.AssetName
	asset.CategoryID = req.CategoryID
	asset.SerialNo = req.SerialNo
	asset.Barcode = req.Barcode
	asset.Manufacturer = req.Manufacturer
	asset.VendorID = req.VendorID
	asset.PurchaseCost = req.PurchaseCost
	asset.SalvageValue = req.SalvageValue
	asset.AcquisitionDate = utils.ParseDate(req.AcquisitionDate)
	asset.UsefulLifeYears = req.UsefulLifeYears
	asset.DepreciationMethod = req.DepreciationMethod
	asset.DepreciationRate = req.DepreciationRate
	asset.AccumulatedDepreciation = req.AccumulatedDepreciation
	asset.BookValue = req.BookValue
	asset.WarrantyEndDate = utils.ParseDate(req.WarrantyEndDate)
	asset.CurrentLocation = req.CurrentLocation
	asset.Status = req.Status
	asset.Loanable = req.Loanable
	asset.Comments = req.Comments

	return db.DB.Save(&asset).Error
}

func (s *AssetService) DeleteAsset(id string) error {
	return db.DB.Delete(&models.Asset{}, id).Error
}
