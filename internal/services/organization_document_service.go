package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationDocumentService struct{}

func NewOrganizationDocumentService() *OrganizationDocumentService {
	return &OrganizationDocumentService{}
}

func (s *OrganizationDocumentService) CreateDocument(req dtos.CreateOrganizationDocumentRequest, userID uint64) (*models.OrganizationDocument, error) {
	document := &models.OrganizationDocument{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		AstraID:      req.AstraID,
		DocumentType: req.DocumentType,
		Document:     req.Document,
		Submitted:    req.Submitted,
	}
	if err := db.DB.Create(document).Error; err != nil {
		return nil, err
	}
	return document, nil
}

func (s *OrganizationDocumentService) GetDocuments(page, limit int) ([]dtos.OrganizationDocumentResponse, int64, error) {
	var results []dtos.OrganizationDocumentResponse
	var total int64
	db.DB.Model(&models.OrganizationDocument{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.OrganizationDocument{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *OrganizationDocumentService) GetDocument(id string) (*dtos.OrganizationDocumentResponse, error) {
	var result dtos.OrganizationDocumentResponse
	err := db.DB.Model(&models.OrganizationDocument{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *OrganizationDocumentService) UpdateDocument(id string, req dtos.UpdateOrganizationDocumentRequest, userID uint64) error {
	var document models.OrganizationDocument
	if err := db.DB.First(&document, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"astra_id":      req.AstraID,
		"document_type": req.DocumentType,
		"document":      req.Document,
		"submitted":     req.Submitted,
		"updated_by":    userID,
	}

	return db.DB.Model(&document).Updates(updates).Error
}

func (s *OrganizationDocumentService) DeleteDocument(id string, userID uint64) error {
	var document models.OrganizationDocument
	if err := db.DB.First(&document, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&document).Update("updated_by", userID).Delete(&document).Error
}

func (s *OrganizationDocumentService) GetDocumentsByAstraID(astraID string) ([]dtos.OrganizationDocumentResponse, error) {
	var results []dtos.OrganizationDocumentResponse
	err := db.DB.Model(&models.OrganizationDocument{}).Where("astra_id = ?", astraID).Find(&results).Error
	return results, err
}
