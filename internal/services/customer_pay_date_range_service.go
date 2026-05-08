package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type CustomerPayDateRangeService struct{}

func NewCustomerPayDateRangeService() *CustomerPayDateRangeService {
	return &CustomerPayDateRangeService{}
}

func (s *CustomerPayDateRangeService) CreateCustomerPayDateRange(req dtos.CreateCustomerPayDateRangeRequest) (*models.CustomerPayDateRange, error) {
	dateRange := &models.CustomerPayDateRange{
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

func (s *CustomerPayDateRangeService) GetCustomerPayDateRanges() ([]models.CustomerPayDateRange, int64, error) {
	var ranges []models.CustomerPayDateRange
	var total int64
	db.DB.Model(&models.CustomerPayDateRange{}).Count(&total)
	err := db.DB.Find(&ranges).Error
	return ranges, total, err
}

func (s *CustomerPayDateRangeService) GetCustomerPayDateRange(id string) (*models.CustomerPayDateRange, error) {
	var dateRange models.CustomerPayDateRange
	if err := db.DB.First(&dateRange, id).Error; err != nil {
		return nil, err
	}
	return &dateRange, nil
}

func (s *CustomerPayDateRangeService) UpdateCustomerPayDateRange(id string, req dtos.UpdateCustomerPayDateRangeRequest) error {
	var dateRange models.CustomerPayDateRange
	if err := db.DB.First(&dateRange, id).Error; err != nil {
		return err
	}

	dateRange.Name = req.Name
	dateRange.StartDate = utils.ParseDate(req.StartDate)
	dateRange.EndDate = utils.ParseDate(req.EndDate)
	dateRange.PayMonth = req.PayMonth
	dateRange.PayYear = req.PayYear

	return db.DB.Save(&dateRange).Error
}

func (s *CustomerPayDateRangeService) DeleteCustomerPayDateRange(id string) error {
	return db.DB.Delete(&models.CustomerPayDateRange{}, id).Error
}
