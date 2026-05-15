package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type CustomerInvoiceService struct{}

func NewCustomerInvoiceService() *CustomerInvoiceService {
	return &CustomerInvoiceService{}
}

func (s *CustomerInvoiceService) CreateInvoice(req dtos.CreateCustomerInvoiceRequest, userID uint64) (*models.CustomerInvoice, error) {
	invoice := &models.CustomerInvoice{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		CustomerID:  req.CustomerID,
		BillingID:   req.BillingID,
		InvoiceNo:   fmt.Sprintf("INV-%d", time.Now().Unix()),
		InvoiceDate: utils.ParseDate(req.InvoiceDate),
		DueDate:     utils.ParseDate(req.DueDate),
		GrossAmount: req.GrossAmount,
		TaxAmount:   req.TaxAmount,
		TotalAmount: req.TotalAmount,
		Balance:     req.TotalAmount,
		Status:      "UNPAID",
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		// Update billing status
		return tx.Model(&models.CustomerBilling{}).Where("id = ?", req.BillingID).Updates(map[string]interface{}{
			"status":     "invoiced",
			"invoice_id": invoice.ID,
		}).Error
	})

	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *CustomerInvoiceService) GetInvoices(page, limit int) ([]dtos.CustomerInvoiceResponse, int64, error) {
	var results []dtos.CustomerInvoiceResponse
	var total int64
	db.DB.Model(&models.CustomerInvoice{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT ci.*, c.full_names as customer_name
		FROM customer_invoices ci
		LEFT JOIN customers c ON ci.customer_id = c.id
		WHERE ci.deleted_at IS NULL
		ORDER BY ci.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerInvoiceService) GetInvoice(id string) (*dtos.CustomerInvoiceResponse, error) {
	var result dtos.CustomerInvoiceResponse
	query := `
		SELECT ci.*, c.full_names as customer_name
		FROM customer_invoices ci
		LEFT JOIN customers c ON ci.customer_id = c.id
		WHERE ci.id = ? AND ci.deleted_at IS NULL LIMIT 1
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

func (s *CustomerInvoiceService) DeleteInvoice(id string) error {
	return db.DB.Delete(&models.CustomerInvoice{}, id).Error
}
