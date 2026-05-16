package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeDocumentService struct{}

func NewEmployeeDocumentService() *EmployeeDocumentService {
	return &EmployeeDocumentService{}
}

func (s *EmployeeDocumentService) CreateEmployeeDocument(req dtos.CreateEmployeeDocumentRequest, userID uint64) (*models.EmployeeDocument, error) {
	document := &models.EmployeeDocument{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		EmployeeID:      req.EmployeeID,
		DocumentTypeID:  req.DocumentTypeID,
		FileName:        req.FileName,
		FileDescription: req.FileDescription,
	}
	if err := db.DB.Create(document).Error; err != nil {
		return nil, err
	}
	return document, nil
}

func (s *EmployeeDocumentService) GetEmployeeDocuments(page, limit int) ([]dtos.EmployeeDocumentResponse, int64, error) {
	var results []dtos.EmployeeDocumentResponse
	var total int64
	db.DB.Model(&models.EmployeeDocument{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ed.id, ed.employee_id, ed.document_type_id, dt.document_type,
			ed.file_name, ed.file_description, ed.created_at, ed.updated_at
		FROM employee_documents ed
		LEFT JOIN document_types dt ON ed.document_type_id = dt.id 
		WHERE ed.deleted_at IS NULL
		ORDER BY ed.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeDocumentService) GetEmployeeDocument(id string) (*dtos.EmployeeDocumentResponse, error) {
	var result dtos.EmployeeDocumentResponse
	query := `
		SELECT 
			ed.id, ed.employee_id, ed.document_type_id, dt.document_type,
			ed.file_name, ed.file_description, ed.created_at, ed.updated_at
		FROM employee_documents ed
		LEFT JOIN document_types dt ON ed.document_type_id = dt.id
		WHERE ed.id = ? AND ed.deleted_at IS NULL
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

func (s *EmployeeDocumentService) UpdateEmployeeDocument(id string, req dtos.UpdateEmployeeDocumentRequest, userID uint64) error {
	var document models.EmployeeDocument
	if err := db.DB.First(&document, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"document_type_id": req.DocumentTypeID,
		"file_name":        req.FileName,
		"file_description": req.FileDescription,
		"updated_by":       userID,
	}
	return db.DB.Model(&document).Updates(updates).Error
}

func (s *EmployeeDocumentService) DeleteEmployeeDocument(id string, userID uint64) error {
	var document models.EmployeeDocument
	if err := db.DB.First(&document, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&document).Update("updated_by", userID).Delete(&document).Error
}
