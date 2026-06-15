package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type CashTransactionService struct{}

func NewCashTransactionService() *CashTransactionService {
	return &CashTransactionService{}
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

	if err := db.DB.Create(transaction).Error; err != nil {
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
	return db.DB.Delete(&models.CashTransaction{}, id).Error
}
