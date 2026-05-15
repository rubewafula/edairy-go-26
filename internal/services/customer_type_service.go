package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type CustomerTypeService struct{}

func NewCustomerTypeService() *CustomerTypeService {
	return &CustomerTypeService{}
}

func (s *CustomerTypeService) CreateCustomerType(req dtos.CreateCustomerTypeRequest, userID uint64) (*models.CustomerType, error) {
	customerType := &models.CustomerType{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(customerType).Error; err != nil {
		return nil, err
	}
	return customerType, nil
}

func (s *CustomerTypeService) GetCustomerTypes(page, limit int) ([]dtos.CustomerTypeResponse, int64, error) {
	var results []dtos.CustomerTypeResponse
	var total int64
	db.DB.Model(&models.CustomerType{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.CustomerType{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *CustomerTypeService) GetCustomerType(id string) (*dtos.CustomerTypeResponse, error) {
	var result dtos.CustomerTypeResponse
	err := db.DB.Model(&models.CustomerType{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *CustomerTypeService) UpdateCustomerType(id string, req dtos.UpdateCustomerTypeRequest, userID uint64) error {
	var customerType models.CustomerType
	if err := db.DB.First(&customerType, id).Error; err != nil {
		return err
	}

	customerType.Name = req.Name
	customerType.Description = req.Description
	customerType.UpdatedBy = userID

	return db.DB.Save(&customerType).Error
}

func (s *CustomerTypeService) DeleteCustomerType(id string) error {
	return db.DB.Delete(&models.CustomerType{}, id).Error
}
