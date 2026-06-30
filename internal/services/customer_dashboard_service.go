package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type CustomerDashboardService struct{}

func NewCustomerDashboardService() *CustomerDashboardService {
	return &CustomerDashboardService{}
}

func (s *CustomerDashboardService) getBillingStat() dtos.CustomerDashboardBillingStat {
	var result struct {
		TotalInvoicesThisMonth int64
		TotalBillingsThisMonth float64
	}

	db.DB.Raw(`
		SELECT
			(SELECT COUNT(*)
			 FROM customer_invoices
			 WHERE invoice_date >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			   AND deleted_at IS NULL) AS total_invoices_this_month,
			COALESCE((
				SELECT SUM(IFNULL(total_amount, 0))
				FROM customer_billings
				WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
				  AND deleted_at IS NULL
			), 0) AS total_billings_this_month
	`).Scan(&result)

	return dtos.CustomerDashboardBillingStat{
		TotalInvoicesThisMonth: result.TotalInvoicesThisMonth,
		TotalBillingsThisMonth: result.TotalBillingsThisMonth,
		Currency:               "KES",
	}
}

func (s *CustomerDashboardService) getInvoiceStat() dtos.CustomerDashboardInvoiceStat {
	var result dtos.CustomerDashboardInvoiceStat

	db.DB.Raw(`
		SELECT
			SUM(CASE WHEN IFNULL(balance, 0) > 0 THEN 1 ELSE 0 END) AS pending_invoices,
			SUM(CASE WHEN IFNULL(balance, 0) = 0 THEN 1 ELSE 0 END) AS settled_invoices
		FROM customer_invoices
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *CustomerDashboardService) getPaymentStat() dtos.CustomerDashboardPaymentStat {
	var result dtos.CustomerDashboardPaymentStat

	db.DB.Raw(`
		SELECT COALESCE(SUM(IFNULL(amount, 0)), 0) AS total_payments_this_month
		FROM customer_payments
		WHERE payment_date >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
		  AND deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *CustomerDashboardService) getBillingTrend() []dtos.CustomerDashboardTrendPoint {
	var results []dtos.CustomerDashboardTrendPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(billings.total, 0) AS amount
		FROM (
			SELECT DATE_SUB(CURDATE(), INTERVAL n DAY) AS date
			FROM (
				SELECT 1 AS n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) days
		) dates
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(total_amount, 0)) AS total
			FROM customer_billings
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			  AND deleted_at IS NULL
			GROUP BY DATE(created_at)
		) billings ON billings.date = dates.date
		ORDER BY dates.date ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.CustomerDashboardTrendPoint{}
	}

	return results
}

func (s *CustomerDashboardService) GetDashboard() dtos.CustomerDashboardResponse {
	return dtos.CustomerDashboardResponse{
		CustomerBillingStat: s.getBillingStat(),
		InvoiceStat:         s.getInvoiceStat(),
		CustomerPaymentStat: s.getPaymentStat(),
		BillingTrend:        s.getBillingTrend(),
	}
}
