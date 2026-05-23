package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type CustomerTypeService struct{}

func NewCustomerTypeService() *CustomerTypeService {
	return &CustomerTypeService{}
}

func (s *CustomerTypeService) CreateCustomerType(req dtos.CreateCustomerTypeRequest) (*dtos.CustomerTypeResponse, error) {
	customerType := &models.CustomerType{
		Name:        req.TypeCode,
		Description: req.Description,
	}

	if err := db.DB.Create(customerType).Error; err != nil {
		return nil, err
	}
	return s.GetCustomerType(fmt.Sprintf("%d", customerType.ID))
}

func (s *CustomerTypeService) GetCustomerTypes(page, limit int) ([]dtos.CustomerTypeResponse, int64, error) {
	var results []dtos.CustomerTypeResponse
	var total int64
	db.DB.Model(&models.CustomerType{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.CustomerType{}).
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Scan(&results).Error

	return results, total, err
}

func (s *CustomerTypeService) GetCustomerType(id string) (*dtos.CustomerTypeResponse, error) {
	var customerType dtos.CustomerTypeResponse
	if err := db.DB.Model(&models.CustomerType{}).First(&customerType, id).Error; err != nil {
		return nil, err
	}
	return &customerType, nil
}

func (s *CustomerTypeService) UpdateCustomerType(id string, req dtos.UpdateCustomerTypeRequest) error {
	var customerType models.CustomerType
	if err := db.DB.First(&customerType, id).Error; err != nil {
		return err
	}

	customerType.Name = req.Name
	customerType.Description = req.Description

	return db.DB.Save(&customerType).Error
}

func (s *CustomerTypeService) DeleteCustomerType(id string) error {
	return db.DB.Delete(&models.CustomerType{}, id).Error
}
