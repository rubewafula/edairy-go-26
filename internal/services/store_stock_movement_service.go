package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreStockMovementService struct{}

func NewStoreStockMovementService() *StoreStockMovementService {
	return &StoreStockMovementService{}
}

func (s *StoreStockMovementService) CreateMovement(req dtos.CreateStoreStockMovementRequest, userID uint64) (*models.StoreStockMovement, error) {
	transactionDate := utils.ParseDate(req.TransactionDate)
	var movement models.StoreStockMovement

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 0. Fetch Movement Type and Direction
		var mType models.StoreStockMovementType
		if err := tx.First(&mType, req.MovementTypeID).Error; err != nil {
			return fmt.Errorf("movement type with ID %d not found: %w", req.MovementTypeID, err)
		}

		// 1. Create Master Transaction Header
		transaction := &models.Transaction{
			Reference:       fmt.Sprintf("SSM-%s-%d", transactionDate.Format("20060102"), req.StoreID),
			TransactionName: "Store Stock Movement",
			TransactionType: "STORES",
			TransactionDate: transactionDate,
			Description:     fmt.Sprintf("Stock movement type %s for Store ID: %d. %s", mType.MovementName, req.StoreID, req.Remarks),
			Status:          "POSTED",
			BaseModel: models.BaseModel{
				CreatedBy: userID,
			},
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 2. Create Store Stock Movement Header
		movement = models.StoreStockMovement{
			BaseModel: models.BaseModel{
				CreatedBy: userID,
			},
			TransactionID:   transaction.ID,
			TransactionDate: transactionDate,
			StoreID:         req.StoreID,
			MovementTypeID:  req.MovementTypeID,
			Remarks:         req.Remarks,
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}

		// 3. Process each item in the movement
		for _, itemReq := range req.Items {
			var stock models.StoreStock
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("item_id = ? AND store_id = ?", itemReq.ItemID, req.StoreID).
				First(&stock).Error; err != nil {
				return fmt.Errorf("item %d not found in store %d: %w", itemReq.ItemID, req.StoreID, err)
			}

			var newQuantity float64
			var glRule string
			var glAmount float64

			if mType.Direction == "IN" {
				newQuantity = stock.Quantity + itemReq.Quantity
				glRule = "STOCK_MOVEMENT_IN"
				glAmount = itemReq.Quantity * stock.BuyingPrice // Value of incoming stock
			} else if mType.Direction == "OUT" {
				if stock.Quantity < itemReq.Quantity {
					return fmt.Errorf("insufficient stock for item %d. Available: %.2f, Requested: %.2f", itemReq.ItemID, stock.Quantity, itemReq.Quantity)
				}
				newQuantity = stock.Quantity - itemReq.Quantity
				glRule = "STOCK_MOVEMENT_OUT"
				glAmount = itemReq.Quantity * stock.BuyingPrice // Value of outgoing stock (COGS)
			} else {
				return fmt.Errorf("unsupported movement direction: %s for type %s", mType.Direction, mType.MovementName)
			}

			// Update StoreStock
			stock.Quantity = newQuantity
			if err := tx.Save(&stock).Error; err != nil {
				return err
			}

			// Create StoreStockMovementItem
			movementItem := models.StoreStockMovementItem{
				BaseModel: models.BaseModel{
					CreatedBy: userID,
				},
				StoreStockMovementID: movement.ID,
				ItemID:               itemReq.ItemID,
				Quantity:             itemReq.Quantity,
				UnitCost:             itemReq.UnitCost,
				SellingPrice:         itemReq.SellingPrice,
				BalanceAfter:         newQuantity,
			}
			if err := tx.Create(&movementItem).Error; err != nil {
				return err
			}

			// Post GL Entries
			if glAmount > 0 {
				// Debit side of the rule
				if err := s.postGLEntry(tx, transaction.ID, glRule, true, glAmount, fmt.Sprintf("Stock movement (%s) for item %d", mType.MovementName, itemReq.ItemID), transactionDate, userID); err != nil {
					return err
				}
				// Credit side of the rule
				if err := s.postGLEntry(tx, transaction.ID, glRule, false, glAmount, fmt.Sprintf("Stock movement (%s) for item %d", mType.MovementName, itemReq.ItemID), transactionDate, userID); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &movement, nil
}

func (s *StoreStockMovementService) GetMovements(page, limit int) ([]dtos.StoreStockMovementResponse, int64, error) {
	var results []dtos.StoreStockMovementResponse
	var total int64
	db.DB.Model(&models.StoreStockMovement{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT
			ssm.id, ssm.transaction_date, ssm.store_id, s.store AS store_name,
			ssmt.movement_name AS movement_type, 
			ssm.remarks,
			ssm.created_at, ssm.updated_at
		FROM store_stock_movements ssm
		LEFT JOIN stores s ON ssm.store_id = s.id
		LEFT JOIN store_stock_movement_types ssmt ON ssm.movement_type_id = ssmt.id
		ORDER BY ssm.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range results {
		items, err := s.GetMovementItems(fmt.Sprintf("%d", results[i].ID))
		if err != nil {
			return nil, 0, err
		}
		results[i].Items = items
	}
	return results, total, err
}

func (s *StoreStockMovementService) GetMovement(id string) (*dtos.StoreStockMovementResponse, error) {
	var result dtos.StoreStockMovementResponse
	query := `
		SELECT 
			ssm.id, ssm.transaction_date, ssm.store_id, s.store AS store_name,
			ssmt.movement_name AS movement_type,  
			ssm.remarks,
			ssm.created_at, ssm.updated_at
		FROM store_stock_movements ssm
		LEFT JOIN stores s ON ssm.store_id = s.id
		LEFT JOIN store_stock_movement_types ssmt ON ssm.movement_type_id = ssmt.id
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
	items, err := s.GetMovementItems(id)
	if err != nil {
		return nil, err
	}
	result.Items = items
	return &result, nil
}

func (s *StoreStockMovementService) UpdateMovement(id string, req dtos.UpdateStoreStockMovementRequest, userID uint64) error {
	// This method would require complex logic to reverse previous GL entries and stock updates,
	// then re-apply new ones. For simplicity and to avoid unintended side effects in a real system,
	// stock movements are often immutable or reversed by new offsetting movements.
	// For now, we'll return an error indicating that direct updates are not supported for this complex transaction.
	return fmt.Errorf("direct update of stock movements is not supported due to complex accounting implications. Please create a new offsetting movement if adjustment is needed")
}

func (s *StoreStockMovementService) DeleteMovement(id string, userID uint64) error {
	// Similar to Update, deleting a stock movement has significant accounting implications.
	// A reversal movement should typically be created.
	return fmt.Errorf("direct deletion of stock movements is not supported. Please create a reversal movement if needed")
}

func (s *StoreStockMovementService) GetMovementItems(movementID string) ([]dtos.StoreStockMovementItemResponse, error) {
	var results []dtos.StoreStockMovementItemResponse
	query := `
		SELECT ssmi.*, si.item_name
		FROM store_stock_movement_items ssmi
		LEFT JOIN store_items si ON ssmi.item_id = si.id
		WHERE ssmi.store_stock_movement_id = ? AND ssmi.deleted_at IS NULL
		ORDER BY ssmi.id ASC
	`
	err := db.DB.Raw(query, movementID).Scan(&results).Error
	return results, err
}

func (s *StoreStockMovementService) postGLEntry(tx *gorm.DB, txID uint64, ruleType string, isDebit bool, amount float64, desc string, transDate time.Time, userID uint64) error {
	if amount <= 0 {
		return nil
	}

	var rule models.TransactionPostingRule
	if err := tx.Where("transaction_type = ?", ruleType).First(&rule).Error; err != nil {
		return fmt.Errorf("posting rule not found for %s: %w", ruleType, err)
	}

	entry := models.GeneralLedgerEntry{
		TransactionID:   txID,
		TransactionDate: transDate,
		Description:     desc,
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
	}

	if isDebit {
		entry.AccountID = rule.DebitAccountID
		entry.SubAccountID = rule.DebitSubAccountID
		entry.Debit = amount
	} else {
		entry.AccountID = rule.CreditAccountID
		entry.SubAccountID = rule.CreditSubAccountID
		entry.Credit = amount
	}

	return tx.Create(&entry).Error
}
