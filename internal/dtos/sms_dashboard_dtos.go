package dtos

type SMSDashboardResponse struct {
	SMSSummary     SMSDashboardSummary          `json:"sms_summary"`
	SMSQueueStat   SMSDashboardQueueStat        `json:"sms_queue_stat"`
	ContactStat    SMSDashboardContactStat      `json:"contact_stat"`
	CampaignStat   SMSDashboardCampaignStat     `json:"campaign_stat"`
	CreditStat     SMSDashboardCreditStat       `json:"credit_stat"`
	DeliveryTrend  []SMSDashboardDeliveryPoint  `json:"delivery_trend"`
	MonthlySMSUsage []SMSDashboardMonthlyUsage  `json:"monthly_sms_usage"`
}

type SMSDashboardSummary struct {
	TotalSMSSentThisMonth       int64   `json:"total_sms_sent_this_month"`
	TotalDeliveredThisMonth     int64   `json:"total_delivered_this_month"`
	TotalFailedThisMonth        int64   `json:"total_failed_this_month"`
	DeliveryRateThisMonth       float64 `json:"delivery_rate_this_month"`
	TotalSpendThisMonth         float64 `json:"total_spend_this_month"`
	CostPerSMS                  float64 `json:"cost_per_sms"`
	AverageDeliveryTimeSeconds  float64 `json:"average_delivery_time_seconds"`
}

type SMSDashboardQueueStat struct {
	PendingMessages int64 `json:"pending_messages"`
	QueuedMessages  int64 `json:"queued_messages"`
}

type SMSDashboardContactStat struct {
	TotalContacts int64 `json:"total_contacts"`
	OptOutCount   int64 `json:"opt_out_count"`
}

type SMSDashboardCampaignStat struct {
	ActiveCampaigns    int64   `json:"active_campaigns"`
	ScheduledCampaigns int64   `json:"scheduled_campaigns"`
	OpenRate           float64 `json:"open_rate"`
}

type SMSDashboardCreditStat struct {
	SMSCreditsRemaining int64 `json:"sms_credits_remaining"`
	CreditsUsedToday    int64 `json:"credits_used_today"`
}

type SMSDashboardDeliveryPoint struct {
	Date      string `json:"date"`
	Sent      int64  `json:"sent"`
	Delivered int64  `json:"delivered"`
	Failed    int64  `json:"failed"`
}

type SMSDashboardMonthlyUsage struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}
