package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type StoreItemService struct{}

func NewStoreItemService() *StoreItemService {
	return &StoreItemService{}
}

func (s *StoreItemService) CreateStoreItem(req dtos.CreateStoreItemRequest) (*models.StoreItem, error) {
	item := &models.StoreItem{
		ItemName:                  req.ItemName,
		Description:               req.Description,
		ReorderPoint:              req.ReorderPoint,
		DefaultBuyingPrice:        req.DefaultBuyingPrice,
		DefaultSellingPrice:       req.DefaultSellingPrice,
		DefaultSellingPriceCredit: req.DefaultSellingPriceCredit,
		Status:                    req.Status,
		Thumbnail:                 req.Thumbnail,
		SKU:                       req.SKU,
		Barcode:                   req.Barcode,
		UnitID:                    req.UnitID,
		StoreInventoryID:          req.StoreInventoryID,
	}

	if err := db.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (s *StoreItemService) GetStoreItems(page, limit int) ([]dtos.StoreItemResponse, int64, error) {
	var results []dtos.StoreItemResponse
	var total int64
	db.DB.Model(&models.StoreItem{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			si.*, inv.inventory_name
		FROM store_items si
		LEFT JOIN store_inventories inv ON si.store_inventory_id = inv.id
		WHERE si.deleted_at IS NULL
		ORDER BY si.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreItemService) GetStoreItem(id string) (*dtos.StoreItemResponse, error) {
	var result dtos.StoreItemResponse
	query := `
		SELECT 
			si.*, inv.inventory_name
		FROM store_items si
		LEFT JOIN store_inventories inv ON si.store_inventory_id = inv.id
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

func (s *StoreItemService) UpdateStoreItem(id string, req dtos.UpdateStoreItemRequest) error {
	var item models.StoreItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	item.ItemName = req.ItemName
	item.Description = req.Description
	item.ReorderPoint = req.ReorderPoint
	item.DefaultBuyingPrice = req.DefaultBuyingPrice
	item.DefaultSellingPrice = req.DefaultSellingPrice
	item.DefaultSellingPriceCredit = req.DefaultSellingPriceCredit
	item.Status = req.Status
	item.Thumbnail = req.Thumbnail
	item.SKU = req.SKU
	item.Barcode = req.Barcode
	item.UnitID = req.UnitID
	item.StoreInventoryID = req.StoreInventoryID

	return db.DB.Save(&item).Error
}

func (s *StoreItemService) DeleteStoreItem(id string) error {
	return db.DB.Delete(&models.StoreItem{}, id).Error
}
