package dtos

type TransporterDashboardResponse struct {
	MilkDeliveryStat        TransporterDashboardMilkDeliveryStat   `json:"milk_delivery_stat"`
	MilkRejectStat          TransporterDashboardMilkRejectStat     `json:"milk_reject_stat"`
	MilkCollectionStat      TransporterDashboardMilkCollectionStat `json:"milk_collection_stat"`
	TransportOperationsTrend []TransporterDashboardTrendPoint       `json:"transport_operations_trend"`
}

type TransporterDashboardMilkDeliveryStat struct {
	Today float64 `json:"today"`
	Month float64 `json:"month"`
}

type TransporterDashboardMilkRejectStat struct {
	Today float64 `json:"today"`
	Month float64 `json:"month"`
}

type TransporterDashboardMilkCollectionStat struct {
	MilkToday     float64 `json:"milk_today"`
	MilkThisMonth float64 `json:"milk_this_month"`
}

type TransporterDashboardTrendPoint struct {
	Date       string  `json:"date"`
	Deliveries float64 `json:"deliveries"`
	Rejects    float64 `json:"rejects"`
}
