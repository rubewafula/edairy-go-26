package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type ShareTransactionService struct{}

func NewShareTransactionService() *ShareTransactionService {
	return &ShareTransactionService{}
}

func (s *ShareTransactionService) CreateShareTransaction(req dtos.CreateShareTransactionRequest) (*models.ShareTransaction, error) {
	transaction := &models.ShareTransaction{
		TransactionID:   req.TransactionID,
		ShareAccountID:  req.ShareAccountID,
		MemberID:        req.MemberID,
		TransactionType: req.TransactionType,
		ShareUnits:      req.ShareUnits,
		UnitPrice:       req.UnitPrice,
		Debit:           req.Debit,
		Credit:          req.Credit,
		BalanceAfter:    req.BalanceAfter,
		TransactionDate: utils.ParseDate(req.TransactionDate),
	}

	if err := db.DB.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *ShareTransactionService) GetShareTransactions() ([]dtos.ShareTransactionResponse, int64, error) {
	var results []dtos.ShareTransactionResponse
	var total int64
	db.DB.Model(&models.ShareTransaction{}).Count(&total)

	query := `
		SELECT 
			st.id, st.transaction_id, st.share_account_id, st.member_id, 
			m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			st.transaction_type, st.share_units, st.unit_price, st.debit, st.credit, 
			st.balance_after, st.transaction_date, st.created_at, st.updated_at
		FROM share_transactions st
		LEFT JOIN member_registrations m ON st.member_id = m.id
		WHERE st.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ShareTransactionService) GetShareTransaction(id string) (*dtos.ShareTransactionResponse, error) {
	var result dtos.ShareTransactionResponse
	query := `
		SELECT 
			st.id, st.transaction_id, st.share_account_id, st.member_id, 
			m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			st.transaction_type, st.share_units, st.unit_price, st.debit, st.credit, 
			st.balance_after, st.transaction_date, st.created_at, st.updated_at
		FROM share_transactions st
		LEFT JOIN member_registrations m ON st.member_id = m.id
		WHERE st.id = ? AND st.deleted_at IS NULL
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

func (s *ShareTransactionService) UpdateShareTransaction(id string, req dtos.UpdateShareTransactionRequest) error {
	var transaction models.ShareTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return err
	}

	transaction.TransactionID = req.TransactionID
	transaction.ShareAccountID = req.ShareAccountID
	transaction.TransactionType = req.TransactionType
	transaction.ShareUnits = req.ShareUnits
	transaction.UnitPrice = req.UnitPrice
	transaction.Debit = req.Debit
	transaction.Credit = req.Credit
	transaction.BalanceAfter = req.BalanceAfter
	transaction.TransactionDate = utils.ParseDate(req.TransactionDate)

	return db.DB.Save(&transaction).Error
}

func (s *ShareTransactionService) DeleteShareTransaction(id string) error {
	return db.DB.Delete(&models.ShareTransaction{}, id).Error
}
