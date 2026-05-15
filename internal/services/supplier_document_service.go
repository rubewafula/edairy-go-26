package services

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type SupplierDocumentService struct{}

func NewSupplierDocumentService() *SupplierDocumentService {
	return &SupplierDocumentService{}
}

func (s *SupplierDocumentService) CreateDocument(req dtos.CreateSupplierDocumentRequest, userID uint64) (*models.SupplierDocument, error) {
	doc := &models.SupplierDocument{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		SupplierID:     req.SupplierID,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		Notes:          req.Notes,
		Verified:       "no",
	}

	if req.IssueDate != "" {
		t := utils.ParseDate(req.IssueDate)
		doc.IssueDate = &t
	}
	if req.ExpiryDate != "" {
		t := utils.ParseDate(req.ExpiryDate)
		doc.ExpiryDate = &t
	}

	if err := db.DB.Create(doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *SupplierDocumentService) GetDocuments(page, limit int) ([]dtos.SupplierDocumentResponse, int64, error) {
	var results []dtos.SupplierDocumentResponse
	var total int64
	db.DB.Model(&models.SupplierDocument{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sd.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM supplier_documents sd
		LEFT JOIN suppliers s ON sd.supplier_id = s.id
		WHERE sd.deleted_at IS NULL
		ORDER BY sd.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplierDocumentService) GetDocument(id string) (*dtos.SupplierDocumentResponse, error) {
	var result dtos.SupplierDocumentResponse
	query := `
		SELECT 
			sd.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM supplier_documents sd
		LEFT JOIN suppliers s ON sd.supplier_id = s.id
		WHERE sd.id = ? AND sd.deleted_at IS NULL
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

func (s *SupplierDocumentService) UpdateDocument(id string, req dtos.UpdateSupplierDocumentRequest, userID uint64) error {
	var doc models.SupplierDocument
	if err := db.DB.First(&doc, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"document_type":   req.DocumentType,
		"document_number": req.DocumentNumber,
		"notes":           req.Notes,
		"updated_by":      userID,
	}

	if req.IssueDate != "" {
		updates["issue_date"] = utils.ParseDate(req.IssueDate)
	}
	if req.ExpiryDate != "" {
		updates["expiry_date"] = utils.ParseDate(req.ExpiryDate)
	}

	return db.DB.Model(&doc).Updates(updates).Error
}

func (s *SupplierDocumentService) DeleteDocument(id string, userID uint64) error {
	var doc models.SupplierDocument
	if err := db.DB.First(&doc, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&doc).Update("updated_by", userID).Delete(&doc).Error
}

func (s *SupplierDocumentService) VerifyDocument(id string, req dtos.VerifySupplierDocumentRequest, userID uint64) error {
	now := time.Now()
	updates := map[string]interface{}{
		"verified":    req.Verified,
		"verified_by": userID,
		"verified_at": &now,
		"updated_by":  userID,
	}
	return db.DB.Model(&models.SupplierDocument{}).Where("id = ?", id).Updates(updates).Error
}

func (s *SupplierDocumentService) GetDocumentsBySupplier(supplierID string) ([]dtos.SupplierDocumentResponse, error) {
	var results []dtos.SupplierDocumentResponse
	query := `
		SELECT 
			sd.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM supplier_documents sd
		LEFT JOIN suppliers s ON sd.supplier_id = s.id
		WHERE sd.supplier_id = ? AND sd.deleted_at IS NULL
		ORDER BY sd.id DESC
	`
	err := db.DB.Raw(query, supplierID).Scan(&results).Error
	return results, err
}
