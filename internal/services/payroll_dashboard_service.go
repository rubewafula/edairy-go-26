package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type PayrollDashboardService struct{}

func NewPayrollDashboardService() *PayrollDashboardService {
	return &PayrollDashboardService{}
}

func (s *PayrollDashboardService) getPayrollSummary() dtos.PayrollDashboardSummary {
	var result struct {
		TotalPayrollCostCurrentMonth float64
		NetPayCurrentMonth           float64
		GrossPayCurrentMonth         float64
		TotalMembersPaid             int64
		PendingApprovals             int64
		OvertimeCostCurrentMonth     float64
		TotalDeductionsCurrentMonth  float64
		BenefitsCostCurrentMonth     float64
		PayrollRunDate               *string
	}

	db.DB.Raw(`
		SELECT
			COALESCE(SUM(IFNULL(mp.gross_pay, 0)), 0) AS total_payroll_cost_current_month,
			COALESCE(SUM(IFNULL(mp.net_pay, 0)), 0) AS net_pay_current_month,
			COALESCE(SUM(IFNULL(mp.gross_pay, 0)), 0) AS gross_pay_current_month,
			(SELECT COUNT(*)
			 FROM member_payslips mps
			 INNER JOIN member_payrolls mp2 ON mp2.id = mps.payroll_id
			 WHERE mp2.created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			   AND mp2.deleted_at IS NULL
			   AND mps.deleted_at IS NULL) AS total_members_paid,
			SUM(CASE WHEN UPPER(mp.status) IN ('DRAFT', 'PROCESSING', 'CONFIRMED') THEN 1 ELSE 0 END) AS pending_approvals,
			COALESCE(SUM(IFNULL(mp.transport_cost, 0)), 0) AS overtime_cost_current_month,
			COALESCE(SUM(IFNULL(mp.total_deductions, 0)), 0) AS total_deductions_current_month,
			0 AS benefits_cost_current_month,
			(SELECT DATE_FORMAT(MAX(date_opened), '%Y-%m-%d')
			 FROM member_payrolls
			 WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			   AND deleted_at IS NULL) AS payroll_run_date
		FROM member_payrolls mp
		WHERE mp.created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
		  AND mp.deleted_at IS NULL
	`).Scan(&result)

	averageMemberPay := float64(0)
	if result.TotalMembersPaid > 0 {
		averageMemberPay = result.NetPayCurrentMonth / float64(result.TotalMembersPaid)
	}

	payrollRunDate := ""
	if result.PayrollRunDate != nil {
		payrollRunDate = *result.PayrollRunDate
	}

	return dtos.PayrollDashboardSummary{
		TotalPayrollCostCurrentMonth: result.TotalPayrollCostCurrentMonth,
		NetPayCurrentMonth:           result.NetPayCurrentMonth,
		GrossPayCurrentMonth:         result.GrossPayCurrentMonth,
		TotalMembersPaid:             result.TotalMembersPaid,
		PendingApprovals:             result.PendingApprovals,
		OvertimeCostCurrentMonth:     result.OvertimeCostCurrentMonth,
		TotalDeductionsCurrentMonth:  result.TotalDeductionsCurrentMonth,
		BenefitsCostCurrentMonth:     result.BenefitsCostCurrentMonth,
		AverageMemberPay:             averageMemberPay,
		PayrollRunDate:               payrollRunDate,
		Currency:                     "KES",
	}
}

func (s *PayrollDashboardService) getProcessingStatus() dtos.PayrollDashboardProcessingStatus {
	var result dtos.PayrollDashboardProcessingStatus

	db.DB.Raw(`
		SELECT
			(SELECT COUNT(*)
			 FROM member_payslips mps
			 INNER JOIN member_payrolls mp ON mp.id = mps.payroll_id
			 WHERE mp.created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			   AND UPPER(mps.status) NOT IN ('INCOMPLETE', 'REJECTED', 'CANCELLED')
			   AND mp.deleted_at IS NULL
			   AND mps.deleted_at IS NULL) AS processed,
			(SELECT COUNT(*)
			 FROM member_payroll_generation_errors
			 WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			   AND deleted_at IS NULL) AS errors
	`).Scan(&result)

	return result
}

func (s *PayrollDashboardService) getPayrollGrowth() dtos.PayrollDashboardGrowth {
	var result dtos.PayrollDashboardGrowth

	db.DB.Raw(`
		SELECT
			COALESCE((
				SELECT SUM(IFNULL(gross_pay, 0))
				FROM member_payrolls
				WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
				  AND deleted_at IS NULL
			), 0) AS current_month_pay,
			COALESCE((
				SELECT SUM(IFNULL(gross_pay, 0))
				FROM member_payrolls
				WHERE created_at >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01')
				  AND created_at < DATE_FORMAT(CURDATE(), '%Y-%m-01')
				  AND deleted_at IS NULL
			), 0) AS previous_month_pay
	`).Scan(&result)

	return result
}

func (s *PayrollDashboardService) getPayrollTrend() []dtos.PayrollDashboardTrendPoint {
	var results []dtos.PayrollDashboardTrendPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(month_start, '%b') AS month,
			COALESCE(payroll.gross_pay, 0) AS gross_pay,
			COALESCE(payroll.net_pay, 0) AS net_pay
		FROM (
			SELECT DATE_FORMAT(DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL n MONTH), '%Y-%m-01') AS month_start
			FROM (
				SELECT 0 AS n UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) months
		) months
		LEFT JOIN (
			SELECT
				DATE_FORMAT(created_at, '%Y-%m-01') AS month_start,
				SUM(IFNULL(gross_pay, 0)) AS gross_pay,
				SUM(IFNULL(net_pay, 0)) AS net_pay
			FROM member_payrolls
			WHERE created_at >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 5 MONTH), '%Y-%m-01')
			  AND deleted_at IS NULL
			GROUP BY DATE_FORMAT(created_at, '%Y-%m-01')
		) payroll ON payroll.month_start = months.month_start
		ORDER BY months.month_start ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.PayrollDashboardTrendPoint{}
	}

	return results
}

func (s *PayrollDashboardService) GetDashboard() dtos.PayrollDashboardResponse {
	return dtos.PayrollDashboardResponse{
		PayrollSummary:          s.getPayrollSummary(),
		PayrollProcessingStatus: s.getProcessingStatus(),
		PayrollGrowth:           s.getPayrollGrowth(),
		PayrollTrend:            s.getPayrollTrend(),
	}
}
