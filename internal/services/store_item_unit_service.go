package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type StoreItemUnitService struct{}

func NewStoreItemUnitService() *StoreItemUnitService {
	return &StoreItemUnitService{}
}

func (s *StoreItemUnitService) CreateUnit(req dtos.CreateStoreItemUnitRequest, userID uint64) (*models.StoreItemUnit, error) {
	unit := &models.StoreItemUnit{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		Name:        req.Name,
		Symbol:      req.Symbol,
		Description: req.Description,
	}

	if err := db.DB.Create(unit).Error; err != nil {
		return nil, err
	}
	return unit, nil
}

func (s *StoreItemUnitService) GetUnits(page, limit int) ([]dtos.StoreItemUnitResponse, int64, error) {
	var results []dtos.StoreItemUnitResponse
	var total int64
	db.DB.Model(&models.StoreItemUnit{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.StoreItemUnit{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *StoreItemUnitService) GetUnit(id string) (*models.StoreItemUnit, error) {
	var unit models.StoreItemUnit
	if err := db.DB.First(&unit, id).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (s *StoreItemUnitService) UpdateUnit(id string, req dtos.UpdateStoreItemUnitRequest, userID uint64) error {
	var unit models.StoreItemUnit
	if err := db.DB.First(&unit, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&unit).Updates(map[string]interface{}{
		"name":        req.Name,
		"symbol":      req.Symbol,
		"description": req.Description,
		"updated_by":  userID,
	}).Error
}

func (s *StoreItemUnitService) DeleteUnit(id string) error {
	var unit models.StoreItemUnit
	if err := db.DB.First(&unit, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&unit).Error
}
