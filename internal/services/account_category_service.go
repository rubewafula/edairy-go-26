package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type AccountCategoryService struct{}

func NewAccountCategoryService() *AccountCategoryService {
	return &AccountCategoryService{}
}

func (s *AccountCategoryService) CreateAccountCategory(req dtos.CreateAccountCategoryRequest, userID uint64) (*models.AccountCategory, error) {
	category := &models.AccountCategory{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		Name:          req.Name,
		Description:   req.Description,
		AccountTypeID: req.AccountTypeID,
	}

	if err := db.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *AccountCategoryService) GetAccountCategories(page, limit int) ([]dtos.AccountCategoryResponse, int64, error) {
	var results []dtos.AccountCategoryResponse
	var total int64

	db.DB.Model(&models.AccountCategory{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ac.id, ac.name, ac.description, ac.account_type_id, 
			at.name AS account_type_name,
			ac.created_at, ac.updated_at, ac.created_by, ac.updated_by
		FROM account_categories ac
		LEFT JOIN account_types at ON ac.account_type_id = at.id
		WHERE ac.deleted_at IS NULL
		ORDER BY ac.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *AccountCategoryService) GetAccountCategory(id string) (*dtos.AccountCategoryResponse, error) {
	var result dtos.AccountCategoryResponse
	query := `
		SELECT 
			ac.id, ac.name, ac.description, ac.account_type_id, 
			at.name AS account_type_name,
			ac.created_at, ac.updated_at, ac.created_by, ac.updated_by
		FROM account_categories ac
		LEFT JOIN account_types at ON ac.account_type_id = at.id
		WHERE ac.id = ? AND ac.deleted_at IS NULL
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

func (s *AccountCategoryService) UpdateAccountCategory(id string, req dtos.UpdateAccountCategoryRequest, userID uint64) error {
	var category models.AccountCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":            req.Name,
		"description":     req.Description,
		"account_type_id": req.AccountTypeID,
		"updated_by":      userID,
	}

	return db.DB.Model(&category).Updates(updates).Error
}

func (s *AccountCategoryService) DeleteAccountCategory(id string, userID uint64) error {
	return db.DB.Model(&models.AccountCategory{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.AccountCategory{}).Error
}
