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
	db.DB.Model(&models.Bank{}).Count(&total)

	offset := (page - 1) * limit

	err := db.DB.Model(&models.Bank{}).
		Where("deleted_at IS NULL").
		Limit(limit).Offset(offset).Order("id DESC").
		Scan(&banks).Error

	return banks, total, err
}

func (s *BankService) GetBank(id string) (*dtos.BankResponse, error) {
	var bank dtos.BankResponse

	err := db.DB.Model(&models.Bank{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&bank).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
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
