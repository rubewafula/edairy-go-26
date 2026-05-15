package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type PurchaseOrderService struct{}

func NewPurchaseOrderService() *PurchaseOrderService {
	return &PurchaseOrderService{}
}

func (s *PurchaseOrderService) CreatePO(req dtos.CreatePurchaseOrderRequest, userID uint64) (*models.PurchaseOrder, error) {
	po := &models.PurchaseOrder{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		SupplierID:      req.SupplierID,
		SupplierQuoteID: req.SupplierQuoteID,
		PoNumber:        req.PoNumber,
		PoDate:          utils.ParseDate(req.PoDate),
		Status:          "draft",
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		for _, item := range req.Items {
			totalAmount += item.Quantity * item.UnitPrice
		}
		po.TotalAmount = totalAmount

		if err := tx.Create(po).Error; err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			item := &models.PurchaseOrderItem{
				PurchaseOrderID: po.ID,
				ItemID:          itemReq.ItemID,
				Description:     itemReq.Description,
				Quantity:        itemReq.Quantity,
				UnitPrice:       itemReq.UnitPrice,
				TotalPrice:      itemReq.Quantity * itemReq.UnitPrice,
			}
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return po, err
}

func (s *PurchaseOrderService) CreateRequisition(req dtos.CreatePurchaseRequisitionRequest, userID uint64) (*models.PurchaseRequisition, error) {
	requisition := &models.PurchaseRequisition{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		RequisitionNo:   req.RequisitionNo,
		RequisitionDate: utils.ParseDate(req.RequisitionDate),
		Description:     req.Description,
		Status:          "draft",
	}

	if err := db.DB.Create(requisition).Error; err != nil {
		return nil, err
	}
	return requisition, nil
}

func (s *PurchaseOrderService) GetRequisitions(page, limit int) ([]dtos.PurchaseRequisitionResponse, int64, error) {
	var results []dtos.PurchaseRequisitionResponse
	var total int64
	db.DB.Model(&models.PurchaseRequisition{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.PurchaseRequisition{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *PurchaseOrderService) GetRequisition(id string) (*dtos.PurchaseRequisitionResponse, error) {
	var result dtos.PurchaseRequisitionResponse
	err := db.DB.Model(&models.PurchaseRequisition{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *PurchaseOrderService) UpdateRequisition(id string, req dtos.UpdatePurchaseRequisitionRequest, userID uint64) (*models.PurchaseRequisition, error) {
	var requisition models.PurchaseRequisition
	if err := db.DB.First(&requisition, id).Error; err != nil {
		return nil, err
	}

	requisition.RequisitionNo = req.RequisitionNo
	requisition.RequisitionDate = utils.ParseDate(req.RequisitionDate)
	requisition.Description = req.Description
	requisition.Status = req.Status
	requisition.UpdatedBy = userID

	if err := db.DB.Save(&requisition).Error; err != nil {
		return nil, err
	}

	return &requisition, nil
}

func (s *PurchaseOrderService) DeleteRequisition(id string, userID uint64) error {
	var requisition models.PurchaseRequisition
	if err := db.DB.First(&requisition, id).Error; err != nil {
		return err
	}
	// Update user ID for audit before soft deletion
	return db.DB.Model(&requisition).Update("updated_by", userID).Delete(&requisition).Error
}

func (s *PurchaseOrderService) GetPOs(page, limit int) ([]dtos.PurchaseOrderResponse, int64, error) {
	var results []dtos.PurchaseOrderResponse
	var total int64
	db.DB.Model(&models.PurchaseOrder{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			po.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM purchase_orders po
		LEFT JOIN suppliers s ON po.supplier_id = s.id
		WHERE po.deleted_at IS NULL
		ORDER BY po.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *PurchaseOrderService) CreateRequisitionItem(req dtos.CreatePurchaseRequisitionItemRequest, userID uint64) (*models.PurchaseRequisitionItem, error) {
	item := &models.PurchaseRequisitionItem{
		BaseModel:             models.BaseModel{CreatedBy: userID},
		PurchaseRequisitionID: &req.PurchaseRequisitionID,
		ItemID:                &req.ItemID,
		Quantity:              req.Quantity,
		Status:                "pending",
	}
	err := db.DB.Create(item).Error
	return item, err
}

func (s *PurchaseOrderService) GetRequisitionItems(reqID string) ([]dtos.PurchaseRequisitionItemResponse, error) {
	var items []dtos.PurchaseRequisitionItemResponse
	query := `
		SELECT pri.id, si.item_name, pri.quantity, pri.status, pri.created_at
		FROM purchase_requisition_items pri
		LEFT JOIN store_items si ON pri.item_id = si.id
		WHERE pri.purchase_requisition_id = ? AND pri.deleted_at IS NULL
	`
	err := db.DB.Raw(query, reqID).Scan(&items).Error
	return items, err
}

func (s *PurchaseOrderService) GetRequisitionItem(id string) (*dtos.PurchaseRequisitionItemResponse, error) {
	var result dtos.PurchaseRequisitionItemResponse
	query := `
		SELECT pri.id, si.item_name, pri.quantity, pri.status, pri.created_at
		FROM purchase_requisition_items pri
		LEFT JOIN store_items si ON pri.item_id = si.id
		WHERE pri.id = ? AND pri.deleted_at IS NULL
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

func (s *PurchaseOrderService) UpdateRequisitionItem(id string, req dtos.UpdatePurchaseRequisitionItemRequest, userID uint64) error {
	var item models.PurchaseRequisitionItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	item.ItemID = &req.ItemID
	item.Quantity = req.Quantity
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedBy = userID

	return db.DB.Save(&item).Error
}

func (s *PurchaseOrderService) DeleteRequisitionItem(id string, userID uint64) error {
	var item models.PurchaseRequisitionItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&item).Update("updated_by", userID).Delete(&item).Error
}

func (s *PurchaseOrderService) GetPOItems(poID string) ([]dtos.PurchaseOrderItemResponse, error) {
	var items []dtos.PurchaseOrderItemResponse
	query := `
		SELECT poi.*, si.item_name
		FROM purchase_order_items poi
		LEFT JOIN store_items si ON poi.item_id = si.id
		WHERE poi.purchase_order_id = ?
	`
	err := db.DB.Raw(query, poID).Scan(&items).Error
	return items, err
}

func (s *PurchaseOrderService) GetPO(id string) (*dtos.PurchaseOrderResponse, error) {
	var result dtos.PurchaseOrderResponse
	query := `
		SELECT 
			po.*, 
			CASE WHEN s.company_name != '' THEN s.company_name ELSE CONCAT(s.first_name, ' ', s.last_name) END as supplier_name
		FROM purchase_orders po
		LEFT JOIN suppliers s ON po.supplier_id = s.id
		WHERE po.id = ? AND po.deleted_at IS NULL
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

func (s *PurchaseOrderService) UpdatePO(id string, req dtos.UpdatePurchaseOrderRequest, userID uint64) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	if err := db.DB.First(&po, id).Error; err != nil {
		return nil, err
	}

	// Update fields from the request
	po.SupplierID = req.SupplierID
	po.SupplierQuoteID = req.SupplierQuoteID
	po.PoNumber = req.PoNumber
	po.PoDate = utils.ParseDate(req.PoDate)
	po.Status = req.Status
	po.TotalAmount = req.TotalAmount // Assuming total amount can be updated manually or recalculated elsewhere
	po.UpdatedBy = userID

	// Save the updated purchase order
	if err := db.DB.Save(&po).Error; err != nil {
		return nil, err
	}

	return &po, nil
}

func (s *PurchaseOrderService) DeletePO(id string, userID uint64) error {
	var po models.PurchaseOrder
	if err := db.DB.First(&po, id).Error; err != nil {
		return err
	}
	// Update user ID for audit before soft deletion
	return db.DB.Model(&po).Update("updated_by", userID).Delete(&po).Error
}

func (s *PurchaseOrderService) GetPOItem(id string) (*dtos.PurchaseOrderItemResponse, error) {
	var result dtos.PurchaseOrderItemResponse
	query := `
		SELECT poi.*, si.item_name
		FROM purchase_order_items poi
		LEFT JOIN store_items si ON poi.item_id = si.id
		WHERE poi.id = ?
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

func (s *PurchaseOrderService) UpdatePOItem(id string, req dtos.UpdatePurchaseOrderItemRequest, userID uint64) error {
	var item models.PurchaseOrderItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		item.ItemID = req.ItemID
		item.Description = req.Description
		item.Quantity = req.Quantity
		item.UnitPrice = req.UnitPrice
		item.TotalPrice = req.Quantity * req.UnitPrice

		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return s.recalculatePOTotals(tx, item.PurchaseOrderID, userID)
	})

	return err
}

func (s *PurchaseOrderService) DeletePOItem(id string, userID uint64) error {
	var item models.PurchaseOrderItem
	if err := db.DB.First(&item, id).Error; err != nil {
		return err
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&item).Error; err != nil {
			return err
		}

		return s.recalculatePOTotals(tx, item.PurchaseOrderID, userID)
	})

	return err
}

func (s *PurchaseOrderService) recalculatePOTotals(tx *gorm.DB, poID uint64, userID uint64) error {
	var totalAmount float64
	err := tx.Model(&models.PurchaseOrderItem{}).
		Where("purchase_order_id = ?", poID).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&totalAmount).Error
	if err != nil {
		return err
	}

	return tx.Model(&models.PurchaseOrder{}).Where("id = ?", poID).Updates(map[string]interface{}{
		"total_amount": totalAmount,
		"updated_by":   userID,
	}).Error
}
