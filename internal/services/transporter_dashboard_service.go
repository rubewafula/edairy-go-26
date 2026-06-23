package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type TransporterDashboardService struct{}

func NewTransporterDashboardService() *TransporterDashboardService {
	return &TransporterDashboardService{}
}

func (s *TransporterDashboardService) getMilkDeliveryStat() dtos.TransporterDashboardMilkDeliveryStat {
	var result dtos.TransporterDashboardMilkDeliveryStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity_accepted, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity_accepted, 0)), 0) AS month
		FROM milk_deliveries
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *TransporterDashboardService) getMilkRejectStat() dtos.TransporterDashboardMilkRejectStat {
	var result dtos.TransporterDashboardMilkRejectStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity, 0)), 0) AS month
		FROM milk_rejects
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *TransporterDashboardService) getMilkCollectionStat() dtos.TransporterDashboardMilkCollectionStat {
	var result dtos.TransporterDashboardMilkCollectionStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity, 0), 0)), 0) AS milk_today,
			COALESCE(SUM(IFNULL(quantity, 0)), 0) AS milk_this_month
		FROM milk_journal_entries
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *TransporterDashboardService) getTransportOperationsTrend() []dtos.TransporterDashboardTrendPoint {
	var results []dtos.TransporterDashboardTrendPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(deliveries.total, 0) AS deliveries,
			COALESCE(rejects.total, 0) AS rejects
		FROM (
			SELECT DATE_SUB(CURDATE(), INTERVAL n DAY) AS date
			FROM (
				SELECT 1 AS n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) days
		) dates
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(quantity_accepted, 0)) AS total
			FROM milk_deliveries
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			GROUP BY DATE(created_at)
		) deliveries ON deliveries.date = dates.date
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(quantity, 0)) AS total
			FROM milk_rejects
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			GROUP BY DATE(created_at)
		) rejects ON rejects.date = dates.date
		ORDER BY dates.date ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.TransporterDashboardTrendPoint{}
	}

	return results
}

func (s *TransporterDashboardService) GetDashboard() dtos.TransporterDashboardResponse {
	return dtos.TransporterDashboardResponse{
		MilkDeliveryStat:         s.getMilkDeliveryStat(),
		MilkRejectStat:           s.getMilkRejectStat(),
		MilkCollectionStat:       s.getMilkCollectionStat(),
		TransportOperationsTrend: s.getTransportOperationsTrend(),
	}
}
