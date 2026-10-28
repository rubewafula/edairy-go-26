package dtos

import "github.com/gin-gonic/gin"

type MilkJournalListFilters struct {
	DateFrom string
	DateTo   string
	RouteID  string
	ShiftID  string
	MemberNo string
	// Entry-level filters (milk-journal-entries only).
	MemberID string
	Status   string
}

func ParseMilkJournalListFilters(ctx *gin.Context) MilkJournalListFilters {
	return MilkJournalListFilters{
		DateFrom: ctx.Query("date[from]"),
		DateTo:   ctx.Query("date[to]"),
		RouteID:  ctx.Query("route_id"),
		ShiftID:  ctx.Query("milk_delivery_shift_id"),
		MemberNo: ctx.Query("member_no"),
		MemberID: ctx.Query("member_id"),
		Status:   ctx.Query("status"),
	}
}
