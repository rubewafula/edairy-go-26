package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type EmployeeDeductionTypeService struct{}

func NewEmployeeDeductionTypeService() *EmployeeDeductionTypeService {
	return &EmployeeDeductionTypeService{}
}

func (s *EmployeeDeductionTypeService) CreateDeductionType(req dtos.CreateEmployeeDeductionTypeRequest, userID uint64) (*models.EmployeeDeductionType, error) {
	dtype := &models.EmployeeDeductionType{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		Name:        req.Name,
		Description: req.Description,
		IsStatutory: req.IsStatutory,
	}
	err := db.DB.Create(dtype).Error
	return dtype, err
}

func (s *EmployeeDeductionTypeService) GetDeductionTypes(page, limit int) ([]models.EmployeeDeductionType, int64, error) {
	var dtypes []models.EmployeeDeductionType
	var total int64
	db.DB.Model(&models.EmployeeDeductionType{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&dtypes).Error
	return dtypes, total, err
}

func (s *EmployeeDeductionTypeService) GetDeductionType(id string) (*models.EmployeeDeductionType, error) {
	var dtype models.EmployeeDeductionType
	if err := db.DB.First(&dtype, id).Error; err != nil {
		return nil, err
	}
	return &dtype, nil
}

func (s *EmployeeDeductionTypeService) UpdateDeductionType(id string, req dtos.UpdateEmployeeDeductionTypeRequest, userID uint64) error {
	var dtype models.EmployeeDeductionType
	if err := db.DB.First(&dtype, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"name":         req.Name,
		"description":  req.Description,
		"is_statutory": req.IsStatutory,
		"updated_by":   userID,
	}
	return db.DB.Model(&dtype).Updates(updates).Error
}

func (s *EmployeeDeductionTypeService) DeleteDeductionType(id string, userID uint64) error {
	var dtype models.EmployeeDeductionType
	if err := db.DB.First(&dtype, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&dtype).Update("updated_by", userID).Delete(&dtype).Error
}
