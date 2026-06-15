package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type AccountTypeService struct{}

func NewAccountTypeService() *AccountTypeService {
	return &AccountTypeService{}
}

func (s *AccountTypeService) CreateAccountType(req dtos.CreateAccountTypeRequest, userID uint64) (*models.AccountType, error) {
	accountType := &models.AccountType{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Name:      req.Name,
	}

	if err := db.DB.Create(accountType).Error; err != nil {
		return nil, err
	}
	return accountType, nil
}

func (s *AccountTypeService) GetAccountTypes(page, limit int) ([]dtos.AccountTypeResponse, int64, error) {
	var results []dtos.AccountTypeResponse
	var total int64

	db.DB.Model(&models.AccountType{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.AccountType{}).
		Select("id, name, created_at, updated_at, created_by, updated_by").
		Where("deleted_at IS NULL").
		Order("id DESC").
		Limit(limit).Offset(offset).
		Scan(&results).Error

	return results, total, err
}

func (s *AccountTypeService) GetAccountType(id string) (*dtos.AccountTypeResponse, error) {
	var result dtos.AccountTypeResponse
	err := db.DB.Model(&models.AccountType{}).
		Select("id, name, created_at, updated_at, created_by, updated_by").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AccountTypeService) UpdateAccountType(id string, req dtos.UpdateAccountTypeRequest, userID uint64) error {
	var accountType models.AccountType
	if err := db.DB.First(&accountType, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"updated_by": userID,
	}

	return db.DB.Model(&accountType).Updates(updates).Error
}

func (s *AccountTypeService) DeleteAccountType(id string, userID uint64) error {
	var accountType models.AccountType
	if err := db.DB.First(&accountType, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&accountType).
		Update("updated_by", userID).
		Delete(&accountType).Error
}
