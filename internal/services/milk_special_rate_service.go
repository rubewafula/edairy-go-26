package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type MilkSpecialRateService struct{}

func NewMilkSpecialRateService() *MilkSpecialRateService {
	return &MilkSpecialRateService{}
}

func (s *MilkSpecialRateService) Create(req dtos.CreateMilkSpecialRateRequest, userID uint64) (*models.MilkSpecialRate, error) {
	rate := &models.MilkSpecialRate{
		BaseModel:             models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		MonthlyPayDateRangeID: req.MonthlyPayDateRangeID,
		Rate:                  req.Rate,
		Confirmed:             0,
	}

	if req.MemberID != 0 {
		memberID := req.MemberID
		rate.MemberID = &memberID
	}

	if req.RouteID != 0 {
		routeID := req.RouteID
		rate.RouteID = &routeID
	}

	if err := db.DB.Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *MilkSpecialRateService) List(payDateRangeID string, page, limit int) ([]dtos.MilkSpecialRateResponse, int64, error) {
	var results []dtos.MilkSpecialRateResponse
	var total int64

	queryBuilder := db.DB.Model(&models.MilkSpecialRate{}).Where("deleted_at IS NULL")
	if payDateRangeID != "" {
		queryBuilder = queryBuilder.Where("pay_date_range_id = ?", payDateRangeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			msr.*, 
			pdr.name as pay_date_range_name,
			m.member_no, 
			CONCAT(m.first_name, ' ', m.last_name) as member_name,
			r.route_name,
			msr.confirmed
		FROM milk_special_rates msr
		LEFT JOIN member_pay_date_ranges pdr ON msr.pay_date_range_id = pdr.id
		LEFT JOIN member_registrations m ON msr.member_id = m.id
		LEFT JOIN routes r ON msr.route_id = r.id
		WHERE msr.deleted_at IS NULL AND (? = '' OR msr.pay_date_range_id = ?)
		ORDER BY msr.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, payDateRangeID, payDateRangeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *MilkSpecialRateService) Get(id string) (*dtos.MilkSpecialRateResponse, error) {
	var result dtos.MilkSpecialRateResponse
	query := `
		SELECT 
			msr.*, 
			pdr.name as pay_date_range_name,
			m.member_no, 
			CONCAT(m.first_name, ' ', m.last_name) as member_name,
			r.route_name,
			msr.confirmed
		FROM milk_special_rates msr
		LEFT JOIN member_pay_date_ranges pdr ON msr.pay_date_range_id = pdr.id
		LEFT JOIN member_registrations m ON msr.member_id = m.id
		LEFT JOIN routes r ON msr.route_id = r.id
		WHERE msr.id = ? AND msr.deleted_at IS NULL
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

func (s *MilkSpecialRateService) Update(id string, req dtos.UpdateMilkSpecialRateRequest, userID uint64) error {
	var rate models.MilkSpecialRate
	if err := db.DB.First(&rate, id).Error; err != nil {
		return err
	}

	if rate.Confirmed != 0 {
		return fmt.Errorf("cannot update a confirmed milk special rate")
	}

	updates := map[string]interface{}{
		"pay_date_range_id": req.PayDateRangeID,
		"rate":              req.Rate,
		"updated_by":        userID,
	}

	if req.MemberID != 0 {
		memberID := req.MemberID
		updates["member_id"] = &memberID
	} else {
		updates["member_id"] = nil
	}

	if req.RouteID != 0 {
		routeID := req.RouteID
		updates["route_id"] = &routeID
	} else {
		updates["route_id"] = nil
	}

	if req.Confirmed != 0 {

		updates["confirmed"] = req.Confirmed
	}

	return db.DB.Model(&rate).Updates(updates).Error
}

func (s *MilkSpecialRateService) Delete(id string) error {
	return db.DB.Delete(&models.MilkSpecialRate{}, id).Error
}
