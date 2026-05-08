package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type AssetDepreciationService struct{}

func NewAssetDepreciationService() *AssetDepreciationService {
	return &AssetDepreciationService{}
}

func (s *AssetDepreciationService) CreateEntry(req dtos.CreateAssetDepreciationRequest) (*models.AssetDepreciationEntry, error) {
	entry := &models.AssetDepreciationEntry{
		AssetID:                 req.AssetID,
		DepreciationDate:        utils.ParseDate(req.DepreciationDate),
		DepreciationAmount:      req.DepreciationAmount,
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		BookValue:               req.BookValue,
		TransactionID:           req.TransactionID,
	}

	if err := db.DB.Create(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *AssetDepreciationService) GetEntries() ([]dtos.AssetDepreciationResponse, int64, error) {
	var results []dtos.AssetDepreciationResponse
	var total int64
	db.DB.Model(&models.AssetDepreciationEntry{}).Count(&total)

	query := `
		SELECT 
			ade.id, fa.id AS asset_id, fa.asset_name, fa.asset_code,
			ade.depreciation_date, ade.depreciation_amount, ade.accumulated_depreciation,
			ade.book_value, ade.transaction_id, ade.created_at
		FROM asset_depreciation_entries ade
		LEFT JOIN fixed_assets fa ON ade.asset_id = fa.id
		WHERE ade.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *AssetDepreciationService) GetEntry(id string) (*dtos.AssetDepreciationResponse, error) {
	var result dtos.AssetDepreciationResponse
	query := `
		SELECT 
			ade.id, ade.asset_id, fa.asset_name, fa.asset_code,
			ade.depreciation_date, ade.depreciation_amount, ade.accumulated_depreciation,
			ade.book_value, ade.transaction_id, ade.created_at
		FROM asset_depreciation_entries ade
		LEFT JOIN fixed_assets fa ON ade.asset_id = fa.id
		WHERE ade.id = ? AND ade.deleted_at IS NULL
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

func (s *AssetDepreciationService) DeleteEntry(id string) error {
	return db.DB.Delete(&models.AssetDepreciationEntry{}, id).Error
}
