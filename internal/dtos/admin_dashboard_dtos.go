package dtos

type AdminDashboardResponse struct {
	SupplyStat         AdminDashboardSuppliesStat       `json:"supply_stats"`
	StoreSaleStat      AdminDashboardStoreSaleStat      `json:"store_sale_stat"`
	MemberStat         AdminDashboardMemberStat         `json:"member_stats"`
	MilkCollectionStat AdminDashboardMilkCollectionStat `json:"milk_collection_stat"`
	MilkDevliveryStat  AdminDashboardMilkDeliveryStat   `json:"milk_delivery_stat"`
	MilkRejectStat     AdminDashboardMilkRejectStat     `json:"milk_reject_stat"`
	LoanStat           AdminDashboardLoanStat           `json:"loan_stat"`
	TopRoutes          []AdminDashboardRouteStat        `json:"top_routes"`
}

type AdminDashboardRouteStat struct {
	Route string  `json:"route"`
	Total float64 `json:"total"`
}

type AdminDashboardDailyMilkStat struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type AdminDashboardMilkDeliveryStat struct {
	Today float64 `json:"today"`
	Month float64 `json:"month"`
}

type AdminDashboardMilkRejectStat struct {
	Today float64 `json:"today"`
	Month float64 `json:"month"`
}

type AdminDashboardMilkCollectionStat struct {
	MilkToday     float64                       `json:"milk_today"`
	MilkThisMonth float64                       `json:"milk_this_month"`
	MilkTrend     []AdminDashboardDailyMilkStat `json:"milk_trend"`
}

type AdminDashboardSuppliesStat struct {
	SupplyItemsToday     int64   `json:"supplied_items_today"`
	SupplyItemsThisMonth int64   `json:"supplied_items_this_month"`
	SuppliesToday        float64 `json:"supplies_today"`
	SuppliesThisMonth    float64 `json:"supplies_this_month"`
}

type AdminDashboardMemberStat struct {
	Total  int64 `json:"total_members"`
	Active int64 `json:"active_members"`
	Male   int64 `json:"male"`
	Female int64 `json:"female"`
}

type AdminDashboardLoanStat struct {
	Total              int64   `json:"total_loans"`
	Pending            int64   `json:"pending_loans"`
	TotalAmount        float64 `json:"total_amount"`
	TotalAmountPending float64 `json:"total_amount_pending"`
}

type AdminDashboardStoreSaleStat struct {
	CashStoreSaleToday      float64 `json:"cash_store_sale_today"`
	CashStoreSaleThisMonth  float64 `json:"cash_store_sale_this_month"`
	TotalStoreSaleToday     float64 `json:"total_store_sale_today"`
	TotalStoreSaleThisMonth float64 `json:"total_store_sale_this_month"`
}
