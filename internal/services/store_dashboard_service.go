package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type StoreDashboardService struct{}

func NewStoreDashboardService() *StoreDashboardService {
	return &StoreDashboardService{}
}

func (s *StoreDashboardService) getStoreSaleStat() dtos.StoreDashboardStoreSaleStat {
	var result struct {
		CashStoreSaleToday              float64
		CreditStoreSaleToday            float64
		TotalStoreSaleToday             float64
		CashStoreSaleThisMonth          float64
		CreditStoreSaleThisMonth        float64
		TotalStoreSaleThisMonth         float64
		TotalSalesTransactionsThisMonth int64
	}

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(amount_paid, 0), 0)), 0) AS cash_store_sale_today,
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(amount_due, 0), 0)), 0) AS credit_store_sale_today,
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(total_amount, 0), 0)), 0) AS total_store_sale_today,
			COALESCE(SUM(IFNULL(amount_paid, 0)), 0) AS cash_store_sale_this_month,
			COALESCE(SUM(IFNULL(amount_due, 0)), 0) AS credit_store_sale_this_month,
			COALESCE(SUM(IFNULL(total_amount, 0)), 0) AS total_store_sale_this_month,
			COUNT(*) AS total_sales_transactions_this_month
		FROM store_sales
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
		  AND deleted_at IS NULL
	`).Scan(&result)

	salesTrend := s.getSalesTrend()

	return dtos.StoreDashboardStoreSaleStat{
		CashStoreSaleToday:              result.CashStoreSaleToday,
		CreditStoreSaleToday:            result.CreditStoreSaleToday,
		TotalStoreSaleToday:             result.TotalStoreSaleToday,
		CashStoreSaleThisMonth:          result.CashStoreSaleThisMonth,
		CreditStoreSaleThisMonth:        result.CreditStoreSaleThisMonth,
		TotalStoreSaleThisMonth:         result.TotalStoreSaleThisMonth,
		TotalSalesTransactionsThisMonth: result.TotalSalesTransactionsThisMonth,
		Currency:                        "KES",
		SalesTrend:                      salesTrend,
	}
}

func (s *StoreDashboardService) getSalesTrend() []dtos.StoreDashboardSalesTrendPoint {
	var results []dtos.StoreDashboardSalesTrendPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(sales.total, 0) AS total_sales
		FROM (
			SELECT DATE_SUB(CURDATE(), INTERVAL n DAY) AS date
			FROM (
				SELECT 1 AS n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) days
		) dates
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
		return []dtos.StoreDashboardSalesTrendPoint{}
	}

	return results
}

func (s *StoreDashboardService) getSupplyStats() dtos.StoreDashboardSupplyStats {
	var result dtos.StoreDashboardSupplyStats

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN DATE(supplied_date) = CURDATE() THEN 1 ELSE 0 END), 0) AS supplies_today,
			COALESCE(SUM(CASE WHEN supplied_date >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END), 0) AS supplies_this_month
		FROM supplies
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *StoreDashboardService) GetDashboard() dtos.StoreDashboardResponse {
	return dtos.StoreDashboardResponse{
		StoreSaleStat: s.getStoreSaleStat(),
		SupplyStats:   s.getSupplyStats(),
	}
}
