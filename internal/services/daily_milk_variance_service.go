package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/debuglog"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type DailyMilkVarianceService struct{}

func NewDailyMilkVarianceService() *DailyMilkVarianceService {
	return &DailyMilkVarianceService{}
}

func (s *DailyMilkVarianceService) GetDailyVariances(page, limit int) ([]dtos.DailyMilkVarianceResponse, int64, error) {
	// #region agent log
	debuglog.Log("daily_milk_variance_service.go:GetDailyVariances", "entry", "A", "pre-fix", map[string]any{
		"page": page, "limit": limit,
	})
	// #endregion

	var rows []models.DailyMilkVariance
	var total int64
	db.DB.Model(&models.DailyMilkVariance{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.DailyMilkVariance{}).
		Select("id", "transporter", "transporter_id", "day", "month",
			"field_collections", "mcc", "cash_sales", "credit_sales", "rejects", "balance").
		Limit(limit).Offset(offset).Order("day DESC").
		Find(&rows).Error
	if err != nil {
		// #region agent log
		debuglog.Log("daily_milk_variance_service.go:GetDailyVariances", "find_failed", "A", "pre-fix", map[string]any{
			"error": err.Error(),
		})
		// #endregion
		return nil, 0, err
	}

	results := make([]dtos.DailyMilkVarianceResponse, len(rows))
	for i, r := range rows {
		results[i] = dtos.DailyMilkVarianceResponse{
			ID:               r.ID,
			Transporter:      r.Transporter,
			TransporterID:    r.TransporterID,
			Day:              r.Day,
			Month:            r.Month,
			FieldCollections: r.FieldCollections,
			MCC:              r.MCC,
			CashSales:        r.CashSales,
			CreditSales:      r.CreditSales,
			Rejects:          r.Rejects,
			Balance:          r.Balance,
		}
	}

	// #region agent log
	debuglog.Log("daily_milk_variance_service.go:GetDailyVariances", "success", "A", "pre-fix", map[string]any{
		"total": total, "returned": len(results),
	})
	// #endregion

	return results, total, nil
}
