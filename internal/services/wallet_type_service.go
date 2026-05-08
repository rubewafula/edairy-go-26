package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type WalletTypeService struct{}

func NewWalletTypeService() *WalletTypeService {
	return &WalletTypeService{}
}

func (s *WalletTypeService) CreateWalletType(req dtos.CreateWalletTypeRequest) (*models.WalletType, error) {
	walletType := &models.WalletType{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(walletType).Error; err != nil {
		return nil, err
	}
	return walletType, nil
}

func (s *WalletTypeService) GetWalletTypes() ([]models.WalletType, int64, error) {
	var walletTypes []models.WalletType
	var total int64
	db.DB.Model(&models.WalletType{}).Count(&total)
	err := db.DB.Find(&walletTypes).Error
	return walletTypes, total, err
}

func (s *WalletTypeService) GetWalletType(id string) (*models.WalletType, error) {
	var walletType models.WalletType
	if err := db.DB.First(&walletType, id).Error; err != nil {
		return nil, err
	}
	return &walletType, nil
}

func (s *WalletTypeService) UpdateWalletType(id string, req dtos.UpdateWalletTypeRequest) error {
	var walletType models.WalletType
	if err := db.DB.First(&walletType, id).Error; err != nil {
		return err
	}

	walletType.Code = req.Code
	walletType.Name = req.Name
	walletType.Description = req.Description

	return db.DB.Save(&walletType).Error
}

func (s *WalletTypeService) DeleteWalletType(id string) error {
	return db.DB.Delete(&models.WalletType{}, id).Error
}
