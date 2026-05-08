package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type BankService struct{}

func NewBankService() *BankService {
	return &BankService{}
}

func (s *BankService) CreateBank(req dtos.CreateBankRequest) (*models.Bank, error) {
	bank := &models.Bank{
		Name:        req.Name,
		SwiftCode:   req.SwiftCode,
		Description: req.Description,
	}

	if err := db.DB.Create(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (s *BankService) GetBanks() ([]models.Bank, int64, error) {
	var banks []models.Bank
	var total int64
	db.DB.Model(&models.Bank{}).Count(&total)
	err := db.DB.Find(&banks).Error
	return banks, total, err
}

func (s *BankService) GetBank(id string) (*models.Bank, error) {
	var bank models.Bank
	if err := db.DB.First(&bank, id).Error; err != nil {
		return nil, err
	}
	return &bank, nil
}

func (s *BankService) UpdateBank(id string, req dtos.UpdateBankRequest) error {
	var bank models.Bank
	if err := db.DB.First(&bank, id).Error; err != nil {
		return err
	}

	bank.Name = req.Name
	bank.SwiftCode = req.SwiftCode
	bank.Description = req.Description

	return db.DB.Save(&bank).Error
}

func (s *BankService) DeleteBank(id string) error {
	return db.DB.Delete(&models.Bank{}, id).Error
}
