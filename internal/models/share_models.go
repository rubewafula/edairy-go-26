package models

import "time"

type ShareType struct {
	BaseModel
	ShareCode         string  `gorm:"column:share_code"`
	ShareType         string  `gorm:"column:share_type"`
	Description       string  `gorm:"column:description"`
	Rate              float64 `gorm:"column:rate"`
	Mandatory         int     `gorm:"column:madatory"` // Matches 'madatory' in schema
	HasShareValue     string  `gorm:"column:has_share_value"`
	RepayMethod       string  `gorm:"column:repay_method"`
	CalculatingMethod string  `gorm:"column:calculating_method"`
	ShareValue        float64 `gorm:"column:share_value"`
	DeductionTypeID   uint64  `gorm:"column:deduction_type_id"`
	Priority          int     `gorm:"column:priority"`
}

func (ShareType) TableName() string {
	return "share_types"
}

type ShareDividend struct {
	BaseModel
	DeclarationID int64   `gorm:"column:declaration_id"`
	MemberID      uint64  `gorm:"column:member_id"`
	FiscalYear    int     `gorm:"column:fiscal_year"`
	Period        int     `gorm:"column:period"`
	ShareUnits    float64 `gorm:"column:share_units"`
	Status        string  `gorm:"column:status;default:CALCULATED"`
	RatePerShare  float64 `gorm:"column:rate_per_share"`
	TaxAmount     float64 `gorm:"column:tax_amount"`
	NetAmount     float64 `gorm:"column:net_amount"`
	TransactionID int64   `gorm:"column:transaction_id"`
}

func (ShareDividend) TableName() string {
	return "share_dividends"
}

type SharePayment struct {
	BaseModel
	TransactionID   uint64     `gorm:"column:transaction_id"`
	MemberID        uint64     `gorm:"column:member_id"`
	ShareAccountID  uint64     `gorm:"column:share_account_id"`
	AmountPaid      float64    `gorm:"column:amount_paid"`
	ShareUnits      float64    `gorm:"column:share_units"`
	PaymentModeID   uint64     `gorm:"column:payment_mode_id"`
	Description     string     `gorm:"column:description"`
	Status          string     `gorm:"column:status;default:PENDING"`
	TransactionDate time.Time  `gorm:"column:transaction_date"`
	ApprovedBy      uint64     `gorm:"column:approved_by"`
	DateApproved    *time.Time `gorm:"column:date_approved"`
}

func (SharePayment) TableName() string {
	return "share_payments"
}

type ShareTransaction struct {
	BaseModel
	TransactionID   uint64    `gorm:"column:transaction_id"`
	ShareAccountID  uint64    `gorm:"column:share_account_id"`
	MemberID        uint64    `gorm:"column:member_id"`
	TransactionType string    `gorm:"column:transaction_type"`
	ShareUnits      float64   `gorm:"column:share_units;default:0.0000"`
	UnitPrice       float64   `gorm:"column:unit_price;default:0.00"`
	Debit           float64   `gorm:"column:debit;default:0.00"`
	Credit          float64   `gorm:"column:credit;default:0.00"`
	BalanceAfter    float64   `gorm:"column:balance_after;default:0.00"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
}

func (ShareTransaction) TableName() string {
	return "share_transactions"
}

type ShareTransfer struct {
	BaseModel
	TransactionID   uint64    `gorm:"column:transaction_id"`
	FromMemberID    uint64    `gorm:"column:from_member_id"`
	ToMemberID      uint64    `gorm:"column:to_member_id"`
	ShareUnits      float64   `gorm:"column:share_units"`
	TransferAmount  float64   `gorm:"column:transfer_amount"`
	Status          string    `gorm:"column:status;default:PENDING"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	ApprovedBy      uint64    `gorm:"column:approved_by"`
	DateApproved    time.Time `gorm:"column:date_approved"`
}

func (ShareTransfer) TableName() string {
	return "share_transfers"
}

type ShareAccount struct {
	BaseModel
	MemberID    uint64    `gorm:"column:member_id"`
	ShareTypeID uint64    `gorm:"column:share_type_id"`
	Status      string    `gorm:"column:status;default:ACTIVE"`
	OpenedAt    time.Time `gorm:"column:opened_at;default:CURRENT_TIMESTAMP"`
}

func (ShareAccount) TableName() string {
	return "share_accounts"
}

type DividendDeclaration struct {
	BaseModel
	FiscalYear      int       `gorm:"column:fiscal_year"`
	Period          int       `gorm:"column:period"`
	TotalPool       float64   `gorm:"column:total_pool"`
	RatePerShare    float64   `gorm:"column:rate_per_share"`
	CalculationType string    `gorm:"column:calculation_type"`
	Status          string    `gorm:"column:status;default:DRAFT"`
	ApprovedBy      uint64    `gorm:"column:approved_by"`
	ApprovedAt      time.Time `gorm:"column:approved_at"`
}

func (DividendDeclaration) TableName() string {
	return "dividend_declarations"
}
