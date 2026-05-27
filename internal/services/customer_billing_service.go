package services

import (
	"fmt"
	"log"
	"math"

	"runtime"
	"sync"
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

type aggregatedDelivery struct {
	CustomerID     uint64
	ProductGradeID uint64
	TotalQty       float64
	DefaultRate    float64
}

// GetNextUnprocessedRange finds the oldest cycle that hasn't been billed yet.
func (s *CustomerBillingService) GetNextUnprocessedRange() (uint64, error) {
	var pdr models.CustomerPayDateRange
	err := db.DB.Where("status = ? AND deleted_at IS NULL", "pending").Order("end_date ASC").First(&pdr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no pending billing cycles found")
		}
		return 0, fmt.Errorf("failed to retrieve unprocessed pay date range: %w", err)
	}
	return pdr.ID, nil
}

func (s *CustomerBillingService) CreateBilling(userID uint64) error {
	log.Printf("[CustomerBillingService.CreateBilling] User %d initiated billing catch-up and generation", userID)
	// 1. Generate all missing ranges up to today
	pdrService := NewCustomerPayDateRangeService()
	if err := pdrService.CatchUpRanges(userID); err != nil {
		return fmt.Errorf("failed to catch up billing ranges: %w", err)
	}

	// 2. Trigger background processing for all ranges that are not yet 'processed'
	go s.processAllPendingCycles(userID)

	return nil
}

func (s *CustomerBillingService) processAllPendingCycles(userID uint64) {
	log.Printf("[CustomerBillingService.processAllPendingCycles] Starting background processing for pending cycles...")
	var ranges []models.CustomerPayDateRange
	// Process pending or incomplete ranges in chronological order
	db.DB.Where("status IN ? AND deleted_at IS NULL", []string{"pending", "incomplete"}).
		Order("end_date ASC").Find(&ranges)

	for _, pdr := range ranges {
		// IDEMPOTENCY CHECK: If any billing in this range is confirmed or invoiced, skip the range
		var nonPendingCount int64
		db.DB.Model(&models.CustomerBilling{}).
			Where("pay_date_range_id = ? AND status != ?", pdr.ID, "pending").
			Count(&nonPendingCount)

		if nonPendingCount > 0 {
			log.Printf("[CustomerBillingService] Range %d has confirmed/invoiced billings. Skipping regeneration.", pdr.ID)
			continue
		}

		// IDEMPOTENCY PREPARATION: Wipe existing pending data to start fresh
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			var billingIDs []uint64
			tx.Model(&models.CustomerBilling{}).Where("pay_date_range_id = ? AND status = ?", pdr.ID, "pending").Pluck("id", &billingIDs)

			if len(billingIDs) > 0 {
				tx.Where("customer_billing_id IN ?", billingIDs).Delete(&models.CustomerBillingItem{})
				tx.Where("id IN ?", billingIDs).Delete(&models.CustomerBilling{})
			}

			// Reset 'processed' flag for all deliveries in this window to allow the SQL aggregation to see them again
			tx.Model(&models.MilkDelivery{}).
				Where("transaction_date BETWEEN ? AND ? AND customer_id != 0", pdr.StartDate, pdr.EndDate).
				Update("processed", false)

			return tx.Model(&pdr).Update("status", "processing").Error
		})

		if err != nil {
			log.Printf("[CustomerBillingService] Failed to initialize processing for range %d: %v", pdr.ID, err)
			continue
		}

		log.Printf("[CustomerBillingService.processAllPendingCycles] Dispatching generation for range %d (%s to %s)", pdr.ID, pdr.StartDate.Format("2006-01-02"), pdr.EndDate.Format("2006-01-02"))
		// Call the heavy lifting generation logic synchronously within this background goroutine
		s.generateBillingInBackground(pdr, userID)
	}
}

