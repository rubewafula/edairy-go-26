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

type CustomerPayDateRangeService struct{}

func NewCustomerPayDateRangeService() *CustomerPayDateRangeService {
	return &CustomerPayDateRangeService{}
}

func (s *CustomerPayDateRangeService) CreateCustomerPayDateRange(req dtos.CreateCustomerPayDateRangeRequest, userID uint64) (*models.CustomerPayDateRange, error) {

	dateRange, err := s.GenerateNextPayDateRange(userID)
	if err != nil {
		return nil, err
	}
	return dateRange, nil
}

// CatchUpRanges ensures that billing cycles are generated up to the current date.
// It creates missing 7-day ranges sequentially until the latest range encompasses or nears today.
func (s *CustomerPayDateRangeService) CatchUpRanges(userID uint64) error {
	now := time.Now()
	for {
		var lastRange models.CustomerPayDateRange
		err := db.DB.Where("deleted_at IS NULL").Order("end_date DESC").First(&lastRange).Error

		var nextStart time.Time
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// No previous range, start from the 1st of the current month
				nextStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			} else {
				return err
			}
		} else {
			nextStart = lastRange.EndDate.AddDate(0, 0, 1)
		}

		if nextStart.After(now) {
			break
		}

		if _, err := s.GenerateNextPayDateRange(userID); err != nil {
			return err
		}
	}
	return nil
}

// GenerateNextPayDateRange creates the next customer pay date range based on the last one.
// It ensures cycles are 7 days long and do not overlap months.
func (s *CustomerPayDateRangeService) GenerateNextPayDateRange(userID uint64) (*models.CustomerPayDateRange, error) {
	var lastRange models.CustomerPayDateRange
	// Find the latest processed pay date range
	err := db.DB.Where("deleted_at IS NULL").Order("end_date DESC").First(&lastRange).Error

	var newStartDate time.Time
	if err != nil && err == gorm.ErrRecordNotFound {
		// No previous range, start from the 1st of the current month
		now := time.Now()
		newStartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve last pay date range: %w", err)
	} else {
		// Start date is the day after the last range's end date
		newStartDate = lastRange.EndDate.AddDate(0, 0, 1)
	}

	// Calculate potential end date (7 days from new start date)
	potentialEndDate := newStartDate.AddDate(0, 0, 6)

	// Ensure the cycle does not overlap months
	newEndDate := potentialEndDate
	if potentialEndDate.Month() != newStartDate.Month() {
		// If it overlaps, set end date to the last day of the current month
		lastDayOfMonth := time.Date(newStartDate.Year(), newStartDate.Month()+1, 0, 0, 0, 0, 0, newStartDate.Location())
		newEndDate = lastDayOfMonth
	}

	dateRange := &models.CustomerPayDateRange{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Name:      fmt.Sprintf("Billing Cycle %s - %s", newStartDate.Format("Jan 02"), newEndDate.Format("Jan 02, 2006")),
		StartDate: newStartDate,
		EndDate:   newEndDate,
		Status:    "pending",
	}

	if err := db.DB.Create(dateRange).Error; err != nil {
		return nil, err
	}
	return dateRange, nil
}

func (s *CustomerPayDateRangeService) GetCustomerPayDateRanges(page, limit int) ([]dtos.CustomerPayDateRangeResponse, int64, error) {
	var results []dtos.CustomerPayDateRangeResponse
	var total int64
	db.DB.Model(&models.CustomerPayDateRange{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.CustomerPayDateRange{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *CustomerPayDateRangeService) GetCustomerPayDateRange(id string) (*dtos.CustomerPayDateRangeResponse, error) {
	var result dtos.CustomerPayDateRangeResponse
	err := db.DB.Model(&models.CustomerPayDateRange{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *CustomerPayDateRangeService) UpdateCustomerPayDateRange(id string, req dtos.UpdateCustomerPayDateRangeRequest, userID uint64) error {
	var dateRange models.CustomerPayDateRange
	if err := db.DB.First(&dateRange, id).Error; err != nil {
		return err
	}

	dateRange.Name = req.Name
	dateRange.StartDate = utils.ParseDate(req.StartDate)
	dateRange.EndDate = utils.ParseDate(req.EndDate)
	dateRange.UpdatedBy = userID

	return db.DB.Save(&dateRange).Error
}

func (s *CustomerPayDateRangeService) DeleteCustomerPayDateRange(id string) error {
	return db.DB.Delete(&models.CustomerPayDateRange{}, id).Error
}
