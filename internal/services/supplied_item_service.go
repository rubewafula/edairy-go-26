package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type SuppliedItemService struct{}

func NewSuppliedItemService() *SuppliedItemService {
	return &SuppliedItemService{}
}

func (s *SuppliedItemService) GetSuppliedItem(id string) (*dtos.SuppliedItemResponse, error) {
	var result dtos.SuppliedItemResponse
	query := `
		SELECT si.*, i.item_name
		FROM supplied_items si
		LEFT JOIN store_items i ON si.item_id = i.id
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

func (s *SuppliedItemService) UpdateSuppliedItem(id string, req dtos.UpdateSuppliedItemRequest, userID uint64) error {
	var item models.SuppliedItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		item.ItemID = req.ItemID
		item.Quantity = req.Quantity
		item.UnitPrice = req.UnitPrice
		item.TotalPrice = float64(req.Quantity) * req.UnitPrice
		item.UpdatedBy = userID

		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return s.recalculateSupplyTotals(tx, item.SupplyID)
	})

	return err
}

func (s *SuppliedItemService) DeleteSuppliedItem(id string, userID uint64) error {
	var item models.SuppliedItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&item).Update("updated_by", userID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&item).Error; err != nil {
			return err
		}

		return s.recalculateSupplyTotals(tx, item.SupplyID)
	})

	return err
}

func (s *SuppliedItemService) recalculateSupplyTotals(tx *gorm.DB, supplyID uint64) error {
	var stats struct {
		TotalAmount float64
		ItemCount   int64
	}

	err := tx.Model(&models.SuppliedItem{}).
		Select("COALESCE(SUM(total_price), 0) as total_amount, COUNT(id) as item_count").
		Where("supply_id = ? AND deleted_at IS NULL", supplyID).
		Scan(&stats).Error

	if err != nil {
		return err
	}

	// If item count is 0, the total amount should be 0.
	// Update the parent Supply header
	return tx.Model(&models.Supply{}).Where("id = ?", supplyID).Updates(map[string]interface{}{
		"total_amount": stats.TotalAmount,
		"item_count":   stats.ItemCount,
	}).Error
}

func (s *SuppliedItemService) GetSuppliedItemsBySupply(supplyID string) ([]dtos.SuppliedItemResponse, error) {
	var results []dtos.SuppliedItemResponse
	query := "SELECT si.*, i.item_name FROM supplied_items si LEFT JOIN store_items i ON si.item_id = i.id WHERE si.supply_id = ? AND si.deleted_at IS NULL"
	err := db.DB.Raw(query, supplyID).Scan(&results).Error
	return results, err
}
