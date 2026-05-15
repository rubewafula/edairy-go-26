package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type SupplyService struct{}

func NewSupplyService() *SupplyService {
	return &SupplyService{}
}

func (s *SupplyService) CreateSupply(req dtos.CreateSupplyRequest, userID uint64) (*models.Supply, error) {
	supply := &models.Supply{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		VendorID:        req.VendorID,
		PaymentTypeID:   req.PaymentTypeID,
		PurchaseOrderID: req.PurchaseOrderID,
		Reference:       req.Reference,
		Activity:        req.Activity,
		SuppliedDate:    utils.ParseDate(req.SuppliedDate),
		StoreID:         req.StoreID,
		PaymentTermID:   req.PaymentTermID,
		ItemCount:       uint64(len(req.Items)),
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		for _, itemReq := range req.Items {
			totalAmount += float64(itemReq.Quantity) * itemReq.UnitPrice
		}
		supply.TotalAmount = totalAmount

		if err := tx.Create(supply).Error; err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			item := &models.SuppliedItem{
				BaseModel:  models.BaseModel{CreatedBy: userID},
				SupplyID:   supply.ID,
				ItemID:     itemReq.ItemID,
				Quantity:   itemReq.Quantity,
				UnitPrice:  itemReq.UnitPrice,
				TotalPrice: float64(itemReq.Quantity) * itemReq.UnitPrice,
			}
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return supply, err
}

func (s *SupplyService) GetSupplies(page, limit int) ([]dtos.SupplyResponse, int64, error) {
	var results []dtos.SupplyResponse
	var total int64
	db.DB.Model(&models.Supply{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sup.*, 
			CASE WHEN v.company_name != '' THEN v.company_name ELSE CONCAT(v.first_name, ' ', v.last_name) END as vendor_name,
			st.store as store_name
		FROM supplies sup
		LEFT JOIN suppliers v ON sup.vendor_id = v.id
		LEFT JOIN stores st ON sup.store_id = st.id
		WHERE sup.deleted_at IS NULL
		ORDER BY sup.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplyService) GetSuppliedItems(supplyID string) ([]dtos.SuppliedItemResponse, error) {
	var results []dtos.SuppliedItemResponse
	query := `
		SELECT si.*, i.item_name
		FROM supplied_items si
		LEFT JOIN store_items i ON si.item_id = i.id
		WHERE si.supply_id = ? AND si.deleted_at IS NULL
	`
	err := db.DB.Raw(query, supplyID).Scan(&results).Error
	return results, err
}

func (s *SupplyService) GetSupply(id string) (*models.Supply, error) {
	var supply models.Supply
	err := db.DB.First(&supply, id).Error
	return &supply, err
}
