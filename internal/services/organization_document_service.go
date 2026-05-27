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
		BaseModel:      models.BaseModel{CreatedBy: userID},
		AstraID:        req.AstraID,
		DocumentTypeID: req.DocumentTypeID,
		DocumentName:   req.DocumentName,
		Submitted:      req.Submitted,
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

	// Join with document_types to get the document type name
	query := `
		SELECT 
			od.id, od.astra_id, od.document_type_id, dt.document_type as document_type_name,
			od.document_name, od.document as document_url, od.submitted,
			od.created_at, od.updated_at
		FROM organization_documents od
		LEFT JOIN document_types dt ON od.document_type_id = dt.id
		WHERE od.deleted_at IS NULL
		ORDER BY od.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *OrganizationDocumentService) GetDocument(id string) (*dtos.OrganizationDocumentResponse, error) {
	var result dtos.OrganizationDocumentResponse
	// Join with document_types to get the document type name
	query := `
		SELECT 
			od.id, od.astra_id, od.document_type_id, dt.document_type as document_type_name,
			od.document_name, od.document as document_url, od.submitted,
			od.created_at, od.updated_at
		FROM organization_documents od
		LEFT JOIN document_types dt ON od.document_type_id = dt.id
		WHERE od.id = ? AND od.deleted_at IS NULL
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

func (s *OrganizationDocumentService) UpdateDocument(id string, req dtos.UpdateOrganizationDocumentRequest, userID uint64) error {
	var document models.OrganizationDocument
	if err := db.DB.First(&document, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"astra_id":         req.AstraID,
		"document_type_id": req.DocumentTypeID,
		"document_name":    req.DocumentName,
		"submitted":        req.Submitted,
		"updated_by":       userID,
		// Only update document URL if a new file was uploaded
		"document": document.Document,
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
