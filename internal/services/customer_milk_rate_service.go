package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type CustomerMilkRateService struct{}

func NewCustomerMilkRateService() *CustomerMilkRateService {
	return &CustomerMilkRateService{}
}

func (s *CustomerMilkRateService) CreateCustomerMilkRate(req dtos.CreateCustomerMilkRateRequest, userID uint64) (*models.CustomerMilkRate, error) {
	rate := &models.CustomerMilkRate{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		CustomerID:   req.CustomerID,
		Rate:         req.Rate,
		GradeID:      req.GradeID,
		PayDateRange: req.PayDateRange,
	}

	if err := db.DB.Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *CustomerMilkRateService) GetCustomerMilkRates(page, limit int) ([]dtos.CustomerMilkRateResponse, int64, error) {
	var results []dtos.CustomerMilkRateResponse
	var total int64
	db.DB.Model(&models.CustomerMilkRate{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			cmr.*, c.full_names AS customer_name, pg.name AS grade_name, cpdr.name AS pay_date_range_name
		FROM customer_milk_rates cmr
		LEFT JOIN customers c ON cmr.customer_id = c.id
		LEFT JOIN product_grades pg ON cmr.grade_id = pg.id
		LEFT JOIN customer_pay_date_ranges cpdr ON cmr.customer_pay_date_range_id = cpdr.id
		WHERE cmr.deleted_at IS NULL
		ORDER BY cmr.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerMilkRateService) GetCustomerMilkRate(id string) (*dtos.CustomerMilkRateResponse, error) {
	var result dtos.CustomerMilkRateResponse
	query := `
		SELECT 
			cmr.*, c.full_names AS customer_name, pg.name AS grade_name, cpdr.name AS pay_date_range_name
		FROM customer_milk_rates cmr
		LEFT JOIN customers c ON cmr.customer_id = c.id
		LEFT JOIN product_grades pg ON cmr.grade_id = pg.id
		LEFT JOIN customer_pay_date_ranges cpdr ON cmr.customer_pay_date_range_id = cpdr.id
		WHERE cmr.id = ? AND cmr.deleted_at IS NULL
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

func (s *CustomerMilkRateService) UpdateCustomerMilkRate(id string, req dtos.UpdateCustomerMilkRateRequest, userID uint64) error {
	var rate models.CustomerMilkRate
	if err := db.DB.First(&rate, id).Error; err != nil {
		return err
	}

	rate.CustomerID = req.CustomerID
	rate.Rate = req.Rate
	rate.GradeID = req.GradeID
	rate.PayDateRange = req.PayDateRange
	rate.UpdatedBy = userID

	return db.DB.Save(&rate).Error
}

func (s *CustomerMilkRateService) DeleteCustomerMilkRate(id string) error {
	return db.DB.Delete(&models.CustomerMilkRate{}, id).Error
}
