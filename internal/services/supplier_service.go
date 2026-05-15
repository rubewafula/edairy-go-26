package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type SupplierService struct{}

func NewSupplierService() *SupplierService {
	return &SupplierService{}
}

func (s *SupplierService) CreateSupplier(req dtos.CreateSupplierRequest, userID uint64) (*models.Supplier, error) {
	supplier := &models.Supplier{
		BaseModel:          models.BaseModel{CreatedBy: userID},
		SupplierCategoryID: req.SupplierCategoryID,
		SupplierCode:       req.SupplierCode,
		SupplierType:       req.SupplierType,
		CompanyName:        req.CompanyName,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		PhoneNo:            req.PhoneNo,
		EmailAddress:       req.EmailAddress,
		KraPin:             req.KraPin,
		OpeningBalance:     req.OpeningBalance,
		CurrentBalance:     req.OpeningBalance,
		CreditLimit:        req.CreditLimit,
		PaymentTermsDays:   req.PaymentTermsDays,
		Status:             req.Status,
		Notes:              req.Notes,
	}
	if err := db.DB.Create(supplier).Error; err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) GetSuppliers(page, limit int) ([]dtos.SupplierResponse, int64, error) {
	var results []dtos.SupplierResponse
	var total int64
	db.DB.Model(&models.Supplier{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			s.id, s.supplier_code, s.supplier_type, s.company_name, 
			CONCAT(COALESCE(s.first_name,''), ' ', COALESCE(s.last_name,'')) as full_name,
			sc.category_name, s.email_address, s.phone_no, s.current_balance, s.status, s.created_at
		FROM suppliers s
		LEFT JOIN supplier_categories sc ON s.supplier_category_id = sc.id
		WHERE s.deleted_at IS NULL
		ORDER BY s.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplierService) GetSupplier(id string) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := db.DB.First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (s *SupplierService) CreateContact(req dtos.CreateSupplierContactRequest, userID uint64) (*models.SupplierContact, error) {
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
	err := db.DB.Create(contact).Error
	return contact, err
}

func (s *SupplierService) CreateBankAccount(req dtos.CreateSupplierBankAccountRequest, userID uint64) (*models.SupplierBankAccount, error) {
	account := &models.SupplierBankAccount{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		SupplierID:    req.SupplierID,
		BankID:        req.BankID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		AccountType:   req.AccountType,
		CurrencyCode:  req.CurrencyCode,
		IsDefault:     req.IsDefault,
	}
	err := db.DB.Create(account).Error
	return account, err
}

func (s *SupplierService) GetSupplierContacts(supplierID string) ([]dtos.SupplierContactResponse, error) {
	var contacts []dtos.SupplierContactResponse
	err := db.DB.Model(&models.SupplierContact{}).
		Where("supplier_id = ?", supplierID).
		Find(&contacts).Error
	return contacts, err
}

func (s *SupplierService) GetSupplierBankAccounts(supplierID string) ([]dtos.SupplierBankAccountResponse, error) {
	var results []dtos.SupplierBankAccountResponse
	query := `
		SELECT sba.*, b.name as bank_name 
		FROM supplier_bank_accounts sba LEFT JOIN banks b ON sba.bank_id = b.id
		WHERE sba.supplier_id = ? AND sba.deleted_at IS NULL`
	err := db.DB.Raw(query, supplierID).Scan(&results).Error
	return results, err
}
