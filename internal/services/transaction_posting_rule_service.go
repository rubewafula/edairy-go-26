package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type TransactionPostingRuleService struct{}

func NewTransactionPostingRuleService() *TransactionPostingRuleService {
	return &TransactionPostingRuleService{}
}

func (s *TransactionPostingRuleService) CreateTransactionPostingRule(req dtos.CreateTransactionPostingRuleRequest, userID uint64) (*models.TransactionPostingRule, error) {
	rule := &models.TransactionPostingRule{
		BaseModel:          models.BaseModel{CreatedBy: userID},
		TransactionType:    req.TransactionType,
		DebitAccountID:     req.DebitAccountID,
		DebitSubAccountID:  req.DebitSubAccountID,
		CreditAccountID:    req.CreditAccountID,
		CreditSubAccountID: req.CreditSubAccountID,
		Description:        req.Description,
	}

	if err := db.DB.Create(rule).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *TransactionPostingRuleService) GetTransactionPostingRules(page, limit int) ([]dtos.TransactionPostingRuleResponse, int64, error) {
	var results []dtos.TransactionPostingRuleResponse
	var total int64

	db.DB.Model(&models.TransactionPostingRule{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			tpr.id, tpr.transaction_type, tpr.debit_account_id, da.name as debit_account_name,
			tpr.credit_account_id, ca.name as credit_account_name, tpr.description,
			tpr.created_at, tpr.updated_at, tpr.created_by, tpr.updated_by
		FROM transaction_posting_rules tpr
		LEFT JOIN accounts da ON tpr.debit_account_id = da.id
		LEFT JOIN accounts ca ON tpr.credit_account_id = ca.id
		WHERE tpr.deleted_at IS NULL
		ORDER BY tpr.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *TransactionPostingRuleService) GetTransactionPostingRule(id string) (*dtos.TransactionPostingRuleResponse, error) {
	var result dtos.TransactionPostingRuleResponse
	query := `
		SELECT 
			tpr.id, tpr.transaction_type, tpr.debit_account_id, da.name as debit_account_name,
			tpr.credit_account_id, ca.name as credit_account_name, tpr.description,
			tpr.created_at, tpr.updated_at, tpr.created_by, tpr.updated_by
		FROM transaction_posting_rules tpr
		LEFT JOIN accounts da ON tpr.debit_account_id = da.id
		LEFT JOIN accounts ca ON tpr.credit_account_id = ca.id
		WHERE tpr.id = ? AND tpr.deleted_at IS NULL
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

func (s *TransactionPostingRuleService) UpdateTransactionPostingRule(id string, req dtos.UpdateTransactionPostingRuleRequest, userID uint64) error {
	var rule models.TransactionPostingRule
	if err := db.DB.First(&rule, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"transaction_type":      req.TransactionType,
		"debit_account_id":      req.DebitAccountID,
		"debit_sub_account_id":  req.DebitSubAccountID,
		"credit_account_id":     req.CreditAccountID,
		"credit_sub_account_id": req.CreditSubAccountID,
		"description":           req.Description,
		"updated_by":            userID,
	}

	return db.DB.Model(&rule).Updates(updates).Error
}

func (s *TransactionPostingRuleService) DeleteTransactionPostingRule(id string, userID uint64) error {
	return db.DB.Model(&models.TransactionPostingRule{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.TransactionPostingRule{}).Error
}
