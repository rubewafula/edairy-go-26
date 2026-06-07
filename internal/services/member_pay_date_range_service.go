package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MemberPayDateRangeService struct{}

func NewMemberPayDateRangeService() *MemberPayDateRangeService {
	return &MemberPayDateRangeService{}
}

func (s *MemberPayDateRangeService) Create(req dtos.CreateMemberPayDateRangeRequest, userID uint64) (*models.MemberPayDateRange, error) {
	// Validate that the provided start_date matches the next expected sequence
	expectedStart, _, err := s.GetNextPayDateRange()
	if err != nil {
		return nil, err
	}

	providedStart := utils.ParseDate(req.StartDate)

	// Compare using formatted strings to avoid timezone/location instant mismatch issues
	if providedStart.Format("2006-01-02") != expectedStart.Format("2006-01-02") {
		log.Printf("Start date mismatch start: %s expected: %s", providedStart.Format("2006-01-02"), expectedStart.Format("2006-01-02"))
		return nil, fmt.Errorf("invalid start date: the next pay cycle must start on %s", expectedStart.Format("2006-01-02"))
	}

	dateRange := &models.MemberPayDateRange{
		BaseModel: models.BaseModel{
			CreatedBy: userID, UpdatedBy: userID,
		},
		Name:      req.Name,
		StartDate: providedStart,
		EndDate:   utils.ParseDate(req.EndDate),
		Confirmed: 0,
	}

	if err := db.DB.Create(dateRange).Error; err != nil {
		return nil, err
	}
	return dateRange, nil
}

func (s *MemberPayDateRangeService) List(page, limit int) ([]dtos.MemberPayDateRangeResponse, int64, error) {
	var results []dtos.MemberPayDateRangeResponse
	var total int64
	db.DB.Model(&models.MemberPayDateRange{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.MemberPayDateRange{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *MemberPayDateRangeService) Get(id string) (*dtos.MemberPayDateRangeResponse, error) {
	var result dtos.MemberPayDateRangeResponse
	err := db.DB.Model(&models.MemberPayDateRange{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *MemberPayDateRangeService) Update(id string, req dtos.UpdateMemberPayDateRangeRequest, userID uint64) error {
	var dateRange models.MemberPayDateRange
	if err := db.DB.First(&dateRange, id).Error; err != nil {
		return err
	}

	if dateRange.Confirmed == 1 {
		return fmt.Errorf("Connot edit confirmed date range")
	}
	updates := map[string]interface{}{
		"name":       req.Name,
		"updated_by": userID,
	}

	if req.EndDate != "" {
		updates["end_date"] = utils.ParseDate(req.EndDate)
	}
	if req.Confirmed == 1 {
		updates["confirmed"] = req.Confirmed
	}

	return db.DB.Model(&dateRange).Updates(updates).Error
}

func (s *MemberPayDateRangeService) Delete(id string) error {
	return db.DB.Delete(&models.MemberPayDateRange{}, id).Error
}

// GetNextPayDateRange calculates the next start and end date for a member pay date range based on
// latest range or earliest journal, and the default payment period.
func (s *MemberPayDateRangeService) GetNextPayDateRange() (time.Time, time.Time, error) {
	var startDate, endDate time.Time

	// Get latest processed range
	var lastRange models.MemberPayDateRange
	err := db.DB.
		Where("deleted_at IS NULL").
		Order("end_date DESC").
		First(&lastRange).Error

	switch {
	case err == nil:
		startDate = lastRange.EndDate.AddDate(0, 0, 1)

	case errors.Is(err, gorm.ErrRecordNotFound):
		var earliestDate time.Time

		err = db.DB.
			Table("milk_journals").
			Where("deleted_at IS NULL").
			Select("MIN(DATE(journal_date))").
			Row().
			Scan(&earliestDate)

		if err != nil {
			return startDate, endDate, err
		}

		if earliestDate.IsZero() {
			return startDate, endDate, fmt.Errorf("no milk journals found")
		}

		startDate = earliestDate

	default:
		return startDate, endDate, err
	}

	var period models.PaymentPeriod
	if err := db.DB.
		Where("default_period = ?", true).
		First(&period).Error; err != nil {
		return startDate, endDate,
			fmt.Errorf("default payment period not configured")
	}

	switch strings.ToUpper(strings.TrimSpace(period.Name)) {

	case "WEEKLY":
		endDate = startDate.AddDate(0, 0, 6)

	case "BI-WEEKLY":
		endDate = startDate.AddDate(0, 0, 13)

	case "MONTHLY":
		// Rolling monthly period
		endDate = startDate.AddDate(0, 1, -1)

	default:
		return startDate, endDate,
			fmt.Errorf("unsupported payment period: %s", period.Name)
	}

	return startDate, endDate, nil
}
