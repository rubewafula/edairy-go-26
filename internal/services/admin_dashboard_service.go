package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type AdminDashboardService struct{}

func NewAdminDashboardService() *AdminDashboardService {
	return &AdminDashboardService{}
}

func (s *AdminDashboardService) getMilkSummary() (float64, float64) {
	var result struct {
		Today float64
		Month float64
	}

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity, 0)), 0) AS month
		FROM milk_journal_entries
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result.Today, result.Month
}

func (s *AdminDashboardService) getMilkDeliverySummary() dtos.AdminDashboardMilkDeliveryStat {
	var result dtos.AdminDashboardMilkDeliveryStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity_accepted, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity_accepted, 0)), 0) AS month
		FROM milk_delivery_acceptance
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *AdminDashboardService) getMilkRejectsSummary() dtos.AdminDashboardMilkRejectStat {
	var result dtos.AdminDashboardMilkRejectStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IF(DATE(created_at) = CURDATE(), IFNULL(quantity, 0), 0)), 0) AS today,
			COALESCE(SUM(IFNULL(quantity, 0)), 0) AS month
		FROM milk_rejects
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *AdminDashboardService) getSupplyStats() dtos.AdminDashboardSuppliesStat {
	var result dtos.AdminDashboardSuppliesStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(if(DATE(created_at) = CURDATE(), quantity, 0)), 0) AS supplied_items_today,
			COALESCE(SUM(quantity), 0) AS items_supplied_this_month,

			COALESCE(SUM(CASE 
				WHEN DATE(created_at) = CURDATE() 
				THEN IFNULL(total_price, 0) 
				ELSE 0 
			END), 0) AS cash_store_sale_today,
			
			COALESCE(SUM(IFNULL(total_price, 0)), 0) AS total_supplied_this_month
		FROM supplied_items
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *AdminDashboardService) getStoreSaleStats() dtos.AdminDashboardStoreSaleStat {
	var result dtos.AdminDashboardStoreSaleStat

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(CASE 
				WHEN DATE(created_at) = CURDATE() 
				THEN IFNULL(amount_paid, 0) 
				ELSE 0 
			END), 0) AS cash_store_sale_today,
			COALESCE(SUM(IFNULL(amount_paid, 0)), 0) AS cash_store_sale_this_month,
			COALESCE(SUM(CASE 
				WHEN DATE(created_at) = CURDATE() 
				THEN IFNULL(total_amount, 0) 
				ELSE 0 
			END), 0) AS total_store_sale_today,
			COALESCE(SUM(IFNULL(total_amount, 0)), 0) AS total_store_sale_this_month
		FROM store_sales
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *AdminDashboardService) getMembers() dtos.AdminDashboardMemberStat {
	var result dtos.AdminDashboardMemberStat

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) AS active,
			COUNT(IF(gender = 'MALE', id, NULL)) AS male,
			COUNT(IF(gender = 'FEMALE', id, NULL)) AS female
		FROM member_registrations;
	`).Scan(&result)

	return result
}

func (s *AdminDashboardService) getLoanStats() dtos.AdminDashboardLoanStat {
	var result dtos.AdminDashboardLoanStat
	db.DB.Raw(`
		SELECT  
			count(*) as total_loans,
			count(if(status = 'pending', id, null)) as pending,
			COALESCE(SUM(ifnull(amount, 0)), 0) as total_amount,
			COALESCE(SUM(if(status = 'pending', ifnull(amount, 0), 0)), 0)as total_amount_pending
		FROM member_loans
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
	`).Scan(&result)

	return result
}

func (s *AdminDashboardService) getTopRoutes() []dtos.AdminDashboardRouteStat {
	var results []dtos.AdminDashboardRouteStat

	db.DB.Raw(`
		SELECT r.route_name AS route,  COALESCE(SUM(ifnull(m.quantity, 0)), 0) AS total
		FROM milk_journal_entries m
		INNER JOIN routes r ON r.id = m.route_id
		WHERE DATE(m.created_at) = CURDATE()
		GROUP BY r.id
		ORDER BY total DESC
		LIMIT 10
	`).Scan(&results)

	return results
}

func (s *AdminDashboardService) getMilkTrend() []dtos.AdminDashboardDailyMilkStat {
	var results []dtos.AdminDashboardDailyMilkStat

	db.DB.Raw(`
		SELECT DATE(created_at) as date,
		        COALESCE(SUM(ifnull(quantity, 0)), 0) as total
		FROM milk_journal_entries
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`).Scan(&results)

	return results
}

func (s *AdminDashboardService) GetDashboard() dtos.AdminDashboardResponse {

	milkToday, milkThisMonth := s.getMilkSummary()
	return dtos.AdminDashboardResponse{

		TopRoutes:      s.getTopRoutes(),
		SupplyStat:     s.getSupplyStats(),
		StoreSaleStat:  s.getStoreSaleStats(),
		MemberStat:     s.getMembers(),
		MilkRejectStat: s.getMilkRejectsSummary(),
		LoanStat:       s.getLoanStats(),
		MilkCollectionStat: dtos.AdminDashboardMilkCollectionStat{
			MilkToday:     milkToday,
			MilkThisMonth: milkThisMonth,
			MilkTrend:     s.getMilkTrend(),
		},
		MilkDevliveryStat: s.getMilkDeliverySummary(),
	}
}
