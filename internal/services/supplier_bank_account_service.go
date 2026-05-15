package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type SupplierBankAccountService struct{}

func NewSupplierBankAccountService() *SupplierBankAccountService {
	return &SupplierBankAccountService{}
}

func (s *SupplierBankAccountService) CreateBankAccount(req dtos.CreateSupplierBankAccountRequest, userID uint64) (*models.SupplierBankAccount, error) {
	account := &models.SupplierBankAccount{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		SupplierID:    req.SupplierID,
		BankID:        req.BankID,
		BankBranchID:  req.BankBranchID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		AccountType:   req.AccountType,
		CurrencyCode:  req.CurrencyCode,
		IsDefault:     req.IsDefault,
		Status:        "active", // Default status
	}
	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *SupplierBankAccountService) GetBankAccounts(page, limit int) ([]dtos.SupplierBankAccountResponse, int64, error) {
	var results []dtos.SupplierBankAccountResponse
	var total int64
	db.DB.Model(&models.SupplierBankAccount{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sba.*, b.name as bank_name, bb.name as bank_branch_name
		FROM supplier_bank_accounts sba
		LEFT JOIN banks b ON sba.bank_id = b.id
		LEFT JOIN bank_branches bb ON sba.bank_branch_id = bb.id
		WHERE sba.deleted_at IS NULL
		ORDER BY sba.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplierBankAccountService) GetBankAccount(id string) (*dtos.SupplierBankAccountResponse, error) {
	var result dtos.SupplierBankAccountResponse
	query := `
		SELECT 
			sba.*, b.name as bank_name, bb.name as bank_branch_name
		FROM supplier_bank_accounts sba
		LEFT JOIN banks b ON sba.bank_id = b.id
		LEFT JOIN bank_branches bb ON sba.bank_branch_id = bb.id
		WHERE sba.id = ? AND sba.deleted_at IS NULL
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

func (s *SupplierBankAccountService) UpdateBankAccount(id string, req dtos.UpdateSupplierBankAccountRequest, userID uint64) error {
	var account models.SupplierBankAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"bank_id":         req.BankID,
		"bank_branch_id":  req.BankBranchID,
		"account_name":    req.AccountName,
		"account_number":  req.AccountNumber,
		"account_type":    req.AccountType,
		"currency_code":   req.CurrencyCode,
		"swift_code":      req.SwiftCode,
		"mobile_money_no": req.MobileMoneyNo,
		"is_default":      req.IsDefault,
		"status":          req.Status,
		"updated_by":      userID,
	}

	return db.DB.Model(&account).Updates(updates).Error
}

func (s *SupplierBankAccountService) DeleteBankAccount(id string, userID uint64) error {
	var account models.SupplierBankAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&account).Update("updated_by", userID).Delete(&account).Error
}
