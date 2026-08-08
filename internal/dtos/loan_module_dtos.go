package dtos

import "time"

// --- Loan Products ---

type CreateLoanProductRequest struct {
	ProductCode           string   `json:"product_code" validate:"required"`
	ProductName           string   `json:"product_name" validate:"required"`
	Description           string   `json:"description"`
	InterestMethod        string   `json:"interest_method" validate:"required"`
	InterestRate          float64  `json:"interest_rate" validate:"required"`
	MinAmount             float64  `json:"min_amount"`
	MaxAmount             float64  `json:"max_amount" validate:"required"`
	RepaymentPeriodMonths int      `json:"repayment_period_months" validate:"required"`
	GracePeriodDays       int      `json:"grace_period_days"`
	ProcessingFeeRate     float64  `json:"processing_fee_rate"`
	InsuranceFeeRate      float64  `json:"insurance_fee_rate"`
	FeeCollectionMethod   string   `json:"fee_collection_method"`
	PenaltyRateDaily      float64  `json:"penalty_rate_daily"`
	PenaltyRateMonthly    float64  `json:"penalty_rate_monthly"`
	PenaltyFixedAmount    float64  `json:"penalty_fixed_amount"`
	AllocationPriority    []string `json:"allocation_priority"`
	RecoveryPriority      int      `json:"recovery_priority"`
	IsActive              *bool    `json:"is_active"`
}

type UpdateLoanProductRequest struct {
	ProductName           *string  `json:"product_name"`
	Description           *string  `json:"description"`
	InterestMethod        *string  `json:"interest_method"`
	InterestRate          *float64 `json:"interest_rate"`
	MinAmount             *float64 `json:"min_amount"`
	MaxAmount             *float64 `json:"max_amount"`
	RepaymentPeriodMonths *int     `json:"repayment_period_months"`
	GracePeriodDays       *int     `json:"grace_period_days"`
	ProcessingFeeRate     *float64 `json:"processing_fee_rate"`
	InsuranceFeeRate      *float64 `json:"insurance_fee_rate"`
	FeeCollectionMethod   *string  `json:"fee_collection_method"`
	PenaltyRateDaily      *float64 `json:"penalty_rate_daily"`
	PenaltyRateMonthly    *float64 `json:"penalty_rate_monthly"`
	PenaltyFixedAmount    *float64 `json:"penalty_fixed_amount"`
	AllocationPriority    []string `json:"allocation_priority"`
	RecoveryPriority      *int     `json:"recovery_priority"`
	IsActive              *bool    `json:"is_active"`
}

// --- Disbursement Channels ---

type CreateDisbursementChannelRequest struct {
	ChannelCode         string `json:"channel_code" validate:"required"`
	ChannelName         string `json:"channel_name" validate:"required"`
	Description         string `json:"description"`
	ChannelType         string `json:"channel_type" validate:"required"`
	Provider            string `json:"provider"`
	IsAsync             *bool  `json:"is_async"`
	RequiresExternalRef *bool  `json:"requires_external_ref"`
	GLRuleType          string `json:"gl_rule_type"`
	ConfigSchema        string `json:"config_schema"`
	SortOrder           *int   `json:"sort_order"`
	IsActive            *bool  `json:"is_active"`
}

type UpdateDisbursementChannelRequest struct {
	ChannelName         *string `json:"channel_name"`
	Description         *string `json:"description"`
	ChannelType         *string `json:"channel_type"`
	Provider            *string `json:"provider"`
	IsAsync             *bool   `json:"is_async"`
	RequiresExternalRef *bool   `json:"requires_external_ref"`
	GLRuleType          *string `json:"gl_rule_type"`
	ConfigSchema        *string `json:"config_schema"`
	SortOrder           *int    `json:"sort_order"`
	IsActive            *bool   `json:"is_active"`
}

type DisbursementChannelConfigItem struct {
	Key       string `json:"key"`
	Label     string `json:"label"`
	Value     string `json:"value"`
	HasValue  bool   `json:"has_value"`
	Source    string `json:"source"`
	IsSecret  bool   `json:"is_secret"`
	InputType string `json:"input_type"`
}

type DisbursementChannelConfigUpdateItem struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value"`
}

type UpdateDisbursementChannelConfigRequest struct {
	Items []DisbursementChannelConfigUpdateItem `json:"items" validate:"required,dive"`
}

// --- Applications ---

type LoanGuarantorInput struct {
	MemberID         uint64  `json:"member_id" validate:"required"`
	GuaranteedAmount float64 `json:"guaranteed_amount" validate:"required"`
	Relationship     string  `json:"relationship"`
}

type CreateLoanApplicationRequest struct {
	MemberID                 uint64               `json:"member_id" validate:"required"`
	LoanProductID            uint64               `json:"loan_product_id" validate:"required"`
	RequestedAmount          float64              `json:"requested_amount" validate:"required"`
	RequestedTermMonths      int                  `json:"requested_term_months" validate:"required"`
	Purpose                  string               `json:"purpose"`
	ExpectedDisbursementDate string               `json:"expected_disbursement_date"`
	Guarantors               []LoanGuarantorInput `json:"guarantors"`
}

