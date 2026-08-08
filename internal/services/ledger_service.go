package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

var (
	ErrLedgerImmutability = errors.New("posted ledger entries cannot be modified; use reversal")
	ErrLedgerImbalance    = errors.New("ledger entries are not balanced")
	ErrDuplicateReference = errors.New("transaction reference already exists")
	ErrPostingRuleMissing = errors.New("posting rule not found")
)

// LedgerService is the single entry point for double-entry GL posting.
type LedgerService struct{}

func NewLedgerService() *LedgerService {
	return &LedgerService{}
}

// PostFromRuleRequest posts a balanced journal from transaction_posting_rules.
type PostFromRuleRequest struct {
	Tx              *gorm.DB
	UserID          uint64
	Reference       string
	IdempotencyKey  string
	TransactionName string
	HeaderType      string // transactions.transaction_type (e.g. SHARE, LOAN)
	RuleType        string // transaction_posting_rules.transaction_type
	Amount          float64
	TransactionDate time.Time
	Description     string
	Status          string
	SwapDebitCredit bool
}

// PostFromRuleResult contains the created header and GL lines.
type PostFromRuleResult struct {
	Transaction *models.Transaction
	Entries     []models.GeneralLedgerEntry
}

// PostManualPairRequest posts explicit debit/credit accounts without a rule row.
type PostManualPairRequest struct {
	Tx              *gorm.DB
	UserID          uint64
	Reference       string
	IdempotencyKey  string
	TransactionName string
	HeaderType      string
	DebitAccountID  uint64
	CreditAccountID uint64
	DebitSubID      *uint64
	CreditSubID     *uint64
	Amount          float64
	TransactionDate time.Time
	Description     string
	Status          string
}

func (s *LedgerService) roundAmount(v float64) float64 {
	return math.Round(v*100) / 100
}

func (s *LedgerService) isPostedStatus(status string) bool {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "POSTED", "APPROVED", "CONFIRMED", "COMPLETED":
		return true
	default:
		return false
	}
}

// FindExistingByReference returns an existing transaction for idempotent retries.
func (s *LedgerService) FindExistingByReference(tx *gorm.DB, reference string) (*models.Transaction, error) {
	if reference == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var t models.Transaction
	err := tx.Where("reference = ? AND deleted_at IS NULL", reference).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *LedgerService) FindExistingByIdempotencyKey(tx *gorm.DB, key string) (*models.Transaction, error) {
	if key == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var t models.Transaction
	err := tx.Where("idempotency_key = ? AND deleted_at IS NULL", key).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *LedgerService) loadRule(tx *gorm.DB, ruleType string) (*models.TransactionPostingRule, error) {
	var rule models.TransactionPostingRule
	if err := tx.Where("transaction_type = ? AND deleted_at IS NULL", ruleType).
		Order("id ASC").First(&rule).Error; err != nil {
		return nil, fmt.Errorf("%w: %s", ErrPostingRuleMissing, ruleType)
	}
	return &rule, nil
}

func (s *LedgerService) createHeader(tx *gorm.DB, userID uint64, reference, idempotencyKey, name, headerType string, date time.Time, description, status string, reversalOf *uint64) (*models.Transaction, error) {
	if status == "" {
		status = "POSTED"
	}
	header := &models.Transaction{
		BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		Reference:       reference,
		TransactionName: name,
		TransactionType: headerType,
		TransactionDate: date,
		Description:     description,
		Status:          status,
		ReversalOfID:    reversalOf,
	}
	if idempotencyKey != "" {
		header.IdempotencyKey = idempotencyKey
	}
	if err := tx.Create(header).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, ErrDuplicateReference
		}
		return nil, err
	}
	return header, nil
}

func (s *LedgerService) insertBalancedPair(
	tx *gorm.DB,
	userID uint64,
	transactionID uint64,
	debitAccountID, creditAccountID uint64,
	debitSubID, creditSubID *uint64,
	amount float64,
	transactionDate time.Time,
	description string,
	isReversal bool,
) ([]models.GeneralLedgerEntry, error) {
	amount = s.roundAmount(amount)
	if amount <= 0 {
		return nil, fmt.Errorf("posting amount must be positive")
	}

	now := time.Now()
	debit := models.GeneralLedgerEntry{
		BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		TransactionID:   transactionID,
		AccountID:       debitAccountID,
		SubAccountID:    debitSubID,
		Debit:           amount,
		Credit:          0,
		TransactionDate: transactionDate,
		Description:     description,
		PostedAt:        now,
		IsReversal:      isReversal,
	}
	credit := models.GeneralLedgerEntry{
		BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		TransactionID:   transactionID,
		AccountID:       creditAccountID,
		SubAccountID:    creditSubID,
		Debit:           0,
		Credit:          amount,
		TransactionDate: transactionDate,
		Description:     description,
		PostedAt:        now,
		IsReversal:      isReversal,
	}
	if err := tx.Create(&debit).Error; err != nil {
		return nil, err
	}
	if err := tx.Create(&credit).Error; err != nil {
		return nil, err
	}
	return []models.GeneralLedgerEntry{debit, credit}, nil
}