func (s *CustomerBillingService) generateBillingInBackground(pdr models.CustomerPayDateRange, userID uint64) {
	payDateRangeID := pdr.ID
	log.Printf("[CustomerBillingService.CreateBilling] Initiating billing generation for PayDateRangeID: %d", payDateRangeID)

	// 1. Fetch aggregated delivery data joined with customer default rates
	var summaries []aggregatedDelivery
	query := `
		SELECT 
			md.customer_id, 
			md.product_grade_id, 
			SUM(md.quantity_accepted) AS total_qty, 
			c.rate AS default_rate
		FROM milk_deliveries md
		INNER JOIN customers c ON md.customer_id = c.id
		WHERE md.transaction_date BETWEEN ? AND ?
		  AND md.customer_id != 0
		  AND md.processed = ?
		  AND md.deleted_at IS NULL
		  AND c.deleted_at IS NULL
		GROUP BY md.customer_id, md.product_grade_id, c.rate
	`

	if err := db.DB.Raw(query, pdr.StartDate, pdr.EndDate, false).Scan(&summaries).Error; err != nil {
		log.Printf("[CustomerBillingService.generateBillingInBackground] Error fetching summaries for range %d: %v", payDateRangeID, err)
		db.DB.Model(&models.CustomerPayDateRange{}).Where("id = ?", payDateRangeID).Update("status", "pending")
		return
	}

	if len(summaries) == 0 {
		db.DB.Model(&models.CustomerPayDateRange{}).Where("id = ?", payDateRangeID).Update("status", "processed")
		return
	}

	// 2. Group aggregated rows by CustomerID
	customerGroups := make(map[uint64][]aggregatedDelivery)
	var customerIDs []uint64
	for _, s := range summaries {
		if _, exists := customerGroups[s.CustomerID]; !exists {
			customerIDs = append(customerIDs, s.CustomerID)
		}
		customerGroups[s.CustomerID] = append(customerGroups[s.CustomerID], s)
	}

	var milkRates []models.CustomerMilkRate
	db.DB.Where("customer_id IN ? AND customer_pay_date_range_id = ?", customerIDs, payDateRangeID).Find(&milkRates)
	rateMap := make(map[uint64]map[uint64]float64) // CustomerID -> GradeID -> Rate
	for _, r := range milkRates {
		if rateMap[r.CustomerID] == nil {
			rateMap[r.CustomerID] = make(map[uint64]float64)
		}
		rateMap[r.CustomerID][r.GradeID] = r.Rate
	}

	// Concurrently process billing for each customer
	var wg sync.WaitGroup
	errorChan := make(chan error, len(customerGroups))
	failedCustomerCount := 0

	numWorkers := runtime.NumCPU() * 2 // Use a reasonable number of workers
	if numWorkers < 1 {
		numWorkers = 1
	}
	jobs := make(chan uint64, len(customerGroups))

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for custID := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[CustomerBillingService] Worker panicked during billing generation for customer %d: %v", custID, r)
							db.DB.Create(&models.CustomerBillingGenerationError{
								BaseModel:      models.BaseModel{CreatedBy: userID},
								CustomerID:     custID,
								PayDateRangeID: payDateRangeID,
								Error:          fmt.Sprintf("Panic during generation: %v", r),
							})
							errorChan <- fmt.Errorf("panic during billing generation for customer %d: %v", custID, r)
						}
					}()

					err := s.processSingleCustomerBilling(custID, payDateRangeID, userID, customerGroups[custID], rateMap, pdr.StartDate, pdr.EndDate)
					if err != nil {
						log.Printf("[CustomerBillingService.CreateBilling] Error processing billing for customer %d: %v", custID, err)
						// Log the error to a dedicated error table
						db.DB.Create(&models.CustomerBillingGenerationError{
							BaseModel:      models.BaseModel{CreatedBy: userID},
							CustomerID:     custID,
							PayDateRangeID: payDateRangeID,
							Error:          err.Error(),
						})
						errorChan <- err
					}
				}()
			}
		}()
	}

	for custID := range customerGroups {
		jobs <- custID
	}
	close(jobs)
	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			failedCustomerCount++
		}
	}

	// Final update of Pay Date Range status
	finalStatus := "processed"
	if failedCustomerCount > 0 {
		finalStatus = "incomplete"
		log.Printf("[CustomerBillingService.CreateBilling] Completed billing generation for PayDateRangeID %d with %d failures.", payDateRangeID, failedCustomerCount)
	} else {
		log.Printf("[CustomerBillingService.CreateBilling] Completed billing generation for PayDateRangeID %d successfully.", payDateRangeID)
	}

	db.DB.Model(&models.CustomerPayDateRange{}).Where("id = ?", payDateRangeID).Update("status", finalStatus)
}

