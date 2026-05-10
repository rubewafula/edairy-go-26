package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type StoreInventoryService struct{}

func NewStoreInventoryService() *StoreInventoryService {
	return &StoreInventoryService{}
}

func (s *StoreInventoryService) CreateInventory(req dtos.CreateStoreInventoryRequest) (*models.StoreInventory, error) {
	inventory := &models.StoreInventory{
		InventoryName: req.InventoryName,
		CategoryID:    req.CategoryID,
		IsActive:      req.IsActive,
		Description:   req.Description,
	}

	if err := db.DB.Create(inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *StoreInventoryService) GetInventories(page, limit int) ([]dtos.StoreInventoryResponse, int64, error) {
	var results []dtos.StoreInventoryResponse
	var total int64
	db.DB.Model(&models.StoreInventory{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			si.id, si.inventory_name, si.category_id, ic.name AS category_name,
			si.is_active, si.description, si.created_at, si.updated_at
		FROM store_inventories si
		LEFT JOIN item_categories ic ON si.category_id = ic.id
		WHERE si.deleted_at IS NULL
		ORDER BY si.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreInventoryService) GetInventory(id string) (*dtos.StoreInventoryResponse, error) {
	var result dtos.StoreInventoryResponse
	query := `
		SELECT 
			si.id, si.inventory_name, si.category_id, ic.name AS category_name,
			si.is_active, si.description, si.created_at, si.updated_at
		FROM store_inventories si
		LEFT JOIN item_categories ic ON si.category_id = ic.id
		WHERE si.id = ? AND si.deleted_at IS NULL
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

func (s *StoreInventoryService) UpdateInventory(id string, req dtos.UpdateStoreInventoryRequest) error {
	var inventory models.StoreInventory
	if err := db.DB.First(&inventory, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&inventory).Updates(map[string]interface{}{
		"inventory_name": req.InventoryName,
		"category_id":    req.CategoryID,
		"is_active":      req.IsActive,
		"description":    req.Description,
	}).Error
}

func (s *StoreInventoryService) DeleteInventory(id string) error {
	var inventory models.StoreInventory
	if err := db.DB.First(&inventory, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&inventory).Error
}