// ValidateTransactionBalance ensures debits equal credits for a transaction.
func (s *LedgerService) ValidateTransactionBalance(tx *gorm.DB, transactionID uint64) error {
	var totals struct {
		TotalDebit  float64
		TotalCredit float64
	}
	err := tx.Model(&models.GeneralLedgerEntry{}).
		Select("COALESCE(SUM(debit),0) as total_debit, COALESCE(SUM(credit),0) as total_credit").
		Where("transaction_id = ? AND deleted_at IS NULL", transactionID).
		Scan(&totals).Error
	if err != nil {
		return err
	}
	if s.roundAmount(totals.TotalDebit) != s.roundAmount(totals.TotalCredit) {
		return fmt.Errorf("%w: debits=%.2f credits=%.2f", ErrLedgerImbalance, totals.TotalDebit, totals.TotalCredit)
	}
	return nil
}

// PostFromRule creates a transaction header and balanced GL entries from a posting rule.
func (s *LedgerService) PostFromRule(req PostFromRuleRequest) (*PostFromRuleResult, error) {
	if req.Tx == nil {
		return nil, errors.New("database transaction is required")
	}
	if req.Reference == "" {
		return nil, errors.New("reference is required")
	}

	if existing, err := s.FindExistingByReference(req.Tx, req.Reference); err == nil {
		var entries []models.GeneralLedgerEntry
		req.Tx.Where("transaction_id = ? AND deleted_at IS NULL", existing.ID).Find(&entries)
		return &PostFromRuleResult{Transaction: existing, Entries: entries}, nil
	}
	if req.IdempotencyKey != "" {
		if existing, err := s.FindExistingByIdempotencyKey(req.Tx, req.IdempotencyKey); err == nil {
			var entries []models.GeneralLedgerEntry
			req.Tx.Where("transaction_id = ? AND deleted_at IS NULL", existing.ID).Find(&entries)
			return &PostFromRuleResult{Transaction: existing, Entries: entries}, nil
		}
	}

	rule, err := s.loadRule(req.Tx, req.RuleType)
	if err != nil {
		return nil, err
	}

	debitAccountID := rule.DebitAccountID
	creditAccountID := rule.CreditAccountID
	debitSubID := rule.DebitSubAccountID
	creditSubID := rule.CreditSubAccountID
	if req.SwapDebitCredit {
		debitAccountID, creditAccountID = creditAccountID, debitAccountID
		debitSubID, creditSubID = creditSubID, debitSubID
	}

	desc := req.Description
	if desc == "" && rule.Description != "" {
		desc = rule.Description
	}

	header, err := s.createHeader(req.Tx, req.UserID, req.Reference, req.IdempotencyKey,
		req.TransactionName, req.HeaderType, req.TransactionDate, desc, req.Status, nil)
	if err != nil {
		return nil, err
	}

	entries, err := s.insertBalancedPair(req.Tx, req.UserID, header.ID,
		debitAccountID, creditAccountID, debitSubID, creditSubID,
		req.Amount, req.TransactionDate, desc, false)
	if err != nil {
		return nil, err
	}
	if err := s.ValidateTransactionBalance(req.Tx, header.ID); err != nil {
		return nil, err
	}
	return &PostFromRuleResult{Transaction: header, Entries: entries}, nil
}

// PostManualPair posts with explicit account IDs (payroll multi-line helper).
func (s *LedgerService) PostManualPair(req PostManualPairRequest) (*PostFromRuleResult, error) {
	if req.Tx == nil {
		return nil, errors.New("database transaction is required")
	}
	if req.Reference == "" {
		return nil, errors.New("reference is required")
	}

	if existing, err := s.FindExistingByReference(req.Tx, req.Reference); err == nil {
		var entries []models.GeneralLedgerEntry
		req.Tx.Where("transaction_id = ? AND deleted_at IS NULL", existing.ID).Find(&entries)
		return &PostFromRuleResult{Transaction: existing, Entries: entries}, nil
	}

	header, err := s.createHeader(req.Tx, req.UserID, req.Reference, req.IdempotencyKey,
		req.TransactionName, req.HeaderType, req.TransactionDate, req.Description, req.Status, nil)
	if err != nil {
		return nil, err
	}

	entries, err := s.insertBalancedPair(req.Tx, req.UserID, header.ID,
		req.DebitAccountID, req.CreditAccountID, req.DebitSubID, req.CreditSubID,
		req.Amount, req.TransactionDate, req.Description, false)
	if err != nil {
		return nil, err
	}
	if err := s.ValidateTransactionBalance(req.Tx, header.ID); err != nil {
		return nil, err
	}
	return &PostFromRuleResult{Transaction: header, Entries: entries}, nil
}

