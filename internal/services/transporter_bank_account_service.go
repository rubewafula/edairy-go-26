package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type TransporterBankAccountService struct{}

func NewTransporterBankAccountService() *TransporterBankAccountService {
	return &TransporterBankAccountService{}
}

func (s *TransporterBankAccountService) CreateAccount(req dtos.CreateTransporterBankAccountRequest) (*models.TransporterBankAccount, error) {
	account := &models.TransporterBankAccount{
		TransporterID: req.TransporterID,
		BankID:        req.BankID,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *TransporterBankAccountService) GetAccounts() ([]dtos.TransporterBankAccountResponse, int64, error) {
	var results []dtos.TransporterBankAccountResponse
	var total int64
	db.DB.Model(&models.TransporterBankAccount{}).Count(&total)

	query := `
		SELECT 
			tba.id, tba.transporter_id, t.transporter_no,
			tba.bank_id, b.name AS bank_name,
			tba.account_number, tba.account_name,
			tba.created_at, tba.updated_at
		FROM transporter_bank_accounts tba
		LEFT JOIN transporters t ON tba.transporter_id = t.id
		LEFT JOIN banks b ON tba.bank_id = b.id
		WHERE tba.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *TransporterBankAccountService) GetAccount(id string) (*dtos.TransporterBankAccountResponse, error) {
	var result dtos.TransporterBankAccountResponse
	query := `
		SELECT 
			tba.id, tba.transporter_id, t.transporter_no,
			tba.bank_id, b.name AS bank_name,
			tba.account_number, tba.account_name,
			tba.created_at, tba.updated_at
		FROM transporter_bank_accounts tba
		LEFT JOIN transporters t ON tba.transporter_id = t.id
		LEFT JOIN banks b ON tba.bank_id = b.id
		WHERE tba.id = ? AND tba.deleted_at IS NULL
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

func (s *TransporterBankAccountService) UpdateAccount(id string, req dtos.UpdateTransporterBankAccountRequest) error {
	return db.DB.Model(&models.TransporterBankAccount{}).Where("id = ?", id).Updates(map[string]interface{}{
		"bank_id":        req.BankID,
		"account_number": req.AccountNumber,
		"account_name":   req.AccountName,
	}).Error
}

func (s *TransporterBankAccountService) DeleteAccount(id string) error {
	return db.DB.Delete(&models.TransporterBankAccount{}, id).Error
}
