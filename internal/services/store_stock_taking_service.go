package services

import (
	"fmt"
	"math"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreStockTakingService struct{}

func NewStoreStockTakingService() *StoreStockTakingService {
	return &StoreStockTakingService{}
}

func (s *StoreStockTakingService) CreateStockTaking(req dtos.CreateStoreStockTakingRequest, userID uint64) ([]models.StoreStockTaking, error) {
	transactionDate := utils.ParseDate(req.StockTakeDate)
	var createdItems []models.StoreStockTaking

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create Master Transaction Record
		transaction := &models.Transaction{
			Reference:       req.StockTakeNo,
			TransactionName: "STOCK TAKING STOCK ADJUSTMENT",
			TransactionType: "STORES",
			TransactionDate: transactionDate,
			Description:     fmt.Sprintf("Stock taking for Store ID: %d. %s", req.StoreID, req.Remarks),
			Status:          "POSTED",
			BaseModel: models.BaseModel{
				CreatedBy: userID,
			},
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			// 2. Fetch current system quantity and lock row
			var stock models.StoreStock
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("item_id = ? AND store_id = ?", itemReq.ItemID, req.StoreID).
				First(&stock).Error; err != nil {
				return fmt.Errorf("item %d not found in store %d: %w", itemReq.ItemID, req.StoreID, err)
			}

			systemQty := stock.Quantity
			variance := itemReq.PhysicalQuantity - systemQty
			varianceValue := math.Abs(variance) * stock.BuyingPrice

			// 3. Update Inventory Balance
			stock.Quantity = itemReq.PhysicalQuantity
			if err := tx.Save(&stock).Error; err != nil {
				return err
			}

			// 4. Create Stock Taking Record
			stockTaking := models.StoreStockTaking{
				StockTakeNo:      req.StockTakeNo,
				StoreID:          req.StoreID,
				ItemID:           itemReq.ItemID,
				SystemQuantity:   systemQty,
				PhysicalQuantity: itemReq.PhysicalQuantity,
				VarianceQuantity: variance,
				Remarks:          req.Remarks,
				StockTakeDate:    transactionDate,
				CreatedBy:        userID,
			}

			if err := tx.Create(&stockTaking).Error; err != nil {
				return err
			}

			// 5. Post GL Entries if there is a variance
			if variance != 0 {
				rule := "STOCK_TAKE_GAIN"
				if variance < 0 {
					rule = "STOCK_TAKE_LOSS"
				}

				// Debit side of the rule
				if err := s.postGLEntry(
					tx,
					transaction.ID,
					rule,
					true,
					varianceValue,
					"Stock variance adjustment: "+req.StockTakeNo,
					transactionDate,
					userID); err != nil {
					return err
				}
				// Credit side of the rule
				if err := s.postGLEntry(
					tx,
					transaction.ID,
					rule,
					false,
					varianceValue,
					"Stock variance adjustment: "+req.StockTakeNo,
					transactionDate,
					userID); err != nil {
					return err
				}
			}

			createdItems = append(createdItems, stockTaking)
		}

		return nil
	})

	return createdItems, err
}

func (s *StoreStockTakingService) postGLEntry(
	tx *gorm.DB, txID uint64,
	ruleType string,
	isDebit bool,
	amount float64,
	desc string,
	transDate time.Time,
	userID uint64) error {

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
