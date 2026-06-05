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

type StoreSaleService struct{}

func NewStoreSaleService() *StoreSaleService {
	return &StoreSaleService{}
}

func (s *StoreSaleService) CreateSale(req dtos.CreateStoreSaleRequest, userID uint64) (*models.StoreSale, error) {
	// 0. Calculate Amount Due and Validate
	amountDue := req.TotalAmount - req.AmountPaid
	if amountDue < 0 {
		return nil, fmt.Errorf("amount paid (%.2f) exceeds total amount (%.2f)", req.AmountPaid, req.TotalAmount)
	}

	// Parse the transaction date from the request
	transactionDate := utils.ParseDate(req.TransactionDate)

	var sale models.StoreSale
	var totalCOGS float64

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch Deduction Type for Stores
		var deductionType models.DeductionType
		if err := tx.Where("description = ?", "STORES").First(&deductionType).Error; err != nil {
			return fmt.Errorf("deduction type 'STORES' not found: %w", err)
		}

		// 2. Create master Transaction Record
		transaction := &models.Transaction{
			Reference:       req.Reference,
			TransactionName: "STORE SALES",
			TransactionType: "STORES",
			TransactionDate: transactionDate,
			Description:     fmt.Sprintf("Store sale to %s (ID: %d)", req.CustomerType, req.CustomerID),
			Status:          "POSTED",
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 3. Create Store Sale Header
		sale = models.StoreSale{
			BaseModel: models.BaseModel{
				CreatedBy: userID,
			},
			TotalAmount:   req.TotalAmount,
			AmountPaid:    req.AmountPaid,
			AmountDue:     amountDue,
			Reference:     req.Reference,
			StoreID:       req.StoreID,
			SaleType:      req.SaleType,
			CustomerID:    req.CustomerID,
			CustomerType:  req.CustomerType,
			TransactionID: transaction.ID,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return err
		}

		// 4. Process Sale Items and Inventory
		for _, itemReq := range req.Items {
			// Update Stock
			var stock models.StoreStock
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("item_id = ? AND store_id = ?", itemReq.ItemID, req.StoreID).First(&stock).Error; err != nil {
				return fmt.Errorf("item %d not found in store %d: %w", itemReq.ItemID, req.StoreID, err)
			}

			if stock.Quantity < float64(itemReq.Quantity) {
				return fmt.Errorf("insufficient stock for item %d", itemReq.ItemID)
			}

			stock.Quantity -= float64(itemReq.Quantity)
			if err := tx.Save(&stock).Error; err != nil {
				return err
			}

			// Track COGS (using stock buying price)
			totalCOGS += stock.BuyingPrice * float64(itemReq.Quantity)

			// Create Sale Item record
			saleItem := models.StoreSaleItem{
				BaseModel: models.BaseModel{
					CreatedBy: userID,
				},
				ItemID:      itemReq.ItemID,
				Quantity:    itemReq.Quantity,
				UnitPrice:   itemReq.UnitPrice,
				Total:       itemReq.Total,
				StoreSaleID: sale.ID,
			}
			if err := tx.Create(&saleItem).Error; err != nil {
				return err
			}

			// 4.1 Record Store Transaction (Inventory Audit)
			storeTx := models.StoreTransaction{
				BaseModel:       models.BaseModel{CreatedBy: userID},
				TransactionID:   transaction.ID,
				StoreID:         req.StoreID,
				StoreItemID:     itemReq.ItemID,
				TransactionType: "SALES",
				Quantity:        float64(itemReq.Quantity),
				UnitCost:        stock.BuyingPrice,
				SellingPrice:    itemReq.UnitPrice,
				BalanceAfter:    stock.Quantity,
				ReferenceType:   "SALE",
				ReferenceID:     sale.ID,
				Notes:           fmt.Sprintf("Sale to %s (ID: %d)", req.CustomerType, req.CustomerID),
				TransactedAt:    transactionDate,
			}
			if err := tx.Create(&storeTx).Error; err != nil {
				return err
			}
		}

		// 5. Handle Cash Portion
		if req.AmountPaid > 0 {
			// GL: Debit Cash (Asset Up)
			if err := s.postGLEntry(
				tx,
				transaction.ID,
				"STORE_SALE_REVENUE_CASH",
				true,
				req.AmountPaid,
				"Cash received for store sale",
				transactionDate,
				userID); err != nil {
				return err
			}
		}

		// 6. Handle Credit Portion (Recurrent Deduction)
		if req.AmountDue > 0 {
			recurrent := models.RecurrentDeduction{
				CustomerID:      req.CustomerID,
				TotalAmount:     req.AmountDue,
				PaidAmount:      0,
				RecurrentAmount: req.AmountDue, // Adjust if installment logic is needed
				DeductionTypeID: deductionType.ID,
				Reference:       req.Reference,
				CustomerType:    req.CustomerType,
				PrincipalAmount: req.AmountDue,
				TransactionDate: transactionDate,
				BaseModel: models.BaseModel{
					CreatedBy: userID,
				},
			}
			if err := tx.Create(&recurrent).Error; err != nil {
				return err
			}

			// GL: Debit Receivables (Asset Up)
			if err := s.postGLEntry(
				tx,
				transaction.ID,
				"STORE_SALE_REVENUE_CREDIT",
				true,
				req.AmountDue,
				"Credit sales for store sale",
				transactionDate, userID); err != nil {
				return err
			}
		}

		// 7. Core Accounting GL Entries
		// GL: Credit Sales Revenue (Income Up)
		if err := s.postGLEntry(
			tx,
			transaction.ID,
			"STORE_SALE_REVENUE",
			false,
			req.TotalAmount,
			"Store sales revenue",
			transactionDate, userID); err != nil {
			return err
		}

		// GL: Debit COGS (Expense Up)
		if err := s.postGLEntry(
			tx,
			transaction.ID,
			"STORE_SALE_COGS",
			true,
			totalCOGS,
			"Cost of goods sold",
			transactionDate, userID); err != nil {
			return err
		}

		// GL: Credit Inventory (Asset Down)
		if err := s.postGLEntry(
			tx,
			transaction.ID,
			"STORE_SALE_INVENTORY",
			false,
			totalCOGS,
			"Inventory reduction for sale",
			transactionDate,
			userID); err != nil {
			return err
		}

		// 8. Confirm Ledger Balance (Double Entry Check)
		var glBalance struct {
			TotalDebit  float64
			TotalCredit float64
		}
		err := tx.Model(&models.GeneralLedgerEntry{}).
			Select("SUM(debit) as total_debit, SUM(credit) as total_credit").
			Where("transaction_id = ?", transaction.ID).
			Scan(&glBalance).Error
		if err != nil {
			return err
		}

		if glBalance.TotalDebit != glBalance.TotalCredit {
			return fmt.Errorf("ledger imbalance: Total Debits (%.2f) != Total Credits (%.2f)", glBalance.TotalDebit, glBalance.TotalCredit)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (s *StoreSaleService) postGLEntry(
	tx *gorm.DB,
	txID uint64,
	ruleTransactionType string,
	isDebit bool, // Explicit boolean flags eliminate edge case failures on exact zero updates
	amount float64,
	desc string,
	transDate time.Time,
	userID uint64,
) error {
	if amount == 0 {
		return nil // Avoid polluting the ledger with dead lines
	}

	var rule models.TransactionPostingRule
	if err := tx.Where("transaction_type = ?", ruleTransactionType).First(&rule).Error; err != nil {
		return fmt.Errorf("posting rule not found for %s: %w", ruleTransactionType, err)
	}

	var accountID uint64
	var debitAmt, creditAmt float64

	var subAccountID *uint64

	if isDebit {
		accountID = rule.DebitAccountID
		debitAmt = amount

		// If the rule has a sub-account, pass its pointer along.
		// If rule.DebitSubAccountID is already a pointer and is nil, subAccountID stays nil (NULL)
		if rule.DebitSubAccountID != nil && *rule.DebitSubAccountID != 0 {
			subAccountID = rule.DebitSubAccountID
		}
	} else {
		accountID = rule.CreditAccountID
		creditAmt = amount

		if rule.CreditSubAccountID != nil && *rule.CreditSubAccountID != 0 {
			subAccountID = rule.CreditSubAccountID
		}
	}

	entry := models.GeneralLedgerEntry{
		TransactionID:   txID,
		AccountID:       accountID,
		SubAccountID:    subAccountID,
		Debit:           debitAmt,
		Credit:          creditAmt,
		TransactionDate: transDate,
		Description:     desc,
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
	}

	return tx.Create(&entry).Error
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
