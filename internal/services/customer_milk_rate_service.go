package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type CustomerMilkRateService struct{}

func NewCustomerMilkRateService() *CustomerMilkRateService {
	return &CustomerMilkRateService{}
}

func (s *CustomerMilkRateService) CreateCustomerMilkRate(req dtos.CreateCustomerMilkRateRequest) (*models.CustomerMilkRate, error) {
	rate := &models.CustomerMilkRate{
		CustomerID:   req.CustomerID,
		Rate:         req.Rate,
		GradeID:      req.GradeID,
		PayDateRange: req.PayDateRange,
	}

	if err := db.DB.Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *CustomerMilkRateService) GetCustomerMilkRates() ([]models.CustomerMilkRate, int64, error) {
	var rates []models.CustomerMilkRate
	var total int64
	db.DB.Model(&models.CustomerMilkRate{}).Count(&total)
	err := db.DB.Find(&rates).Error
	return rates, total, err
}

func (s *CustomerMilkRateService) GetCustomerMilkRate(id string) (*models.CustomerMilkRate, error) {
	var rate models.CustomerMilkRate
	if err := db.DB.First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (s *CustomerMilkRateService) UpdateCustomerMilkRate(id string, req dtos.UpdateCustomerMilkRateRequest) error {
	var rate models.CustomerMilkRate
	if err := db.DB.First(&rate, id).Error; err != nil {
		return err
	}

	rate.CustomerID = req.CustomerID
	rate.Rate = req.Rate
	rate.GradeID = req.GradeID
	rate.PayDateRange = req.PayDateRange

	return db.DB.Save(&rate).Error
}

func (s *CustomerMilkRateService) DeleteCustomerMilkRate(id string) error {
	return db.DB.Delete(&models.CustomerMilkRate{}, id).Error
}
