package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type CustomerPaymentService struct{}

func NewCustomerPaymentService() *CustomerPaymentService {
	return &CustomerPaymentService{}
}

func (s *CustomerPaymentService) CreatePayment(req dtos.CreateCustomerPaymentRequest, userID uint64) (*models.CustomerPayment, error) {
	payment := &models.CustomerPayment{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		CustomerID:    req.CustomerID,
		InvoiceID:     req.InvoiceID,
		ReceiptNumber: fmt.Sprintf("RCP-%d", utils.Now().Unix()),
		PaymentDate:   utils.ParseDate(req.PaymentDate),
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		ReferenceNo:   req.ReferenceNo,
		Notes:         req.Notes,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// If linked to an invoice, update invoice balance and create allocation
		if req.InvoiceID != nil {
			if err := tx.Model(&models.CustomerInvoice{}).Where("id = ?", *req.InvoiceID).UpdateColumn("balance", gorm.Expr("balance - ?", req.Amount)).Error; err != nil {
				return err
			}

			allocation := models.CustomerPaymentAllocation{
				BaseModel: models.BaseModel{
					CreatedBy: userID,
				},
				InvoiceID:         *req.InvoiceID,
				CustomerPaymentID: payment.ID,
				AllocatedAmount:   req.Amount,
			}
			if err := tx.Create(&allocation).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *CustomerPaymentService) GetPayments(page, limit int) ([]dtos.CustomerPaymentResponse, int64, error) {
	var results []dtos.CustomerPaymentResponse
	var total int64
	db.DB.Model(&models.CustomerPayment{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT cp.*, c.full_names as customer_name, ci.invoice_no
		FROM customer_payments cp
		LEFT JOIN customers c ON cp.customer_id = c.id
		LEFT JOIN customer_invoices ci ON cp.invoice_id = ci.id
		WHERE cp.deleted_at IS NULL
		ORDER BY cp.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerPaymentService) GetPayment(id string) (*dtos.CustomerPaymentResponse, error) {
	var result dtos.CustomerPaymentResponse
	query := `
		SELECT cp.*, c.full_names as customer_name, ci.invoice_no
		FROM customer_payments cp
		LEFT JOIN customers c ON cp.customer_id = c.id
		LEFT JOIN customer_invoices ci ON cp.invoice_id = ci.id
		WHERE cp.id = ? AND cp.deleted_at IS NULL LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, err
}
