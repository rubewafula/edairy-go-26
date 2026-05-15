package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type CustomerPayDateRangeService struct{}

func NewCustomerPayDateRangeService() *CustomerPayDateRangeService {
	return &CustomerPayDateRangeService{}
}

func (s *CustomerPayDateRangeService) CreateCustomerPayDateRange(req dtos.CreateCustomerPayDateRangeRequest, userID uint64) (*models.CustomerPayDateRange, error) {
	dateRange := &models.CustomerPayDateRange{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		Name:      req.Name,
		StartDate: utils.ParseDate(req.StartDate),
		EndDate:   utils.ParseDate(req.EndDate),
	}

	if err := db.DB.Create(dateRange).Error; err != nil {
		return nil, err
	}
	return dateRange, nil
}

func (s *CustomerPayDateRangeService) GetCustomerPayDateRanges(page, limit int) ([]dtos.CustomerPayDateRangeResponse, int64, error) {
	var results []dtos.CustomerPayDateRangeResponse
	var total int64
	db.DB.Model(&models.CustomerPayDateRange{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.CustomerPayDateRange{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *CustomerPayDateRangeService) GetCustomerPayDateRange(id string) (*dtos.CustomerPayDateRangeResponse, error) {
	var result dtos.CustomerPayDateRangeResponse
	err := db.DB.Model(&models.CustomerPayDateRange{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *CustomerPayDateRangeService) UpdateCustomerPayDateRange(id string, req dtos.UpdateCustomerPayDateRangeRequest, userID uint64) error {
	var dateRange models.CustomerPayDateRange
	if err := db.DB.First(&dateRange, id).Error; err != nil {
		return err
	}

	dateRange.Name = req.Name
	dateRange.StartDate = utils.ParseDate(req.StartDate)
	dateRange.EndDate = utils.ParseDate(req.EndDate)
	dateRange.UpdatedBy = userID

	return db.DB.Save(&dateRange).Error
}

func (s *CustomerPayDateRangeService) DeleteCustomerPayDateRange(id string) error {
	return db.DB.Delete(&models.CustomerPayDateRange{}, id).Error
}
