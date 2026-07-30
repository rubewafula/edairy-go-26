package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type FinancialReportService struct{}

func NewFinancialReportService() *FinancialReportService {
	return &FinancialReportService{}
}

func (s *FinancialReportService) parseDateRange(from, to string) (time.Time, time.Time, error) {
	var start, end time.Time
	if from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			return start, end, fmt.Errorf("invalid from date")
		}
		start = t
	}
	if to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			return start, end, fmt.Errorf("invalid to date")
		}
		end = t.Add(24*time.Hour - time.Nanosecond)
	}
	return start, end, nil
}

func (s *FinancialReportService) glDateClause(from, to time.Time) (string, []interface{}) {
	clauses := []string{"gle.deleted_at IS NULL"}
	args := []interface{}{}
	if !from.IsZero() {
		clauses = append(clauses, "gle.transaction_date >= ?")
		args = append(args, from)
	}
	if !to.IsZero() {
		clauses = append(clauses, "gle.transaction_date <= ?")
		args = append(args, to)
	}
	return strings.Join(clauses, " AND "), args
}

func (s *FinancialReportService) GetGeneralLedger(from, to string, accountID uint64, page, limit int) (*dtos.GeneralLedgerResponse, error) {
	start, end, err := s.parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	where, args := s.glDateClause(start, end)
	if accountID > 0 {
		where += " AND gle.account_id = ?"
		args = append(args, accountID)
	}

	var total int64
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM general_ledger_entries gle WHERE %s`, where)
	if err := db.DB.Raw(countQ, args...).Scan(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT gle.id, gle.transaction_id, t.reference, t.transaction_name, t.status,
			gle.account_id, a.account_code, a.name as account_name,
			gle.debit, gle.credit, gle.transaction_date, gle.description, gle.is_reversal
		FROM general_ledger_entries gle
		INNER JOIN accounts a ON a.id = gle.account_id
		LEFT JOIN transactions t ON t.id = gle.transaction_id
		WHERE %s
		ORDER BY gle.transaction_date DESC, gle.id DESC
		LIMIT ? OFFSET ?`, where)
	args = append(args, limit, offset)

	var items []dtos.GeneralLedgerItem
	if err := db.DB.Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, err
	}
	return &dtos.GeneralLedgerResponse{Items: items, Total: total}, nil
}

func (s *FinancialReportService) GetCashFlowStatement(from, to string) (*dtos.CashFlowResponse, error) {
	start, end, err := s.parseDateRange(from, to)
	if err != nil {
		return nil, err
	}
	where, args := s.glDateClause(start, end)

	// Cash accounts: codes 1000, 1050, 1055 and names containing 'Cash' or 'Bank'
	query := fmt.Sprintf(`
		SELECT a.account_code, a.name as account_name,
			SUM(gle.debit) as total_inflow,
			SUM(gle.credit) as total_outflow,
			SUM(gle.debit - gle.credit) as net_change
		FROM general_ledger_entries gle
		INNER JOIN accounts a ON a.id = gle.account_id
		WHERE %s AND (a.account_code IN ('1000','1050','1055') OR a.name LIKE '%%Cash%%' OR a.name LIKE '%%Bank%%')
		GROUP BY a.id, a.account_code, a.name
		ORDER BY a.account_code`, where)

	var items []dtos.CashFlowItem
	if err := db.DB.Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, err
	}

	resp := &dtos.CashFlowResponse{Items: items}
	for _, it := range items {
		resp.TotalInflow += it.TotalInflow
		resp.TotalOutflow += it.TotalOutflow
		resp.NetCashChange += it.NetChange
	}
	return resp, nil
}

