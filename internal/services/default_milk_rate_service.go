package services

import (
	"log"
	"fmt"
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
		Rate:        req.Rate,
		AdvanceRate: positiveOrNil(req.AdvanceRate.Ptr()),
	}

	if req.RouteID != 0 {
		routeID := req.RouteID
		rate.RouteID = &routeID
	}

	if req.MemberID != 0 {
		memberID := req.MemberID
		rate.MemberID = &memberID
	}

	if err := db.DB.Create(rate).Error; err != nil {
		return nil, err
	}
	recalcCreditAfterMilkRateChange(fmt.Sprintf("Default milk rate added: KES %.2f", rate.Rate))
	return rate, nil
}

func (s *DefaultMilkRateService) GetDefaultMilkRates(page, limit int) ([]dtos.DefaultMilkRateResponse, int64, error) {
	var results []dtos.DefaultMilkRateResponse
	var total int64
	db.DB.Model(&models.DefaultMilkRate{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			dmr.id, dmr.rate, dmr.advance_rate, COALESCE(dmr.route_id, 0) AS route_id, COALESCE(dmr.member_id, 0) AS member_id, r.route_name,
			concat(member_no, ' ', first_name, ' ', last_name) as member_name,
			dmr.created_at, dmr.updated_at
		FROM default_milk_rates dmr
		LEFT JOIN routes r ON dmr.route_id = r.id
		LEFT JOIN member_registrations m ON m.id = dmr.member_id
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
			dmr.id, dmr.rate, dmr.advance_rate, COALESCE(dmr.route_id, 0) AS route_id, COALESCE(dmr.member_id, 0) AS member_id, r.route_name, 
			concat(member_no, ' ', first_name, ' ', last_name) as member_name,
			dmr.created_at, dmr.updated_at
		FROM default_milk_rates dmr
		LEFT JOIN routes r ON dmr.route_id = r.id
		LEFT JOIN member_registrations m ON m.id = dmr.member_id
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
	rate.AdvanceRate = positiveOrNil(req.AdvanceRate.Ptr())

	if req.RouteID != 0 {
		routeID := req.RouteID
		rate.RouteID = &routeID
	} else {
		rate.RouteID = nil
	}

	if req.MemberID != 0 {
		memberID := req.MemberID
		rate.MemberID = &memberID
	} else {
		rate.MemberID = nil
	}

	if err := db.DB.Save(&rate).Error; err != nil {
		return err
	}
	recalcCreditAfterMilkRateChange(fmt.Sprintf("Default milk rate changed: KES %.2f", rate.Rate))
	return nil
}

func (s *DefaultMilkRateService) DeleteDefaultMilkRate(id string) error {
	if err := db.DB.Delete(&models.DefaultMilkRate{}, id).Error; err != nil {
		return err
	}
	recalcCreditAfterMilkRateChange("Default milk rate removed")
	return nil
}

func positiveOrNil(v *float64) *float64 {
	if v == nil || *v <= 0 {
		return nil
	}
	x := *v
	return &x
}

// EnsureAdvanceRateSchema adds default_milk_rates.advance_rate.
func EnsureAdvanceRateSchema() {
	if !db.DB.Migrator().HasColumn(&models.DefaultMilkRate{}, "advance_rate") {
		if err := db.DB.Migrator().AddColumn(&models.DefaultMilkRate{}, "AdvanceRate"); err != nil {
			log.Printf("[MilkRates] add advance_rate: %v", err)
		}
	}
}
