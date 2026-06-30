package dtos

type FinanceDashboardResponse struct {
	StoreSaleStat        FinanceDashboardStoreSaleStat        `json:"store_sale_stat"`
	LoanStat             FinanceDashboardLoanStat             `json:"loan_stat"`
	LoansByStatus        []FinanceDashboardStatusCount        `json:"loans_by_status"`
	FinanceOverviewTrend []FinanceDashboardOverviewTrendPoint `json:"finance_overview_trend"`
	LoanAmountSplit      FinanceDashboardLoanAmountSplit      `json:"loan_amount_split"`
}

type FinanceDashboardStoreSaleStat struct {
	TotalStoreSaleThisMonth float64 `json:"total_store_sale_this_month"`
	TotalStoreSaleToday     float64 `json:"total_store_sale_today"`
	CashStoreSaleToday      float64 `json:"cash_store_sale_today"`
	Currency                string  `json:"currency"`
}

type FinanceDashboardLoanStat struct {
	TotalLoans     int64   `json:"total_loans"`
	PendingLoans   int64   `json:"pending_loans"`
	ApprovedLoans  int64   `json:"approved_loans"`
	TotalAmount    float64 `json:"total_amount"`
	ApprovedAmount float64 `json:"approved_amount"`
}

type FinanceDashboardStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type FinanceDashboardOverviewTrendPoint struct {
	Date       string  `json:"date"`
	LoanAmount float64 `json:"loan_amount"`
	StoreSales float64 `json:"store_sales"`
}

type FinanceDashboardLoanAmountSplit struct {
	TotalAmount    float64 `json:"total_amount"`
	ApprovedAmount float64 `json:"approved_amount"`
}
