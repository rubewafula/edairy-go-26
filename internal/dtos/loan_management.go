package dtos

import "time"

// LoanAccount DTOs
type CreateLoanAccountRequest struct {
	MemberID      uint64  `json:"member_id" validate:"required"`
	AccountNumber string  `json:"account_number" validate:"required"`
	Balance       float64 `json:"balance"`
	Status        string  `json:"status"`
}

type UpdateLoanAccountRequest struct {
	MemberID      uint64  `json:"member_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Status        string  `json:"status"`
}

type LoanAccountResponse struct {
	ID            uint64    `json:"ID"`
	MemberID      uint64    `json:"MemberID"`
	AccountNumber string    `json:"AccountNumber"`
	Balance       float64   `json:"Balance"`
	Status        string    `json:"Status"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}

// LoanCallback DTOs
type CreateLoanCallbackRequest struct {
	Detail string `json:"detail" validate:"required"`
	LoanID uint64 `json:"loan_id" validate:"required"`
	Type   string `json:"type" validate:"required"`
}

type UpdateLoanCallbackRequest struct {
	Detail string `json:"detail"`
	LoanID uint64 `json:"loan_id"`
	Type   string `json:"type"`
}

type LoanCallbackResponse struct {
	ID        uint64    `json:"ID"`
	Detail    string    `json:"Detail"`
	LoanID    uint64    `json:"LoanID"`
	Type      string    `json:"Type"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// LoanOrganizationProfile DTOs
type CreateLoanOrganizationProfileRequest struct {
	NextLevel       string `json:"next_level"`
	AstraID         string `json:"astra_id"`
	LinkStatus      string `json:"link_status"`
	UUID            string `json:"uuid"`
	Version         string `json:"version"`
	ProductID       string `json:"product_id"`
	CompanyDetailID uint64 `json:"company_detail_id"`
	ManuallyRatify  bool   `json:"manually_ratify"`
}

type UpdateLoanOrganizationProfileRequest struct {
	NextLevel       string `json:"next_level"`
	AstraID         string `json:"astra_id"`
	LinkStatus      string `json:"link_status"`
	UUID            string `json:"uuid"`
	Version         string `json:"version"`
	ProductID       string `json:"product_id"`
	CompanyDetailID uint64 `json:"company_detail_id"`
	ManuallyRatify  bool   `json:"manually_ratify"`
}

type LoanOrganizationProfileResponse struct {
	ID              uint64    `json:"ID"`
	NextLevel       string    `json:"NextLevel"`
	AstraID         string    `json:"AstraID"`
	LinkStatus      string    `json:"LinkStatus"`
	UUID            string    `json:"UUID"`
	Version         string    `json:"Version"`
	ProductID       string    `json:"ProductID"`
	CompanyDetailID uint64    `json:"CompanyDetailID"`
	ManuallyRatify  bool      `json:"ManuallyRatify"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}

// LoanOriginationCallbackLog DTOs
type CreateLoanOriginationCallbackLogRequest struct {
	AstraDetail string `json:"astra_detail" validate:"required"`
	SyncAttempt uint64 `json:"sync_attempt"`
}

type UpdateLoanOriginationCallbackLogRequest struct {
	AstraDetail string `json:"astra_detail"`
	SyncAttempt uint64 `json:"sync_attempt"`
}

type LoanOriginationCallbackLogResponse struct {
	ID          uint64    `json:"ID"`
	AstraDetail string    `json:"AstraDetail"`
	SyncAttempt uint64    `json:"SyncAttempt"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

// LoanTransaction DTOs
type CreateLoanTransactionRequest struct {
	LoanID      uint64  `json:"loan_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,min=0"`
	Type        string  `json:"type" validate:"required,oneof=DEBIT CREDIT"`
	Reference   string  `json:"reference" validate:"required"`
	Description string  `json:"description"`
	Date        string  `json:"date" validate:"required,datetime"`
}

type UpdateLoanTransactionRequest struct {
	LoanID      uint64  `json:"loan_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Reference   string  `json:"reference"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
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
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

// MemberLoan DTOs
type CreateMemberLoanRequest struct {
	MemberID     uint64  `json:"member_id" validate:"required"`
	LoanType     string  `json:"loan_type" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,min=0"`
	InterestRate float64 `json:"interest_rate"`
	Status       string  `json:"status"`
	DisbursedAt  string  `json:"disbursed_at" validate:"omitempty,datetime"`
}

type UpdateMemberLoanRequest struct {
	LoanType     string  `json:"loan_type"`
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
	Status       string  `json:"status"`
	DisbursedAt  string  `json:"disbursed_at"`
}

type MemberLoanResponse struct {
	ID           uint64     `json:"ID"`
	MemberID     uint64     `json:"MemberID"`
	LoanType     string     `json:"LoanType"`
	Amount       float64    `json:"Amount"`
	InterestRate float64    `json:"InterestRate"`
	Status       string     `json:"Status"`
	DisbursedAt  *time.Time `json:"DisbursedAt"`
	CreatedAt    time.Time  `json:"CreatedAt"`
	UpdatedAt    time.Time  `json:"UpdatedAt"`
}
