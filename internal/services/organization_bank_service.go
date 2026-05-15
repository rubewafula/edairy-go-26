package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationBankService struct{}

func NewOrganizationBankService() *OrganizationBankService {
	return &OrganizationBankService{}
}

func (s *OrganizationBankService) CreateBank(req dtos.CreateOrganizationBankRequest, userID uint64) (*models.OrganizationBank, error) {
	bank := &models.OrganizationBank{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Name:      req.Name,
	}
	if err := db.DB.Create(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (s *OrganizationBankService) GetBanks(page, limit int) ([]dtos.OrganizationBankResponse, int64, error) {
	var results []dtos.OrganizationBankResponse
	var total int64
	db.DB.Model(&models.OrganizationBank{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.OrganizationBank{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *OrganizationBankService) GetBank(id string) (*dtos.OrganizationBankResponse, error) {
	var result dtos.OrganizationBankResponse
	err := db.DB.Model(&models.OrganizationBank{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *OrganizationBankService) UpdateBank(id string, req dtos.UpdateOrganizationBankRequest, userID uint64) error {
	var bank models.OrganizationBank
	if err := db.DB.First(&bank, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"updated_by": userID,
	}

	return db.DB.Model(&bank).Updates(updates).Error
}

func (s *OrganizationBankService) DeleteBank(id string, userID uint64) error {
	var bank models.OrganizationBank
	if err := db.DB.First(&bank, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&bank).Update("updated_by", userID).Delete(&bank).Error
}
