package models

import (
	"time"

	"gorm.io/datatypes"
)

// --- Loan Products ---

type LoanProduct struct {
	BaseModel
	ProductCode           string         `gorm:"column:product_code;uniqueIndex;size:32;not null" json:"product_code"`
	ProductName           string         `gorm:"column:product_name;not null" json:"product_name"`
	Description           string         `gorm:"column:description;type:text" json:"description"`
	InterestMethod        string         `gorm:"column:interest_method;not null;default:REDUCING_BALANCE" json:"interest_method"`
	InterestRate          float64        `gorm:"column:interest_rate;type:decimal(8,4);not null" json:"interest_rate"`
	MinAmount             float64        `gorm:"column:min_amount;type:decimal(15,2);not null" json:"min_amount"`
	MaxAmount             float64        `gorm:"column:max_amount;type:decimal(15,2);not null" json:"max_amount"`
	RepaymentPeriodMonths int            `gorm:"column:repayment_period_months;not null" json:"repayment_period_months"`
	GracePeriodDays       int            `gorm:"column:grace_period_days;not null;default:0" json:"grace_period_days"`
	ProcessingFeeRate     float64        `gorm:"column:processing_fee_rate;type:decimal(8,4)" json:"processing_fee_rate"`
	InsuranceFeeRate      float64        `gorm:"column:insurance_fee_rate;type:decimal(8,4)" json:"insurance_fee_rate"`
	FeeCollectionMethod   string         `gorm:"column:fee_collection_method;size:32;not null;default:DEDUCT_FROM_PROCEEDS" json:"fee_collection_method"`
	PenaltyRateDaily      float64        `gorm:"column:penalty_rate_daily;type:decimal(8,4)" json:"penalty_rate_daily"`
	PenaltyRateMonthly    float64        `gorm:"column:penalty_rate_monthly;type:decimal(8,4)" json:"penalty_rate_monthly"`
	PenaltyFixedAmount    float64        `gorm:"column:penalty_fixed_amount;type:decimal(15,2)" json:"penalty_fixed_amount"`
	AllocationPriority    datatypes.JSON `gorm:"column:allocation_priority;type:json" json:"allocation_priority"`
	RecoveryPriority      int            `gorm:"column:recovery_priority;default:1" json:"recovery_priority"`
	EligibilityRules      datatypes.JSON `gorm:"column:eligibility_rules;type:json" json:"eligibility_rules"`
	ApprovalWorkflow      datatypes.JSON `gorm:"column:approval_workflow;type:json" json:"approval_workflow"`
	IsActive              bool           `gorm:"column:is_active;default:1" json:"is_active"`
}

func (LoanProduct) TableName() string { return "loan_products" }

// --- Applications ---

type LoanApplication struct {
	BaseModel
	ApplicationNo            string     `gorm:"column:application_no;uniqueIndex;size:64;not null" json:"application_no"`
	MemberID                 uint64     `gorm:"column:member_id;index;not null" json:"member_id"`
	LoanProductID            uint64     `gorm:"column:loan_product_id;index;not null" json:"loan_product_id"`
	RequestedAmount          float64    `gorm:"column:requested_amount;type:decimal(15,2);not null" json:"requested_amount"`
	RequestedTermMonths      int        `gorm:"column:requested_term_months;not null" json:"requested_term_months"`
	Purpose                  string     `gorm:"column:purpose;size:500" json:"purpose"`
	Status                   string     `gorm:"column:status;not null;default:DRAFT" json:"status"`
	ApprovedAmount           *float64   `gorm:"column:approved_amount;type:decimal(15,2)" json:"approved_amount"`
	ApprovedTermMonths       *int       `gorm:"column:approved_term_months" json:"approved_term_months"`
	DateApplied              time.Time  `gorm:"column:date_applied;not null" json:"date_applied"`
	ExpectedDisbursementDate *time.Time `gorm:"column:expected_disbursement_date;type:date" json:"expected_disbursement_date"`
	RejectedReason           string     `gorm:"column:rejected_reason;type:text" json:"rejected_reason"`
	LoanContractID           *uint64    `gorm:"column:loan_contract_id" json:"loan_contract_id"`
	MemberNo                 string     `gorm:"-" json:"member_no,omitempty"`
	MemberName               string     `gorm:"-" json:"member_name,omitempty"`
	LoanProduct              *LoanProduct `gorm:"foreignKey:LoanProductID" json:"loan_product,omitempty"`
	Guarantors               []LoanGuarantor `gorm:"foreignKey:LoanApplicationID" json:"guarantors,omitempty"`
}

