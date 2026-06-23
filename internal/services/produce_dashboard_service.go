package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type ProduceDashboardService struct{}

func NewProduceDashboardService() *ProduceDashboardService {
	return &ProduceDashboardService{}
}

func (s *ProduceDashboardService) getMilkCollectionStat() dtos.ProduceDashboardQuantityStat {
	var result dtos.ProduceDashboardQuantityStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity, 0)), 0) AS month
		FROM milk_journal_entries
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *ProduceDashboardService) getMilkDeliveryStat() dtos.ProduceDashboardQuantityStat {
	var result dtos.ProduceDashboardQuantityStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity_accepted, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity_accepted, 0)), 0) AS month
		FROM milk_deliveries
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *ProduceDashboardService) getMilkRejectStat() dtos.ProduceDashboardQuantityStat {
	var result dtos.ProduceDashboardQuantityStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity, 0)), 0) AS month
		FROM milk_rejects
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *ProduceDashboardService) getCollectionVsRejectTrend() []dtos.ProduceDashboardTrendPoint {
	var results []dtos.ProduceDashboardTrendPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(collections.total, 0) AS collections,
			COALESCE(rejects.total, 0) AS rejects
		FROM (
			SELECT DISTINCT activity_date AS date
			FROM (
				SELECT DATE(created_at) AS activity_date
				FROM milk_journal_entries
				WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
				UNION
				SELECT DATE(created_at) AS activity_date
				FROM milk_rejects
				WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			) activity
		) dates
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(quantity, 0)) AS total
			FROM milk_journal_entries
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			GROUP BY DATE(created_at)
		) collections ON collections.date = dates.date
		LEFT JOIN (
			SELECT DATE(created_at) AS date, SUM(IFNULL(quantity, 0)) AS total
			FROM milk_rejects
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			GROUP BY DATE(created_at)
		) rejects ON rejects.date = dates.date
		ORDER BY dates.date ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.ProduceDashboardTrendPoint{}
	}

	return results
}

func (s *ProduceDashboardService) GetDashboard() dtos.ProduceDashboardResponse {
	return dtos.ProduceDashboardResponse{
		MilkCollectionStat:      s.getMilkCollectionStat(),
		MilkDeliveryStat:        s.getMilkDeliveryStat(),
		MilkRejectStat:          s.getMilkRejectStat(),
		CollectionVsRejectTrend: s.getCollectionVsRejectTrend(),
	}
}