type ApproveLoanApplicationRequest struct {
	ApprovalRole       string  `json:"approval_role" validate:"required"`
	ApprovedAmount     float64 `json:"approved_amount"`
	ApprovedTermMonths int     `json:"approved_term_months"`
	Comments           string  `json:"comments"`
	Conditions         string  `json:"conditions"`
	FinalApproval      bool    `json:"final_approval"`
}

type RejectLoanApplicationRequest struct {
	Comments string `json:"comments" validate:"required"`
}

// --- Contracts ---

type LoanDisbursementQuote struct {
	ApprovedAmount  float64 `json:"approved_amount"`
	ProcessingFee   float64 `json:"processing_fee"`
	InsuranceFee    float64 `json:"insurance_fee"`
	TotalFees       float64 `json:"total_fees"`
	NetDisbursement float64 `json:"net_disbursement"`
	AlreadyDisbursed float64 `json:"already_disbursed"`
	FeesDeducted    float64 `json:"fees_deducted"`
	RemainingNet    float64 `json:"remaining_net"`
}

type DisburseLoanRequest struct {
	Amount              float64                `json:"amount"`
	ChannelCode         string                 `json:"channel_code"`
	ChannelPayload      map[string]interface{} `json:"channel_payload"`
	Method              string                 `json:"method"` // legacy; used if channel_code omitted
	Reference           string                 `json:"reference" validate:"required"`
	IdempotencyKey      string                 `json:"idempotency_key"`
	DisbursementDate    string                 `json:"disbursement_date"`
	PostToGL            bool                   `json:"post_to_gl"`
	CreateMilkDeduction bool                   `json:"create_milk_deduction"`
	InstallmentAmount   float64                `json:"installment_amount"`
	Notes               string                 `json:"notes"`
}

type ConfirmLoanDisbursementRequest struct {
	Success           bool   `json:"success"`
	ExternalReference string `json:"external_reference"`
	FailureCode       string `json:"failure_code"`
	FailureMessage    string `json:"failure_message"`
	PostToGL          bool   `json:"post_to_gl"`
	CreateMilkDeduction bool `json:"create_milk_deduction"`
	OTP               string `json:"otp"`
	ScaIntentID       string `json:"sca_intent_id"`
}

type RecordLoanRepaymentRequest struct {
	Amount          float64 `json:"amount" validate:"required"`
	Channel         string  `json:"channel" validate:"required"`
	Reference       string  `json:"reference" validate:"required"`
	IdempotencyKey  string  `json:"idempotency_key"`
	PaymentDate     string  `json:"payment_date"`
	PostToGL        bool    `json:"post_to_gl"`
	MemberPayrollID *uint64 `json:"member_payroll_id"`
	MemberPayslipID *uint64 `json:"member_payslip_id"`
	Notes           string  `json:"notes"`
}

type RestructureLoanRequest struct {
	RestructureType string  `json:"restructure_type" validate:"required"`
	NewTermMonths   *int    `json:"new_term_months"`
	NewRate         *float64 `json:"new_rate"`
	Reason          string  `json:"reason" validate:"required"`
	EffectiveDate   string  `json:"effective_date" validate:"required"`
}

type WriteOffLoanRequest struct {
	Amount       float64 `json:"amount" validate:"required"`
	Reason       string  `json:"reason" validate:"required"`
	WriteOffDate string  `json:"write_off_date"`
	PostToGL     bool    `json:"post_to_gl"`
}

type AccrueInterestRequest struct {
	AsOfDate        string   `json:"as_of_date" validate:"required"`
	LoanContractIDs []uint64 `json:"loan_contract_ids"`
	PostToGL        bool     `json:"post_to_gl"`
}

type SettleLoanRequest struct {
	SettlementDate string `json:"settlement_date"`
	QuoteOnly      bool   `json:"quote_only"`
	PostToGL       bool   `json:"post_to_gl"`
}

// --- Reports ---

type LoanPortfolioReport struct {
	TotalContracts       int64   `json:"total_contracts"`
	ActiveContracts      int64   `json:"active_contracts"`
	OverdueContracts     int64   `json:"overdue_contracts"`
	TotalOutstanding     float64 `json:"total_outstanding"`
	TotalPrincipal       float64 `json:"total_principal"`
	TotalInterestDue     float64 `json:"total_interest_due"`
	TotalDisbursed       float64 `json:"total_disbursed"`
	TotalRepaid          float64 `json:"total_repaid"`
	AverageInterestRate  float64 `json:"average_interest_rate"`
}

type LoanAgingBucket struct {
	Bucket      string  `json:"bucket"`
	Count       int64   `json:"count"`
	Outstanding float64 `json:"outstanding"`
}

