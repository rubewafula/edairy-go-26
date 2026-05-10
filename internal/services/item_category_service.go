package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type ItemCategoryService struct{}

func NewItemCategoryService() *ItemCategoryService {
	return &ItemCategoryService{}
}

func (s *ItemCategoryService) CreateCategory(req dtos.CreateItemCategoryRequest) (*models.ItemCategory, error) {
	category := &models.ItemCategory{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryID: req.ParentCategoryID,
	}

	if err := db.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *ItemCategoryService) GetCategories(page, limit int) ([]dtos.ItemCategoryResponse, int64, error) {
	var results []dtos.ItemCategoryResponse
	var total int64
	db.DB.Model(&models.ItemCategory{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ic.id, ic.name, ic.description, ic.parent_category_id, 
			p.name AS parent_category_name,
			ic.created_at, ic.updated_at
		FROM item_categories ic
		LEFT JOIN item_categories p ON ic.parent_category_id = p.id
		WHERE ic.deleted_at IS NULL
		ORDER BY ic.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *ItemCategoryService) GetCategory(id string) (*dtos.ItemCategoryResponse, error) {
	var result dtos.ItemCategoryResponse
	query := `
		SELECT 
			ic.id, ic.name, ic.description, ic.parent_category_id, 
			p.name AS parent_category_name,
			ic.created_at, ic.updated_at
		FROM item_categories ic
		LEFT JOIN item_categories p ON ic.parent_category_id = p.id
		WHERE ic.id = ? AND ic.deleted_at IS NULL
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

func (s *ItemCategoryService) UpdateCategory(id string, req dtos.UpdateItemCategoryRequest) error {
	return db.DB.Model(&models.ItemCategory{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":               req.Name,
		"description":        req.Description,
		"parent_category_id": req.ParentCategoryID,
	}).Error
}

func (s *ItemCategoryService) DeleteCategory(id string) error {
	return db.DB.Delete(&models.ItemCategory{}, id).Error
}
