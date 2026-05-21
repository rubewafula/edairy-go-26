package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type DocumentTypeService struct{}

func NewDocumentTypeService() *DocumentTypeService {
	return &DocumentTypeService{}
}

func (s *DocumentTypeService) Create(req dtos.CreateDocumentTypeRequest, userID uint64) (*models.DocumentType, error) {
	docType := &models.DocumentType{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		DocumentType: req.DocumentType,
		Description:  req.Description,
	}

	if err := db.DB.Create(docType).Error; err != nil {
		return nil, err
	}
	return docType, nil
}

func (s *DocumentTypeService) List(page, limit int) ([]dtos.DocumentTypeResponse, int64, error) {
	var results []dtos.DocumentTypeResponse
	var total int64

	db.DB.Model(&models.DocumentType{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.DocumentType{}).
		Limit(limit).Offset(offset).Order("id DESC").
		Scan(&results).Error

	return results, total, err
}

func (s *DocumentTypeService) Get(id string) (*dtos.DocumentTypeResponse, error) {
	var result dtos.DocumentTypeResponse
	err := db.DB.Model(&models.DocumentType{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (s *DocumentTypeService) Update(id string, req dtos.UpdateDocumentTypeRequest, userID uint64) error {
	var docType models.DocumentType
	if err := db.DB.First(&docType, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"document_type": req.DocumentType,
		"description":   req.Description,
		"updated_by":    userID,
	}

	return db.DB.Model(&docType).Updates(updates).Error
}

func (s *DocumentTypeService) Delete(id string, userID uint64) error {
	var docType models.DocumentType
	if err := db.DB.First(&docType, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&docType).Update("updated_by", userID).Delete(&docType).Error
}
