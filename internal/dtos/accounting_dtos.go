package dtos

import "time"

// Account DTOs
type CreateAccountRequest struct {
	AccountCode       string  `json:"account_code" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	Description       string  `json:"description"`
	AccountCategoryID uint64  `json:"account_category_id" validate:"required"`
	ParentAccountID   *uint64 `json:"parent_account_id"`
	IsPostable        bool    `json:"is_postable"`
	IsActive          bool    `json:"is_active"`
}

type UpdateAccountRequest struct {
	AccountCode       string  `json:"account_code"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	AccountCategoryID uint64  `json:"account_category_id"`
	ParentAccountID   *uint64 `json:"parent_account_id"`
	IsPostable        bool    `json:"is_postable"`
	IsActive          bool    `json:"is_active"`
}

type AccountResponse struct {
	ID                  uint64    `json:"id"`
	AccountCode         string    `json:"account_code"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	AccountCategoryID   uint64    `json:"account_category_id"`
	AccountCategoryName string    `json:"account_category_name"` // Joined field
	ParentAccountID     *uint64   `json:"parent_account_id"`
	ParentAccountName   *string   `json:"parent_account_name"` // Joined field
	IsPostable          bool      `json:"is_postable"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedBy           *uint64   `json:"created_by"`
	UpdatedBy           *uint64   `json:"updated_by"`
}

// AccountCategory DTOs
type CreateAccountCategoryRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	AccountTypeID uint64 `json:"account_type_id" validate:"required"`
}

type UpdateAccountCategoryRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	AccountTypeID uint64 `json:"account_type_id" validate:"required"`
}

type AccountCategoryResponse struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AccountTypeID   uint64    `json:"account_type_id"`
	AccountTypeName string    `json:"account_type_name"` // Joined field
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       *uint64   `json:"created_by"`
	UpdatedBy       *uint64   `json:"updated_by"`
}

// AccountType DTOs
type CreateAccountTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateAccountTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type AccountTypeResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *uint64   `json:"created_by"`
	UpdatedBy *uint64   `json:"updated_by"`
}

// AccountSubAccount DTOs
type CreateAccountSubAccountRequest struct {
	SubAccountCode string `json:"sub_account_code" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	AccountID      uint64 `json:"account_id" validate:"required"`
}

type UpdateAccountSubAccountRequest struct {
	SubAccountCode string `json:"sub_account_code"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	AccountID      uint64 `json:"account_id"`
}

type AccountSubAccountResponse struct {
	ID             uint64    `json:"id"`
	SubAccountCode string    `json:"sub_account_code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	AccountID      uint64    `json:"account_id"`
	AccountName    string    `json:"account_name"` // Joined field
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *uint64   `json:"created_by"`
	UpdatedBy      *uint64   `json:"updated_by"`
}

// TransactionPostingRule DTOs
type CreateTransactionPostingRuleRequest struct {
	TransactionType    string  `json:"transaction_type" validate:"required"`
	DebitAccountID     uint64  `json:"debit_account_id" validate:"required"`
	DebitSubAccountID  *uint64 `json:"debit_sub_account_id"`
	CreditAccountID    uint64  `json:"credit_account_id" validate:"required"`
	CreditSubAccountID *uint64 `json:"credit_sub_account_id"`
	Description        string  `json:"description"`
}

type UpdateTransactionPostingRuleRequest struct {
	TransactionType    string  `json:"transaction_type"`
	DebitAccountID     uint64  `json:"debit_account_id"`
	DebitSubAccountID  *uint64 `json:"debit_sub_account_id"`
	CreditAccountID    uint64  `json:"credit_account_id"`
	CreditSubAccountID *uint64 `json:"credit_sub_account_id"`
	Description        string  `json:"description"`
}

type TransactionPostingRuleResponse struct {
	ID                uint64    `json:"id"`
	TransactionType   string    `json:"transaction_type"`
	DebitAccountID    uint64    `json:"debit_account_id"`
	DebitAccountName  string    `json:"debit_account_name"`
	CreditAccountID   uint64    `json:"credit_account_id"`
	CreditAccountName string    `json:"credit_account_name"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         *uint64   `json:"created_by"`
	UpdatedBy         *uint64   `json:"updated_by"`
}

// Trial Balance DTOs
type TrialBalanceItem struct {
	AccountID   uint64  `json:"account_id"`
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	Balance     float64 `json:"balance"`
}

type TrialBalanceResponse struct {
	Items        []TrialBalanceItem `json:"items"`
	TotalDebits  float64            `json:"total_debits"`
	TotalCredits float64            `json:"total_credits"`
}

// Profit and Loss DTOs
type ProfitLossItem struct {
	AccountID    uint64  `json:"account_id"`
	AccountCode  string  `json:"account_code"`
	AccountName  string  `json:"account_name"`
	CategoryName string  `json:"category_name"`
	TypeName     string  `json:"account_type_name"`
	Amount       float64 `json:"amount"`
}

