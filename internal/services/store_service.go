package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type StoreService struct {
	notificationService *UINotificationService
}

func NewStoreService() *StoreService {
	return &StoreService{
		notificationService: NewUINotificationService(),
	}
}

func (s *StoreService) CreateStore(req dtos.CreateStoreRequest) (*models.Store, error) {
	store := &models.Store{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(store).Error; err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StoreService) GetStores() ([]models.Store, int64, error) {
	var stores []models.Store
	var total int64
	db.DB.Model(&models.Store{}).Count(&total)
	err := db.DB.Find(&stores).Error
	return stores, total, err
}

func (s *StoreService) GetStore(id string) (*models.Store, error) {
	var store models.Store
	if err := db.DB.First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *StoreService) UpdateStore(id string, req dtos.UpdateStoreRequest) error {
	var store models.Store
	if err := db.DB.First(&store, id).Error; err != nil {
		return err
	}

	store.Name = req.Name
	store.Description = req.Description

	return db.DB.Save(&store).Error
}

func (s *StoreService) DeleteStore(id string) error {
	return db.DB.Delete(&models.Store{}, id).Error
}

// ImportStoreStock handles bulk import of stock from CSV/XLS/XLSX
func (s *StoreService) ImportStoreStock(file *multipart.FileHeader, userID uint64) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var data [][]string

	if ext == ".csv" {
		reader := csv.NewReader(src)
		data, err = reader.ReadAll()
	} else if ext == ".xlsx" || ext == ".xls" {
		f, err := excelize.OpenReader(src)
		if err != nil {
			return err
		}
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return fmt.Errorf("no sheets found in excel file")
		}
		data, err = f.GetRows(sheets[0])
	} else {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return err
	}

	go s.processStoreStockRowsInBackground(data, userID)
	return nil
}

