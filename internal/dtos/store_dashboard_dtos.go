package dtos

type StoreDashboardResponse struct {
	StoreSaleStat StoreDashboardStoreSaleStat `json:"store_sale_stat"`
	SupplyStats   StoreDashboardSupplyStats   `json:"supply_stats"`
}

type StoreDashboardStoreSaleStat struct {
	CashStoreSaleToday              float64                         `json:"cash_store_sale_today"`
	CreditStoreSaleToday            float64                         `json:"credit_store_sale_today"`
	TotalStoreSaleToday             float64                         `json:"total_store_sale_today"`
	CashStoreSaleThisMonth          float64                         `json:"cash_store_sale_this_month"`
	CreditStoreSaleThisMonth        float64                         `json:"credit_store_sale_this_month"`
	TotalStoreSaleThisMonth         float64                         `json:"total_store_sale_this_month"`
	TotalSalesTransactionsThisMonth int64                           `json:"total_sales_transactions_this_month"`
	Currency                        string                          `json:"currency"`
	SalesTrend                      []StoreDashboardSalesTrendPoint `json:"sales_trend"`
}

type StoreDashboardSalesTrendPoint struct {
	Date       string  `json:"date"`
	TotalSales float64 `json:"total_sales"`
}

type StoreDashboardSupplyStats struct {
	SuppliesToday     int64 `json:"supplies_today"`
	SuppliesThisMonth int64 `json:"supplies_this_month"`
}
