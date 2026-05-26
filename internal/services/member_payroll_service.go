package services

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MemberPayrollService struct{}

func NewMemberPayrollService() *MemberPayrollService {
	return &MemberPayrollService{}
}

func (s *MemberPayrollService) Create(req dtos.CreateMemberPayrollRequest, userID uint64) (*models.MemberPayroll, error) {
	// 1. Verify Date Range
	var pdr models.MemberPayDateRange
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", req.PayDateRangeID).First(&pdr).Error; err != nil {
		return nil, fmt.Errorf("invalid pay date range: %w", err)
	}

	// Check if payroll already exists for this period
	var existing models.MemberPayroll
	err := db.DB.Where("pay_date_range_id = ? AND deleted_at IS NULL", req.PayDateRangeID).First(&existing).Error
	if err == nil && existing.Status != "draft" {
		return nil, fmt.Errorf("payroll for this period is already %s and cannot be regenerated", existing.Status)
	}

	dateOpened := utils.ParseDate(req.DateOpened)
	var payroll models.MemberPayroll

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Clean up existing draft if it exists
		if existing.ID != 0 {
			if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.MemberPayslip{}).Error; err != nil {
				return err
			}
			if err := tx.Where("payroll_id = ?", existing.ID).Delete(&models.MemberPayrollDeduction{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		}

		payroll = models.MemberPayroll{
			BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			DateOpened:     &dateOpened,
			Description:    req.Description,
			Status:         "processing",
			PayDateRangeID: &req.PayDateRangeID,
			PhysicalPeriod: req.PhysicalPeriod,
		}
		return tx.Create(&payroll).Error
	})

	if err != nil {
		return nil, err
	}

	// Defer heavy processing to background
	go s.generatePayrollInBackground(payroll.ID, pdr, req, userID)

	return &payroll, nil
}

