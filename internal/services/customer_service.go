package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type CustomerService struct{}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) CreateCustomer(req dtos.CreateCustomerRequest) (*models.Customer, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	customer := &models.Customer{
		ClassID:       req.ClassID,
		FullNames:     req.FullNames,
		Phone:         req.Phone,
		EmailAddress:  req.EmailAddress,
		CustomerNo:    req.CustomerNo,
		KraPin:        req.KraPin,
		Status:        status,
		CreditLimit:   req.CreditLimit,
		PostalAddress: req.PostalAddress,
		PostalCode:    req.PostalCode,
		PostalTown:    req.PostalTown,
		SiteID:        req.SiteID,
		Terms:         req.Terms,
		Rate:          req.Rate,
	}

	if err := db.DB.Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetCustomers() ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64
	db.DB.Model(&models.Customer{}).Count(&total)
	err := db.DB.Find(&customers).Error
	return customers, total, err
}

func (s *CustomerService) GetCustomer(id string) (*models.Customer, error) {
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) UpdateCustomer(id string, req dtos.UpdateCustomerRequest) error {
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		return err
	}

	customer.ClassID = req.ClassID
	customer.FullNames = req.FullNames
	customer.Phone = req.Phone
	customer.EmailAddress = req.EmailAddress
	customer.KraPin = req.KraPin
	customer.Status = req.Status
	customer.CreditLimit = req.CreditLimit
	customer.PostalAddress = req.PostalAddress
	customer.PostalCode = req.PostalCode
	customer.PostalTown = req.PostalTown
	customer.SiteID = req.SiteID
	customer.Terms = req.Terms
	customer.Rate = req.Rate

	return db.DB.Save(&customer).Error
}

func (s *CustomerService) DeleteCustomer(id string) error {
	return db.DB.Delete(&models.Customer{}, id).Error
}
