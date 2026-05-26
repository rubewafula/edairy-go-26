package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type CustomerService struct{}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) CreateCustomer(req dtos.CreateCustomerRequest) (*dtos.CustomerResponse, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	customer := &models.Customer{
		CustomerTypeID: req.CustomerTypeID,
		FullNames:      req.FullNames,
		Phone:          req.Phone,
		EmailAddress:   req.EmailAddress,
		CustomerNo:     req.CustomerNo,
		KraPin:         req.KraPin,
		Status:         status,
		CreditLimit:    req.CreditLimit,
		PostalAddress:  req.PostalAddress,
		PostalCode:     req.PostalCode,
		PostalTown:     req.PostalTown,
		Terms:          req.Terms,
		Rate:           req.Rate,
	}

	if err := db.DB.Create(customer).Error; err != nil {
		return nil, err
	}
	return s.GetCustomer(fmt.Sprintf("%d", customer.ID))
}

func (s *CustomerService) GetCustomers(page, limit int) ([]dtos.CustomerResponse, int64, error) {
	var results []dtos.CustomerResponse
	var total int64
	db.DB.Model(&models.Customer{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			c.*, ct.description as customer_type_name
		FROM customers c
		LEFT JOIN customer_types ct ON c.customer_type_id = ct.id
		WHERE c.deleted_at IS NULL
		ORDER BY c.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerService) GetCustomer(id string) (*dtos.CustomerResponse, error) {
	var result dtos.CustomerResponse
	query := `
		SELECT 
			c.*, ct.description as customer_type_name
		FROM customers c
		LEFT JOIN customer_types ct ON c.customer_type_id = ct.id
		WHERE c.id = ? AND c.deleted_at IS NULL
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

func (s *CustomerService) UpdateCustomer(id string, req dtos.UpdateCustomerRequest) error {
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		return err
	}

	customer.CustomerTypeID = req.CustomerTypeID
	customer.FullNames = req.FullNames
	customer.Phone = req.Phone
	customer.EmailAddress = req.EmailAddress
	customer.KraPin = req.KraPin
	customer.Status = req.Status
	customer.CreditLimit = req.CreditLimit
	customer.PostalAddress = req.PostalAddress
	customer.PostalCode = req.PostalCode
	customer.PostalTown = req.PostalTown
	customer.Terms = req.Terms
	customer.Rate = req.Rate

	return db.DB.Save(&customer).Error
}

func (s *CustomerService) DeleteCustomer(id string) error {
	return db.DB.Delete(&models.Customer{}, id).Error
}
