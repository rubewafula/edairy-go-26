package models

import "time"

type Loan struct {
	BaseModel
	MemberID              uint64    `gorm:"index;column:member_id"`
	Amount                float64   `gorm:"column:amount"`
	Interest              float64   `gorm:"column:interest"`
	TotalPayable          float64   `gorm:"column:total_payable"`
	Status                string    `gorm:"column:status"`
	ApprovedAmt           float64   `gorm:"column:approved_amount"`
	ProcessedBy           uint64    `gorm:"column:processed_by"`
	LoanLimitBy           uint64    `gorm:"column:loan_limit_by"`
	CreditLimit           uint64    `gorm:"column:credit_limit"`
	ReviewAccepted        bool      `gorm:"column:review_accepted"`
	UUID                  string    `gorm:"uniqueIndex;column:uuid"`
	DisbursedAt           time.Time `gorm:"column:disbursed_at"`
	ProcessedAt           time.Time `gorm:"column:processed_at"`
	RequestID             string    `gorm:"column:request_id"`
	TotalDue              float64   `gorm:"column:total_due"`
	RepaymentAmount       float64   `gorm:"column:repayment_amount"`
	WithdrawalRequestUUID string    `gorm:"column:withdrawal_request_uuid"`
}

type LoanRepayment struct {
	BaseModel
	LoanID uint64    `gorm:"index;column:loan_id"`
	Amount float64   `gorm:"column:amount"`
	Date   time.Time `gorm:"index;column:date"`
}

type LoanCallback struct {
	BaseModel
	Detail string `gorm:"column:detail"`
	LoanID uint64 `gorm:"index;column:loan_id"`
	Type   string `gorm:"column:type"`
}

type LoanOrganizationProfile struct {
	BaseModel
	NextLevel       string `gorm:"column:next_level"`
	AstraID         string `gorm:"column:astra_id"`
	LinkStatus      string `gorm:"column:link_status"`
	UUID            string `gorm:"uniqueIndex;column:uuid"`
	Version         string `gorm:"column:version"`
	ProductID       string `gorm:"column:product_id"`
	CompanyDetailID uint64 `gorm:"column:company_detail_id"`
	ManuallyRatify  bool   `gorm:"column:manually_ratify"`
}

type LoanOriginationCallbackLog struct {
	BaseModel
	AstraDetail string `gorm:"column:astra_detail"`
	SyncAttempt uint64 `gorm:"column:sync_attempt"`
}

type LoanAccount struct {
	BaseModel
	ManuallyRatify bool                   `gorm:"column:manually_ratify;default:0"`
	NextLevel      string                 `gorm:"column:next_level;default:detail"`
	Status         string                 `gorm:"column:status;default:ACTIVE"`
	AstraID        *string                `gorm:"column:astra_id"`
	CreditLimit    *uint64                `gorm:"column:credit_limit"`
	LinkStatus     string                 `gorm:"column:link_status;default:1"`
	LivenessPassed bool                   `gorm:"column:liveness_passed;default:0"`
	AstraRemarks   map[string]interface{} `gorm:"column:astra_remarks;serializer:json"`
	UUID           string                 `gorm:"column:uuid;default:uuid();uniqueIndex"`
	AuthCreated    bool                   `gorm:"column:auth_created;default:0"`
	Locale         string                 `gorm:"column:locale;default:KE"`
	CustomerID     uint64                 `gorm:"column:customer_id;index"`
	CustomerType   string                 `gorm:"column:customer_type;type:enum('MEMBER','TRANSPORTER','EMPLOYEE','VENDOR')"`
}

type LoanTransaction struct {
	BaseModel
	LoanID      uint64    `gorm:"column:loan_id"`
	Amount      float64   `gorm:"column:amount"`
	Type        string    `gorm:"column:type"` // DEBIT, CREDIT
	Reference   string    `gorm:"column:reference;uniqueIndex"`
	Description string    `gorm:"column:description"`
	Date        time.Time `gorm:"column:transaction_date"`
}
