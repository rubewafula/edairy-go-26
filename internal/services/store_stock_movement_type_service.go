package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type StoreStockMovementTypeService struct{}

func NewStoreStockMovementTypeService() *StoreStockMovementTypeService {
	return &StoreStockMovementTypeService{}
}

func (s *StoreStockMovementTypeService) CreateMovementType(req dtos.CreateStoreStockMovementTypeRequest) (*models.StoreStockMovementType, error) {
	movementType := &models.StoreStockMovementType{
		MovementCode: req.MovementCode,
		MovementName: req.MovementName,
		Direction:    req.Direction,
		AffectsStock: req.AffectsStock,
		Description:  req.Description,
		IsSystem:     req.IsSystem,
	}

	if err := db.DB.Create(movementType).Error; err != nil {
		return nil, err
	}
	return movementType, nil
}

func (s *StoreStockMovementTypeService) GetMovementTypes(page, limit int) ([]dtos.StoreStockMovementTypeResponse, int64, error) {
	var results []dtos.StoreStockMovementTypeResponse
	var total int64
	db.DB.Model(&models.StoreStockMovementType{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.StoreStockMovementType{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *StoreStockMovementTypeService) GetMovementType(id string) (*models.StoreStockMovementType, error) {
	var movementType models.StoreStockMovementType
	if err := db.DB.First(&movementType, id).Error; err != nil {
		return nil, err
	}
	return &movementType, nil
}

func (s *StoreStockMovementTypeService) UpdateMovementType(id string, req dtos.UpdateStoreStockMovementTypeRequest) error {
	var movementType models.StoreStockMovementType
	if err := db.DB.First(&movementType, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&movementType).Updates(map[string]interface{}{
		"movement_code": req.MovementCode,
		"movement_name": req.MovementName,
		"direction":     req.Direction,
		"affects_stock": req.AffectsStock,
		"description":   req.Description,
		"is_system":     req.IsSystem,
	}).Error
}

func (s *StoreStockMovementTypeService) DeleteMovementType(id string) error {
	var movementType models.StoreStockMovementType
	if err := db.DB.First(&movementType, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&movementType).Error
}
