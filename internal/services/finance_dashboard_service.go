package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type FinanceDashboardService struct{}

func NewFinanceDashboardService() *FinanceDashboardService {
	return &FinanceDashboardService{}
}

func (s *FinanceDashboardService) getStoreSaleStat() dtos.FinanceDashboardStoreSaleStat {
	var result dtos.FinanceDashboardStoreSaleStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IFNULL(total_amount, 0)), 0) AS total_store_sale_this_month,
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(total_amount, 0), 0)), 0) AS total_store_sale_today,
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(amount_paid, 0), 0)), 0) AS cash_store_sale_today
		FROM store_sales
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
		  AND deleted_at IS NULL
	`).Scan(&result)

	result.Currency = "KES"
	return result
}

func (s *FinanceDashboardService) getLoanStat() dtos.FinanceDashboardLoanStat {
	var result dtos.FinanceDashboardLoanStat

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_loans,
			SUM(CASE WHEN UPPER(status) = 'PENDING' THEN 1 ELSE 0 END) AS pending_loans,
			SUM(CASE WHEN UPPER(status) = 'APPROVED' THEN 1 ELSE 0 END) AS approved_loans,
			COALESCE(SUM(IFNULL(amount, 0)), 0) AS total_amount,
			COALESCE(SUM(IFNULL(approved_amount, 0)), 0) AS approved_amount
		FROM loans
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *FinanceDashboardService) getLoansByStatus() []dtos.FinanceDashboardStatusCount {
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

	return []dtos.FinanceDashboardStatusCount{
		{Status: "Approved", Count: counts.Approved},
		{Status: "Rejected", Count: counts.Rejected},
		{Status: "Active", Count: counts.Active},
		{Status: "Pending", Count: counts.Pending},
	}
}

func (s *FinanceDashboardService) getFinanceOverviewTrend() []dtos.FinanceDashboardOverviewTrendPoint {
	var results []dtos.FinanceDashboardOverviewTrendPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(loans.total, 0) AS loan_amount,
			COALESCE(sales.total, 0) AS store_sales
		FROM (
			SELECT DATE_SUB(CURDATE(), INTERVAL n DAY) AS date
			FROM (
				SELECT 1 AS n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) days
		) dates
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(amount, 0)) AS total
			FROM loans
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			  AND deleted_at IS NULL
			GROUP BY DATE(created_at)
		) loans ON loans.date = dates.date
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(total_amount, 0)) AS total
			FROM store_sales
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			  AND deleted_at IS NULL
			GROUP BY DATE(created_at)
		) sales ON sales.date = dates.date
		ORDER BY dates.date ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.FinanceDashboardOverviewTrendPoint{}
	}

	return results
}

func (s *FinanceDashboardService) getLoanAmountSplit() dtos.FinanceDashboardLoanAmountSplit {
	var result dtos.FinanceDashboardLoanAmountSplit

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IFNULL(amount, 0)), 0) AS total_amount,
			COALESCE(SUM(IFNULL(approved_amount, 0)), 0) AS approved_amount
		FROM loans
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *FinanceDashboardService) GetDashboard() dtos.FinanceDashboardResponse {
	return dtos.FinanceDashboardResponse{
		StoreSaleStat:        s.getStoreSaleStat(),
		LoanStat:             s.getLoanStat(),
		LoansByStatus:        s.getLoansByStatus(),
		FinanceOverviewTrend: s.getFinanceOverviewTrend(),
		LoanAmountSplit:      s.getLoanAmountSplit(),
	}
}
