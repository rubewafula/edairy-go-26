package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type MemberDashboardService struct{}

func NewMemberDashboardService() *MemberDashboardService {
	return &MemberDashboardService{}
}

func (s *MemberDashboardService) getMemberStats() dtos.MemberDashboardMemberStats {
	var result struct {
		TotalMembers int64
		NewMembers   int64
		Male         int64
		Female       int64
		Active       int64
		Inactive     int64
		Pending      int64
	}

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_members,
			SUM(CASE WHEN created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS new_members,
			COUNT(IF(gender = 'MALE', id, NULL)) AS male,
			COUNT(IF(gender = 'FEMALE', id, NULL)) AS female,
			SUM(CASE WHEN UPPER(status) = 'ACTIVE' THEN 1 ELSE 0 END) AS active,
			SUM(CASE WHEN UPPER(status) IN ('INACTIVE', 'SUSPENDED') THEN 1 ELSE 0 END) AS inactive,
			SUM(CASE WHEN UPPER(status) = 'PENDING' THEN 1 ELSE 0 END) AS pending
		FROM member_registrations
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return dtos.MemberDashboardMemberStats{
		TotalMembers: result.TotalMembers,
		NewMembers:   result.NewMembers,
		Male:         result.Male,
		Female:       result.Female,
		Composition: []dtos.MemberDashboardStatusCount{
			{Status: "Active", Count: result.Active},
			{Status: "Inactive", Count: result.Inactive},
			{Status: "Pending", Count: result.Pending},
		},
	}
}

func (s *MemberDashboardService) getLoanStats() dtos.MemberDashboardLoanStat {
	var result struct {
		TotalLoans    int64
		PendingLoans  int64
		ApprovedLoans int64
	}

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_loans,
			SUM(CASE WHEN UPPER(status) = 'PENDING' THEN 1 ELSE 0 END) AS pending_loans,
			SUM(CASE WHEN UPPER(status) = 'APPROVED' THEN 1 ELSE 0 END) AS approved_loans
		FROM loans
		WHERE deleted_at IS NULL
	`).Scan(&result)

	return dtos.MemberDashboardLoanStat{
		TotalLoans:    result.TotalLoans,
		PendingLoans:  result.PendingLoans,
		ApprovedLoans: result.ApprovedLoans,
		ByStatus: []dtos.MemberDashboardStatusCount{
			{Status: "Pending", Count: result.PendingLoans},
			{Status: "Approved", Count: result.ApprovedLoans},
		},
	}
}

func (s *MemberDashboardService) GetDashboard() dtos.MemberDashboardResponse {
	return dtos.MemberDashboardResponse{
		MemberStats: s.getMemberStats(),
		LoanStat:    s.getLoanStats(),
	}
}
