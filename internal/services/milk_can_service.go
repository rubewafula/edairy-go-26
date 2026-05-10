package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
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
		Units:      req.Units,
		TareWeight: req.TareWeight,
		RouteID:    req.RouteID,
	}

	if err := db.DB.Create(milkCan).Error; err != nil {
		return nil, err
	}
	return milkCan, nil
}

func (s *MilkCanService) GetMilkCans(page, limit int) ([]dtos.MilkCanResponse, int64, error) {
	var milkCans []dtos.MilkCanResponse
	var total int64
	db.DB.Model(&models.MilkCan{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			mc.id, mc.can_id, mc.can_type, mc.can_size, mc.units, mc.tare_weight, 
			mc.route_id, r.route_name, mc.created_at, mc.updated_at
		FROM milk_cans mc
		LEFT JOIN routes r ON mc.route_id = r.id
		WHERE mc.deleted_at IS NULL
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&milkCans).Error
	return milkCans, total, err
}

func (s *MilkCanService) GetMilkCan(id string) (*dtos.MilkCanResponse, error) {
	var milkCan dtos.MilkCanResponse
	query := `
		SELECT 
			mc.id, mc.can_id, mc.can_type, mc.can_size, mc.units, mc.tare_weight, 
			mc.route_id, r.route_name, mc.created_at, mc.updated_at
		FROM milk_cans mc
		LEFT JOIN routes r ON mc.route_id = r.id
		WHERE mc.id = ? AND mc.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&milkCan).Error
	if err != nil {
		return nil, err
	}
	if milkCan.ID == 0 {
		return nil, gorm.ErrRecordNotFound
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
	milkCan.Units = req.Units
	milkCan.TareWeight = req.TareWeight
	milkCan.RouteID = req.RouteID

	return db.DB.Save(&milkCan).Error
}

func (s *MilkCanService) DeleteMilkCan(id string) error {
	return db.DB.Delete(&models.MilkCan{}, id).Error
}