func (s *FinancialReportService) GetMemberStatement(memberID uint64, from, to string) (*dtos.MemberStatementResponse, error) {
	start, end, err := s.parseDateRange(from, to)
	if err != nil {
		return nil, err
	}

	resp := &dtos.MemberStatementResponse{MemberID: memberID}

	// Shares
	db.DB.Raw(`SELECT COALESCE(SUM(sp.amount_paid),0) FROM share_payments sp
		WHERE sp.member_id = ? AND sp.deleted_at IS NULL`, memberID).Scan(&resp.TotalShareContributions)

	// Milk payslips (net pay proxy)
	payQ := `SELECT COALESCE(SUM(net_pay),0) FROM member_payslips WHERE member_id = ? AND status = 'approved'`
	payArgs := []interface{}{memberID}
	if !start.IsZero() {
		payQ += " AND created_at >= ?"
		payArgs = append(payArgs, start)
	}
	if !end.IsZero() {
		payQ += " AND created_at <= ?"
		payArgs = append(payArgs, end)
	}
	db.DB.Raw(payQ, payArgs...).Scan(&resp.TotalMilkPayments)

	// Loan balance from loan_transactions
	db.DB.Raw(`SELECT COALESCE(SUM(CASE WHEN type='DEBIT' THEN amount ELSE -amount END),0)
		FROM loan_transactions lt INNER JOIN loans l ON l.id = lt.loan_id
		WHERE l.member_id = ? AND lt.deleted_at IS NULL`, memberID).Scan(&resp.LoanBalance)

	// Share transactions sub-ledger
	ledgerQ := `
		SELECT st.transaction_date, st.transaction_type, st.debit, st.credit, st.balance_after, '' as description
		FROM share_transactions st
		WHERE st.member_id = ? AND st.deleted_at IS NULL`
	ledgerArgs := []interface{}{memberID}
	if !start.IsZero() {
		ledgerQ += " AND st.transaction_date >= ?"
		ledgerArgs = append(ledgerArgs, start)
	}
	if !end.IsZero() {
		ledgerQ += " AND st.transaction_date <= ?"
		ledgerArgs = append(ledgerArgs, end)
	}
	ledgerQ += " ORDER BY st.transaction_date DESC LIMIT 100"
	db.DB.Raw(ledgerQ, ledgerArgs...).Scan(&resp.Transactions)

	return resp, nil
}

func (s *FinancialReportService) GetLoanStatement(loanID uint64, from, to string) (*dtos.LoanStatementResponse, error) {
	start, end, err := s.parseDateRange(from, to)
	if err != nil {
		return nil, err
	}

	var loan models.Loan
	if err := db.DB.First(&loan, loanID).Error; err != nil {
		return nil, err
	}

	resp := &dtos.LoanStatementResponse{
		LoanID:       loanID,
		MemberID:     loan.MemberID,
		Principal:    loan.Amount,
		TotalPayable: loan.TotalPayable,
		Status:       loan.Status,
	}

	q := `SELECT id, reference, type, amount, transaction_date, description
		FROM loan_transactions WHERE loan_id = ? AND deleted_at IS NULL`
	args := []interface{}{loanID}
	if !start.IsZero() {
		q += " AND transaction_date >= ?"
		args = append(args, start)
	}
	if !end.IsZero() {
		q += " AND transaction_date <= ?"
		args = append(args, end)
	}
	q += " ORDER BY transaction_date ASC"
	db.DB.Raw(q, args...).Scan(&resp.Transactions)

	var balance float64
	for _, tx := range resp.Transactions {
		if strings.ToUpper(tx.Type) == "DEBIT" {
			balance += tx.Amount
		} else {
			balance -= tx.Amount
		}
	}
	resp.OutstandingBalance = balance
	return resp, nil
}

type FinancialPeriodService struct{}

func NewFinancialPeriodService() *FinancialPeriodService {
	return &FinancialPeriodService{}
}

func (s *FinancialPeriodService) ListPeriods() ([]models.FinancialPeriod, error) {
	var periods []models.FinancialPeriod
	err := db.DB.Where("deleted_at IS NULL").Order("start_date DESC").Find(&periods).Error
	return periods, err
}

func (s *FinancialPeriodService) ClosePeriod(id, userID uint64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var period models.FinancialPeriod
		if err := tx.First(&period, id).Error; err != nil {
			return err
		}
		if period.Status == "CLOSED" {
			return nil
		}
		now := time.Now()
		return tx.Model(&period).Updates(map[string]interface{}{
			"status":     "CLOSED",
			"closed_at":  now,
			"closed_by":  userID,
			"updated_by": userID,
		}).Error
	})
}

func (s *FinancialPeriodService) CreatePeriod(name string, start, end time.Time, userID uint64) (*models.FinancialPeriod, error) {
	p := &models.FinancialPeriod{
		BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		Name:      name,
		StartDate: start,
		EndDate:   end,
		Status:    "OPEN",
	}
	if err := db.DB.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

// BudgetService (P3 stub — stores budget lines for future variance reporting).
type BudgetService struct{}

func NewBudgetService() *BudgetService {
	return &BudgetService{}
}

func (s *BudgetService) GetBudgetVsActual(_ string, _ string) (*dtos.BudgetVsActualResponse, error) {
	return &dtos.BudgetVsActualResponse{Message: "Budget module not yet configured for this dairy"}, nil
}

// BankReconciliationService (P3 stub).
type BankReconciliationService struct{}

func NewBankReconciliationService() *BankReconciliationService {
	return &BankReconciliationService{}
}

func (s *BankReconciliationService) GetReconciliationStatus() (*dtos.BankReconciliationResponse, error) {
	var count int64
	db.DB.Model(&models.BankAccount{}).Count(&count)
	return &dtos.BankReconciliationResponse{
		BankAccountsConfigured: int(count),
		Message:              "Configure bank_accounts and import statements to reconcile",
	}, nil
}
