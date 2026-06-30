package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

type SMSDashboardService struct{}

func NewSMSDashboardService() *SMSDashboardService {
	return &SMSDashboardService{}
}

func (s *SMSDashboardService) getSMSSummary() dtos.SMSDashboardSummary {
	var result struct {
		TotalSMSSentThisMonth      int64
		TotalDeliveredThisMonth    int64
		TotalFailedThisMonth       int64
		AverageDeliveryTimeSeconds float64
	}

	db.DB.Raw(`
		SELECT
			SUM(CASE WHEN UPPER(status) IN ('SENT', 'DELIVERED') THEN 1 ELSE 0 END) AS total_sms_sent_this_month,
			SUM(CASE WHEN UPPER(status) = 'DELIVERED' THEN 1 ELSE 0 END) AS total_delivered_this_month,
			SUM(CASE WHEN UPPER(status) = 'FAILED' THEN 1 ELSE 0 END) AS total_failed_this_month,
			COALESCE(AVG(CASE
				WHEN delivered_at IS NOT NULL AND sent_at IS NOT NULL
				THEN TIMESTAMPDIFF(SECOND, sent_at, delivered_at)
			END), 0) AS average_delivery_time_seconds
		FROM sms_outboxes
		WHERE created_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
		  AND deleted_at IS NULL
	`).Scan(&result)

	deliveryRate := float64(0)
	if result.TotalSMSSentThisMonth > 0 {
		deliveryRate = float64(result.TotalDeliveredThisMonth) / float64(result.TotalSMSSentThisMonth) * 100
	}

	const costPerSMS = 0.25
	totalSpend := float64(result.TotalSMSSentThisMonth) * costPerSMS

	return dtos.SMSDashboardSummary{
		TotalSMSSentThisMonth:      result.TotalSMSSentThisMonth,
		TotalDeliveredThisMonth:    result.TotalDeliveredThisMonth,
		TotalFailedThisMonth:       result.TotalFailedThisMonth,
		DeliveryRateThisMonth:      deliveryRate,
		TotalSpendThisMonth:        totalSpend,
		CostPerSMS:                 costPerSMS,
		AverageDeliveryTimeSeconds: result.AverageDeliveryTimeSeconds,
	}
}

func (s *SMSDashboardService) getSMSQueueStat() dtos.SMSDashboardQueueStat {
	var result dtos.SMSDashboardQueueStat

	db.DB.Raw(`
		SELECT
			(SELECT COUNT(*)
			 FROM sms_outboxes
			 WHERE UPPER(status) = 'PENDING'
			   AND deleted_at IS NULL) AS pending_messages,
			(SELECT COUNT(*)
			 FROM sms_queue
			 WHERE processed = 'no') AS queued_messages
	`).Scan(&result)

	return result
}

func (s *SMSDashboardService) getContactStat() dtos.SMSDashboardContactStat {
	var result dtos.SMSDashboardContactStat

	db.DB.Raw(`
		SELECT
			(SELECT COUNT(*) FROM sms_contacts) AS total_contacts,
			(SELECT COUNT(*)
			 FROM sms_contact_groups
			 WHERE status = 'inactive'
			   AND deleted_at IS NULL) AS opt_out_count
	`).Scan(&result)

	return result
}

func (s *SMSDashboardService) getCampaignStat() dtos.SMSDashboardCampaignStat {
	var result struct {
		ActiveCampaigns    int64
		ScheduledCampaigns int64
		TotalRecipients    int64
		TotalSent          int64
	}

	db.DB.Raw(`
		SELECT
			SUM(CASE WHEN status IN ('running', 'processing') THEN 1 ELSE 0 END) AS active_campaigns,
			SUM(CASE WHEN status = 'scheduled' THEN 1 ELSE 0 END) AS scheduled_campaigns,
			COALESCE(SUM(total_recipients), 0) AS total_recipients,
			COALESCE(SUM(total_sent), 0) AS total_sent
		FROM sms_campaigns
		WHERE deleted_at IS NULL
	`).Scan(&result)

	openRate := float64(0)
	if result.TotalRecipients > 0 {
		openRate = float64(result.TotalSent) / float64(result.TotalRecipients) * 100
	}

	return dtos.SMSDashboardCampaignStat{
		ActiveCampaigns:    result.ActiveCampaigns,
		ScheduledCampaigns: result.ScheduledCampaigns,
		OpenRate:           openRate,
	}
}