func (LoanApplication) TableName() string { return "loan_applications" }

type LoanApproval struct {
	BaseModel
	LoanApplicationID uint64    `gorm:"column:loan_application_id;index;not null" json:"loan_application_id"`
	ApprovalRole      string    `gorm:"column:approval_role;not null" json:"approval_role"`
	ApproverID        uint64    `gorm:"column:approver_id;not null" json:"approver_id"`
	ApprovalDate      time.Time `gorm:"column:approval_date;not null" json:"approval_date"`
	Decision          string    `gorm:"column:decision;not null" json:"decision"`
	ApprovedAmount    *float64  `gorm:"column:approved_amount;type:decimal(15,2)" json:"approved_amount"`
	ApprovedTermMonths *int     `gorm:"column:approved_term_months" json:"approved_term_months"`
	Comments          string    `gorm:"column:comments;type:text" json:"comments"`
	Conditions        string    `gorm:"column:conditions;type:text" json:"conditions"`
	StepOrder         int       `gorm:"column:step_order;default:1" json:"step_order"`
}

func (LoanApproval) TableName() string { return "loan_approvals" }

type LoanGuarantor struct {
	BaseModel
	LoanApplicationID uint64  `gorm:"column:loan_application_id;index;not null" json:"loan_application_id"`
	MemberID          uint64  `gorm:"column:member_id;not null" json:"member_id"`
	GuaranteedAmount  float64 `gorm:"column:guaranteed_amount;type:decimal(15,2);not null" json:"guaranteed_amount"`
	Relationship      string  `gorm:"column:relationship;size:128" json:"relationship"`
}

func (LoanGuarantor) TableName() string { return "loan_guarantors" }

type LoanDocument struct {
	BaseModel
	LoanApplicationID *uint64 `gorm:"column:loan_application_id;index" json:"loan_application_id"`
	LoanContractID    *uint64 `gorm:"column:loan_contract_id;index" json:"loan_contract_id"`
	DocumentType      string  `gorm:"column:document_type;size:64;not null" json:"document_type"`
	FilePath          string  `gorm:"column:file_path;size:512;not null" json:"file_path"`
	FileName          string  `gorm:"column:file_name;size:255;not null" json:"file_name"`
	UploadedBy        *uint64 `gorm:"column:uploaded_by" json:"uploaded_by"`
}

func (LoanDocument) TableName() string { return "loan_documents" }

// --- Contracts ---