type ProfitLossResponse struct {
	RevenueItems  []ProfitLossItem `json:"revenue_items"`
	ExpenseItems  []ProfitLossItem `json:"expense_items"`
	TotalRevenue  float64          `json:"total_revenue"`
	TotalExpenses float64          `json:"total_expenses"`
	NetProfit     float64          `json:"net_profit"`
}

// Balance Sheet DTOs
type BalanceSheetItem struct {
	AccountID    uint64  `json:"account_id"`
	AccountCode  string  `json:"account_code"`
	AccountName  string  `json:"account_name"`
	CategoryName string  `json:"category_name"`
	TypeName     string  `json:"account_type_name"`
	Amount       float64 `json:"amount"`
}

type BalanceSheetResponse struct {
	AssetItems             []BalanceSheetItem `json:"asset_items"`
	LiabilityItems         []BalanceSheetItem `json:"liability_items"`
	EquityItems            []BalanceSheetItem `json:"equity_items"`
	TotalAssets            float64            `json:"total_assets"`
	TotalLiabilities       float64            `json:"total_liabilities"`
	TotalEquity            float64            `json:"total_equity"`
	TotalLiabilitiesEquity float64            `json:"total_liabilities_equity"`
	AsOfDate               string             `json:"as_of_date,omitempty"`
	FromDate               string             `json:"from_date,omitempty"`
	ToDate                 string             `json:"to_date,omitempty"`
}

// DateRangeQuery optional filters for financial reports.
type DateRangeQuery struct {
	From string `form:"from" json:"from"`
	To   string `form:"to" json:"to"`
}

type GeneralLedgerItem struct {
	ID              uint64    `json:"id"`
	TransactionID   uint64    `json:"transaction_id"`
	Reference       string    `json:"reference"`
	TransactionName string    `json:"transaction_name"`
	Status          string    `json:"status"`
	AccountID       uint64    `json:"account_id"`
	AccountCode     string    `json:"account_code"`
	AccountName     string    `json:"account_name"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
	IsReversal      bool      `json:"is_reversal"`
}

type GeneralLedgerResponse struct {
	Items []GeneralLedgerItem `json:"items"`
	Total int64               `json:"total"`
}

type CashFlowItem struct {
	AccountCode  string  `json:"account_code"`
	AccountName  string  `json:"account_name"`
	TotalInflow  float64 `json:"total_inflow"`
	TotalOutflow float64 `json:"total_outflow"`
	NetChange    float64 `json:"net_change"`
}

type CashFlowResponse struct {
	Items         []CashFlowItem `json:"items"`
	TotalInflow   float64        `json:"total_inflow"`
	TotalOutflow  float64        `json:"total_outflow"`
	NetCashChange float64        `json:"net_cash_change"`
}

type MemberStatementLine struct {
	TransactionDate time.Time `json:"transaction_date"`
	TransactionType string    `json:"transaction_type"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	BalanceAfter    float64   `json:"balance_after"`
	Description     string    `json:"description"`
}

type MemberStatementResponse struct {
	MemberID               uint64                `json:"member_id"`
	TotalShareContributions float64              `json:"total_share_contributions"`
	TotalMilkPayments      float64               `json:"total_milk_payments"`
	LoanBalance            float64               `json:"loan_balance"`
	Transactions           []MemberStatementLine `json:"transactions"`
}

type LoanStatementLine struct {
	ID              uint64    `json:"id"`
	Reference       string    `json:"reference"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
}

type LoanStatementResponse struct {
	LoanID             uint64              `json:"loan_id"`
	MemberID           uint64              `json:"member_id"`
	Principal          float64             `json:"principal"`
	TotalPayable       float64             `json:"total_payable"`
	OutstandingBalance float64             `json:"outstanding_balance"`
	Status             string              `json:"status"`
	Transactions       []LoanStatementLine `json:"transactions"`
}

type PostFinancialTransactionRequest struct {
	Reference       string  `json:"reference" validate:"required"`
	IdempotencyKey  string  `json:"idempotency_key"`
	RuleType        string  `json:"rule_type" validate:"required"`
	HeaderType      string  `json:"header_type"`
	Amount          float64 `json:"amount" validate:"required"`
	TransactionDate string  `json:"transaction_date"`
	Description     string  `json:"description"`
	SwapDebitCredit bool    `json:"swap_debit_credit"`
}

type ReverseTransactionRequest struct {
	Reason string `json:"reason"`
}

type CreateFinancialPeriodRequest struct {
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type BudgetVsActualResponse struct {
	Message string `json:"message"`
}

type BankReconciliationResponse struct {
	BankAccountsConfigured int    `json:"bank_accounts_configured"`
	Message                string `json:"message"`
}
