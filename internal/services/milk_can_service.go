package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type MilkCanService struct{}

func NewMilkCanService() *MilkCanService {
	return &MilkCanService{}
}

func (s *MilkCanService) CreateMilkCan(req dtos.CreateMilkCanRequest) (*models.MilkCan, error) {
	milkCan := &models.MilkCan{
		CanID:      req.CanID,
		CanType:    req.CanType,
		CanSize:    req.CanSize,
		TareWeight: req.TareWeight,
		RouteID:    req.RouteID,
	}

	if err := db.DB.Create(milkCan).Error; err != nil {
		return nil, err
	}
	return milkCan, nil
}

func (s *MilkCanService) GetMilkCans() ([]models.MilkCan, int64, error) {
	var milkCans []models.MilkCan
	var total int64
	db.DB.Model(&models.MilkCan{}).Count(&total)
	err := db.DB.Find(&milkCans).Error
	return milkCans, total, err
}

func (s *MilkCanService) GetMilkCan(id string) (*models.MilkCan, error) {
	var milkCan models.MilkCan
	if err := db.DB.First(&milkCan, id).Error; err != nil {
		return nil, err
	}
	return &milkCan, nil
}

func (s *MilkCanService) UpdateMilkCan(id string, req dtos.UpdateMilkCanRequest) error {
	var milkCan models.MilkCan
	if err := db.DB.First(&milkCan, id).Error; err != nil {
		return err
	}

	milkCan.CanID = req.CanID
	milkCan.CanType = req.CanType
	milkCan.CanSize = req.CanSize
	milkCan.TareWeight = req.TareWeight
	milkCan.RouteID = req.RouteID

	return db.DB.Save(&milkCan).Error
}

func (s *MilkCanService) DeleteMilkCan(id string) error {
	return db.DB.Delete(&models.MilkCan{}, id).Error
}
