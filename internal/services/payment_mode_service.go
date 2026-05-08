package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type PaymentModeService struct{}

func NewPaymentModeService() *PaymentModeService {
	return &PaymentModeService{}
}

func (s *PaymentModeService) CreatePaymentMode(req dtos.CreatePaymentModeRequest) (*models.PaymentMode, error) {
	paymentMode := &models.PaymentMode{
		Code: req.Code,
		Name: req.Name,
	}

	if err := db.DB.Create(paymentMode).Error; err != nil {
		return nil, err
	}
	return paymentMode, nil
}

func (s *PaymentModeService) GetPaymentModes() ([]dtos.PaymentModeResponse, int64, error) {
	var results []dtos.PaymentModeResponse
	var total int64
	db.DB.Model(&models.PaymentMode{}).Count(&total)

	query := `
		SELECT 
			id, code, name, created_at, updated_at
		FROM payment_modes
		WHERE deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *PaymentModeService) GetPaymentMode(id string) (*dtos.PaymentModeResponse, error) {
	var result dtos.PaymentModeResponse
	query := `
		SELECT 
			id, code, name, created_at, updated_at
		FROM payment_modes
		WHERE id = ? AND deleted_at IS NULL
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

func (s *PaymentModeService) UpdatePaymentMode(id string, req dtos.UpdatePaymentModeRequest) error {
	var paymentMode models.PaymentMode
	if err := db.DB.First(&paymentMode, id).Error; err != nil {
		return err
	}

	paymentMode.Code = req.Code
	paymentMode.Name = req.Name

	return db.DB.Save(&paymentMode).Error
}

func (s *PaymentModeService) DeletePaymentMode(id string) error {
	return db.DB.Delete(&models.PaymentMode{}, id).Error
}
