package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type SupplierQuoteService struct{}

func NewSupplierQuoteService() *SupplierQuoteService {
	return &SupplierQuoteService{}
}

func (s *SupplierQuoteService) CreateQuote(req dtos.CreateSupplierQuoteRequest, userID uint64) (*models.SupplierQuote, error) {
	quote := &models.SupplierQuote{
		BaseModel:        models.BaseModel{CreatedBy: userID},
		VendorID:         req.VendorID,
		Description:      req.Description,
		Status:           "pending",
		RfqNo:            req.RfqNo,
		SupplierQuoteRef: req.SupplierQuoteRef,
	}
	if err := db.DB.Create(quote).Error; err != nil {
		return nil, err
	}
	return quote, nil
}

func (s *SupplierQuoteService) GetQuotes(page, limit int) ([]dtos.SupplierQuoteResponse, int64, error) {
	var results []dtos.SupplierQuoteResponse
	var total int64
	db.DB.Model(&models.SupplierQuote{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sq.*, 
			CASE 
				WHEN s.company_name IS NOT NULL AND s.company_name != '' THEN s.company_name 
				ELSE CONCAT(s.first_name, ' ', s.last_name) 
			END as vendor_name
		FROM supplier_quotes sq
		LEFT JOIN suppliers s ON sq.supplier_id = s.id
		WHERE sq.deleted_at IS NULL
		ORDER BY sq.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplierQuoteService) CreateQuoteItem(req dtos.CreateSupplierQuoteItemRequest, userID uint64) (*models.SupplierQuoteItem, error) {
	item := &models.SupplierQuoteItem{
		BaseModel:         models.BaseModel{CreatedBy: userID},
		SupplierQuoteID:   req.SupplierQuoteID,
		ItemID:            req.ItemID,
		QuantityRequested: req.QuantityRequested,
		UnitPrice:         req.UnitPrice,
		Notes:             req.Notes,
	}
	if err := db.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SupplierQuoteService) GetQuoteItems(quoteID string) ([]dtos.SupplierQuoteItemResponse, error) {
	var results []dtos.SupplierQuoteItemResponse
	query := `
		SELECT sqi.*, si.item_name
		FROM supplier_quote_items sqi
		LEFT JOIN store_items si ON sqi.item_id = si.id
		WHERE sqi.supplier_quote_id = ? AND sqi.deleted_at IS NULL
	`
	err := db.DB.Raw(query, quoteID).Scan(&results).Error
	return results, err
}

func (s *SupplierQuoteService) UpdateQuoteStatus(id string, status string, userID uint64) error {
	return db.DB.Model(&models.SupplierQuote{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_by": userID,
	}).Error
}

func (s *SupplierQuoteService) GetQuote(id string) (*models.SupplierQuote, error) {
	var quote models.SupplierQuote
	err := db.DB.First(&quote, id).Error
	return &quote, err
}

func (s *SupplierQuoteService) GetQuoteItem(id string) (*dtos.SupplierQuoteItemResponse, error) {
	var result dtos.SupplierQuoteItemResponse
	query := `
		SELECT sqi.*, si.item_name
		FROM supplier_quote_items sqi
		LEFT JOIN store_items si ON sqi.item_id = si.id
		WHERE sqi.id = ? AND sqi.deleted_at IS NULL
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

func (s *SupplierQuoteService) UpdateQuoteItem(id string, req dtos.UpdateSupplierQuoteItemRequest, userID uint64) error {
	var item models.SupplierQuoteItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	item.QuantitySupplied = req.QuantitySupplied
	item.UnitPrice = req.UnitPrice
	item.Notes = req.Notes
	item.UpdatedBy = userID

	return db.DB.Save(&item).Error
}

func (s *SupplierQuoteService) DeleteQuoteItem(id string) error {
	return db.DB.Delete(&models.SupplierQuoteItem{}, id).Error
}
