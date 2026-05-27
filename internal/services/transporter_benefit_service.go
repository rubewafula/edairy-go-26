package services

import (
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

func (s *TransporterBenefitService) CreateBenefit(req dtos.CreateTransporterBenefitRequest, userID uint64) (*models.TransporterBenefit, error) {
	benefit := &models.TransporterBenefit{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		Name:        req.Name,
		MinQuantity: req.MinQuantity,
		Rate:        req.Rate,
		Status:      req.Status,
		StartDate:   utils.ParseDatePtr(req.StartDate),
		EndDate:     utils.ParseDatePtr(req.EndDate),
	}

	if req.RouteID != 0 {
		routeID := req.RouteID
		benefit.RouteID = &routeID
	}

	if err := db.DB.Create(benefit).Error; err != nil {
		return nil, err
	}
	return benefit, nil
}

func (s *TransporterBenefitService) GetBenefits(routeID string, page, limit int) ([]dtos.TransporterBenefitResponse, int64, error) {
	var results []dtos.TransporterBenefitResponse
	var total int64
	query := db.DB.Model(&models.TransporterBenefit{}).Where("transporter_benefits.deleted_at IS NULL")

	if routeID != "" {
		query = query.Where("transporter_benefits.route_id = ?", routeID)
	}

	query.Count(&total)
	offset := (page - 1) * limit

	err := query.Select("transporter_benefits.*, routes.route_name").
		Joins("LEFT JOIN routes ON transporter_benefits.route_id = routes.id").
		Limit(limit).Offset(offset).Order("transporter_benefits.id DESC").Scan(&results).Error
	return results, total, err
}

func (s *TransporterBenefitService) GetBenefit(id string) (*dtos.TransporterBenefitResponse, error) {
	var result dtos.TransporterBenefitResponse
	query := `
		SELECT
			tb.*, r.route_name
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

func (s *TransporterBenefitService) UpdateBenefit(id string, req dtos.UpdateTransporterBenefitRequest, userID uint64) error {
	var benefit models.TransporterBenefit
	if err := db.DB.First(&benefit, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":         req.Name,
		"min_quantity": req.MinQuantity,
		"rate":         req.Rate,
		"status":       req.Status,
		"updated_by":   userID,
	}
	if req.StartDate != "" {
		updates["start_date"] = utils.ParseDatePtr(req.StartDate)
	}
	if req.EndDate != "" {
		updates["end_date"] = utils.ParseDatePtr(req.EndDate)
	}
	if req.RouteID != nil {
		updates["route_id"] = req.RouteID
	} else {
		updates["route_id"] = nil
	}

	return db.DB.Model(&benefit).Updates(updates).Error
}

func (s *TransporterBenefitService) DeleteBenefit(id string) error {
	return db.DB.Delete(&models.TransporterBenefit{}, id).Error
}
