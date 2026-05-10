package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type DailyMilkVarianceService struct{}

func NewDailyMilkVarianceService() *DailyMilkVarianceService {
	return &DailyMilkVarianceService{}
}

func (s *DailyMilkVarianceService) GetDailyVariances(page, limit int) ([]dtos.DailyMilkVarianceResponse, int64, error) {
	var results []dtos.DailyMilkVarianceResponse
	var total int64
	db.DB.Model(&models.DailyMilkVariance{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.DailyMilkVariance{}).
		Limit(limit).Offset(offset).Order("day DESC").Scan(&results).Error
	return results, total, err
}