type LoanContract struct {
	BaseModel
	ContractNo              string     `gorm:"column:contract_no;uniqueIndex;size:64;not null" json:"contract_no"`
	LoanApplicationID       *uint64    `gorm:"column:loan_application_id" json:"loan_application_id"`
	MemberID                uint64     `gorm:"column:member_id;index;not null" json:"member_id"`
	MemberNo                string     `gorm:"-" json:"member_no,omitempty"`
	MemberName              string     `gorm:"-" json:"member_name,omitempty"`
	LoanProductID           uint64     `gorm:"column:loan_product_id;index;not null" json:"loan_product_id"`
	PrincipalAmount         float64    `gorm:"column:principal_amount;type:decimal(15,2);not null" json:"principal_amount"`
	ApprovedAmount          float64    `gorm:"column:approved_amount;type:decimal(15,2);not null" json:"approved_amount"`
	DisbursedAmount            float64    `gorm:"column:disbursed_amount;type:decimal(15,2);default:0" json:"disbursed_amount"`
	FeesDeductedAtDisbursement float64    `gorm:"column:fees_deducted_at_disbursement;type:decimal(15,2);default:0" json:"fees_deducted_at_disbursement"`
	InterestRate            float64    `gorm:"column:interest_rate;type:decimal(8,4);not null" json:"interest_rate"`
	InterestMethod          string     `gorm:"column:interest_method;size:32;not null" json:"interest_method"`
	TermMonths              int        `gorm:"column:term_months;not null" json:"term_months"`
	GracePeriodDays         int        `gorm:"column:grace_period_days;default:0" json:"grace_period_days"`
	Status                  string     `gorm:"column:status;not null;default:PENDING_DISBURSEMENT" json:"status"`
	OutstandingPrincipal    float64    `gorm:"column:outstanding_principal;type:decimal(15,2);default:0" json:"outstanding_principal"`
	OutstandingInterest     float64    `gorm:"column:outstanding_interest;type:decimal(15,2);default:0" json:"outstanding_interest"`
	OutstandingFees         float64    `gorm:"column:outstanding_fees;type:decimal(15,2);default:0" json:"outstanding_fees"`
	OutstandingPenalties    float64    `gorm:"column:outstanding_penalties;type:decimal(15,2);default:0" json:"outstanding_penalties"`
	TotalPaid               float64    `gorm:"column:total_paid;type:decimal(15,2);default:0" json:"total_paid"`
	DaysInArrears           int        `gorm:"column:days_in_arrears;default:0" json:"days_in_arrears"`
	DisbursementDate        *time.Time `gorm:"column:disbursement_date;type:date" json:"disbursement_date"`
	MaturityDate            *time.Time `gorm:"column:maturity_date;type:date" json:"maturity_date"`
	LastPaymentDate         *time.Time `gorm:"column:last_payment_date;type:date" json:"last_payment_date"`
	InstallmentAmount       *float64   `gorm:"column:installment_amount;type:decimal(15,2)" json:"installment_amount"`
	MilkDeductionEnabled    bool       `gorm:"column:milk_deduction_enabled;default:1" json:"milk_deduction_enabled"`
	MaxDeductionPercent     float64    `gorm:"column:max_deduction_percent;type:decimal(5,2);default:50" json:"max_deduction_percent"`
	RecurrentDeductionID    *uint64    `gorm:"column:recurrent_deduction_id" json:"recurrent_deduction_id"`
	LoanProduct             *LoanProduct `gorm:"foreignKey:LoanProductID" json:"loan_product,omitempty"`
	Installments            []LoanScheduleInstallment `gorm:"foreignKey:LoanContractID" json:"installments,omitempty"`
	Disbursements           []LoanDisbursement        `gorm:"foreignKey:LoanContractID" json:"disbursements,omitempty"`
}

func (LoanContract) TableName() string { return "loan_contracts" }

// --- Disbursement Channels ---

type DisbursementChannel struct {
	BaseModel
	ChannelCode         string         `gorm:"column:channel_code;uniqueIndex;size:32;not null" json:"channel_code"`
	ChannelName         string         `gorm:"column:channel_name;size:128;not null" json:"channel_name"`
	Description         string         `gorm:"column:description;size:512" json:"description"`
	ChannelType         string         `gorm:"column:channel_type;size:16;not null" json:"channel_type"`
	Provider            string         `gorm:"column:provider;size:64" json:"provider"`
	IsAsync             bool           `gorm:"column:is_async;default:0" json:"is_async"`
	RequiresExternalRef bool           `gorm:"column:requires_external_ref;default:0" json:"requires_external_ref"`
	GLRuleType          string         `gorm:"column:gl_rule_type;size:64;default:LOAN_DISBURSEMENT" json:"gl_rule_type"`
	ConfigSchema        datatypes.JSON `gorm:"column:config_schema;type:json" json:"config_schema"`
	SortOrder           int            `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsActive            bool           `gorm:"column:is_active;default:1" json:"is_active"`
	ChannelConfigs      []DisbursementChannelConfig `gorm:"foreignKey:DisbursementChannelID" json:"channel_configs,omitempty"`
}

func (DisbursementChannel) TableName() string { return "disbursement_channels" }

type DisbursementChannelConfig struct {
	BaseModel
	DisbursementChannelID uint64 `gorm:"column:disbursement_channel_id;not null;index;uniqueIndex:uk_channel_config_key" json:"disbursement_channel_id"`
	ConfigKey             string `gorm:"column:config_key;size:128;not null;uniqueIndex:uk_channel_config_key" json:"config_key"`
	ConfigValue           string `gorm:"column:config_value;type:text" json:"config_value"`
	IsSecret              bool   `gorm:"column:is_secret;default:0" json:"is_secret"`
}

func (DisbursementChannelConfig) TableName() string { return "disbursement_channel_config" }

type LoanDisbursement struct {
	BaseModel
	LoanContractID        uint64         `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	DisbursementChannelID *uint64        `gorm:"column:disbursement_channel_id;index" json:"disbursement_channel_id"`
	ChannelCode           string         `gorm:"column:channel_code;size:32;index" json:"channel_code"`
	Status                string         `gorm:"column:status;size:32;not null;default:PENDING" json:"status"`
	GrossAmount           float64        `gorm:"column:gross_amount;type:decimal(15,2);default:0" json:"gross_amount"`
	FeesDeducted          float64        `gorm:"column:fees_deducted;type:decimal(15,2);default:0" json:"fees_deducted"`
	Amount                float64        `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"` // net paid via channel
	Method                string         `gorm:"column:method;not null;default:BANK" json:"method"`     // legacy
	Reference             string         `gorm:"column:reference;uniqueIndex;size:128;not null" json:"reference"`
	IdempotencyKey        string         `gorm:"column:idempotency_key;uniqueIndex;size:64" json:"idempotency_key"`
	ExternalReference     string         `gorm:"column:external_reference;size:128" json:"external_reference"`
	FailureCode           string         `gorm:"column:failure_code;size:64" json:"failure_code"`
	FailureMessage        string         `gorm:"column:failure_message;type:text" json:"failure_message"`
	ChannelPayload        datatypes.JSON `gorm:"column:channel_payload;type:json" json:"channel_payload"`
	DisbursementDate      time.Time      `gorm:"column:disbursement_date;type:date;not null" json:"disbursement_date"`
	ProcessedAt           *time.Time     `gorm:"column:processed_at" json:"processed_at"`
	CompletedAt           *time.Time     `gorm:"column:completed_at" json:"completed_at"`
	GLTransactionID       *uint64        `gorm:"column:gl_transaction_id" json:"gl_transaction_id"`
	Notes                 string         `gorm:"column:notes;type:text" json:"notes"`
	DisbursementChannel   *DisbursementChannel `gorm:"foreignKey:DisbursementChannelID" json:"disbursement_channel,omitempty"`
}

