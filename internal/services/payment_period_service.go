package services

import (
	"strconv"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type PaymentPeriodService struct{}

func NewPaymentPeriodService() *PaymentPeriodService {
	return &PaymentPeriodService{}
}

func (s *PaymentPeriodService) Create(req dtos.CreatePaymentPeriodRequest, userID uint64) (*dtos.PaymentPeriodResponse, error) {
	period := &models.PaymentPeriod{
		BaseModel:     models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		Name:          req.Name,
		Description:   req.Description,
		DefaultPeriod: req.DefaultPeriod,
	}
	if err := db.DB.Create(period).Error; err != nil {
		return nil, err
	}
	return s.Get(strconv.FormatUint(period.ID, 10))
}

func (s *PaymentPeriodService) List() ([]dtos.PaymentPeriodResponse, error) {
	var results []dtos.PaymentPeriodResponse
	err := db.DB.Model(&models.PaymentPeriod{}).Where("deleted_at IS NULL").Scan(&results).Error
	return results, err
}

func (s *PaymentPeriodService) Get(id string) (*dtos.PaymentPeriodResponse, error) {
	var result dtos.PaymentPeriodResponse
	if err := db.DB.Model(&models.PaymentPeriod{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *PaymentPeriodService) Update(id string, req dtos.UpdatePaymentPeriodRequest, userID uint64) error {
	var period models.PaymentPeriod
	if err := db.DB.First(&period, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"updated_by": userID,
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if req.DefaultPeriod == 1 {
			// Drop current default status from all other records
			if err := tx.Model(&models.PaymentPeriod{}).Where("deleted_at IS NULL").Update("default_period", 0).Error; err != nil {
				return err
			}
			updates["default_period"] = 1
		}

		return tx.Model(&period).Updates(updates).Error
	})
}

func (s *PaymentPeriodService) Delete(id string, userID uint64) error {
	// Audit the update before soft delete
	return db.DB.Model(&models.PaymentPeriod{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.PaymentPeriod{}).Error
}
