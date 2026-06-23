package dtos

type MemberDashboardResponse struct {
	MemberStats MemberDashboardMemberStats `json:"member_stats"`
	LoanStat    MemberDashboardLoanStat    `json:"loan_stat"`
}

type MemberDashboardMemberStats struct {
	TotalMembers int64                          `json:"total_members"`
	NewMembers   int64                          `json:"new_members"`
	Male         int64                          `json:"male"`
	Female       int64                          `json:"female"`
	Composition  []MemberDashboardStatusCount `json:"composition"`
}

type MemberDashboardLoanStat struct {
	TotalLoans    int64                          `json:"total_loans"`
	PendingLoans  int64                          `json:"pending_loans"`
	ApprovedLoans int64                          `json:"approved_loans"`
	ByStatus      []MemberDashboardStatusCount `json:"by_status"`
}

type MemberDashboardStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
