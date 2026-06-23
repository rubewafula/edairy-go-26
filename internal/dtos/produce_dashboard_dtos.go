package dtos

type ProduceDashboardResponse struct {
	MilkCollectionStat     ProduceDashboardQuantityStat      `json:"milk_collection_stat"`
	MilkDeliveryStat       ProduceDashboardQuantityStat      `json:"milk_delivery_stat"`
	MilkRejectStat         ProduceDashboardQuantityStat      `json:"milk_reject_stat"`
	CollectionVsRejectTrend []ProduceDashboardTrendPoint      `json:"collection_vs_reject_trend"`
}

type ProduceDashboardQuantityStat struct {
	Today float64 `json:"today"`
	Month float64 `json:"month"`
}

type ProduceDashboardTrendPoint struct {
	Date        string  `json:"date"`
	Collections float64 `json:"collections"`
	Rejects     float64 `json:"rejects"`
}
