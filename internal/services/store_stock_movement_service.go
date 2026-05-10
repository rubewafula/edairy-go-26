package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type StoreStockMovementService struct{}

func NewStoreStockMovementService() *StoreStockMovementService {
	return &StoreStockMovementService{}
}

func (s *StoreStockMovementService) CreateMovement(req dtos.CreateStoreStockMovementRequest, userID uint64) (*models.StoreStockMovement, error) {
	movement := &models.StoreStockMovement{
		TransactionDate: utils.ParseDate(req.TransactionDate),
		StoreID:         req.StoreID,
		ItemID:          req.ItemID,
		MovementType:    req.MovementType,
		ReferenceTable:  req.ReferenceTable,
		ReferenceID:     req.ReferenceID,
		QtyIn:           req.QtyIn,
		QtyOut:          req.QtyOut,
		BalanceAfter:    req.BalanceAfter,
		UnitCost:        req.UnitCost,
		SellingPrice:    req.SellingPrice,
		Remarks:         req.Remarks,
		CreatedBy:       userID,
	}

	if err := db.DB.Create(movement).Error; err != nil {
		return nil, err
	}
	return movement, nil
}

func (s *StoreStockMovementService) GetMovements(page, limit int) ([]dtos.StoreStockMovementResponse, int64, error) {
	var results []dtos.StoreStockMovementResponse
	var total int64
	db.DB.Model(&models.StoreStockMovement{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ssm.*, s.name AS store_name, si.item_name
		FROM store_stock_movements ssm
		LEFT JOIN stores s ON ssm.store_id = s.id
		LEFT JOIN store_items si ON ssm.item_id = si.id
		ORDER BY ssm.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreStockMovementService) GetMovement(id string) (*dtos.StoreStockMovementResponse, error) {
	var result dtos.StoreStockMovementResponse
	query := `
		SELECT 
			ssm.*, s.name AS store_name, si.item_name
		FROM store_stock_movements ssm
		LEFT JOIN stores s ON ssm.store_id = s.id
		LEFT JOIN store_items si ON ssm.item_id = si.id
		WHERE ssm.id = ?
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

func (s *StoreStockMovementService) UpdateMovement(id string, req dtos.UpdateStoreStockMovementRequest) error {
	return db.DB.Model(&models.StoreStockMovement{}).Where("id = ?", id).Updates(map[string]interface{}{
		"transaction_date": utils.ParseDate(req.TransactionDate),
		"store_id":         req.StoreID,
		"item_id":          req.ItemID,
		"movement_type":    req.MovementType,
		"reference_table":  req.ReferenceTable,
		"reference_id":     req.ReferenceID,
		"qty_in":           req.QtyIn,
		"qty_out":          req.QtyOut,
		"balance_after":    req.BalanceAfter,
		"unit_cost":        req.UnitCost,
		"selling_price":    req.SellingPrice,
		"remarks":          req.Remarks,
	}).Error
}

func (s *StoreStockMovementService) DeleteMovement(id string) error {
	return db.DB.Delete(&models.StoreStockMovement{}, id).Error
}
