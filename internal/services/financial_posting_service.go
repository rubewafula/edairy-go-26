package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

// FinancialPostingService posts domain transactions through the central ledger.
type FinancialPostingService struct {
	ledger *LedgerService
}

func NewFinancialPostingService() *FinancialPostingService {
	return &FinancialPostingService{ledger: NewLedgerService()}
}

type DomainPostRequest struct {
	UserID          uint64
	Reference       string
	IdempotencyKey  string
	Amount          float64
	TransactionDate time.Time
	Description     string
	HeaderType      string
	RuleType        string
	SwapDebitCredit bool
}

func (s *FinancialPostingService) postRule(req DomainPostRequest) (*PostFromRuleResult, error) {
	var result *PostFromRuleResult
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		res, err := s.ledger.PostFromRule(PostFromRuleRequest{
			Tx:              tx,
			UserID:          req.UserID,
			Reference:       req.Reference,
			IdempotencyKey:  req.IdempotencyKey,
			TransactionName: req.RuleType,
			HeaderType:      req.HeaderType,
			RuleType:        req.RuleType,
			Amount:          req.Amount,
			TransactionDate: req.TransactionDate,
			Description:     req.Description,
			Status:          "POSTED",
			SwapDebitCredit: req.SwapDebitCredit,
		})
		if err != nil {
			return err
		}
		result = res
		return nil
	})
	return result, err
}

func (s *FinancialPostingService) PostLoanDisbursement(userID, loanID uint64, amount float64, date time.Time, desc string, idempotencyKey string) (*PostFromRuleResult, error) {
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: fmt.Sprintf("LOAN-DISB-%d-%d", loanID, date.Unix()),
		IdempotencyKey: idempotencyKey, Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "LOAN", RuleType: "LOAN_DISBURSEMENT",
	})
}

// LoanNetDisbursementAmounts books gross principal to member loans while paying net cash and recognizing fees.
type LoanNetDisbursementAmounts struct {
	GrossPrincipal float64
	NetCash        float64
	ProcessingFee  float64
	InsuranceFee   float64
}

// PostLoanNetDisbursement posts DR Member Loans (gross) / CR Cash (net) / CR Fee Income (fees).
func (s *FinancialPostingService) PostLoanNetDisbursement(userID uint64, reference, idempotencyKey, contractNo string, amounts LoanNetDisbursementAmounts, date time.Time) (*PostFromRuleResult, error) {
	var result *PostFromRuleResult
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if reference == "" {
			return fmt.Errorf("reference is required")
		}
		if existing, err := s.ledger.FindExistingByReference(tx, reference); err == nil {
			var entries []models.GeneralLedgerEntry
			tx.Where("transaction_id = ? AND deleted_at IS NULL", existing.ID).Find(&entries)
			result = &PostFromRuleResult{Transaction: existing, Entries: entries}
			return nil
		}
		if idempotencyKey != "" {
			if existing, err := s.ledger.FindExistingByIdempotencyKey(tx, idempotencyKey); err == nil {
				var entries []models.GeneralLedgerEntry
				tx.Where("transaction_id = ? AND deleted_at IS NULL", existing.ID).Find(&entries)
				result = &PostFromRuleResult{Transaction: existing, Entries: entries}
				return nil
			}
		}

		disburseRule, err := s.ledger.loadRule(tx, "LOAN_DISBURSEMENT")
		if err != nil {
			return err
		}
		feeRule, err := s.ledger.loadRule(tx, "LOAN_PROCESSING_FEE_CAPITALIZED")
		if err != nil {
			return err
		}
		loanAccountID := disburseRule.DebitAccountID
		cashAccountID := disburseRule.CreditAccountID
		feeIncomeAccountID := feeRule.CreditAccountID

		desc := fmt.Sprintf("Loan net disbursement %s", contractNo)
		header, err := s.ledger.createHeader(tx, userID, reference, idempotencyKey,
			"LOAN_NET_DISBURSEMENT", "LOAN", date, desc, "POSTED", nil)
		if err != nil {
			return err
		}

		if amounts.NetCash > 0 {
			if _, err := s.ledger.insertBalancedPair(tx, userID, header.ID,
				loanAccountID, cashAccountID, disburseRule.DebitSubAccountID, disburseRule.CreditSubAccountID,
				amounts.NetCash, date, desc+" (cash)", false); err != nil {
				return err
			}
		}
		if amounts.ProcessingFee > 0 {
			if _, err := s.ledger.insertBalancedPair(tx, userID, header.ID,
				loanAccountID, feeIncomeAccountID, disburseRule.DebitSubAccountID, feeRule.CreditSubAccountID,
				amounts.ProcessingFee, date, desc+" (processing fee)", false); err != nil {
				return err
			}
		}
		if amounts.InsuranceFee > 0 {
			if _, err := s.ledger.insertBalancedPair(tx, userID, header.ID,
				loanAccountID, feeIncomeAccountID, disburseRule.DebitSubAccountID, feeRule.CreditSubAccountID,
				amounts.InsuranceFee, date, desc+" (insurance fee)", false); err != nil {
				return err
			}
		}
		if err := s.ledger.ValidateTransactionBalance(tx, header.ID); err != nil {
			return err
		}
		var entries []models.GeneralLedgerEntry
		tx.Where("transaction_id = ? AND deleted_at IS NULL", header.ID).Find(&entries)
		result = &PostFromRuleResult{Transaction: header, Entries: entries}
		return nil
	})
	return result, err
}