func (s *MemberPayrollService) generatePayrollInBackground(payrollID uint64, pdr models.MemberPayDateRange, req dtos.CreateMemberPayrollRequest, userID uint64) {
	type milkSum struct {
		MemberID uint64
		Kilos    float64
	}

	// 1. Fetch Aggregated Milk Collections
	var milkCollections []milkSum
	db.DB.Table("milk_journal_entries").
		Select("member_id, SUM(quantity) as kilos").
		Joins("JOIN milk_journals ON milk_journal_entries.milk_journal_id = milk_journals.id").
		Where("milk_journals.status = ? AND milk_journals.date_captured BETWEEN ? AND ?", "approved", pdr.StartDate, pdr.EndDate).
		Group("member_id").Scan(&milkCollections)

	if len(milkCollections) == 0 {
		db.DB.Model(&models.MemberPayroll{}).Where("id = ?", payrollID).Update("status", "draft")
		return
	}

	// 2. Fetch Aggregated Rejects
	var milkRejects []milkSum
	db.DB.Table("milk_rejects").
		Select("member_id, SUM(quantity) as kilos").
		Where("reject_date BETWEEN ? AND ?", pdr.StartDate, pdr.EndDate).
		Group("member_id").Scan(&milkRejects)

	rejectMap := make(map[uint64]float64)
	for _, r := range milkRejects {
		rejectMap[r.MemberID] = r.Kilos
	}

	// 3. Pre-fetch Rate Resolution Maps
	var specialRates []models.MilkSpecialRate
	db.DB.Where("monthly_pay_date_range_id = ? AND deleted_at IS NULL", req.PayDateRangeID).Find(&specialRates)
	memberRateMap := make(map[uint64]float64)
	routePeriodRateMap := make(map[uint64]float64)
	for _, r := range specialRates {
		if r.MemberID != 0 {
			memberRateMap[r.MemberID] = r.Rate
		} else if r.RouteID != 0 {
			routePeriodRateMap[r.RouteID] = r.Rate
		}
	}

	var defaultRates []models.DefaultMilkRate
	db.DB.Where("deleted_at IS NULL").Find(&defaultRates)
	defaultRouteRateMap := make(map[uint64]float64)
	var globalDefault float64
	for _, r := range defaultRates {
		if r.RouteID != 0 {
			defaultRouteRateMap[r.RouteID] = r.Rate
		} else {
			globalDefault = r.Rate
		}
	}

	// 4. Pre-fetch Member metadata and Deductions in bulk
	memberIDs := make([]uint64, len(milkCollections))
	for i, c := range milkCollections {
		memberIDs[i] = c.MemberID
	}

	var members []models.Member
	db.DB.Select("id, route_id").Where("id IN ?", memberIDs).Find(&members)
	memberRouteMap := make(map[uint64]uint64)
	for _, m := range members {
		memberRouteMap[m.ID] = m.RouteID
	}

	var allDeductions []models.RecurrentDeduction
	db.DB.Where("customer_id IN ? AND settled = 0 AND customer_type = 'member'", memberIDs).
		Order("created_at ASC").Find(&allDeductions)
	deductionMap := make(map[uint64][]models.RecurrentDeduction)
	for _, d := range allDeductions {
		deductionMap[d.CustomerID] = append(deductionMap[d.CustomerID], d)
	}

	// 5. Worker Pool Setup
	dateOpened := utils.ParseDate(req.DateOpened)
	numWorkers := runtime.NumCPU() * 2
	if sqlDB, err := db.DB.DB(); err == nil {
		stats := sqlDB.Stats()
		if stats.MaxOpenConnections > 0 && numWorkers > stats.MaxOpenConnections/2 {
			numWorkers = stats.MaxOpenConnections / 2
		}
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	var wg sync.WaitGroup
	jobs := make(chan milkSum, len(milkCollections))
	var mu sync.Mutex
	var totalGross, totalNet, totalDeductions float64
	var failedCount int64

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for collection := range jobs {
				var mGross, mNet, mDeductions float64
				var err error

				// Retry logic: try up to 3 times per member
				for attempt := 1; attempt <= 3; attempt++ {
					err = db.DB.Transaction(func(tx *gorm.DB) error {
						mGross, mNet, mDeductions = 0, 0, 0
						rejectKilos := rejectMap[collection.MemberID]
						netKilos := collection.Kilos - rejectKilos
						if netKilos < 0 {
							netKilos = 0
						}

						routeID := memberRouteMap[collection.MemberID]
						var rate float64
						if r, ok := memberRateMap[collection.MemberID]; ok {
							rate = r
						} else if r, ok := routePeriodRateMap[routeID]; ok {
							rate = r
						} else if r, ok := defaultRouteRateMap[routeID]; ok {
							rate = r
						} else {
							rate = globalDefault
						}

						mGross = netKilos * rate
						var payslipDeductions float64
						if mDeds, ok := deductionMap[collection.MemberID]; ok {
							for _, rd := range mDeds {
								remaining := rd.TotalAmount - rd.PaidAmount
								deductAmount := rd.RecurrentAmount
								if deductAmount > remaining {
									deductAmount = remaining
								}
								if payslipDeductions+deductAmount > mGross {
									deductAmount = mGross - payslipDeductions
								}
								if deductAmount <= 0 {
									continue
								}

								mpd := models.MemberPayrollDeduction{
									MemberID:        int64(collection.MemberID),
									PayrollID:       payrollID,
									DeductionTypeID: rd.DeductionTypeID,
									Amount:          fmt.Sprintf("%.2f", deductAmount),
									TransactionDate: dateOpened,
									Reference:       rd.Reference,
								}
								if err := tx.Create(&mpd).Error; err != nil {
									return err
								}
								payslipDeductions += deductAmount
							}
						}

						mDeductions = payslipDeductions
						mNet = mGross - mDeductions
						payslip := models.MemberPayslip{
							MemberID:        collection.MemberID,
							PayrollID:       payrollID,
							DateOpened:      &dateOpened,
							Status:          "draft",
							GrossKilos:      collection.Kilos,
							RejectKilos:     rejectKilos,
							NetKilos:        netKilos,
							GrossPay:        mGross,
							TotalDeductions: mDeductions,
							NetPay:          mNet,
							PhysicalPeriod:  req.PhysicalPeriod,
							PayDateRangeID:  &req.PayDateRangeID,
						}
						return tx.Create(&payslip).Error
					})

					if err == nil {
						break
					}
					time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
				}

				if err != nil {
					mu.Lock()
					failedCount++
					mu.Unlock()
					log.Printf("[MemberPayrollService] Failed to generate payroll for member %d: %v", collection.MemberID, err)
					db.DB.Create(&models.MemberPayrollGenerationError{
						BaseModel: models.BaseModel{CreatedBy: userID},
						MemberID:  collection.MemberID,
						PayrollID: payrollID,
						Error:     err.Error(),
					})
				} else {
					mu.Lock()
					totalGross += mGross
					totalNet += mNet
					totalDeductions += mDeductions
					mu.Unlock()
				}
			}
		}()
	}

	for _, c := range milkCollections {
		jobs <- c
	}
	close(jobs)
	wg.Wait()

	// 6. Finalize Payroll Header
	status := "draft"
	if failedCount > 0 {
		status = "incomplete"
	}

	db.DB.Model(&models.MemberPayroll{}).Where("id = ?", payrollID).Updates(map[string]interface{}{
		"gross_pay":        totalGross,
		"net_pay":          totalNet,
		"total_deductions": totalDeductions,
		"status":           status,
		"updated_by":       userID,
		"updated_at":       time.Now(),
	})
}

func (s *MemberPayrollService) getRawPayroll(id string) (*models.MemberPayroll, error) {
	var payroll models.MemberPayroll
	if err := db.DB.First(&payroll, id).Error; err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (s *MemberPayrollService) List(page, limit int) ([]dtos.MemberPayrollResponse, int64, error) {
	var results []dtos.MemberPayrollResponse
	var total int64
	db.DB.Model(&models.MemberPayroll{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			mp.*, 
			pdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM member_payrolls mp
		LEFT JOIN customer_pay_date_ranges pdr ON mp.pay_date_range_id = pdr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.deleted_at IS NULL
		ORDER BY mp.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *MemberPayrollService) Get(id string) (*dtos.MemberPayrollResponse, error) {
	var result dtos.MemberPayrollResponse
	query := `
		SELECT 
			mp.*, 
			pdr.name as pay_date_range_name,
			u_post.name as posted_by_name,
			u_conf.name as confirmed_by_name,
			u_appr.name as approved_by_name
		FROM member_payrolls mp
		LEFT JOIN customer_pay_date_ranges pdr ON mp.pay_date_range_id = pdr.id
		LEFT JOIN users u_post ON mp.posted_by = u_post.id
		LEFT JOIN users u_conf ON mp.confirmed_by = u_conf.id
		LEFT JOIN users u_appr ON mp.approved_by = u_appr.id
		WHERE mp.id = ? AND mp.deleted_at IS NULL
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
