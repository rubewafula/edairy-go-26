package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type SupplierContactService struct{}

func NewSupplierContactService() *SupplierContactService {
	return &SupplierContactService{}
}

func (s *SupplierContactService) CreateContact(req dtos.CreateSupplierContactRequest, userID uint64) (*models.SupplierContact, error) {
	contact := &models.SupplierContact{
		BaseModel:          models.BaseModel{CreatedBy: userID},
		SupplierID:         req.SupplierID,
		ContactType:        req.ContactType,
		FullName:           req.FullName,
		Designation:        req.Designation,
		PhoneNo:            req.PhoneNo,
		AlternativePhoneNo: req.AlternativePhoneNo,
		EmailAddress:       req.EmailAddress,
		IsDefault:          req.IsDefault,
		Notes:              req.Notes,
	}
	if err := db.DB.Create(contact).Error; err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *SupplierContactService) GetContacts(page, limit int) ([]dtos.SupplierContactResponse, int64, error) {
	var results []dtos.SupplierContactResponse
	var total int64
	db.DB.Model(&models.SupplierContact{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sc.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM supplier_contacts sc
		LEFT JOIN suppliers s ON sc.supplier_id = s.id
		WHERE sc.deleted_at IS NULL
		ORDER BY sc.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplierContactService) GetContact(id string) (*dtos.SupplierContactResponse, error) {
	var result dtos.SupplierContactResponse
	query := `
		SELECT 
			sc.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM supplier_contacts sc
		LEFT JOIN suppliers s ON sc.supplier_id = s.id
		WHERE sc.id = ? AND sc.deleted_at IS NULL
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

func (s *SupplierContactService) UpdateContact(id string, req dtos.UpdateSupplierContactRequest, userID uint64) error {
	var contact models.SupplierContact
	if err := db.DB.First(&contact, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"contact_type":         req.ContactType,
		"full_name":            req.FullName,
		"designation":          req.Designation,
		"phone_no":             req.PhoneNo,
		"alternative_phone_no": req.AlternativePhoneNo,
		"email_address":        req.EmailAddress,
		"is_default":           req.IsDefault,
		"notes":                req.Notes,
		"updated_by":           userID,
	}

	return db.DB.Model(&contact).Updates(updates).Error
}

func (s *SupplierContactService) DeleteContact(id string, userID uint64) error {
	var contact models.SupplierContact
	if err := db.DB.First(&contact, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&contact).Update("updated_by", userID).Delete(&contact).Error
}

func (s *SupplierContactService) GetContactsBySupplier(supplierID string) ([]dtos.SupplierContactResponse, error) {
	var results []dtos.SupplierContactResponse
	query := `
		SELECT 
			sc.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM supplier_contacts sc
		LEFT JOIN suppliers s ON sc.supplier_id = s.id
		WHERE sc.supplier_id = ? AND sc.deleted_at IS NULL
		ORDER BY sc.id DESC
	`
	err := db.DB.Raw(query, supplierID).Scan(&results).Error
	return results, err
}
