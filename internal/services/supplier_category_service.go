package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type SupplierCategoryService struct{}

func NewSupplierCategoryService() *SupplierCategoryService {
	return &SupplierCategoryService{}
}

func (s *SupplierCategoryService) CreateCategory(req dtos.CreateSupplierCategoryRequest, userID uint64) (*models.SupplierCategory, error) {
	category := &models.SupplierCategory{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		CategoryCode: req.CategoryCode,
		CategoryName: req.CategoryName,
		Description:  req.Description,
		Status:       req.Status,
		SiteID:       req.SiteID,
	}
	if err := db.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *SupplierCategoryService) GetCategories(page, limit int) ([]dtos.SupplierCategoryResponse, int64, error) {
	var results []dtos.SupplierCategoryResponse
	var total int64
	db.DB.Model(&models.SupplierCategory{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.SupplierCategory{}).
		Select("id, category_code, category_name, description, status, created_at").
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *SupplierCategoryService) GetCategory(id string) (*models.SupplierCategory, error) {
	var category models.SupplierCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *SupplierCategoryService) UpdateCategory(id string, req dtos.CreateSupplierCategoryRequest, userID uint64) error {
	var category models.SupplierCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"category_code": req.CategoryCode,
		"category_name": req.CategoryName,
		"description":   req.Description,
		"status":        req.Status,
		"site_id":       req.SiteID,
		"updated_by":    userID,
	}

	return db.DB.Model(&category).Updates(updates).Error
}

func (s *SupplierCategoryService) DeleteCategory(id string) error {
	return db.DB.Delete(&models.SupplierCategory{}, id).Error
}