func (s *SMSDashboardService) getCreditStat() dtos.SMSDashboardCreditStat {
	var result struct {
		CreditsUsedToday int64
	}

	db.DB.Raw(`
		SELECT COUNT(*) AS credits_used_today
		FROM sms_outboxes
		WHERE DATE(created_at) = CURDATE()
		  AND UPPER(status) IN ('SENT', 'DELIVERED')
		  AND deleted_at IS NULL
	`).Scan(&result)

	return dtos.SMSDashboardCreditStat{
		SMSCreditsRemaining: 0,
		CreditsUsedToday:    result.CreditsUsedToday,
	}
}

func (s *SMSDashboardService) getDeliveryTrend() []dtos.SMSDashboardDeliveryPoint {
	var results []dtos.SMSDashboardDeliveryPoint

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(dates.date, '%Y-%m-%d') AS date,
			COALESCE(stats.sent, 0) AS sent,
			COALESCE(stats.delivered, 0) AS delivered,
			COALESCE(stats.failed, 0) AS failed
		FROM (
			SELECT DISTINCT DATE(created_at) AS date
			FROM sms_outboxes
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			  AND deleted_at IS NULL
		) dates
		LEFT JOIN (
			SELECT
				DATE(created_at) AS date,
				SUM(CASE WHEN UPPER(status) IN ('SENT', 'DELIVERED') THEN 1 ELSE 0 END) AS sent,
				SUM(CASE WHEN UPPER(status) = 'DELIVERED' THEN 1 ELSE 0 END) AS delivered,
				SUM(CASE WHEN UPPER(status) = 'FAILED' THEN 1 ELSE 0 END) AS failed
			FROM sms_outboxes
			WHERE DATE(created_at) BETWEEN DATE_SUB(CURDATE(), INTERVAL 5 DAY) AND DATE_SUB(CURDATE(), INTERVAL 1 DAY)
			  AND deleted_at IS NULL
			GROUP BY DATE(created_at)
		) stats ON stats.date = dates.date
		ORDER BY dates.date ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.SMSDashboardDeliveryPoint{}
	}

	return results
}

func (s *SMSDashboardService) getMonthlySMSUsage() []dtos.SMSDashboardMonthlyUsage {
	var results []dtos.SMSDashboardMonthlyUsage

	db.DB.Raw(`
		SELECT
			DATE_FORMAT(month_start, '%b') AS month,
			COALESCE(usage.count, 0) AS count
		FROM (
			SELECT DATE_FORMAT(DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL n MONTH), '%Y-%m-01') AS month_start
			FROM (
				SELECT 0 AS n UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
			) months
		) months
		LEFT JOIN (
			SELECT
				DATE_FORMAT(created_at, '%Y-%m-01') AS month_start,
				COUNT(*) AS count
			FROM sms_outboxes
			WHERE created_at >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 5 MONTH), '%Y-%m-01')
			  AND UPPER(status) IN ('SENT', 'DELIVERED')
			  AND deleted_at IS NULL
			GROUP BY DATE_FORMAT(created_at, '%Y-%m-01')
		) usage ON usage.month_start = months.month_start
		ORDER BY months.month_start ASC
	`).Scan(&results)

	if results == nil {
		return []dtos.SMSDashboardMonthlyUsage{}
	}

	return results
}

func (s *SMSDashboardService) GetDashboard() dtos.SMSDashboardResponse {
	return dtos.SMSDashboardResponse{
		SMSSummary:      s.getSMSSummary(),
		SMSQueueStat:    s.getSMSQueueStat(),
		ContactStat:     s.getContactStat(),
		CampaignStat:    s.getCampaignStat(),
		CreditStat:      s.getCreditStat(),
		DeliveryTrend:   s.getDeliveryTrend(),
		MonthlySMSUsage: s.getMonthlySMSUsage(),
	}
}
