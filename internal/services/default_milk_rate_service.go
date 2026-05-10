package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type DefaultMilkRateService struct{}

func NewDefaultMilkRateService() *DefaultMilkRateService {
	return &DefaultMilkRateService{}
}

func (s *DefaultMilkRateService) CreateDefaultMilkRate(req dtos.CreateDefaultMilkRateRequest) (*models.DefaultMilkRate, error) {
	rate := &models.DefaultMilkRate{
		Rate:    req.Rate,
		RouteID: req.RouteID,
	}

	if err := db.DB.Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *DefaultMilkRateService) GetDefaultMilkRates(page, limit int) ([]dtos.DefaultMilkRateResponse, int64, error) {
	var results []dtos.DefaultMilkRateResponse
	var total int64
	db.DB.Model(&models.DefaultMilkRate{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			dmr.id, dmr.rate, dmr.route_id, r.route_name,
			dmr.created_at, dmr.updated_at
		FROM default_milk_rates dmr
		LEFT JOIN routes r ON dmr.route_id = r.id
		WHERE dmr.deleted_at IS NULL
		ORDER BY dmr.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *DefaultMilkRateService) GetDefaultMilkRate(id string) (*dtos.DefaultMilkRateResponse, error) {
	var result dtos.DefaultMilkRateResponse
	query := `
		SELECT 
			dmr.id, dmr.rate, dmr.route_id, r.route_name,
			dmr.created_at, dmr.updated_at
		FROM default_milk_rates dmr
		LEFT JOIN routes r ON dmr.route_id = r.id
		WHERE dmr.id = ? AND dmr.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *DefaultMilkRateService) UpdateDefaultMilkRate(id string, req dtos.UpdateDefaultMilkRateRequest) error {
	var rate models.DefaultMilkRate
	if err := db.DB.First(&rate, id).Error; err != nil {
		return err
	}
	rate.Rate = req.Rate
	rate.RouteID = req.RouteID
	return db.DB.Save(&rate).Error
}

func (s *DefaultMilkRateService) DeleteDefaultMilkRate(id string) error {
	return db.DB.Delete(&models.DefaultMilkRate{}, id).Error
}