// processSingleCustomerBilling handles the billing generation for a single customer.
func (s *CustomerBillingService) processSingleCustomerBilling(
	customerID, payDateRangeID, userID uint64,
	items []aggregatedDelivery,
	rateMap map[uint64]map[uint64]float64,
	startDate, endDate time.Time,
) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Check if billing already exists for this customer and pay date range
		var existingBilling models.CustomerBilling
		err := tx.Where("customer_id = ? AND pay_date_range_id = ? AND deleted_at IS NULL", customerID, payDateRangeID).First(&existingBilling).Error
		if err == nil {
			return fmt.Errorf("billing already exists for customer %d and pay date range ID %d", customerID, payDateRangeID)
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check for existing billing for customer %d: %w", customerID, err)
		}

		// Create CustomerBilling header
		billing := models.CustomerBilling{
			BaseModel:       models.BaseModel{CreatedBy: userID},
			PayDateRangeID:  payDateRangeID,
			CustomerID:      customerID,
			Status:          "pending",
			TotalDeliveries: 0, // Will be summed from items
			TotalAmount:     0, // Will be summed from items
		}
		if err := tx.Create(&billing).Error; err != nil {
			return fmt.Errorf("error creating billing header for customer %d: %w", customerID, err)
		}

		// Create CustomerBillingItems and update billing totals
		var billingItems []models.CustomerBillingItem
		var totalDeliveries float64
		var totalAmount float64

		for _, itm := range items {
			// Resolve Rate: CustomerMilkRate (specific) -> Customer.Rate (fallback)
			rate := itm.DefaultRate
			if r, exists := rateMap[customerID][itm.ProductGradeID]; exists {
				rate = r
			}

			unitPrice := math.Round(rate*100) / 100
			itemTotalAmount := math.Round(itm.TotalQty*rate*100) / 100

			item := models.CustomerBillingItem{
				BaseModel:         models.BaseModel{CreatedBy: userID},
				CustomerBillingID: billing.ID,
				ProductGradeID:    itm.ProductGradeID,
				TotalQuantity:     itm.TotalQty,
				UnitPrice:         unitPrice,
				TotalAmount:       itemTotalAmount,
			}
			billingItems = append(billingItems, item)
			totalDeliveries += itm.TotalQty
			totalAmount += itemTotalAmount
		}

		if len(billingItems) > 0 {
			if err := tx.CreateInBatches(billingItems, 100).Error; err != nil { // Bulk insert items
				return fmt.Errorf("error creating billing items for customer %d: %w", customerID, err)
			}
		}

		// Update CustomerBilling header with final totals
		if err := tx.Model(&billing).Updates(map[string]interface{}{
			"total_deliveries": totalDeliveries,
			"total_amount":     totalAmount,
		}).Error; err != nil {
			return fmt.Errorf("error updating billing totals for billing %d: %w", billing.ID, err)
		}

		// 3. Mark deliveries as processed in a single batch update for this customer and period
		if err := tx.Model(&models.MilkDelivery{}).
			Where("customer_id = ? AND transaction_date BETWEEN ? AND ? AND processed = ?",
				customerID, startDate, endDate, false).
			Update("processed", true).Error; err != nil {

			return fmt.Errorf("error marking deliveries as processed for customer %d: %w", customerID, err)
		}

		log.Printf("[CustomerBillingService.processSingleCustomerBilling] Successfully processed billing for customer %d in range %d", customerID, payDateRangeID)
		return nil
	})
}

func (s *CustomerBillingService) ConfirmBilling(billingID uint64, userID uint64) error {
	var billing models.CustomerBilling
	if err := db.DB.First(&billing, billingID).Error; err != nil {
		return fmt.Errorf("customer billing not found: %w", err)
	}

	if billing.Status != "pending" {
		return fmt.Errorf("billing must be in 'pending' status to be confirmed, current status: %s", billing.Status)
	}

	log.Printf("[CustomerBillingService.ConfirmBilling] Billing %d confirmed by user %d", billingID, userID)
	return db.DB.Model(&billing).Updates(map[string]interface{}{
		"status":     "confirmed",
		"updated_by": userID,
	}).Error
}

func (s *CustomerBillingService) ApproveBilling(billingID uint64, userID uint64) error {
	var billing models.CustomerBilling
	if err := db.DB.First(&billing, billingID).Error; err != nil {
		return fmt.Errorf("customer billing not found: %w", err)
	}

	if billing.Status != "confirmed" {
		return fmt.Errorf("billing must be in 'confirmed' status to be approved, current status: %s", billing.Status)
	}

	log.Printf("[CustomerBillingService.ApproveBilling] Billing %d approved by user %d", billingID, userID)
	return db.DB.Model(&billing).Updates(map[string]interface{}{
		"status":     "approved",
		"updated_by": userID,
	}).Error
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
	log.Printf("[CustomerBillingService.DeleteBilling] Deleting billing record ID: %s", id)
	return db.DB.Delete(&models.CustomerBilling{}, id).Error
}
