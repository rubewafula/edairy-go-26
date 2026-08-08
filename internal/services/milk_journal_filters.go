package services

import (
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

func appendMilkJournalHeaderFilters(
	whereClauses []string,
	args []interface{},
	filters dtos.MilkJournalListFilters,
) ([]string, []interface{}) {
	if filters.DateFrom != "" {
		whereClauses = append(whereClauses, "DATE(mj.journal_date) >= ?")
		args = append(args, filters.DateFrom)
	}
	if filters.DateTo != "" {
		whereClauses = append(whereClauses, "DATE(mj.journal_date) <= ?")
		args = append(args, filters.DateTo)
	}
	if filters.RouteID != "" {
		whereClauses = append(whereClauses, "mj.route_id = ?")
		args = append(args, filters.RouteID)
	}
	if filters.ShiftID != "" {
		whereClauses = append(whereClauses, "mj.milk_delivery_shift_id = ?")
		args = append(args, filters.ShiftID)
	}
	if strings.TrimSpace(filters.MemberNo) != "" {
		whereClauses = append(whereClauses, `EXISTS (
			SELECT 1
			FROM milk_journal_entries mje_filter
			INNER JOIN member_registrations m_filter ON m_filter.id = mje_filter.member_id
			WHERE mje_filter.milk_journal_id = mj.id
			  AND mje_filter.deleted_at IS NULL
			  AND m_filter.member_no LIKE ?
		)`)
		args = append(args, "%"+strings.TrimSpace(filters.MemberNo)+"%")
	}

	return whereClauses, args
}

func appendMilkJournalEntryFilters(
	whereClauses []string,
	args []interface{},
	filters dtos.MilkJournalListFilters,
) ([]string, []interface{}) {
	if filters.DateFrom != "" {
		whereClauses = append(whereClauses, "DATE(mj.journal_date) >= ?")
		args = append(args, filters.DateFrom)
	}
	if filters.DateTo != "" {
		whereClauses = append(whereClauses, "DATE(mj.journal_date) <= ?")
		args = append(args, filters.DateTo)
	}
	if filters.RouteID != "" {
		whereClauses = append(whereClauses, "mj.route_id = ?")
		args = append(args, filters.RouteID)
	}
	if filters.ShiftID != "" {
		whereClauses = append(whereClauses, "mj.milk_delivery_shift_id = ?")
		args = append(args, filters.ShiftID)
	}
	if strings.TrimSpace(filters.MemberNo) != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+strings.TrimSpace(filters.MemberNo)+"%")
	}

	return whereClauses, args
}
