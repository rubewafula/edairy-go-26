package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type LoanDashboardService struct{}

func NewLoanDashboardService() *LoanDashboardService {
	return &LoanDashboardService{}
}

// allLoansSQL is every loan the dashboards summarise: the legacy `loans` table plus loan-module
// applications (desktop and mobile cashouts). Loan-module statuses are mapped to the legacy ones:
// SUBMITTED/UNDER_REVIEW -> PENDING, APPROVED -> APPROVED, DISBURSED -> ACTIVE. DRAFT applications
// (mobile requests whose OTP was never confirmed) are not loans yet and are left out.
const allLoansSQL = `(
	SELECT l.amount AS amount, l.approved_amount AS approved_amount, UPPER(l.status) AS status,
		l.created_at AS created_at, COALESCE(NULLIF(lc.type, ''), 'General Loan') AS type, 'LEGACY' AS source
	FROM loans l
	LEFT JOIN (SELECT loan_id, MIN(type) AS type FROM loan_callbacks WHERE deleted_at IS NULL GROUP BY loan_id) lc
		ON lc.loan_id = l.id
	WHERE l.deleted_at IS NULL
	UNION ALL
	SELECT a.requested_amount, a.approved_amount,
		CASE UPPER(a.status)
			WHEN 'SUBMITTED' THEN 'PENDING' WHEN 'UNDER_REVIEW' THEN 'PENDING'
			WHEN 'DISBURSED' THEN 'ACTIVE' ELSE UPPER(a.status) END,
		a.created_at, COALESCE(NULLIF(p.product_name, ''), 'Loan'), UPPER(COALESCE(NULLIF(a.source, ''), 'BACKOFFICE'))
	FROM loan_applications a
	LEFT JOIN loan_products p ON p.id = a.loan_product_id
	WHERE a.deleted_at IS NULL AND UPPER(a.status) <> 'DRAFT'
) all_loans`

func (s *LoanDashboardService) getLoanStat() dtos.LoanDashboardLoanStat {
	var result dtos.LoanDashboardLoanStat

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_loans,
			COALESCE(SUM(IFNULL(amount, 0)), 0) AS total_amount,
			COALESCE(SUM(IFNULL(approved_amount, 0)), 0) AS approved_amount
		FROM ` + allLoansSQL).Scan(&result)

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
			COALESCE(SUM(CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END), 0) AS approved,
			COALESCE(SUM(CASE WHEN status = 'REJECTED' THEN 1 ELSE 0 END), 0) AS rejected,
			COALESCE(SUM(CASE WHEN status IN ('DISBURSED', 'ACTIVE') THEN 1 ELSE 0 END), 0) AS active,
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending
		FROM ` + allLoansSQL).Scan(&counts)

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
		SELECT type, COALESCE(SUM(COALESCE(approved_amount, amount)), 0) AS amount
		FROM ` + allLoansSQL + `
		WHERE status NOT IN ('REJECTED', 'CANCELLED')
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
		FROM ` + allLoansSQL).Scan(&result)

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
