package services

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type TransporterBenefitService struct{}

func NewTransporterBenefitService() *TransporterBenefitService {
	return &TransporterBenefitService{}
}

func (s *TransporterBenefitService) CreateBenefit(req dtos.CreateTransporterBenefitRequest) (*models.TransporterBenefit, error) {
	status := req.Status
	if status == "" {
		status = "1"
	}

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		t := utils.ParseDate(req.StartDate)
		startDate = &t
	}
	if req.EndDate != "" {
		t := utils.ParseDate(req.EndDate)
		endDate = &t
	}

	benefit := &models.TransporterBenefit{
		Name:        req.Name,
		MinQuantity: req.MinQuantity,
		Rate:        req.Rate,
		RouteID:     req.RouteID,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := db.DB.Create(benefit).Error; err != nil {
		return nil, err
	}
	return benefit, nil
}

func (s *TransporterBenefitService) GetBenefits() ([]dtos.TransporterBenefitResponse, int64, error) {
	var results []dtos.TransporterBenefitResponse
	var total int64
	db.DB.Model(&models.TransporterBenefit{}).Count(&total)

	query := `
		SELECT 
			tb.id, tb.name, tb.min_quantity, tb.rate,
			tb.route_id, r.route_name,
			tb.status, tb.start_date, tb.end_date,
			tb.created_at, tb.updated_at
		FROM transporter_benefits tb
		LEFT JOIN routes r ON tb.route_id = r.id
		WHERE tb.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *TransporterBenefitService) GetBenefit(id string) (*dtos.TransporterBenefitResponse, error) {
	var result dtos.TransporterBenefitResponse
	query := `
		SELECT 
			tb.id, tb.name, tb.min_quantity, tb.rate,
			tb.route_id, r.route_name,
			tb.status, tb.start_date, tb.end_date,
			tb.created_at, tb.updated_at
		FROM transporter_benefits tb
		LEFT JOIN routes r ON tb.route_id = r.id
		WHERE tb.id = ? AND tb.deleted_at IS NULL
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

func (s *TransporterBenefitService) UpdateBenefit(id string, req dtos.UpdateTransporterBenefitRequest) error {
	updates := map[string]interface{}{
		"name":         req.Name,
		"min_quantity": req.MinQuantity,
		"rate":         req.Rate,
		"route_id":     req.RouteID,
		"status":       req.Status,
	}

	if req.StartDate != "" {
		t := utils.ParseDate(req.StartDate)
		updates["start_date"] = &t
	} else {
		updates["start_date"] = nil
	}

	if req.EndDate != "" {
		t := utils.ParseDate(req.EndDate)
		updates["end_date"] = &t
	} else {
		updates["end_date"] = nil
	}

	return db.DB.Model(&models.TransporterBenefit{}).Where("id = ?", id).Updates(updates).Error
}

func (s *TransporterBenefitService) DeleteBenefit(id string) error {
	return db.DB.Delete(&models.TransporterBenefit{}, id).Error
}
