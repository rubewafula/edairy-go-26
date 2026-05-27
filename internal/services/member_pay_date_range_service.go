package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MemberPayDateRangeService struct{}

func NewMemberPayDateRangeService() *MemberPayDateRangeService {
	return &MemberPayDateRangeService{}
}

func (s *MemberPayDateRangeService) Create(req dtos.CreateMemberPayDateRangeRequest, userID uint64) (*models.MemberPayDateRange, error) {
	dateRange := &models.MemberPayDateRange{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		Name:      req.Name,
		StartDate: utils.ParseDate(req.StartDate),
		EndDate:   utils.ParseDate(req.EndDate),
		PayMonth:  req.PayMonth,
		PayYear:   req.PayYear,
	}

	if err := db.DB.Create(dateRange).Error; err != nil {
		return nil, err
	}
	return dateRange, nil
}

func (s *MemberPayDateRangeService) List(page, limit int) ([]dtos.MemberPayDateRangeResponse, int64, error) {
	var results []dtos.MemberPayDateRangeResponse
	var total int64
	db.DB.Model(&models.MemberPayDateRange{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.MemberPayDateRange{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *MemberPayDateRangeService) Get(id string) (*dtos.MemberPayDateRangeResponse, error) {
	var result dtos.MemberPayDateRangeResponse
	err := db.DB.Model(&models.MemberPayDateRange{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *MemberPayDateRangeService) Update(id string, req dtos.UpdateMemberPayDateRangeRequest, userID uint64) error {
	var dateRange models.MemberPayDateRange
	if err := db.DB.First(&dateRange, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"pay_month":  req.PayMonth,
		"pay_year":   req.PayYear,
		"updated_by": userID,
	}
	if req.StartDate != "" {
		updates["start_date"] = utils.ParseDate(req.StartDate)
	}
	if req.EndDate != "" {
		updates["end_date"] = utils.ParseDate(req.EndDate)
	}

	return db.DB.Model(&dateRange).Updates(updates).Error
}

func (s *MemberPayDateRangeService) Delete(id string) error {
	return db.DB.Delete(&models.MemberPayDateRange{}, id).Error
}
