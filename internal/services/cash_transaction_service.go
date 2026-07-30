package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type CashTransactionService struct {
	posting *FinancialPostingService
}

func NewCashTransactionService() *CashTransactionService {
	return &CashTransactionService{posting: NewFinancialPostingService()}
}

func (s *CashTransactionService) Create(req dtos.CreateCashTransactionRequest, userID uint64) (*models.CashTransaction, error) {
	transaction := &models.CashTransaction{
		ReferenceNumber:        req.ReferenceNumber,
		TransactionDescription: req.TransactionDescription,
		TransactionType:        req.TransactionType,
		TransactionDate:        req.TransactionDate,
		PaidBy:                 &req.PaidBy,
		TransactionAmount:      &req.TransactionAmount,
		CustomerType:           &req.CustomerType,
		CustomerID:             &req.CustomerID,
		PaymentModeID:          &req.PaymentModeID,
		PaymentType:            &req.PaymentType,
		TransactionID:          &req.TransactionID,
		CreatedBy:              &userID,
		UpdatedBy:              &userID,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		if !req.PostToGL {
			return nil
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(req.TransactionAmount), 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("invalid transaction amount for GL posting")
		}

		date := utils.ParseFlexibleDate(req.TransactionDate)
		idempotencyKey := req.IdempotencyKey
		if idempotencyKey == "" {
			idempotencyKey = req.ReferenceNumber
		}

		var glResult *PostFromRuleResult
		switch strings.ToUpper(req.TransactionType) {
		case "CASH_IN", "IN", "RECEIPT":
			glResult, err = s.posting.PostCashIn(userID, amount, date, req.TransactionDescription, req.ReferenceNumber, idempotencyKey)
		case "CASH_OUT", "OUT", "PAYMENT":
			glResult, err = s.posting.PostCashOut(userID, amount, date, req.TransactionDescription, req.ReferenceNumber, idempotencyKey)
		default:
			return fmt.Errorf("unsupported cash transaction type for GL: %s", req.TransactionType)
		}
		if err != nil {
			return err
		}

		glID := int64(glResult.Transaction.ID)
		transaction.TransactionID = &glID
		return tx.Model(transaction).Update("transaction_id", glID).Error
	})

	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *CashTransactionService) List(page, limit int) ([]models.CashTransaction, int64, error) {
	var transactions []models.CashTransaction
	var total int64
	db.DB.Model(&models.CashTransaction{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&transactions).Error
	return transactions, total, err
}

func (s *CashTransactionService) Get(id string) (*models.CashTransaction, error) {
	var transaction models.CashTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *CashTransactionService) Update(id string, req dtos.CreateCashTransactionRequest, userID uint64) error {
	var transaction models.CashTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return err
	}
	if transaction.TransactionID != nil && *transaction.TransactionID > 0 {
		return ErrLedgerImmutability
	}

	updates := map[string]interface{}{
		"reference_number":        req.ReferenceNumber,
		"transaction_description": req.TransactionDescription,
		"transaction_type":        req.TransactionType,
		"transaction_date":        req.TransactionDate,
		"paid_by":                 req.PaidBy,
		"transaction_amount":      req.TransactionAmount,
		"customer_type":           req.CustomerType,
		"customer_id":             req.CustomerID,
		"payment_mode_id":         req.PaymentModeID,
		"payment_type":            req.PaymentType,
		"transaction_id":          req.TransactionID,
		"updated_by":              userID,
	}

	return db.DB.Model(&transaction).Updates(updates).Error
}

func (s *CashTransactionService) Delete(id string) error {
	var transaction models.CashTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return err
	}
	if transaction.TransactionID != nil && *transaction.TransactionID > 0 {
		return ErrLedgerImmutability
	}
	return db.DB.Delete(&models.CashTransaction{}, id).Error
}
