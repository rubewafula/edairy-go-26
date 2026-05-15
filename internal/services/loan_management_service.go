package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type LoanManagementService struct{}

func NewLoanManagementService() *LoanManagementService {
	return &LoanManagementService{}
}

// Loan Accounts
func (s *LoanManagementService) CreateAccount(req dtos.CreateLoanAccountRequest, userID uint64) (*models.LoanAccount, error) {
	account := &models.LoanAccount{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		MemberID:      req.MemberID,
		AccountNumber: req.AccountNumber,
		Status:        req.Status,
	}
	if account.Status == "" {
		account.Status = "ACTIVE"
	}
	err := db.DB.Create(account).Error
	return account, err
}

func (s *LoanManagementService) GetAccounts(page, limit int) ([]dtos.LoanAccountResponse, int64, error) {
	var results []dtos.LoanAccountResponse
	var total int64
	db.DB.Model(&models.LoanAccount{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT la.*, CONCAT(m.first_name, ' ', m.last_name) as member_name
		FROM loan_accounts la
		LEFT JOIN member_registrations m ON la.member_id = m.id
		WHERE la.deleted_at IS NULL
		ORDER BY la.id DESC LIMIT ? OFFSET ?`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

// Loan Callbacks
func (s *LoanManagementService) CreateCallback(req dtos.CreateLoanCallbackRequest, userID uint64) (*models.LoanCallback, error) {
	callback := &models.LoanCallback{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Detail:    req.Detail,
		LoanID:    req.LoanID,
		Type:      req.Type,
	}
	err := db.DB.Create(callback).Error
	return callback, err
}

func (s *LoanManagementService) GetCallbacks(page, limit int) ([]dtos.LoanCallbackResponse, int64, error) {
	var results []dtos.LoanCallbackResponse
	var total int64
	db.DB.Model(&models.LoanCallback{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.LoanCallback{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *LoanManagementService) GetCallback(id string) (*dtos.LoanCallbackResponse, error) {
	var result dtos.LoanCallbackResponse
	err := db.DB.Model(&models.LoanCallback{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *LoanManagementService) UpdateCallback(id string, req dtos.UpdateLoanCallbackRequest, userID uint64) error {
	var callback models.LoanCallback
	if err := db.DB.First(&callback, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"detail":     req.Detail,
		"loan_id":    req.LoanID,
		"type":       req.Type,
		"updated_by": userID,
	}

	return db.DB.Model(&callback).Updates(updates).Error
}

func (s *LoanManagementService) DeleteCallback(id string, userID uint64) error {
	var callback models.LoanCallback
	if err := db.DB.First(&callback, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&callback).Update("updated_by", userID).Delete(&callback).Error
}

// Loan Origination Logs
func (s *LoanManagementService) CreateOriginationLog(req dtos.CreateLoanOriginationLogRequest, userID uint64) (*models.LoanOriginationCallbackLog, error) {
	log := &models.LoanOriginationCallbackLog{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		AstraDetail: req.AstraDetail,
		SyncAttempt: req.SyncAttempt,
	}
	err := db.DB.Create(log).Error
	return log, err
}

// Loan Transactions
func (s *LoanManagementService) CreateTransaction(req dtos.CreateLoanTransactionRequest, userID uint64) (*models.LoanTransaction, error) {
	transaction := &models.LoanTransaction{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		LoanID:      req.LoanID,
		Amount:      req.Amount,
		Type:        req.Type,
		Reference:   req.Reference,
		Description: req.Description,
		Date:        utils.ParseDate(req.Date),
	}
	err := db.DB.Create(transaction).Error
	return transaction, err
}

func (s *LoanManagementService) GetTransactions(loanID string) ([]dtos.LoanTransactionResponse, error) {
	var results []dtos.LoanTransactionResponse
	err := db.DB.Model(&models.LoanTransaction{}).Where("loan_id = ?", loanID).Find(&results).Error
	return results, err
}

// Member Loans
func (s *LoanManagementService) CreateMemberLoan(req dtos.CreateMemberLoanRequest, userID uint64) (*models.MemberLoan, error) {
	loan := &models.MemberLoan{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		MemberID:     req.MemberID,
		LoanType:     req.LoanType,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		Status:       req.Status,
	}
	if loan.Status == "" {
		loan.Status = "PENDING"
	}
	err := db.DB.Create(loan).Error
	return loan, err
}

func (s *LoanManagementService) GetMemberLoans(page, limit int) ([]dtos.MemberLoanResponse, int64, error) {
	var results []dtos.MemberLoanResponse
	var total int64
	db.DB.Model(&models.MemberLoan{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT ml.*, CONCAT(m.first_name, ' ', m.last_name) as member_name
		FROM member_loans ml
		LEFT JOIN member_registrations m ON ml.member_id = m.id
		WHERE ml.deleted_at IS NULL
		ORDER BY ml.id DESC LIMIT ? OFFSET ?`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LoanManagementService) UpdateMemberLoan(id string, req dtos.UpdateMemberLoanRequest, userID uint64) error {
	var loan models.MemberLoan
	if err := db.DB.First(&loan, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"amount":        req.Amount,
		"interest_rate": req.InterestRate,
		"status":        req.Status,
		"updated_by":    userID,
	}
	if req.DisbursedAt != "" {
		t := utils.ParseDate(req.DisbursedAt)
		updates["disbursed_at"] = &t
	}

	return db.DB.Model(&loan).Updates(updates).Error
}
