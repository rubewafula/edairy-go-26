package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type StoreSaleService struct{}

func NewStoreSaleService() *StoreSaleService {
	return &StoreSaleService{}
}

func (s *StoreSaleService) CreateSale(req dtos.CreateStoreSaleRequest, userID uint64) (*models.StoreSale, error) {
	sale := &models.StoreSale{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		TotalAmount:   req.TotalAmount,
		AmountPaid:    req.AmountPaid,
		AmountDue:     req.AmountDue,
		Reference:     req.Reference,
		StoreID:       req.StoreID,
		SaleType:      req.SaleType,
		CustomerID:    req.CustomerID,
		CustomerType:  req.CustomerType,
		TransactionID: req.TransactionID,
	}

	if err := db.DB.Create(sale).Error; err != nil {
		return nil, err
	}
	return sale, nil
}

func (s *StoreSaleService) GetSales(page, limit int) ([]dtos.StoreSaleResponse, int64, error) {
	var results []dtos.StoreSaleResponse
	var total int64
	db.DB.Model(&models.StoreSale{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ss.*, s.store AS store_name
		FROM store_sales ss
		LEFT JOIN stores s ON ss.store_id = s.id
		WHERE ss.deleted_at IS NULL
		ORDER BY ss.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StoreSaleService) GetSale(id string) (*dtos.StoreSaleResponse, error) {
	var result dtos.StoreSaleResponse
	query := `
		SELECT 
			ss.*, s.store AS store_name
		FROM store_sales ss
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

func (s *StoreSaleService) UpdateSale(id string, req dtos.UpdateStoreSaleRequest, userID uint64) error {
	var sale models.StoreSale
	if err := db.DB.First(&sale, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&sale).Updates(map[string]interface{}{
		"total_amount":   req.TotalAmount,
		"amount_paid":    req.AmountPaid,
		"amount_due":     req.AmountDue,
		"reference":      req.Reference,
		"store_id":       req.StoreID,
		"sale_type":      req.SaleType,
		"customer_id":    req.CustomerID,
		"customer_type":  req.CustomerType,
		"transaction_id": req.TransactionID,
		"updated_by":     userID,
	}).Error
}

func (s *StoreSaleService) DeleteSale(id string) error {
	return db.DB.Delete(&models.StoreSale{}, id).Error
}
