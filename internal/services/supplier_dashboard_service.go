package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type SupplierDashboardService struct{}

func NewSupplierDashboardService() *SupplierDashboardService {
	return &SupplierDashboardService{}
}

func (s *SupplierDashboardService) getSupplyStats() dtos.SupplierDashboardSupplyStats {
	var result dtos.SupplierDashboardSupplyStats

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN supplied_date >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END), 0) AS total_supplies_this_month,
			COALESCE(SUM(CASE WHEN DATE(supplied_date) = CURDATE() THEN 1 ELSE 0 END), 0) AS supplies_today,
			COALESCE(SUM(CASE WHEN supplied_date >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END), 0) AS supplies_this_month
		FROM supplies
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *SupplierDashboardService) getSupplyVolumeTrend() []dtos.SupplierDashboardVolumePoint {
	var results []dtos.SupplierDashboardVolumePoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(volumes.total, 0) AS volume
		FROM (
			SELECT DATE_SUB(CURDATE(), INTERVAL n DAY) AS date
			FROM (
				SELECT 1 AS n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) days
		) dates
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(total_price, 0)) AS total
			FROM supplied_items
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			  AND deleted_at IS NULL
			GROUP BY DATE(created_at)
		) volumes ON volumes.date = dates.date
		ORDER BY dates.date ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.SupplierDashboardVolumePoint{}
	}

	return results
}

func (s *SupplierDashboardService) getSupplyMixToday() []dtos.SupplierDashboardMixPoint {
	var results []dtos.SupplierDashboardMixPoint

	db.DB.Raw(`
		SELECT
			COALESCE(i.item_name, 'Unknown') AS item,
			COALESCE(SUM(IFNULL(si.total_price, 0)), 0) AS value
		FROM supplied_items si
		LEFT JOIN store_items i ON si.item_id = i.id
		WHERE DATE(si.created_at) = CURDATE()
		  AND si.deleted_at IS NULL
		GROUP BY i.item_name
		ORDER BY value DESC
	`).Scan(&results)

	if results == nil {
		return []dtos.SupplierDashboardMixPoint{}
	}

	return results
}

func (s *SupplierDashboardService) getSalesMixToday() dtos.SupplierDashboardSalesMix {
	var result dtos.SupplierDashboardSalesMix

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IFNULL(amount_paid, 0)), 0) AS cash_sales,
			COALESCE(SUM(IFNULL(amount_due, 0)), 0) AS credit_sales
		FROM store_sales
		WHERE DATE(created_at) = CURDATE()
		  AND deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *SupplierDashboardService) GetDashboard() dtos.SupplierDashboardResponse {
	return dtos.SupplierDashboardResponse{
		SupplyStats:       s.getSupplyStats(),
		SupplyVolumeTrend: s.getSupplyVolumeTrend(),
		SupplyMixToday:    s.getSupplyMixToday(),
		SalesMixToday:     s.getSalesMixToday(),
	}
}
