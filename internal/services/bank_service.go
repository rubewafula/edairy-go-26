package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type BankService struct{}

func NewBankService() *BankService {
	return &BankService{}
}

func (s *BankService) CreateBank(req dtos.CreateBankRequest) (*models.Bank, error) {
	bank := &models.Bank{
		BankName: req.Name,
		BankCode: req.SwiftCode,
	}

	if err := db.DB.Create(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (s *BankService) GetBanks(page, limit int) ([]dtos.BankResponse, int64, error) {
	var banks []dtos.BankResponse
	var total int64
	db.DB.Model(&models.Bank{}).Where("deleted_at IS NULL").Count(&total)

	offset := (page - 1) * limit

	query := `
		SELECT id, bank_name, bank_code, created_at, updated_at
		FROM banks
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&banks).Error

	return banks, total, err
}

func (s *BankService) GetBank(id string) (*dtos.BankResponse, error) {
	var bank dtos.BankResponse

	query := `
		SELECT id, bank_name, bank_code, created_at, updated_at
		FROM banks
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&bank).Error

	if err != nil {
		return nil, err
	}

	if bank.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &bank, nil
}

func (s *BankService) UpdateBank(id string, req dtos.UpdateBankRequest) error {
	var bank models.Bank
	if err := db.DB.First(&bank, id).Error; err != nil {
		return err
	}

	bank.BankName = req.Name
	bank.BankCode = req.SwiftCode

	return db.DB.Save(&bank).Error
}

func (s *BankService) DeleteBank(id string) error {
	return db.DB.Delete(&models.Bank{}, id).Error
}
