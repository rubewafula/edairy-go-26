package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type EmployeeTerminationCategoryService struct{}

func NewEmployeeTerminationCategoryService() *EmployeeTerminationCategoryService {
	return &EmployeeTerminationCategoryService{}
}

func (s *EmployeeTerminationCategoryService) Create(req dtos.CreateEmployeeTerminationCategoryRequest, userID uint64) (*models.EmployeeTerminationCategory, error) {
	category := &models.EmployeeTerminationCategory{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		Name:        req.Name,
		Description: req.Description,
	}
	if err := db.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *EmployeeTerminationCategoryService) List(page, limit int) ([]dtos.EmployeeTerminationCategoryResponse, int64, error) {
	var results []dtos.EmployeeTerminationCategoryResponse
	var total int64

	db.DB.Model(&models.EmployeeTerminationCategory{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.EmployeeTerminationCategory{}).Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *EmployeeTerminationCategoryService) Get(id string) (*dtos.EmployeeTerminationCategoryResponse, error) {
	var result dtos.EmployeeTerminationCategoryResponse
	err := db.DB.Model(&models.EmployeeTerminationCategory{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmployeeTerminationCategoryService) Update(id string, req dtos.UpdateEmployeeTerminationCategoryRequest, userID uint64) error {
	var category models.EmployeeTerminationCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"updated_by":  userID,
	}
	return db.DB.Model(&category).Updates(updates).Error
}

func (s *EmployeeTerminationCategoryService) Delete(id string, userID uint64) error {
	var category models.EmployeeTerminationCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&category).Update("updated_by", userID).Delete(&category).Error
}
