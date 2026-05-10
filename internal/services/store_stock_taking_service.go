package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type StoreStockTakingService struct{}

func NewStoreStockTakingService() *StoreStockTakingService {
	return &StoreStockTakingService{}
}

func (s *StoreStockTakingService) CreateStockTaking(req dtos.CreateStoreStockTakingRequest, userID uint64) (*models.StoreStockTaking, error) {
	stockTaking := &models.StoreStockTaking{
		StockTakeNo:      req.StockTakeNo,
		StoreID:          req.StoreID,
		ItemID:           req.ItemID,
		SystemQuantity:   req.SystemQuantity,
		PhysicalQuantity: req.PhysicalQuantity,
		VarianceQuantity: req.VarianceQuantity,
		Remarks:          req.Remarks,
		StockTakeDate:    utils.ParseDate(req.StockTakeDate),
		CreatedBy:        userID,
	}

	if err := db.DB.Create(stockTaking).Error; err != nil {
		return nil, err
	}
	return stockTaking, nil
}

func (s *StoreStockTakingService) GetStockTakings(page, limit int) ([]dtos.StoreStockTakingResponse, int64, error) {
	var results []dtos.StoreStockTakingResponse
	var total int64
	db.DB.Model(&models.StoreStockTaking{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sst.*, s.store AS store_name, si.item_name
		FROM store_stock_takings sst
		LEFT JOIN stores s ON sst.store_id = s.id
		LEFT JOIN store_items si ON sst.item_id = si.id
		ORDER BY sst.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreStockTakingService) GetStockTaking(id string) (*dtos.StoreStockTakingResponse, error) {
	var result dtos.StoreStockTakingResponse
	query := `
		SELECT 
			sst.*, s.store AS store_name, si.item_name
		FROM store_stock_takings sst
		LEFT JOIN stores s ON sst.store_id = s.id
		LEFT JOIN store_items si ON sst.item_id = si.id
		WHERE sst.id = ?
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

func (s *StoreStockTakingService) UpdateStockTaking(id string, req dtos.UpdateStoreStockTakingRequest) error {
	return db.DB.Model(&models.StoreStockTaking{}).Where("id = ?", id).Updates(map[string]interface{}{
		"stock_take_no":     req.StockTakeNo,
		"store_id":          req.StoreID,
		"item_id":           req.ItemID,
		"system_quantity":   req.SystemQuantity,
		"physical_quantity": req.PhysicalQuantity,
		"variance_quantity": req.VarianceQuantity,
		"remarks":           req.Remarks,
		"stock_take_date":   utils.ParseDate(req.StockTakeDate),
	}).Error
}

func (s *StoreStockTakingService) DeleteStockTaking(id string) error {
	return db.DB.Delete(&models.StoreStockTaking{}, id).Error
}
