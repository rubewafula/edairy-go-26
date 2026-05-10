package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type StoreStockService struct{}

func NewStoreStockService() *StoreStockService {
	return &StoreStockService{}
}

func (s *StoreStockService) CreateStock(req dtos.CreateStoreStockRequest) (*models.StoreStock, error) {
	stock := &models.StoreStock{
		ItemID:             req.ItemID,
		StoreID:            req.StoreID,
		Quantity:           req.Quantity,
		Unit:               req.Unit,
		BuyingPrice:        req.BuyingPrice,
		SellingPrice:       req.SellingPrice,
		CreditSellingPrice: req.CreditSellingPrice,
	}

	if stock.Unit == "" {
		stock.Unit = "KG"
	}

	if err := db.DB.Create(stock).Error; err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *StoreStockService) GetStocks(page, limit int) ([]dtos.StoreStockResponse, int64, error) {
	var results []dtos.StoreStockResponse
	var total int64
	db.DB.Model(&models.StoreStock{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ss.*, si.item_name, s.store AS store_name
		FROM store_stocks ss
		LEFT JOIN store_items si ON ss.item_id = si.id
		LEFT JOIN stores s ON ss.store_id = s.id
		WHERE ss.deleted_at IS NULL
		ORDER BY ss.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreStockService) GetStock(id string) (*dtos.StoreStockResponse, error) {
	var result dtos.StoreStockResponse
	query := `
		SELECT 
			ss.*, si.item_name, s.store AS store_name
		FROM store_stocks ss
		LEFT JOIN store_items si ON ss.item_id = si.id
		LEFT JOIN stores s ON ss.store_id = s.id
		WHERE ss.id = ? AND ss.deleted_at IS NULL
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

func (s *StoreStockService) UpdateStock(id string, req dtos.UpdateStoreStockRequest) error {
	var stock models.StoreStock
	if err := db.DB.First(&stock, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&stock).Updates(map[string]interface{}{
		"item_id":              req.ItemID,
		"store_id":             req.StoreID,
		"quantity":             req.Quantity,
		"unit":                 req.Unit,
		"buying_price":         req.BuyingPrice,
		"selling_price":        req.SellingPrice,
		"credit_selling_price": req.CreditSellingPrice,
	}).Error
}

func (s *StoreStockService) DeleteStock(id string) error {
	return db.DB.Delete(&models.StoreStock{}, id).Error
}