func (LoanDisbursement) TableName() string { return "loan_disbursements" }

type LoanScheduleInstallment struct {
	BaseModel
	LoanContractID uint64    `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	InstallmentNo  int       `gorm:"column:installment_no;not null" json:"installment_no"`
	DueDate        time.Time `gorm:"column:due_date;type:date;index;not null" json:"due_date"`
	PrincipalDue   float64   `gorm:"column:principal_due;type:decimal(15,2)" json:"principal_due"`
	InterestDue    float64   `gorm:"column:interest_due;type:decimal(15,2)" json:"interest_due"`
	FeeDue         float64   `gorm:"column:fee_due;type:decimal(15,2)" json:"fee_due"`
	InsuranceDue   float64   `gorm:"column:insurance_due;type:decimal(15,2)" json:"insurance_due"`
	PenaltyDue     float64   `gorm:"column:penalty_due;type:decimal(15,2)" json:"penalty_due"`
	TotalDue       float64   `gorm:"column:total_due;type:decimal(15,2)" json:"total_due"`
	PrincipalPaid  float64   `gorm:"column:principal_paid;type:decimal(15,2)" json:"principal_paid"`
	InterestPaid   float64   `gorm:"column:interest_paid;type:decimal(15,2)" json:"interest_paid"`
	FeePaid        float64   `gorm:"column:fee_paid;type:decimal(15,2)" json:"fee_paid"`
	InsurancePaid  float64   `gorm:"column:insurance_paid;type:decimal(15,2)" json:"insurance_paid"`
	PenaltyPaid    float64   `gorm:"column:penalty_paid;type:decimal(15,2)" json:"penalty_paid"`
	BalanceAfter   float64   `gorm:"column:balance_after;type:decimal(15,2)" json:"balance_after"`
	Status         string    `gorm:"column:status;default:PENDING" json:"status"`
}

func (LoanScheduleInstallment) TableName() string { return "loan_schedule_installments" }

type LoanRepaymentRecord struct {
	BaseModel
	LoanContractID   uint64    `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	Amount           float64   `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`
	Channel          string    `gorm:"column:channel;not null" json:"channel"`
	Reference        string    `gorm:"column:reference;uniqueIndex;size:128;not null" json:"reference"`
	IdempotencyKey   string    `gorm:"column:idempotency_key;uniqueIndex;size:64" json:"idempotency_key"`
	PaymentDate      time.Time `gorm:"column:payment_date;type:date;index;not null" json:"payment_date"`
	MemberPayrollID  *uint64   `gorm:"column:member_payroll_id" json:"member_payroll_id"`
	MemberPayslipID  *uint64   `gorm:"column:member_payslip_id" json:"member_payslip_id"`
	GLTransactionID  *uint64   `gorm:"column:gl_transaction_id" json:"gl_transaction_id"`
	Notes            string    `gorm:"column:notes;type:text" json:"notes"`
	Allocations      []LoanRepaymentAllocation `gorm:"foreignKey:LoanRepaymentID" json:"allocations,omitempty"`
}

