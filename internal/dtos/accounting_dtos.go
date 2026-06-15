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
}
