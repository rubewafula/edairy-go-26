package dtos

type CreateLoanRequest struct {
	MemberID              uint64  `json:"member_id" validate:"required"`
	Amount                float64 `json:"amount" validate:"required,min=1"`
	Interest              float64 `json:"interest" validate:"min=0"`
	TotalPayable          float64 `json:"total_payable"`
	Status                string  `json:"status" validate:"omitempty,oneof=PENDING APPROVED DISBURSED REJECTED"`
	RequestID             string  `json:"request_id"`
	RepaymentAmount       float64 `json:"repayment_amount"`
	WithdrawalRequestUUID string  `json:"withdrawal_request_uuid"`
}

type UpdateLoanRequest struct {
	Amount                float64 `json:"amount" validate:"required,min=1"`
	Interest              float64 `json:"interest" validate:"min=0"`
	TotalPayable          float64 `json:"total_payable"`
	Status                string  `json:"status" validate:"required,oneof=PENDING APPROVED DISBURSED REJECTED"`
	ApprovedAmt           float64 `json:"approved_amount"`
	ProcessedBy           uint64  `json:"processed_by"`
	LoanLimitBy           uint64  `json:"loan_limit_by"`
	CreditLimit           uint64  `json:"credit_limit"`
	ReviewAccepted        bool    `json:"review_accepted"`
	DisbursedAt           string  `json:"disbursed_at" validate:"omitempty,datetime"`
	ProcessedAt           string  `json:"processed_at" validate:"omitempty,datetime"`
	RequestID             string  `json:"request_id"`
	TotalDue              float64 `json:"total_due"`
	RepaymentAmount       float64 `json:"repayment_amount"`
	WithdrawalRequestUUID string  `json:"withdrawal_request_uuid"`
}
