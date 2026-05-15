package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type CustomerBillingService struct{}

func NewCustomerBillingService() *CustomerBillingService {
	return &CustomerBillingService{}
}

func (s *CustomerBillingService) GetBillings(page, limit int) ([]dtos.CustomerBillingResponse, int64, error) {
	var results []dtos.CustomerBillingResponse
	var total int64
	db.DB.Model(&models.CustomerBilling{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT cb.*, cpdr.name as pay_date_range_name
		FROM customer_billings cb
		LEFT JOIN customer_pay_date_ranges cpdr ON cb.pay_date_range_id = cpdr.id
		WHERE cb.deleted_at IS NULL
		ORDER BY cb.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerBillingService) GetBilling(id string) (*dtos.CustomerBillingResponse, error) {
	var result dtos.CustomerBillingResponse
	query := `
		SELECT cb.*, cpdr.name as pay_date_range_name
		FROM customer_billings cb
		LEFT JOIN customer_pay_date_ranges cpdr ON cb.pay_date_range_id = cpdr.id
		WHERE cb.id = ? AND cb.deleted_at IS NULL LIMIT 1
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

func (s *CustomerBillingService) GetBillingItems(billingID string) ([]dtos.CustomerBillingItemResponse, error) {
	var items []dtos.CustomerBillingItemResponse
	query := `
		SELECT cbi.*, pg.name as grade_name
		FROM customer_billing_items cbi
		LEFT JOIN product_grades pg ON cbi.product_grade_id = pg.id
		WHERE cbi.customer_billing_id = ? AND cbi.deleted_at IS NULL
	`
	err := db.DB.Raw(query, billingID).Scan(&items).Error
	return items, err
}
