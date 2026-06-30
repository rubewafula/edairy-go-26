package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type HRDashboardService struct{}

func NewHRDashboardService() *HRDashboardService {
	return &HRDashboardService{}
}

func (s *HRDashboardService) getEmployeeStats() dtos.HRDashboardEmployeeStats {
	var result struct {
		TotalEmployees    int64
		NewEmployees      int64
		ExistingEmployees int64
		Male              int64
		Female            int64
	}

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_employees,
			SUM(CASE WHEN created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS new_employees,
			SUM(CASE WHEN created_at < DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS existing_employees,
			COUNT(IF(UPPER(gender) = 'MALE', id, NULL)) AS male,
			COUNT(IF(UPPER(gender) = 'FEMALE', id, NULL)) AS female
		FROM employees
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return dtos.HRDashboardEmployeeStats{
		TotalEmployees:    result.TotalEmployees,
		NewEmployees:      result.NewEmployees,
		ExistingEmployees: result.ExistingEmployees,
		Male:              result.Male,
		Female:            result.Female,
	}
}

func (s *HRDashboardService) getLeaveStats() dtos.HRDashboardLeaveStats {
	var result dtos.HRDashboardLeaveStats

	db.DB.Raw(`
		SELECT
			SUM(CASE WHEN UPPER(status) = 'PENDING' THEN 1 ELSE 0 END) AS pending_applications,
			SUM(CASE
				WHEN approved = 1
				 AND CURDATE() BETWEEN DATE(start_date) AND DATE(end_date)
				THEN 1 ELSE 0
			END) AS employees_on_leave
		FROM employee_leave_applications
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return result
}

func (s *HRDashboardService) getEmployeeCategories() []dtos.HRDashboardCategoryCount {
	var results []dtos.HRDashboardCategoryCount

	db.DB.Raw(`
		SELECT
			COALESCE(NULLIF(ecd.contract_type, ''), 'Unclassified') AS category,
			COUNT(DISTINCT e.id) AS count
		FROM employees e
		LEFT JOIN employee_contract_details ecd ON ecd.employee_id = e.id AND ecd.deleted_at IS NULL
		WHERE e.deleted_at IS NULL
		GROUP BY ecd.contract_type
		ORDER BY count DESC
	`).Scan(&results)

	if results == nil {
		return []dtos.HRDashboardCategoryCount{}
	}

	return results
}

func (s *HRDashboardService) GetDashboard() dtos.HRDashboardResponse {
	return dtos.HRDashboardResponse{
		EmployeeStats:      s.getEmployeeStats(),
		LeaveStats:         s.getLeaveStats(),
		EmployeeCategories: s.getEmployeeCategories(),
	}
}
