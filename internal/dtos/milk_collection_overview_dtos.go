package dtos

// MilkCollectionOverviewResponse powers the mobile dashboard "Milk Delivery Overview" card.
// Quantities are litres collected from members (milk_journal_entries).
type MilkCollectionOverviewResponse struct {
	PeriodDays            int                           `json:"period_days"`
	StartDate             string                        `json:"start_date"`
	EndDate               string                        `json:"end_date"`
	TotalQuantity         float64                       `json:"total_quantity"`
	DailyAverage          float64                       `json:"daily_average"`
	TodayQuantity         float64                       `json:"today_quantity"`
	PreviousTotalQuantity float64                       `json:"previous_total_quantity"`
	ChangePercent         *float64                      `json:"change_percent"`
	RejectedQuantity      float64                       `json:"rejected_quantity"`
	ActiveMembers         int64                         `json:"active_members"`
	Trend                 []MilkCollectionOverviewPoint `json:"trend"`
	// Scope is "all" (staff), "member" (the logged-in member's own milk) or "transporter".
	Scope     string `json:"scope"`
	ScopeID   uint64 `json:"scope_id,omitempty"`
	ScopeName string `json:"scope_name,omitempty"`
}

type MilkCollectionOverviewPoint struct {
	Date     string  `json:"date"`
	Quantity float64 `json:"quantity"`
}
