package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type AssetCategoryService struct{}

func NewAssetCategoryService() *AssetCategoryService {
	return &AssetCategoryService{}
}

func (s *AssetCategoryService) CreateCategory(req dtos.CreateAssetCategoryRequest) (*models.AssetCategory, error) {
	category := &models.AssetCategory{

		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *AssetCategoryService) GetCategories() ([]dtos.AssetCategoryResponse, int64, error) {
	var results []dtos.AssetCategoryResponse

	var total int64
	db.DB.Model(&models.AssetCategory{}).Count(&total)

	query := `
		SELECT ac.id, ac.name, ac.description, ac.created_at, ac.updated_at,
		(SELECT COUNT(*) FROM fixed_assets WHERE asset_category_id = ac.id AND deleted_at IS NULL) AS asset_count
		FROM asset_categories ac
		WHERE ac.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *AssetCategoryService) GetCategory(id string) (*dtos.AssetCategoryResponse, error) {
	var result dtos.AssetCategoryResponse
	query := `
		SELECT ac.id, ac.name, ac.description, ac.created_at, ac.updated_at,
		(SELECT COUNT(*) FROM fixed_assets WHERE asset_category_id = ac.id AND deleted_at IS NULL) AS asset_count
		FROM asset_categories ac
		WHERE ac.id = ? AND ac.deleted_at IS NULL
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

func (s *AssetCategoryService) UpdateCategory(id string, req dtos.UpdateAssetCategoryRequest) error {
	var category models.AssetCategory
	if err := db.DB.First(&category, "id = ?", id).Error; err != nil {

		return err
	}

	category.Name = req.Name
	category.Description = req.Description

	return db.DB.Save(&category).Error

}

func (s *AssetCategoryService) DeleteCategory(id string) error {
	return db.DB.Delete(&models.AssetCategory{}, "id = ?", id).Error
}
