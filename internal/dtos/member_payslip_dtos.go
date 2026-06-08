package dtos

import "time"

type MemberPayslipResponse struct {
	ID              uint64     `json:"id"`
	MemberID        uint64     `json:"member_id"`
	MemberNo        string     `json:"member_no"`
	MemberName      string     `json:"member_name"`
	PayrollID       uint64     `json:"payroll_id"`
	PayDateRangeID  *uint64    `json:"pay_date_range_id"`
	DateOpened      *time.Time `json:"date_opened"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`
	RateUsed        string     `json:"rate_used"`
	PostedAt        *time.Time `json:"posted_at"`
	PostedByName    string     `json:"posted_by_name"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
	ConfirmedByName string     `json:"confirmed_by_name"`
	ApprovedAt      *time.Time `json:"approved_at"`
	ApprovedByName  string     `json:"approved_by_name"`
	GrossKilos      float64    `json:"gross_kilos"`
	RejectKilos     float64    `json:"reject_kilos"`
	NetKilos        float64    `json:"net_kilos"`
	GrossPay        float64    `json:"gross_pay"`
	TotalDeductions float64    `json:"total_deductions"`
	NetPay          float64    `json:"net_pay"`
	FiscalPeriod    string     `json:"fiscal_period"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