type LoanAgingReport struct {
	AsOfDate time.Time         `json:"as_of_date"`
	Buckets  []LoanAgingBucket `json:"buckets"`
}

type LoanPARReport struct {
	AsOfDate      time.Time `json:"as_of_date"`
	PAR30         float64   `json:"par_30"`
	PAR60         float64   `json:"par_60"`
	PAR90         float64   `json:"par_90"`
	NPLRatio      float64   `json:"npl_ratio"`
	RecoveryRate  float64   `json:"recovery_rate"`
	TotalPortfolio float64  `json:"total_portfolio"`
}

type LoanRegisterRow struct {
	ContractID           uint64     `json:"contract_id"`
	ContractNo           string     `json:"contract_no"`
	MemberID             uint64     `json:"member_id"`
	MemberNo             string     `json:"member_no"`
	MemberName           string     `json:"member_name"`
	ProductName          string     `json:"product_name"`
	ApprovedAmount       float64    `json:"approved_amount"`
	DisbursedAmount      float64    `json:"disbursed_amount"`
	OutstandingPrincipal float64    `json:"outstanding_principal"`
	OutstandingInterest  float64    `json:"outstanding_interest"`
	OutstandingFees      float64    `json:"outstanding_fees"`
	OutstandingPenalties float64    `json:"outstanding_penalties"`
	TotalOutstanding     float64    `json:"total_outstanding"`
	Status               string     `json:"status"`
	InterestRate         float64    `json:"interest_rate"`
	TermMonths           int        `json:"term_months"`
	DaysInArrears        int        `json:"days_in_arrears"`
	DisbursementDate     *time.Time `json:"disbursement_date"`
}

type LoanRegisterReport struct {
	AsOfDate time.Time          `json:"as_of_date"`
	Rows     []LoanRegisterRow  `json:"rows"`
	Total    int64              `json:"total"`
}

type LoanApplicationStatusSummary struct {
	Status         string  `json:"status"`
	Count          int64   `json:"count"`
	TotalRequested float64 `json:"total_requested"`
}

type LoanApplicationsPipelineReport struct {
	AsOfDate          time.Time                      `json:"as_of_date"`
	Items             []LoanApplicationStatusSummary `json:"items"`
	TotalApplications int64                          `json:"total_applications"`
	TotalRequested    float64                        `json:"total_requested"`
}

type LoanDisbursementReportRow struct {
	ID                uint64     `json:"id"`
	ContractNo        string     `json:"contract_no"`
	MemberNo          string     `json:"member_no"`
	MemberName        string     `json:"member_name"`
	ChannelCode       string     `json:"channel_code"`
	Status            string     `json:"status"`
	Amount            float64    `json:"amount"`
	GrossAmount       float64    `json:"gross_amount"`
	FeesDeducted      float64    `json:"fees_deducted"`
	DisbursementDate  time.Time  `json:"disbursement_date"`
	ExternalReference string     `json:"external_reference"`
	Reference         string     `json:"reference"`
}

type LoanDisbursementsReport struct {
	AsOfDate time.Time                   `json:"as_of_date"`
	Rows     []LoanDisbursementReportRow `json:"rows"`
	Total    int64                       `json:"total"`
	TotalAmount float64                  `json:"total_amount"`
}

type LoanSettlementQuote struct {
	ContractID           uint64    `json:"contract_id"`
	SettlementDate       time.Time `json:"settlement_date"`
	OutstandingPrincipal float64   `json:"outstanding_principal"`
	OutstandingInterest  float64   `json:"outstanding_interest"`
	OutstandingFees      float64   `json:"outstanding_fees"`
	OutstandingPenalties float64   `json:"outstanding_penalties"`
	TotalSettlement      float64   `json:"total_settlement"`
}

type LoanModuleStatement struct {
	ContractID     uint64                   `json:"contract_id"`
	ContractNo     string                   `json:"contract_no"`
	MemberID       uint64                   `json:"member_id"`
	OpeningBalance float64                  `json:"opening_balance"`
	ClosingBalance float64                  `json:"closing_balance"`
	Disbursements  []LoanModuleStatementLine      `json:"disbursements"`
	Repayments     []LoanModuleStatementLine      `json:"repayments"`
	Accruals       []LoanModuleStatementLine      `json:"accruals"`
	Schedule       []LoanScheduleLine       `json:"schedule"`
}

type LoanModuleStatementLine struct {
	Date        time.Time `json:"date"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Debit       float64   `json:"debit"`
	Credit      float64   `json:"credit"`
	Balance     float64   `json:"balance"`
}

type LoanScheduleLine struct {
	InstallmentNo int       `json:"installment_no"`
	DueDate       time.Time `json:"due_date"`
	PrincipalDue  float64   `json:"principal_due"`
	InterestDue   float64   `json:"interest_due"`
	TotalDue      float64   `json:"total_due"`
	TotalPaid     float64   `json:"total_paid"`
	Status        string    `json:"status"`
}
