package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type DeductionTypeService struct{}

func NewDeductionTypeService() *DeductionTypeService {
	return &DeductionTypeService{}
}

func (s *DeductionTypeService) CreateDeductionType(req dtos.CreateDeductionTypeRequest) (*models.DeductionType, error) {
	deductionType := &models.DeductionType{
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
		IsStatutory: req.IsStatutory,
	}

	if err := db.DB.Create(deductionType).Error; err != nil {
		return nil, err
	}
	return deductionType, nil
}

func (s *DeductionTypeService) GetDeductionTypes(page, limit int) ([]dtos.DeductionTypeResponse, int64, error) {
	var results []dtos.DeductionTypeResponse
	var total int64
	db.DB.Model(&models.DeductionType{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.DeductionType{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *DeductionTypeService) GetDeductionType(id string) (*dtos.DeductionTypeResponse, error) {
	var deductionTypeModel models.DeductionType
	if err := db.DB.First(&deductionTypeModel, id).Error; err != nil {
		return nil, err // gorm.ErrRecordNotFound will be returned here if not found
	}

	// Map model to DTO for snake_case JSON output
	response := &dtos.DeductionTypeResponse{
		ID:          deductionTypeModel.ID,
		Code:        deductionTypeModel.Code,
		Description: deductionTypeModel.Description,
		Status:      deductionTypeModel.Status,
		IsStatutory: deductionTypeModel.IsStatutory,
		CreatedAt:   deductionTypeModel.CreatedAt,
		UpdatedAt:   deductionTypeModel.UpdatedAt,
	}
	return response, nil
}

func (s *DeductionTypeService) UpdateDeductionType(id string, req dtos.UpdateDeductionTypeRequest) error {
	var deductionType models.DeductionType
	if err := db.DB.First(&deductionType, id).Error; err != nil {
		return err
	}

	deductionType.Code = req.Code
	deductionType.Description = req.Description
	deductionType.Status = req.Status
	deductionType.IsStatutory = req.IsStatutory

	return db.DB.Save(&deductionType).Error
}

func (s *DeductionTypeService) DeleteDeductionType(id string) error {
	return db.DB.Delete(&models.DeductionType{}, id).Error
}