func (LoanRepaymentRecord) TableName() string { return "loan_repayments" }

type LoanRepaymentAllocation struct {
	BaseModel
	LoanRepaymentID           uint64  `gorm:"column:loan_repayment_id;index;not null" json:"loan_repayment_id"`
	LoanScheduleInstallmentID *uint64 `gorm:"column:loan_schedule_installment_id" json:"loan_schedule_installment_id"`
	AllocationType            string  `gorm:"column:allocation_type;not null" json:"allocation_type"`
	Amount                    float64 `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`
}

func (LoanRepaymentAllocation) TableName() string { return "loan_repayment_allocations" }

type LoanChargeRecord struct {
	BaseModel
	LoanContractID  uint64    `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	ChargeType      string    `gorm:"column:charge_type;not null" json:"charge_type"`
	Amount          float64   `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`
	Capitalized     bool      `gorm:"column:capitalized;default:0" json:"capitalized"`
	ChargedDate     time.Time `gorm:"column:charged_date;type:date;not null" json:"charged_date"`
	GLTransactionID *uint64   `gorm:"column:gl_transaction_id" json:"gl_transaction_id"`
}

func (LoanChargeRecord) TableName() string { return "loan_charges" }

type LoanPenaltyRecord struct {
	BaseModel
	LoanContractID          uint64     `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	LoanScheduleInstallmentID *uint64  `gorm:"column:loan_schedule_installment_id" json:"loan_schedule_installment_id"`
	PenaltyType             string     `gorm:"column:penalty_type;not null" json:"penalty_type"`
	Amount                  float64    `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`
	PenaltyDate             time.Time  `gorm:"column:penalty_date;type:date;not null" json:"penalty_date"`
	Waived                  bool       `gorm:"column:waived;default:0" json:"waived"`
	WaivedBy                *uint64    `gorm:"column:waived_by" json:"waived_by"`
	WaivedAt                *time.Time `gorm:"column:waived_at" json:"waived_at"`
	GLTransactionID         *uint64    `gorm:"column:gl_transaction_id" json:"gl_transaction_id"`
}

func (LoanPenaltyRecord) TableName() string { return "loan_penalties" }

type LoanInterestAccrual struct {
	BaseModel
	LoanContractID  uint64    `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	AccrualDate     time.Time `gorm:"column:accrual_date;type:date;not null" json:"accrual_date"`
	PrincipalBase   float64   `gorm:"column:principal_base;type:decimal(15,2);not null" json:"principal_base"`
	InterestAmount  float64   `gorm:"column:interest_amount;type:decimal(15,2);not null" json:"interest_amount"`
	Posted          bool      `gorm:"column:posted;default:0" json:"posted"`
	GLTransactionID *uint64   `gorm:"column:gl_transaction_id" json:"gl_transaction_id"`
}

func (LoanInterestAccrual) TableName() string { return "loan_interest_accruals" }

type LoanRestructuring struct {
	BaseModel
	LoanContractID     uint64     `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	RestructureType    string     `gorm:"column:restructure_type;not null" json:"restructure_type"`
	PreviousTermMonths *int       `gorm:"column:previous_term_months" json:"previous_term_months"`
	NewTermMonths      *int       `gorm:"column:new_term_months" json:"new_term_months"`
	PreviousRate       *float64   `gorm:"column:previous_rate;type:decimal(8,4)" json:"previous_rate"`
	NewRate            *float64   `gorm:"column:new_rate;type:decimal(8,4)" json:"new_rate"`
	Reason             string     `gorm:"column:reason;type:text" json:"reason"`
	EffectiveDate      time.Time  `gorm:"column:effective_date;type:date;not null" json:"effective_date"`
}

func (LoanRestructuring) TableName() string { return "loan_restructurings" }

