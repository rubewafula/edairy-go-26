package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type CustomerClassService struct{}

func NewCustomerClassService() *CustomerClassService {
	return &CustomerClassService{}
}

func (s *CustomerClassService) CreateCustomerClass(req dtos.CreateCustomerClassRequest) (*models.CustomerClass, error) {
	class := &models.CustomerClass{
		ClassCode:   req.ClassCode,
		Description: req.Description,
	}

	if err := db.DB.Create(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}

func (s *CustomerClassService) GetCustomerClasses() ([]models.CustomerClass, int64, error) {
	var classes []models.CustomerClass
	var total int64
	db.DB.Model(&models.CustomerClass{}).Count(&total)
	err := db.DB.Find(&classes).Error
	return classes, total, err
}

func (s *CustomerClassService) GetCustomerClass(id string) (*models.CustomerClass, error) {
	var class models.CustomerClass
	if err := db.DB.First(&class, id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (s *CustomerClassService) UpdateCustomerClass(id string, req dtos.UpdateCustomerClassRequest) error {
	var class models.CustomerClass
	if err := db.DB.First(&class, id).Error; err != nil {
		return err
	}

	class.ClassCode = req.ClassCode
	class.Description = req.Description

	return db.DB.Save(&class).Error
}

func (s *CustomerClassService) DeleteCustomerClass(id string) error {
	return db.DB.Delete(&models.CustomerClass{}, id).Error
}
