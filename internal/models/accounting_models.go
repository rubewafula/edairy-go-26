package models

import "time"

// Accounting
type AccountCategory struct {
	BaseModel
	Name          string `gorm:"column:name"`
	Description   string `gorm:"column:description"`
	AccountTypeID uint64 `gorm:"column:account_type_id"`
}

type AccountType struct {
	BaseModel
	Name string `gorm:"column:name"`
}

type Account struct {
	BaseModel
	AccountCode       string  `gorm:"uniqueIndex;column:account_code"`
	Name              string  `gorm:"column:name"`
	Description       string  `gorm:"column:description"`
	AccountCategoryID uint64  `gorm:"column:account_category_id"`
	ParentAccountID   uint64  `gorm:"column:parent_account_id"`
	IsPostable        bool    `gorm:"column:is_postable"`
	IsActive          bool    `gorm:"default:true;column:is_active"`
	Balance           float64 `gorm:"column:balance"`
}

type AccountSubAccount struct {
	BaseModel
	SubAccountCode string `gorm:"uniqueIndex;column:sub_account_code"`
	Name           string `gorm:"column:name"`
	Description    string `gorm:"column:description"`
	AccountID      uint64 `gorm:"index;column:account_id"`
}

type Transaction struct {
	BaseModel
	Reference       string    `gorm:"index;column:reference"`
	TransactionName string    `gorm:"column:transaction_name"`
	TransactionType string    `gorm:"column:transaction_type"`
	TransactionDate time.Time `gorm:"index;column:transaction_date"`
	Description     string    `gorm:"column:description"`
	Status          string    `gorm:"index;column:status"`
}

type LedgerEntry struct {
	BaseModel
	TransactionID uint64  `gorm:"index;column:transaction_id"`
	AccountID     uint64  `gorm:"index;column:account_id"`
	SubAccountID  uint64  `gorm:"index;column:sub_account_id"`
	Debit         float64 `gorm:"column:debit"`
	Credit        float64 `gorm:"column:credit"`
}

type WalletType struct {
	BaseModel
	Code        string `gorm:"uniqueIndex;column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

type Wallet struct {
	BaseModel
	WalletID      string  `gorm:"uniqueIndex;column:wallet_id"`
	WalletName    string  `gorm:"column:wallet_name"`
	MemberID      uint64  `gorm:"index;column:member_id"`
	AccountNumber string  `gorm:"column:account_number"`
	Balance       float64 `gorm:"column:balance"`
	UUID          string  `gorm:"uniqueIndex;column:uuid"`
	WalletTypeID  string  `gorm:"column:walletTypeId"`
}

type MoneyTransfer struct {
	BaseModel
	Type     string  `gorm:"column:type"` // mpesa/wallet_transfer
	MemberID uint64  `gorm:"index;column:member_id"`
	Amount   float64 `gorm:"column:amount"`
	Status   string  `gorm:"column:status"` // pending/success/failed
	Remarks  string  `gorm:"column:remarks"`
}

type WalletWithdrawal struct {
	BaseModel
	WithdrawalUUID string `gorm:"uniqueIndex;column:withdrawal_uuid"`
	Status         string `gorm:"column:status"`
	LoanID         uint64 `gorm:"index;column:loan_id"`
	MemberID       uint64 `gorm:"index;column:member_id"`
}

type TransactionPostingRule struct {
	BaseModel
	TransactionType    string  `gorm:"column:transaction_type;not null"`
	DebitAccountID     uint64  `gorm:"column:debit_account_id;not null"`
	DebitSubAccountID  *uint64 `gorm:"column:debit_sub_account_id"`
	CreditAccountID    uint64  `gorm:"column:credit_account_id;not null"`
	CreditSubAccountID *uint64 `gorm:"column:credit_sub_account_id"`
	Description        string  `gorm:"column:description"`
}

func (TransactionPostingRule) TableName() string {
	return "transaction_posting_rules"
}

type GeneralLedgerEntry struct {
	BaseModel
	TransactionID   uint64    `gorm:"column:transaction_id;not null;index"`
	AccountID       uint64    `gorm:"column:account_id;not null;index"`
	SubAccountID    *uint64   `gorm:"column:sub_account_id;index"`
	Debit           float64   `gorm:"column:debit;type:decimal(10,2);not null;default:0.00"`
	Credit          float64   `gorm:"column:credit;type:decimal(10,2);not null;default:0.00"`
	TransactionDate time.Time `gorm:"column:transaction_date;not null"`
	Description     string    `gorm:"column:description"`
}

func (GeneralLedgerEntry) TableName() string {
	return "general_ledger_entries"
}

type StatutoryDeductionConfiguration struct {
	BaseModel
	DeductionID           uint64  `gorm:"column:deduction_id"` // Links to employee_deduction_types.id
	EmployeeDeductionRate float64 `gorm:"column:employee_deduction_rate"`
	EmployerDeductionRate float64 `gorm:"column:employer_deduction_rate"`
	MinAmount             float64 `gorm:"column:min_amount"`
	MaxAmount             float64 `gorm:"column:max_amount"`
	FixedAmount           float64 `gorm:"column:fixed_amount"`
	BandLowerLimitAmount  float64 `gorm:"column:band_lower_limit_amount"`
	BandUpperLimitAmount  float64 `gorm:"column:band_upper_limit_amount"`
	MinApplicableAmount   float64 `gorm:"column:min_applicable_amount"`
	// Add other fields as needed from the schema, e.g., 'is_percentage', 'is_fixed'
}

func (StatutoryDeductionConfiguration) TableName() string {
	return "statutory_deductions_configurations"
}
