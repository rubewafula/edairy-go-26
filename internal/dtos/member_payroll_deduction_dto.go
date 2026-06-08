package dtos

import "time"

// CreateMemberPayrollDeductionRequest defines the structure for creating a new member payroll deduction.
type CreateMemberPayrollDeductionRequest struct {
	MemberID        uint64 `json:"member_id" binding:"required"`
	DeductionMonth  string `json:"deduction_month" binding:"required"`
	FiscalYear      int    `json:"fiscal_year" binding:"required"`
	DeductionTypeID uint64 `json:"deduction_type_id" binding:"required"`
	Amount          string `json:"amount" binding:"required"` // Keep as string for now, as per DB
	Priority        int    `json:"priority"`
	Settled         string `json:"settled"` // "0" or "1"
	TransactionDate string `json:"transaction_date" binding:"required"`
	DateCaptured    string `json:"date_captured"`
	Confirmed       string `json:"confirmed"` // "0" or "1"
	PayrollID       uint64 `json:"payroll_id"`
	Reference       string `json:"reference"`
	SettlementType  string `json:"settlement_type"`
}

// UpdateMemberPayrollDeductionRequest defines the structure for updating an existing member payroll deduction.
type UpdateMemberPayrollDeductionRequest struct {
	MemberID        uint64 `json:"member_id"`
	DeductionMonth  string `json:"deduction_month"`
	FiscalYear      int    `json:"fiscal_year"`
	DeductionTypeID uint64 `json:"deduction_type_id"`
	Amount          string `json:"amount"`
	Priority        int    `json:"priority"`
	Settled         string `json:"settled"`
	TransactionDate string `json:"transaction_date"`
	DateCaptured    string `json:"date_captured"`
	Confirmed       string `json:"confirmed"`
	PayrollID       uint64 `json:"payroll_id"`
	Reference       string `json:"reference"`
	SettlementType  string `json:"settlement_type"`
}

// MemberPayrollDeductionResponse defines the structure for a member payroll deduction response.
type MemberPayrollDeductionResponse struct {
	ID                uint64     `json:"id"`
	MemberID          uint64     `json:"member_id"`
	MemberNo          string     `json:"member_no"`
	MemberName        string     `json:"member_name"`
	DeductionMonth    string     `json:"deduction_month"`
	FiscalYear        int        `json:"fiscal_year"`
	DeductionTypeID   uint64     `json:"deduction_type_id"`
	DeductionTypeName string     `json:"deduction_type_name"`
	Amount            string     `json:"amount"`
	Priority          int        `json:"priority"`
	Settled           string     `json:"settled"`
	TransactionDate   *time.Time `json:"transaction_date"`
	DateCaptured      *time.Time `json:"date_captured"`
	Confirmed         string     `json:"confirmed"`
	PayrollID         uint64     `json:"payroll_id"`
	FiscalPeriod      string     `json:"fiscal_period"`
	Reference         string     `json:"reference"`
	SettlementType    string     `json:"settlement_type"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
