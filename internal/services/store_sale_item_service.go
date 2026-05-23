package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type StoreSaleItemService struct{}

func NewStoreSaleItemService() *StoreSaleItemService {
	return &StoreSaleItemService{}
}

func (s *StoreSaleItemService) CreateSaleItem(req dtos.CreateStoreSaleItemRequest, userID uint64) (*models.StoreSaleItem, error) {
	item := &models.StoreSaleItem{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		ItemID:      req.ItemID,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		Total:       req.Total,
		StoreSaleID: req.StoreSaleID,
	}

	if err := db.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (s *StoreSaleItemService) GetSaleItems(page, limit int) ([]dtos.StoreSaleItemResponse, int64, error) {
	var results []dtos.StoreSaleItemResponse
	var total int64
	db.DB.Model(&models.StoreSaleItem{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ssi.id, ssi.item_id, si.item_name, ssi.quantity, 
			CAST(ssi.unit_price AS DECIMAL(10,2)) as unit_price, 
			CAST(ssi.total AS DECIMAL(10,2)) as total, 
			ssi.store_sale_id, ss.reference as sale_reference,
			ssi.created_at, ssi.updated_at
		FROM store_sale_items ssi
		LEFT JOIN store_items si ON ssi.item_id = si.id
		LEFT JOIN store_sales ss ON ssi.store_sale_id = ss.id
		WHERE ssi.deleted_at IS NULL
		ORDER BY ssi.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreSaleItemService) GetSaleItem(id string) (*dtos.StoreSaleItemResponse, error) {
	var result dtos.StoreSaleItemResponse
	query := `
		SELECT 
			ssi.id, ssi.item_id, si.item_name, ssi.quantity, 
			CAST(ssi.unit_price AS DECIMAL(10,2)) as unit_price, 
			CAST(ssi.total AS DECIMAL(10,2)) as total, 
			ssi.store_sale_id, ss.reference as sale_reference,
			ssi.created_at, ssi.updated_at
		FROM store_sale_items ssi
		LEFT JOIN store_items si ON ssi.item_id = si.id
		LEFT JOIN store_sales ss ON ssi.store_sale_id = ss.id
		WHERE ssi.id = ? AND ssi.deleted_at IS NULL
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

func (s *StoreSaleItemService) UpdateSaleItem(id string, req dtos.UpdateStoreSaleItemRequest, userID uint64) error {
	var item models.StoreSaleItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&item).Updates(map[string]interface{}{
		"item_id":       req.ItemID,
		"quantity":      req.Quantity,
		"unit_price":    fmt.Sprintf("%.2f", req.UnitPrice),
		"total":         fmt.Sprintf("%.2f", req.Total),
		"store_sale_id": req.StoreSaleID,
		"updated_by":    userID,
	}).Error
}

func (s *StoreSaleItemService) DeleteSaleItem(id string) error {
	return db.DB.Delete(&models.StoreSaleItem{}, id).Error
}
