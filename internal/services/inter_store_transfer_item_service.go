package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type InterStoreTransferItemService struct{}

func NewInterStoreTransferItemService() *InterStoreTransferItemService {
	return &InterStoreTransferItemService{}
}

func (s *InterStoreTransferItemService) CreateTransferItem(req dtos.CreateInterStoreTransferItemRequest, userID uint64) (*models.InterStoreTransferItem, error) {
	item := &models.InterStoreTransferItem{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		TransferID: req.TransferID,
		ItemID:     req.ItemID,
		Quantity:   fmt.Sprintf("%.2f", req.Quantity),
		StockID:    req.StockID,
	}

	if err := db.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (s *InterStoreTransferItemService) GetTransferItems(page, limit int) ([]dtos.InterStoreTransferItemResponse, int64, error) {
	var results []dtos.InterStoreTransferItemResponse
	var total int64
	db.DB.Model(&models.InterStoreTransferItem{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			iti.id, ist.reference, iti.item_id, si.item_name, 
			CAST(iti.quantity AS DECIMAL(10,2)) as quantity, 
			iti.inter_store_transfer_id as transfer_id, iti.stock_id,
			iti.created_at, iti.updated_at
		FROM inter_store_transfer_items iti
		LEFT JOIN store_items si ON iti.item_id = si.id
		LEFT JOIN inter_store_transfers ist ON iti.inter_store_transfer_id = ist.id
		WHERE iti.deleted_at IS NULL
		ORDER BY iti.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *InterStoreTransferItemService) GetTransferItem(id string) (*dtos.InterStoreTransferItemResponse, error) {
	var result dtos.InterStoreTransferItemResponse
	query := `
		SELECT 
			iti.id, ist.reference, iti.item_id, si.item_name, 
			CAST(iti.quantity AS DECIMAL(10,2)) as quantity, 
			iti.inter_store_transfer_id as transfer_id, iti.stock_id,
			iti.created_at, iti.updated_at
		FROM inter_store_transfer_items iti
		LEFT JOIN store_items si ON iti.item_id = si.id
		LEFT JOIN inter_store_transfers ist ON iti.inter_store_transfer_id = ist.id
		WHERE iti.id = ? AND iti.deleted_at IS NULL
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

func (s *InterStoreTransferItemService) UpdateTransferItem(id string, req dtos.UpdateInterStoreTransferItemRequest, userID uint64) error {
	var item models.InterStoreTransferItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&item).Updates(map[string]interface{}{
		"inter_store_transfer_id": req.TransferID,
		"item_id":                 req.ItemID,
		"quantity":                fmt.Sprintf("%.2f", req.Quantity),
		"stock_id":                req.StockID,
		"updated_by":              userID,
	}).Error
}

func (s *InterStoreTransferItemService) DeleteTransferItem(id string) error {
	var item models.InterStoreTransferItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&item).Error
}
