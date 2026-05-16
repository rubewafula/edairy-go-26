package dtos

import "time"

// LoanAccount DTOs
type CreateLoanAccountRequest struct {
	CustomerID     uint64                 `json:"customer_id" validate:"required"`
	CustomerType   string                 `json:"customer_type" validate:"required,oneof=MEMBER TRANSPORTER EMPLOYEE VENDOR"`
	ManuallyRatify bool                   `json:"manually_ratify"`
	NextLevel      string                 `json:"next_level"`
	Status         string                 `json:"status"`
	AstraID        string                 `json:"astra_id"`
	CreditLimit    uint64                 `json:"credit_limit"`
	LinkStatus     string                 `json:"link_status"`
	LivenessPassed bool                   `json:"liveness_passed"`
	AstraRemarks   map[string]interface{} `json:"astra_remarks"`
	AuthCreated    bool                   `json:"auth_created"`
	Locale         string                 `json:"locale"`
}

type UpdateLoanAccountRequest struct {
	ManuallyRatify *bool                  `json:"manually_ratify"`
	NextLevel      string                 `json:"next_level"`
	Status         string                 `json:"status"`
	AstraID        string                 `json:"astra_id"`
	CreditLimit    uint64                 `json:"credit_limit"`
	LinkStatus     string                 `json:"link_status"`
	LivenessPassed *bool                  `json:"liveness_passed"`
	AstraRemarks   map[string]interface{} `json:"astra_remarks"`
	AuthCreated    *bool                  `json:"auth_created"`
	Locale         string                 `json:"locale"`
}

type LoanAccountResponse struct {
	ID             uint64                 `json:"id"`
	CustomerID     uint64                 `json:"customer_id"`
	CustomerType   string                 `json:"customer_type"`
	ManuallyRatify bool                   `json:"manually_ratify"`
	NextLevel      string                 `json:"next_level"`
	Status         string                 `json:"status"`
	AstraID        *string                `json:"astra_id"`
	CreditLimit    *uint64                `json:"credit_limit"`
	LinkStatus     string                 `json:"link_status"`
	LivenessPassed bool                   `json:"liveness_passed"`
	AstraRemarks   map[string]interface{} `json:"astra_remarks"`
	UUID           string                 `json:"uuid"`
	AuthCreated    bool                   `json:"auth_created"`
	Locale         string                 `json:"locale"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
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
	ID              uint64    `json:"id"`
	NextLevel       string    `json:"next_level"`
	AstraID         string    `json:"astra_id"`
	LinkStatus      string    `json:"link_status"`
	UUID            string    `json:"uuid"`
	Version         string    `json:"version"`
	ProductID       string    `json:"product_id"`
	CompanyDetailID uint64    `json:"company_detail_id"`
	ManuallyRatify  bool      `json:"manually_ratify"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
	ID          uint64    `json:"id"`
	AstraDetail string    `json:"astra_detail"`
	SyncAttempt uint64    `json:"sync_attempt"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	ID          uint64    `json:"id"`
	LoanID      uint64    `json:"loan_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	ID           uint64     `json:"id"`
	MemberID     uint64     `json:"member_id"`
	LoanType     string     `json:"loan_type"`
	Amount       float64    `json:"amount"`
	InterestRate float64    `json:"interest_rate"`
	Status       string     `json:"status"`
	DisbursedAt  *time.Time `json:"disbursed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
