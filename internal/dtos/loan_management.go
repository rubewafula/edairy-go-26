package dtos

import "time"

// Loan Account
type CreateLoanAccountRequest struct {
	MemberID      uint64 `json:"member_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	Status        string `json:"status"`
}

type LoanAccountResponse struct {
	ID            uint64    `json:"ID"`
	MemberID      uint64    `json:"MemberID"`
	MemberName    string    `json:"MemberName"`
	AccountNumber string    `json:"AccountNumber"`
	Balance       float64   `json:"Balance"`
	Status        string    `json:"Status"`
	CreatedAt     time.Time `json:"CreatedAt"`
}

// Loan Callback
type CreateLoanCallbackRequest struct {
	Detail string `json:"detail" validate:"required"`
	LoanID uint64 `json:"loan_id" validate:"required"`
	Type   string `json:"type" validate:"required"`
}

type UpdateLoanCallbackRequest struct {
	Detail string `json:"detail" validate:"required"`
	LoanID uint64 `json:"loan_id" validate:"required"`
	Type   string `json:"type" validate:"required"`
}

type LoanCallbackResponse struct {
	ID        uint64    `json:"ID"`
	LoanID    uint64    `json:"LoanID"`
	Detail    string    `json:"Detail"`
	Type      string    `json:"Type"`
	CreatedAt time.Time `json:"CreatedAt"`
}

// Loan Origination Log
type CreateLoanOriginationLogRequest struct {
	AstraDetail string `json:"astra_detail" validate:"required"`
	SyncAttempt uint64 `json:"sync_attempt"`
}

type LoanOriginationLogResponse struct {
	ID          uint64    `json:"ID"`
	AstraDetail string    `json:"AstraDetail"`
	SyncAttempt uint64    `json:"SyncAttempt"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

// Loan Transaction
type CreateLoanTransactionRequest struct {
	LoanID      uint64  `json:"loan_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=DEBIT CREDIT"`
	Reference   string  `json:"reference" validate:"required"`
	Description string  `json:"description"`
	Date        string  `json:"transaction_date" validate:"required,datetime"`
}

type LoanTransactionResponse struct {
	ID          uint64    `json:"ID"`
	LoanID      uint64    `json:"LoanID"`
	Amount      float64   `json:"Amount"`
	Type        string    `json:"Type"`
	Reference   string    `json:"Reference"`
	Description string    `json:"Description"`
	Date        time.Time `json:"Date"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

// Member Loan
type CreateMemberLoanRequest struct {
	MemberID     uint64  `json:"member_id" validate:"required"`
	LoanType     string  `json:"loan_type" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	InterestRate float64 `json:"interest_rate"`
	Status       string  `json:"status"`
}

type UpdateMemberLoanRequest struct {
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
	Status       string  `json:"status"`
	DisbursedAt  string  `json:"disbursed_at" validate:"omitempty,datetime"`
}

type MemberLoanResponse struct {
	ID           uint64     `json:"ID"`
	MemberID     uint64     `json:"MemberID"`
	MemberName   string     `json:"MemberName"`
	LoanType     string     `json:"LoanType"`
	Amount       float64    `json:"Amount"`
	InterestRate float64    `json:"InterestRate"`
	Status       string     `json:"Status"`
	DisbursedAt  *time.Time `json:"DisbursedAt"`
	CreatedAt    time.Time  `json:"CreatedAt"`
}