type LoanWriteOffRecord struct {
	BaseModel
	LoanContractID  uint64    `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	Amount          float64   `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`
	Reason          string    `gorm:"column:reason;type:text" json:"reason"`
	WriteOffDate    time.Time `gorm:"column:write_off_date;type:date;not null" json:"write_off_date"`
	GLTransactionID *uint64   `gorm:"column:gl_transaction_id" json:"gl_transaction_id"`
	ApprovedBy      *uint64   `gorm:"column:approved_by" json:"approved_by"`
}

func (LoanWriteOffRecord) TableName() string { return "loan_write_offs" }

type LoanAuditLog struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	EntityType string         `gorm:"column:entity_type;size:64;not null;index:idx_loan_audit_entity,priority:1" json:"entity_type"`
	EntityID   uint64         `gorm:"column:entity_id;not null;index:idx_loan_audit_entity,priority:2" json:"entity_id"`
	Action     string         `gorm:"column:action;size:64;not null" json:"action"`
	OldValues  datatypes.JSON `gorm:"column:old_values;type:json" json:"old_values"`
	NewValues  datatypes.JSON `gorm:"column:new_values;type:json" json:"new_values"`
	UserID     *uint64        `gorm:"column:user_id" json:"user_id"`
	IPAddress  string         `gorm:"column:ip_address;size:45" json:"ip_address"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null" json:"created_at"`
}

func (LoanAuditLog) TableName() string { return "loan_audit_logs" }

// Loan module status constants
const (
	LoanAppDraft        = "DRAFT"
	LoanAppSubmitted    = "SUBMITTED"
	LoanAppUnderReview  = "UNDER_REVIEW"
	LoanAppApproved     = "APPROVED"
	LoanAppRejected     = "REJECTED"
	LoanAppCancelled    = "CANCELLED"
	LoanAppDisbursed    = "DISBURSED"

	LoanContractPending    = "PENDING_DISBURSEMENT"
	LoanContractActive     = "ACTIVE"
	LoanContractOverdue    = "OVERDUE"
	LoanContractRestructured = "RESTRUCTURED"
	LoanContractSettled    = "SETTLED"
	LoanContractWrittenOff = "WRITTEN_OFF"
	LoanContractClosed     = "CLOSED"

	LoanChannelMilkDeduction = "MILK_DEDUCTION"
	LoanChannelCash          = "CASH"
	LoanChannelBank          = "BANK"
	LoanChannelMobileMoney   = "MOBILE_MONEY"
	LoanChannelJournal       = "JOURNAL"

	LoanInterestFlat              = "FLAT"
	LoanInterestReducingBalance   = "REDUCING_BALANCE"
	LoanInterestEqualInstallments = "EQUAL_INSTALLMENTS"
	LoanInterestOnly              = "INTEREST_ONLY"
	LoanInterestBalloon           = "BALLOON"

	DefaultAllocationOrder = "PENALTY,FEE,INTEREST,PRINCIPAL,INSURANCE"

	LoanFeeDeductFromProceeds      = "DEDUCT_FROM_PROCEEDS"
	LoanFeePayUpfrontCash          = "PAY_UPFRONT_CASH"
	LoanFeeSpreadOverInstallments  = "SPREAD_OVER_INSTALLMENTS"

	// Disbursement channel types
	DisbursementChannelCash    = "CASH"
	DisbursementChannelBank    = "BANK"
	DisbursementChannelMobile  = "MOBILE"
	DisbursementChannelItem    = "ITEM"
	DisbursementChannelWallet  = "WALLET"

	// Disbursement channel codes
	DisbursementCodeCash        = "CASH"
	DisbursementCodeItemPurchase = "ITEM_PURCHASE"
	DisbursementCodeMobileMpesa  = "MOBILE_MPESA"
	DisbursementCodeMobileDtbMpesa = "MOBILE_DTB_MPESA"
	DisbursementCodeMobileAirtel = "MOBILE_AIRTEL"
	DisbursementCodeMobileEquitel = "MOBILE_EQUITEL"
	DisbursementCodeBankCoop     = "BANK_COOP"
	DisbursementCodeBankEquity   = "BANK_EQUITY"
	DisbursementCodeBankKCB      = "BANK_KCB"
	DisbursementCodeWalletMember = "WALLET_MEMBER"

	// Disbursement attempt status
	DisbursementStatusPending    = "PENDING"
	DisbursementStatusProcessing = "PROCESSING"
	DisbursementStatusSuccess    = "SUCCESS"
	DisbursementStatusFailed     = "FAILED"
	DisbursementStatusCancelled  = "CANCELLED"
)
