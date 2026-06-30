package dtos

type LoanDashboardResponse struct {
	LoanStat        LoanDashboardLoanStat        `json:"loan_stat"`
	LoansByStatus   []LoanDashboardStatusCount   `json:"loans_by_status"`
	LoanPortfolio   []LoanDashboardPortfolioItem `json:"loan_portfolio"`
	LoanAmountSplit LoanDashboardAmountSplit     `json:"loan_amount_split"`
}

type LoanDashboardLoanStat struct {
	TotalLoans     int64   `json:"total_loans"`
	TotalAmount    float64 `json:"total_amount"`
	ApprovedAmount float64 `json:"approved_amount"`
	Currency       string  `json:"currency"`
}

type LoanDashboardStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type LoanDashboardPortfolioItem struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

type LoanDashboardAmountSplit struct {
	TotalAmount    float64 `json:"total_amount"`
	ApprovedAmount float64 `json:"approved_amount"`
}
