package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type MilkDeliveryShiftService struct{}

func NewMilkDeliveryShiftService() *MilkDeliveryShiftService {
	return &MilkDeliveryShiftService{}
}

func (s *MilkDeliveryShiftService) CreateShift(req dtos.CreateMilkDeliveryShiftRequest) (*models.MilkDeliveryShift, error) {
	shift := &models.MilkDeliveryShift{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(shift).Error; err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *MilkDeliveryShiftService) GetShifts(page, limit int) ([]dtos.MilkDeliveryShiftResponse, int64, error) {
	var shifts []dtos.MilkDeliveryShiftResponse
	var total int64
	db.DB.Model(&models.MilkDeliveryShift{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.MilkDeliveryShift{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&shifts).Error
	return shifts, total, err
}

func (s *MilkDeliveryShiftService) GetShift(id string) (*dtos.MilkDeliveryShiftResponse, error) {
	var shift dtos.MilkDeliveryShiftResponse
	if err := db.DB.Model(&models.MilkDeliveryShift{}).First(&shift, id).Error; err != nil {
		return nil, err
	}
	return &shift, nil
}

func (s *MilkDeliveryShiftService) UpdateShift(id string, req dtos.UpdateMilkDeliveryShiftRequest) error {
	var shift models.MilkDeliveryShift
	if err := db.DB.First(&shift, id).Error; err != nil {
		return err
	}

	shift.Name = req.Name
	shift.Description = req.Description

	return db.DB.Save(&shift).Error
}

func (s *MilkDeliveryShiftService) DeleteShift(id string) error {
	return db.DB.Delete(&models.MilkDeliveryShift{}, id).Error
}
