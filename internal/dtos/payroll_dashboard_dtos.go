package dtos

type PayrollDashboardResponse struct {
	PayrollSummary          PayrollDashboardSummary          `json:"payroll_summary"`
	PayrollProcessingStatus PayrollDashboardProcessingStatus `json:"payroll_processing_status"`
	PayrollGrowth           PayrollDashboardGrowth           `json:"payroll_growth"`
	PayrollTrend            []PayrollDashboardTrendPoint     `json:"payroll_trend"`
}

type PayrollDashboardSummary struct {
	TotalPayrollCostCurrentMonth float64 `json:"total_payroll_cost_current_month"`
	NetPayCurrentMonth           float64 `json:"net_pay_current_month"`
	GrossPayCurrentMonth         float64 `json:"gross_pay_current_month"`
	TotalMembersPaid             int64   `json:"total_members_paid"`
	PendingApprovals             int64   `json:"pending_approvals"`
	OvertimeCostCurrentMonth     float64 `json:"overtime_cost_current_month"`
	TotalDeductionsCurrentMonth  float64 `json:"total_deductions_current_month"`
	BenefitsCostCurrentMonth     float64 `json:"benefits_cost_current_month"`
	AverageMemberPay             float64 `json:"average_member_pay"`
	PayrollRunDate               string  `json:"payroll_run_date"`
	Currency                     string  `json:"currency"`
}

type PayrollDashboardProcessingStatus struct {
	Processed int64 `json:"processed"`
	Errors    int64 `json:"errors"`
}

type PayrollDashboardGrowth struct {
	CurrentMonthPay  float64 `json:"current_month_pay"`
	PreviousMonthPay float64 `json:"previous_month_pay"`
}

type PayrollDashboardTrendPoint struct {
	Month    string  `json:"month"`
	GrossPay float64 `json:"gross_pay"`
	NetPay   float64 `json:"net_pay"`
}
