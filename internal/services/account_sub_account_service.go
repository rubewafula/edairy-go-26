package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type AccountSubAccountService struct{}

func NewAccountSubAccountService() *AccountSubAccountService {
	return &AccountSubAccountService{}
}

func (s *AccountSubAccountService) CreateAccountSubAccount(req dtos.CreateAccountSubAccountRequest, userID uint64) (*models.AccountSubAccount, error) {
	subAccount := &models.AccountSubAccount{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		SubAccountCode: req.SubAccountCode,
		Name:           req.Name,
		Description:    req.Description,
		AccountID:      req.AccountID,
	}

	if err := db.DB.Create(subAccount).Error; err != nil {
		return nil, err
	}
	return subAccount, nil
}

func (s *AccountSubAccountService) GetAccountSubAccounts(accountID string, page, limit int) ([]dtos.AccountSubAccountResponse, int64, error) {
	var results []dtos.AccountSubAccountResponse
	var total int64

	queryBuilder := db.DB.Model(&models.AccountSubAccount{})
	if accountID != "" {
		queryBuilder = queryBuilder.Where("account_id = ?", accountID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			asa.id, asa.sub_account_code, asa.name, asa.description, asa.account_id,
			a.name as account_name, asa.created_at, asa.updated_at, asa.created_by, asa.updated_by
		FROM account_sub_accounts asa
		LEFT JOIN accounts a ON asa.account_id = a.id
		WHERE asa.deleted_at IS NULL
		AND (? = '' OR asa.account_id = ?)
		ORDER BY asa.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, accountID, accountID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *AccountSubAccountService) GetAccountSubAccount(id string) (*dtos.AccountSubAccountResponse, error) {
	var result dtos.AccountSubAccountResponse
	query := `
		SELECT 
			asa.id, asa.sub_account_code, asa.name, asa.description, asa.account_id,
			a.name as account_name, asa.created_at, asa.updated_at, asa.created_by, asa.updated_by
		FROM account_sub_accounts asa
		LEFT JOIN accounts a ON asa.account_id = a.id
		WHERE asa.id = ? AND asa.deleted_at IS NULL
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

func (s *AccountSubAccountService) UpdateAccountSubAccount(id string, req dtos.UpdateAccountSubAccountRequest, userID uint64) error {
	var subAccount models.AccountSubAccount
	if err := db.DB.First(&subAccount, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"sub_account_code": req.SubAccountCode,
		"name":             req.Name,
		"description":      req.Description,
		"account_id":       req.AccountID,
		"updated_by":       userID,
	}

	return db.DB.Model(&subAccount).Updates(updates).Error
}

func (s *AccountSubAccountService) DeleteAccountSubAccount(id string, userID uint64) error {
	return db.DB.Model(&models.AccountSubAccount{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.AccountSubAccount{}).Error
}
