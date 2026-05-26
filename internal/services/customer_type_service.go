package services

import (
	"fmt"

	"gorm.io/gorm"

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
		Name:        req.Name,
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
	db.DB.Model(&models.CustomerType{}).Where("deleted_at IS NULL").Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM customer_types
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error

	return results, total, err
}

func (s *CustomerTypeService) GetCustomerType(id string) (*dtos.CustomerTypeResponse, error) {
	var result dtos.CustomerTypeResponse
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM customer_types
		WHERE id = ? AND deleted_at IS NULL
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
