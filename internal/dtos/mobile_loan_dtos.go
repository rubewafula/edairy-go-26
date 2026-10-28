package dtos

import "time"

type MobileLoanRequest struct {
	MemberID      uint64  `json:"member_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	LoanProductID uint64  `json:"loan_product_id"`
	Purpose       string  `json:"purpose"`
}

type VerifyMobileLoanOTPRequest struct {
	OTP string `json:"otp" validate:"required,len=6,numeric"`
}

type MobileLoanRequestResponse struct {
	ApplicationID     uint64    `json:"application_id"`
	ApplicationNo     string    `json:"application_no"`
	MemberID          uint64    `json:"member_id"`
	Amount            float64   `json:"amount"`
	TotalRepayable    float64   `json:"total_repayable"`
	TermMonths        int       `json:"term_months"`
	LoanProductID     uint64    `json:"loan_product_id"`
	MaskedPhone       string    `json:"masked_phone"`
	OTPExpiresAt      time.Time `json:"otp_expires_at"`
	ResendAvailableAt time.Time `json:"resend_available_at"`
	CreditLimit       float64   `json:"credit_limit"`
}

type MobileLoanStatusResponse struct {
	ApplicationID         uint64     `json:"application_id"`
	ApplicationNo         string     `json:"application_no"`
	MemberID              uint64     `json:"member_id"`
	Amount                float64    `json:"amount"`
	ApplicationStatus     string     `json:"application_status"`
	OTPStatus             string     `json:"otp_status"`
	OTPExpiresAt          *time.Time `json:"otp_expires_at,omitempty"`
	OTPAttemptsLeft       int        `json:"otp_attempts_left"`
	ContractID            *uint64    `json:"contract_id,omitempty"`
	ContractNo            string     `json:"contract_no,omitempty"`
	ContractStatus        string     `json:"contract_status,omitempty"`
	TotalRepayable        float64    `json:"total_repayable"`
	DisbursementID        *uint64    `json:"disbursement_id,omitempty"`
	DisbursementStatus    string     `json:"disbursement_status,omitempty"`
	DisbursementReference string     `json:"disbursement_reference,omitempty"`
	MpesaReference        string     `json:"mpesa_reference,omitempty"`
	FailureMessage        string     `json:"failure_message,omitempty"`
	CreditLimit           float64    `json:"credit_limit"`
}

// MobileLoanListItem is one row of GET /api/loan-module/mobile/loan-requests?member_id=…
type MobileLoanListItem struct {
	MobileLoanStatusResponse
	// DisplayStatus is a short label for the app: Awaiting OTP, Expired, Queued, Sending, Disbursed,
	// Cancelled, Rejected, Under review or Approved.
	DisplayStatus string    `json:"display_status"`
	StatusMessage string    `json:"status_message"`
	CreatedAt     time.Time `json:"created_at"`
}
