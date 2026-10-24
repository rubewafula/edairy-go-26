package models

import "time"

const (
	LoanAppSourceBackoffice = "BACKOFFICE"
	LoanAppSourceMobile     = "MOBILE"
)

const (
	LoanOTPStatusPending  = "PENDING"
	LoanOTPStatusVerified = "VERIFIED"
	LoanOTPStatusExpired  = "EXPIRED"
	LoanOTPStatusLocked   = "LOCKED"
)

// Credit limit change actions recorded in credit_limit_change_log.action.
const (
	CreditActionMilkConfirmed      = "MILK_CONFIRMED"
	CreditActionMilkUnconfirmed    = "MILK_UNCONFIRMED"
	CreditActionMilkEntryChanged   = "MILK_ENTRY_CHANGED"
	CreditActionMilkRejected       = "MILK_REJECTED"
	CreditActionMilkRejectReversed = "MILK_REJECT_REVERSED"
	CreditActionLoanReserved       = "LOAN_RESERVED"
	CreditActionLoanDisbursed      = "LOAN_DISBURSED"
	CreditActionLoanReleased       = "LOAN_RELEASED"
	CreditActionLoanRepayment      = "LOAN_REPAYMENT"
	CreditActionLoanSettled        = "LOAN_SETTLED"
	CreditActionLoanWrittenOff     = "LOAN_WRITTEN_OFF"
	CreditActionLoanRestructured   = "LOAN_RESTRUCTURED"
	CreditActionLoanInterest       = "LOAN_INTEREST"
	CreditActionLoanPenalty        = "LOAN_PENALTY"
	CreditActionPayrollPaid        = "PAYROLL_PAID"
	CreditActionDeductionChanged   = "DEDUCTION_CHANGED"
	CreditActionRecalculated       = "RECALCULATED"
	CreditActionSeeded             = "SEEDED"
)

type LoanApplicationOTP struct {
	BaseModel
	LoanApplicationID uint64     `gorm:"column:loan_application_id;index;not null" json:"loan_application_id"`
	MemberID          uint64     `gorm:"column:member_id;index;not null" json:"member_id"`
	Phone             string     `gorm:"column:phone;size:20;not null" json:"-"`
	OTPHash           string     `gorm:"column:otp_hash;size:128;not null" json:"-"`
	ExpiresAt         time.Time  `gorm:"column:expires_at;not null" json:"expires_at"`
	Attempts          int        `gorm:"column:attempts;default:0" json:"attempts"`
	MaxAttempts       int        `gorm:"column:max_attempts;default:5" json:"max_attempts"`
	ResendCount       int        `gorm:"column:resend_count;default:0" json:"resend_count"`
	LastSentAt        time.Time  `gorm:"column:last_sent_at;not null" json:"last_sent_at"`
	VerifiedAt        *time.Time `gorm:"column:verified_at" json:"verified_at"`
	Status            string     `gorm:"column:status;size:16;not null;default:PENDING" json:"status"`
}

func (LoanApplicationOTP) TableName() string { return "loan_application_otps" }
