package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type LoanDashboardService struct{}

func NewLoanDashboardService() *LoanDashboardService {
	return &LoanDashboardService{}
}

func (s *LoanDashboardService) getLoanStat() dtos.LoanDashboardLoanStat {
	var result dtos.LoanDashboardLoanStat

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_loans,
			COALESCE(SUM(IFNULL(amount, 0)), 0) AS total_amount,
			COALESCE(SUM(IFNULL(approved_amount, 0)), 0) AS approved_amount
		FROM loans
		WHERE deleted_at IS NULL
	`).Scan(&result)

	result.Currency = "KES"
	return result
}

func (s *LoanDashboardService) getLoansByStatus() []dtos.LoanDashboardStatusCount {
	var counts struct {
		Approved int64
		Rejected int64
		Active   int64
		Pending  int64
	}

	db.DB.Raw(`
		SELECT
			SUM(CASE WHEN UPPER(status) = 'APPROVED' THEN 1 ELSE 0 END) AS approved,
			SUM(CASE WHEN UPPER(status) = 'REJECTED' THEN 1 ELSE 0 END) AS rejected,
			SUM(CASE WHEN UPPER(status) IN ('DISBURSED', 'ACTIVE') THEN 1 ELSE 0 END) AS active,
			SUM(CASE WHEN UPPER(status) = 'PENDING' THEN 1 ELSE 0 END) AS pending
		FROM loans
		WHERE deleted_at IS NULL
	`).Scan(&counts)

	return []dtos.LoanDashboardStatusCount{
		{Status: "Approved", Count: counts.Approved},
		{Status: "Rejected", Count: counts.Rejected},
		{Status: "Active", Count: counts.Active},
		{Status: "Pending", Count: counts.Pending},
	}
}

func (s *LoanDashboardService) getLoanPortfolio() []dtos.LoanDashboardPortfolioItem {
	var results []dtos.LoanDashboardPortfolioItem

	db.DB.Raw(`
		SELECT
			COALESCE(NULLIF(lc.type, ''), 'General Loan') AS type,
			COALESCE(SUM(IFNULL(l.approved_amount, l.amount)), 0) AS amount
		FROM loans l
		LEFT JOIN (
			SELECT loan_id, MIN(type) AS type
			FROM loan_callbacks
			WHERE deleted_at IS NULL
			GROUP BY loan_id
		) lc ON lc.loan_id = l.id
		WHERE l.deleted_at IS NULL
		GROUP BY type
		ORDER BY amount DESC
	`).Scan(&results)

	if results == nil {
		return []dtos.LoanDashboardPortfolioItem{}
	}

	return results
}

func (s *LoanDashboardService) getLoanAmountSplit() dtos.LoanDashboardAmountSplit {
	var result dtos.LoanDashboardAmountSplit

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IFNULL(amount, 0)), 0) AS total_amount,
			COALESCE(SUM(IFNULL(approved_amount, 0)), 0) AS approved_amount
		FROM loans
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *LoanDashboardService) GetDashboard() dtos.LoanDashboardResponse {
	stat := s.getLoanStat()
	return dtos.LoanDashboardResponse{
		LoanStat:        stat,
		LoansByStatus:   s.getLoansByStatus(),
		LoanPortfolio:   s.getLoanPortfolio(),
		LoanAmountSplit: s.getLoanAmountSplit(),
	}
}
