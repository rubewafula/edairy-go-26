package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type CustomerBillingService struct{}

func NewCustomerBillingService() *CustomerBillingService {
	return &CustomerBillingService{}
}

func (s *CustomerBillingService) CreateBilling(payDateRangeID uint64, userID uint64) error {
	var pdr models.CustomerPayDateRange
	if err := db.DB.First(&pdr, payDateRangeID).Error; err != nil {
		return fmt.Errorf("invalid pay date range: %w", err)
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch all confirmed, unprocessed deliveries in range
		var deliveries []models.MilkDelivery
		if err := tx.Where("transaction_date BETWEEN ? AND ? AND customer_id != 0 AND processed = ?",
			pdr.StartDate, pdr.EndDate, false).Find(&deliveries).Error; err != nil {
			return err
		}

		if len(deliveries) == 0 {
			return fmt.Errorf("no pending deliveries found for this period")
		}

		// 2. Group by Customer
		customerGroups := make(map[uint64][]models.MilkDelivery)
		for _, d := range deliveries {
			customerGroups[d.CustomerID] = append(customerGroups[d.CustomerID], d)
		}

		for customerID, delys := range customerGroups {
			var totalQty, totalAmt float64

			// Create Billing Header
			billing := models.CustomerBilling{
				BaseModel:      models.BaseModel{CreatedBy: userID},
				PayDateRangeID: payDateRangeID,
				CustomerID:     customerID,
				Status:         "pending",
			}
			if err := tx.Create(&billing).Error; err != nil {
				return err
			}

			// Group by Product Grade for Items (Simplified: assuming Grade information exists or is default)
			// In a real scenario, you'd join with product grades.
			// Here we aggregate items based on the delivery records.
			for _, d := range delys {
				item := models.CustomerBillingItem{
					BaseModel:         models.BaseModel{CreatedBy: userID},
					CustomerBillingID: billing.ID,
					TotalQuantity:     d.QuantityAccepted,
					UnitPrice:         d.Amount / d.QuantityAccepted,
					TotalAmount:       d.Amount,
				}
				if err := tx.Create(&item).Error; err != nil {
					return err
				}
				totalQty += d.QuantityAccepted
				totalAmt += d.Amount

				// Mark delivery as processed
				if err := tx.Model(&d).Update("processed", true).Error; err != nil {
					return err
				}
			}

			// Update Header totals
			if err := tx.Model(&billing).Updates(map[string]interface{}{
				"total_deliveries": totalQty,
				"total_amount":     totalAmt,
			}).Error; err != nil {
				return err
			}
		}

		// Update Pay Date Range status
		return tx.Model(&pdr).Update("updated_at", time.Now()).Error
	})
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

func (s *CustomerBillingService) DeleteBilling(id string) error {
	return db.DB.Delete(&models.CustomerBilling{}, id).Error
}