// AppendManualPair adds another balanced pair to an existing transaction header.
func (s *LedgerService) AppendManualPair(tx *gorm.DB, userID, transactionID uint64,
	debitAccountID, creditAccountID uint64, debitSubID, creditSubID *uint64,
	amount float64, transactionDate time.Time, description string) error {
	_, err := s.insertBalancedPair(tx, userID, transactionID,
		debitAccountID, creditAccountID, debitSubID, creditSubID,
		amount, transactionDate, description, false)
	if err != nil {
		return err
	}
	return s.ValidateTransactionBalance(tx, transactionID)
}

// GuardGLMutation blocks in-place edits/deletes on posted transactions.
func (s *LedgerService) GuardGLMutation(tx *gorm.DB, transactionID uint64) error {
	var header models.Transaction
	if err := tx.First(&header, transactionID).Error; err != nil {
		return err
	}
	if s.isPostedStatus(header.Status) {
		return ErrLedgerImmutability
	}
	return nil
}

// ReverseTransaction creates an offsetting transaction and marks the original reversed.
func (s *LedgerService) ReverseTransaction(tx *gorm.DB, transactionID, userID uint64, reason string) (*models.Transaction, error) {
	var original models.Transaction
	if err := tx.First(&original, transactionID).Error; err != nil {
		return nil, err
	}
	if strings.EqualFold(original.Status, "REVERSED") {
		return &original, nil
	}

	var entries []models.GeneralLedgerEntry
	if err := tx.Where("transaction_id = ? AND deleted_at IS NULL", transactionID).Find(&entries).Error; err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no ledger entries to reverse for transaction %d", transactionID)
	}

	revRef := original.Reference + "-REV"
	revHeader, err := s.createHeader(tx, userID, revRef, "",
		"Reversal: "+original.TransactionName, original.TransactionType,
		time.Now(), reason, "REVERSED", &transactionID)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		rev := models.GeneralLedgerEntry{
			BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			TransactionID:   revHeader.ID,
			AccountID:       e.AccountID,
			SubAccountID:    e.SubAccountID,
			Debit:           e.Credit,
			Credit:          e.Debit,
			TransactionDate: time.Now(),
			Description:     "Reversal: " + e.Description,
			PostedAt:        time.Now(),
			IsReversal:      true,
			ReversesEntryID: &e.ID,
		}
		if err := tx.Create(&rev).Error; err != nil {
			return nil, err
		}
	}

	if err := s.ValidateTransactionBalance(tx, revHeader.ID); err != nil {
		return nil, err
	}

	if err := tx.Model(&original).Updates(map[string]interface{}{
		"status":     "REVERSED",
		"updated_by": userID,
	}).Error; err != nil {
		return nil, err
	}

	return revHeader, nil
}

// Global ledger service instance for convenience.
var defaultLedger = NewLedgerService()

func Ledger() *LedgerService {
	return defaultLedger
}

// EnsureLedgerSchema adds optional columns/indexes when missing (safe for existing DBs).
func EnsureLedgerSchema() error {
	migrator := db.DB.Migrator()
	if !migrator.HasColumn(&models.Transaction{}, "idempotency_key") {
		if err := migrator.AddColumn(&models.Transaction{}, "IdempotencyKey"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn(&models.Transaction{}, "reversal_of_id") {
		if err := migrator.AddColumn(&models.Transaction{}, "ReversalOfID"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn(&models.GeneralLedgerEntry{}, "posted_at") {
		if err := migrator.AddColumn(&models.GeneralLedgerEntry{}, "PostedAt"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn(&models.GeneralLedgerEntry{}, "is_reversal") {
		if err := migrator.AddColumn(&models.GeneralLedgerEntry{}, "IsReversal"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn(&models.GeneralLedgerEntry{}, "reverses_entry_id") {
		if err := migrator.AddColumn(&models.GeneralLedgerEntry{}, "ReversesEntryID"); err != nil {
			return err
		}
	}
	return db.DB.AutoMigrate(&models.FinancialPeriod{}, &models.BankAccount{})
}