func (s *StoreService) processStoreStockRowsInBackground(data [][]string, userID uint64) {
	totalRows := len(data) - 1
	if totalRows < 0 {
		return
	}

	importID := uint64(time.Now().UnixNano())
	var wg sync.WaitGroup
	jobs := make(chan []string, totalRows)
	errorChan := make(chan error, totalRows)
	numWorkers := runtime.NumCPU() * 2

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[StoreService] Worker panicked during store stock import: %v", r)
							db.DB.Create(&models.ImportError{
								BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic during import: %v", r),
								ImportId:  importID,
							})
							errorChan <- fmt.Errorf("panic during import: %v", r)
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						// Columns: Store(0), Inventory Item(1), Item Name(2), Unit of Measure(3), Quantity(4), Buying Price(5), Selling Price(6), Credit Selling Price(7)
						if len(row) < 7 {
							return fmt.Errorf("insufficient columns (found %d, need at least 7)", len(row))
						}

						storeName := strings.TrimSpace(row[0])
						itemCode := strings.TrimSpace(row[1])
						itemName := strings.TrimSpace(row[2])
						unitName := strings.TrimSpace(row[3])
						qty, _ := utils.ParseFloat(row[4])
						buyingPrice, _ := utils.ParseFloat(row[5])
						sellingPrice, _ := utils.ParseFloat(row[6])
						creditPrice := sellingPrice
						if len(row) > 7 && strings.TrimSpace(row[7]) != "" {
							creditPrice, _ = utils.ParseFloat(row[7])
						}

						// 1. Resolve Store (Table: stores)
						var store models.Store
						if err := tx.Where("store = ?", storeName).First(&store).Error; err != nil {
							store = models.Store{Name: storeName, BaseModel: models.BaseModel{CreatedBy: userID}}
							if err := tx.Create(&store).Error; err != nil {
								return err
							}
						}

						// 2. Resolve Store Inventory (Table: store_inventories)
						var storeInv models.StoreInventory
						if err := tx.Where("inventory_name = ?", storeName).First(&storeInv).Error; err != nil {
							storeInv = models.StoreInventory{InventoryName: storeName, IsActive: true, BaseModel: models.BaseModel{CreatedBy: userID}}
							if err := tx.Create(&storeInv).Error; err != nil {
								return err
							}
						}

						// 3. Resolve Unit of Measure (Table: store_item_units)
						var unit models.StoreItemUnit
						if err := tx.Where("name = ? OR symbol = ?", unitName, unitName).First(&unit).Error; err != nil {
							unit = models.StoreItemUnit{Name: unitName, Symbol: unitName, BaseModel: models.BaseModel{CreatedBy: userID}}
							if err := tx.Create(&unit).Error; err != nil {
								return err
							}
						}

						// 4. Resolve/Update Store Item (Table: store_items)
						var item models.StoreItem
						if err := tx.Where("sku = ?", itemCode).First(&item).Error; err != nil {
							item = models.StoreItem{
								ItemName:                  itemName,
								SKU:                       itemCode,
								UnitID:                    int64(unit.ID),
								StoreInventoryID:          storeInv.ID,
								DefaultBuyingPrice:        buyingPrice,
								DefaultSellingPrice:       sellingPrice,
								DefaultSellingPriceCredit: creditPrice,
								Status:                    "ACTIVE",
								BaseModel:                 models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							}
							if err := tx.Create(&item).Error; err != nil {
								return err
							}
						} else {
							tx.Model(&item).Updates(map[string]interface{}{"item_name": itemName, "unit_id": unit.ID})
						}

						// 5. Create Master Transaction Record
						transaction := models.Transaction{
							Reference:       fmt.Sprintf("STK-IMP-%d-%s", importID, item.SKU),
							TransactionName: "STORE STOCK IMPORT",
							TransactionType: "STORES",
							TransactionDate: time.Now(),
							Description:     fmt.Sprintf("Bulk stock import for %s in %s", itemName, storeName),
							Status:          "POSTED",
							BaseModel:       models.BaseModel{CreatedBy: userID},
						}
						if err := tx.Create(&transaction).Error; err != nil {
							return err
						}

						// 6. Upsert Store Stock (Table: store_stocks)
						var stock models.StoreStock
						var newQty float64
						if err := tx.Where("item_id = ? AND store_id = ?", item.ID, store.ID).First(&stock).Error; err != nil {
							newQty = qty
							stock = models.StoreStock{
								ItemID: item.ID, StoreID: store.ID, Quantity: qty, Unit: unitName,
								BuyingPrice: buyingPrice, SellingPrice: sellingPrice, CreditSellingPrice: creditPrice,
								BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							}
							if err := tx.Create(&stock).Error; err != nil {
								return err
							}
						} else {
							newQty = stock.Quantity + qty
							if err := tx.Model(&stock).Updates(map[string]interface{}{
								"quantity":             newQty,
								"buying_price":         buyingPrice,
								"selling_price":        sellingPrice,
								"credit_selling_price": creditPrice,
								"updated_by":           userID,
							}).Error; err != nil {
								return err
							}
						}

						// 7. Record Store Transaction (Table: store_transactions)
						storeTx := models.StoreTransaction{
							TransactionID:   transaction.ID,
							StoreID:         store.ID,
							StoreItemID:     item.ID,
							TransactionType: "ADJUSTMENT_IN",
							Quantity:        qty,
							UnitCost:        buyingPrice,
							SellingPrice:    sellingPrice,
							BalanceAfter:    newQty,
							ReferenceType:   "IMPORT",
							ReferenceID:     importID,
							Notes:           "Bulk import",
							BaseModel:       models.BaseModel{CreatedBy: userID},
						}
						if err := tx.Create(&storeTx).Error; err != nil {
							return err
						}

						// 8. Post GL Entries (Accounting)
						glAmount := qty * buyingPrice
						if glAmount > 0 {
							// Use generic STOCK_IMPORT rule for inventory valuation vs opening balances/suspense
							rule := "STOCK_IMPORT"
							// Debit side
							if err := s.postGLEntry(tx, transaction.ID, rule, true, glAmount, "Inventory adjustment via import", time.Now(), userID); err != nil {
								return err
							}
							// Credit side
							if err := s.postGLEntry(tx, transaction.ID, rule, false, glAmount, "Inventory adjustment via import", time.Now(), userID); err != nil {
								return err
							}
						}

						return nil
					})
					if err != nil {
						db.DB.Create(&models.ImportError{
							BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							RowData:   strings.Join(row, ","),
							Error:     err.Error(), ImportId: importID,
						})
						errorChan <- err
					}
				}()
			}
		}()
	}

	for i := 1; i < len(data); i++ {
		jobs <- data[i]
	}
	close(jobs)
	wg.Wait()
	close(errorChan)

	failCount := 0
	for range errorChan {
		failCount++
	}

	message := fmt.Sprintf("Store stock import completed. Success: %d, Failed: %d out of %d records.", totalRows-failCount, failCount, totalRows)
	notificationType := "SUCCESS"
	errorLink := ""
	if failCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/store-stocks/import-errors/%d", importID)
	} else if totalRows == 0 {
		message = "Store stock import completed. No records were processed."
		notificationType = "SUCCESS"
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title: "Store Stock Import Status", Message: message, NotificationType: notificationType, ErrorLink: errorLink,
	})
}

func (s *StoreService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

func (s *StoreService) postGLEntry(
	tx *gorm.DB,
	txID uint64,
	ruleType string,
	isDebit bool,
	amount float64,
	desc string,
	transDate time.Time,
	userID uint64,
) error {
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
