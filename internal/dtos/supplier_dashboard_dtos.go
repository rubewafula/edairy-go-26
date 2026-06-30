package dtos

type SupplierDashboardResponse struct {
	SupplyStats       SupplierDashboardSupplyStats       `json:"supply_stats"`
	SupplyVolumeTrend []SupplierDashboardVolumePoint     `json:"supply_volume_trend"`
	SupplyMixToday    []SupplierDashboardMixPoint        `json:"supply_mix_today"`
	SalesMixToday     SupplierDashboardSalesMix          `json:"sales_mix_today"`
}

type SupplierDashboardSupplyStats struct {
	TotalSuppliesThisMonth int64 `json:"total_supplies_this_month"`
	SuppliesToday          int64 `json:"supplies_today"`
	SuppliesThisMonth      int64 `json:"supplies_this_month"`
}

type SupplierDashboardVolumePoint struct {
	Date   string  `json:"date"`
	Volume float64 `json:"volume"`
}

type SupplierDashboardMixPoint struct {
	Item  string  `json:"item"`
	Value float64 `json:"value"`
}

type SupplierDashboardSalesMix struct {
	CashSales   float64 `json:"cash_sales"`
	CreditSales float64 `json:"credit_sales"`
}