func (s *FinancialPostingService) PostLoanRepayment(userID, loanID uint64, amount float64, date time.Time, desc string, idempotencyKey string) (*PostFromRuleResult, error) {
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: fmt.Sprintf("LOAN-REP-%d-%d", loanID, date.Unix()),
		IdempotencyKey: idempotencyKey, Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "LOAN", RuleType: "LOAN_REPAYMENT",
	})
}

func (s *FinancialPostingService) PostCashIn(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("CASH-IN-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "CASH", RuleType: "CASH_IN_PAYMENT",
	})
}

func (s *FinancialPostingService) PostCashOut(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("CASH-OUT-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "CASH", RuleType: "CASH_OUT_PAYMENT",
	})
}

func (s *FinancialPostingService) PostMemberSavingsContribution(userID, memberID uint64, amount float64, date time.Time, desc, idempotencyKey string) (*PostFromRuleResult, error) {
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: fmt.Sprintf("SAV-IN-%d-%s", memberID, date.Format("20060102")),
		IdempotencyKey: idempotencyKey, Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "SAVINGS", RuleType: "MEMBER_SAVINGS_CONTRIBUTION",
	})
}

func (s *FinancialPostingService) PostMemberSavingsWithdrawal(userID, memberID uint64, amount float64, date time.Time, desc, idempotencyKey string) (*PostFromRuleResult, error) {
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: fmt.Sprintf("SAV-OUT-%d-%s", memberID, date.Format("20060102")),
		IdempotencyKey: idempotencyKey, Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "SAVINGS", RuleType: "MEMBER_SAVINGS_WITHDRAWAL",
	})
}

func (s *FinancialPostingService) PostLocalMilkSale(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("MILK-SALE-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "MILK", RuleType: "LOCAL_MILK_SALES",
	})
}

func (s *FinancialPostingService) PostVendorPayment(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("VENDOR-PAY-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "AP", RuleType: "VENDOR_PAYMENTS",
	})
}

func (s *FinancialPostingService) PostGoodReceivedOnCredit(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("GRN-CR-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "INVENTORY", RuleType: "GOOD_RECEIVED_ON_CREDIT",
	})
}

func (s *FinancialPostingService) PostDividendDeclaration(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("DIV-DECL-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "DIVIDEND", RuleType: "DIVIDEND_DECLARATION",
	})
}

func (s *FinancialPostingService) PostDividendPayment(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("DIV-PAY-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "DIVIDEND", RuleType: "DIVIDEND_PAYMENT",
	})
}

func (s *FinancialPostingService) PostBadDebtProvision(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("BDP-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "LOAN", RuleType: "BAD_DEBT_PROVISION",
	})
}

func (s *FinancialPostingService) PostBadDebtWriteOff(userID uint64, amount float64, date time.Time, desc, ref, idempotencyKey string) (*PostFromRuleResult, error) {
	if ref == "" {
		ref = fmt.Sprintf("BDWO-%d", time.Now().UnixNano())
	}
	return s.postRule(DomainPostRequest{
		UserID: userID, Reference: ref, IdempotencyKey: idempotencyKey,
		Amount: amount, TransactionDate: date, Description: desc,
		HeaderType: "LOAN", RuleType: "BAD_DEBT_WRITE_OFF",
	})
}

// PostFinancialTransaction is the generic API-facing posting endpoint.
func (s *FinancialPostingService) PostFinancialTransaction(req dtos.PostFinancialTransactionRequest, userID uint64) (*PostFromRuleResult, error) {
	date := utils.ParseFlexibleDate(req.TransactionDate)
	if date.IsZero() {
		date = time.Now()
	}
	return s.postRule(DomainPostRequest{
		UserID:          userID,
		Reference:       req.Reference,
		IdempotencyKey:  req.IdempotencyKey,
		Amount:          req.Amount,
		TransactionDate: date,
		Description:     req.Description,
		HeaderType:      req.HeaderType,
		RuleType:        req.RuleType,
		SwapDebitCredit: req.SwapDebitCredit,
	})
}

func (s *FinancialPostingService) ReverseByTransactionID(transactionID, userID uint64, reason string) (*models.Transaction, error) {
	var rev *models.Transaction
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		res, err := s.ledger.ReverseTransaction(tx, transactionID, userID, reason)
		if err != nil {
			return err
		}
		rev = res
		return nil
	})
	return rev, err
}
