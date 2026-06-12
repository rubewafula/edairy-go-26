package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"os"
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

type SupplierService struct {
	notificationService *UINotificationService
}

func NewSupplierService() *SupplierService {
	return &SupplierService{
		notificationService: NewUINotificationService(),
	}
}

func (s *SupplierService) CreateSupplier(req dtos.CreateSupplierRequest, userID uint64) (*models.Supplier, error) {
	supplier := &models.Supplier{
		BaseModel:          models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		SupplierCategoryID: req.SupplierCategoryID,
		SupplierCode:       req.SupplierCode,
		SupplierType:       req.SupplierType,
		CompanyName:        req.CompanyName,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		PhoneNo:            req.PhoneNo,
		EmailAddress:       req.EmailAddress,
		KraPin:             req.KraPin,
		CreditLimit:        req.CreditLimit,
		PaymentTermsDays:   req.PaymentTermsDays,
		Status:             req.Status,
		Notes:              req.Notes,
	}

	if req.Gender != "" {
		supplier.Gender = &req.Gender
	}

	if err := db.DB.Create(supplier).Error; err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) GetSuppliers(page, limit int) ([]dtos.SupplierResponse, int64, error) {
	var results []dtos.SupplierResponse
	var total int64
	db.DB.Model(&models.Supplier{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			s.id, s.supplier_code, s.supplier_type, s.company_name, 
			CONCAT(COALESCE(s.first_name,''), ' ', COALESCE(s.last_name,'')) as full_name,
			sc.category_name, s.email_address, s.phone_no, s.current_balance, s.status, s.created_at
		FROM suppliers s
		LEFT JOIN supplier_categories sc ON s.supplier_category_id = sc.id
		WHERE s.deleted_at IS NULL
		ORDER BY s.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplierService) GetSupplier(id string) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := db.DB.First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (s *SupplierService) UpdateSupplier(id string, req dtos.CreateSupplierRequest, userID uint64) error {
	var supplier models.Supplier
	if err := db.DB.First(&supplier, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"supplier_category_id": req.SupplierCategoryID,
		"supplier_type":        req.SupplierType,
		"company_name":         req.CompanyName,
		"first_name":           req.FirstName,
		"last_name":            req.LastName,
		"phone_no":             req.PhoneNo,
		"email_address":        req.EmailAddress,
		"kra_pin":              req.KraPin,
		"credit_limit":         req.CreditLimit,
		"payment_terms_days":   req.PaymentTermsDays,
		"status":               req.Status,
		"notes":                req.Notes,
		"updated_by":           userID,
	}
	return db.DB.Model(&supplier).Updates(updates).Error
}

func (s *SupplierService) DeleteSupplier(id string) error {
	return db.DB.Delete(&models.Supplier{}, id).Error
}

func (s *SupplierService) CreateContact(req dtos.CreateSupplierContactRequest, userID uint64) (*models.SupplierContact, error) {
	contact := &models.SupplierContact{
		BaseModel:          models.BaseModel{CreatedBy: userID},
		SupplierID:         req.SupplierID,
		ContactType:        req.ContactType,
		FullName:           req.FullName,
		Designation:        req.Designation,
		PhoneNo:            req.PhoneNo,
		AlternativePhoneNo: req.AlternativePhoneNo,
		EmailAddress:       req.EmailAddress,
		IsDefault:          req.IsDefault,
		Notes:              req.Notes,
	}
	err := db.DB.Create(contact).Error
	return contact, err
}

func (s *SupplierService) CreateBankAccount(req dtos.CreateSupplierBankAccountRequest, userID uint64) (*models.SupplierBankAccount, error) {
	account := &models.SupplierBankAccount{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		SupplierID:    req.SupplierID,
		BankID:        req.BankID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		AccountType:   req.AccountType,
		CurrencyCode:  req.CurrencyCode,
		IsDefault:     req.IsDefault,
	}
	err := db.DB.Create(account).Error
	return account, err
}

func (s *SupplierService) GetSupplierContacts(supplierID string) ([]dtos.SupplierContactResponse, error) {
	var contacts []dtos.SupplierContactResponse
	err := db.DB.Model(&models.SupplierContact{}).
		Where("supplier_id = ?", supplierID).
		Find(&contacts).Error
	return contacts, err
}

func (s *SupplierService) GetSupplierBankAccounts(supplierID string) ([]dtos.SupplierBankAccountResponse, error) {
	var results []dtos.SupplierBankAccountResponse
	query := `
		SELECT sba.*, b.name as bank_name 
		FROM supplier_bank_accounts sba LEFT JOIN banks b ON sba.bank_id = b.id
		WHERE sba.supplier_id = ? AND sba.deleted_at IS NULL`
	err := db.DB.Raw(query, supplierID).Scan(&results).Error
	return results, err
}

// ExportSuppliers initiates a background process to export supplier data to a CSV file.
func (s *SupplierService) ExportSuppliers(userID uint64, categoryID, supplierType, status string) error {
	go s.processSupplierExportInBackground(userID, categoryID, supplierType, status)
	return nil
}

// supplierExportQueryResult is a helper struct to hold the results of the supplier export query,
// including joined fields from related tables.
type supplierExportQueryResult struct {
	models.Supplier
	CategoryCode  string `gorm:"column:category_code"`
	CategoryName  string `gorm:"column:category_name"`
	BankName      string `gorm:"column:bank_name"`
	BranchName    string `gorm:"column:branch_name"`
	AccountName   string `gorm:"column:account_name"`
	AccountNumber string `gorm:"column:account_number"`
}

// processSupplierExportInBackground performs the actual data export in a separate goroutine.
// It queries the database, generates a CSV file, saves it, and sends a UI notification.
func (s *SupplierService) processSupplierExportInBackground(userID uint64, categoryID, supplierType, status string) {
	var results []supplierExportQueryResult

	// Base SQL query to fetch supplier details along with their category and primary bank account information.
	baseQuery := `
		SELECT
			s.*, 
			sc.category_code, sc.category_name,
			b.bank_name, bb.branch_name, sba.account_name, sba.account_number
		FROM suppliers s
		LEFT JOIN supplier_categories sc ON s.supplier_category_id = sc.id
		LEFT JOIN supplier_bank_accounts sba ON s.id = sba.supplier_id AND sba.is_default = 'yes' AND sba.deleted_at IS NULL
		LEFT JOIN banks b ON sba.bank_id = b.id
		LEFT JOIN bank_branches bb ON sba.branch_id = bb.id
	`

	var args []interface{}
	whereClauses := []string{"s.deleted_at IS NULL"}

	if categoryID != "" {
		whereClauses = append(whereClauses, "s.supplier_category_id = ?")
		args = append(args, categoryID)
	}
	if supplierType != "" {
		whereClauses = append(whereClauses, "s.supplier_type = ?")
		args = append(args, supplierType)
	}
	if status != "" {
		whereClauses = append(whereClauses, "s.status = ?")
		args = append(args, status)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	// Execute the query to retrieve all matching supplier records.
	err := db.DB.Raw(baseQuery+whereSql+" ORDER BY s.id DESC", args...).Scan(&results).Error
	if err != nil {
		log.Printf("[SupplierService.processSupplierExportInBackground] Error querying suppliers for export: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Supplier Export Failed",
			Message:          fmt.Sprintf("Failed to export supplier data: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	fileData, err := s.generateSupplierCSV(results)
	if err != nil {
		log.Printf("[SupplierService.processSupplierExportInBackground] Error generating CSV: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Supplier Export Failed",
			Message:          fmt.Sprintf("Failed to generate export file: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	exportDir := "./storage/exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Printf("[SupplierService.processSupplierExportInBackground] Failed to create export directory: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Supplier Export Failed",
			Message:          fmt.Sprintf("Failed to create export directory: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	filename := fmt.Sprintf("suppliers_export_%d.csv", time.Now().UnixNano())
	filePath := filepath.Join(exportDir, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[SupplierService.processSupplierExportInBackground] Failed to save exported CSV: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Supplier Export Failed",
			Message:          fmt.Sprintf("Failed to save exported file: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	downloadLink := fmt.Sprintf("/api/suppliers/export/download/%s", filename)
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Supplier Export Complete",
		Message:          "Your supplier data export is ready for download.",
		NotificationType: "SUCCESS",
		DownloadLink:     downloadLink,
	})
	log.Printf("[SupplierService.processSupplierExportInBackground] Supplier export completed successfully. File: %s", filename)
}

func (s *SupplierService) generateSupplierCSV(results []supplierExportQueryResult) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Define the CSV header, matching the import columns for consistency.
	headers := []string{
		"category_code", "category_name", "supplier_type", "company_name", "first_name",
		"last_name", "kra_pin", "license_number", "phone_no", "email_address",
		"contact_person", "postal_address", "postal_code", "town", "bank_name",
		"bank_branch", "account_name", "account_number",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Iterate through the query results and write each supplier's data to the CSV.
	for _, supplier := range results {
		row := []string{
			supplier.CategoryCode,
			supplier.CategoryName,
			supplier.SupplierType,
			supplier.CompanyName,
			supplier.FirstName,
			supplier.LastName,
			supplier.KraPin,
			supplier.LicenseNumber,
			supplier.PhoneNo,
			supplier.EmailAddress,
			supplier.ContactPerson,
			supplier.PostalAddress,
			supplier.PostalCode,
			supplier.Town,
			supplier.BankName,
			supplier.BranchName,
			supplier.AccountName,
			supplier.AccountNumber,
		}
		if err := writer.Write(row); err != nil {
			log.Printf("[SupplierService.generateSupplierCSV] Error writing CSV row for supplier %s: %v", supplier.SupplierCode, err)
		}
	}
	writer.Flush()

	return buf.Bytes(), writer.Error()
}

// ImportSuppliers bulk imports suppliers from CSV, XLS, or XLSX files.
func (s *SupplierService) ImportSuppliers(file *multipart.FileHeader, userID uint64) error {
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

	go s.processSupplierRowsInBackground(data, userID)

	return nil
}

func (s *SupplierService) processSupplierRowsInBackground(data [][]string, userID uint64) {
	totalRows := len(data) - 1
	if totalRows < 0 {
		return
	}

	importID := uint64(time.Now().UnixNano())
	var wg sync.WaitGroup
	jobs := make(chan []string, totalRows)
	errorChan := make(chan error, totalRows)

	numWorkers := runtime.NumCPU() * 2
	if numWorkers < 1 {
		numWorkers = 1
	}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							db.DB.Create(&models.ImportError{
								BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic during import: %v", r),
								ImportId:  importID,
							})
							errorChan <- fmt.Errorf("panic during row processing")
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						// Column mapping: [category_code(0), category_name(1), supplier_type(2), company_name(3), first_name(4), last_name(5), kra_pin(6), license_number(7), phone_no(8), email_address(9), contact_person(10), postal_address(11), postal_code(12), town(13), bank_name(14), bank_branch(15), account_name(16), account_number(17)]
						if len(row) < 18 {
							return fmt.Errorf("row has insufficient columns (found %d, need 18)", len(row))
						}

						// 1. Resolve/Create Category
						catCode := strings.TrimSpace(row[0])
						catName := strings.TrimSpace(row[1])
						var category models.SupplierCategory
						if err := tx.Where("category_code = ? OR category_name = ?", catCode, catName).First(&category).Error; err != nil {
							category = models.SupplierCategory{
								BaseModel:    models.BaseModel{CreatedBy: userID},
								CategoryCode: catCode,
								CategoryName: catName,
								Status:       "active",
							}
							if err := tx.Create(&category).Error; err != nil {
								return fmt.Errorf("failed to create category: %w", err)
							}
						}

						// 2. Check for existing supplier (Idempotency)
						email := strings.TrimSpace(row[9])
						kraPin := strings.TrimSpace(row[6])
						var count int64
						tx.Model(&models.Supplier{}).Where("email_address = ? OR kra_pin = ?", email, kraPin).Count(&count)
						if count > 0 {
							return fmt.Errorf("supplier with email %s or KRA PIN %s already exists", email, kraPin)
						}

						// 3. Create Supplier
						supplier := models.Supplier{
							BaseModel:          models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							SupplierCategoryID: category.ID,
							SupplierCode:       utils.GenerateMemberNo(),
							SupplierType:       strings.ToLower(strings.TrimSpace(row[2])),
							CompanyName:        strings.TrimSpace(row[3]),
							FirstName:          strings.TrimSpace(row[4]),
							LastName:           strings.TrimSpace(row[5]),
							KraPin:             kraPin,
							LicenseNumber:      strings.TrimSpace(row[7]),
							PhoneNo:            utils.NormalizePhone(row[8]),
							EmailAddress:       email,
							ContactPerson:      strings.TrimSpace(row[10]),
							PostalAddress:      strings.TrimSpace(row[11]),
							PostalCode:         strings.TrimSpace(row[12]),
							Town:               strings.TrimSpace(row[13]),
							Status:             "active",
						}

						if err := tx.Create(&supplier).Error; err != nil {
							return err
						}

						// 4. Handle Bank Account
						bankNameInput := strings.TrimSpace(row[14])
						branchNameInput := strings.TrimSpace(row[15])
						accountNameInput := strings.TrimSpace(row[16])
						accNo := strings.TrimSpace(row[17])

						if bankNameInput != "" && accNo != "" {
							var bank models.Bank
							if err := tx.Where("bank_name = ?", bankNameInput).First(&bank).Error; err != nil {
								bank = models.Bank{
									BankName: bankNameInput,
								}
								if err := tx.Create(&bank).Error; err != nil {
									return fmt.Errorf("failed to create bank: %w", err)
								}
							}

							var branch models.BankBranch
							if branchNameInput != "" {
								if err := tx.Where("bank_id = ? AND name = ?", bank.ID, branchNameInput).First(&branch).Error; err != nil {
									branch = models.BankBranch{
										BankID: bank.ID,
										Name:   branchNameInput,
									}
									if err := tx.Create(&branch).Error; err != nil {
										return fmt.Errorf("failed to create bank branch: %w", err)
									}
								}
							}

							bankAcc := models.SupplierBankAccount{
								BaseModel:     models.BaseModel{CreatedBy: userID},
								SupplierID:    supplier.ID,
								BankID:        bank.ID,
								BankBranchID:  branch.ID,
								AccountNumber: accNo,
								AccountName:   accountNameInput,
								IsDefault:     "yes",
								Status:        "active",
							}
							if err := tx.Create(&bankAcc).Error; err != nil {
								return fmt.Errorf("failed to create supplier bank account: %w", err)
							}
						}

						return nil
					})

					if err != nil {
						db.DB.Create(&models.ImportError{
							BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							RowData:   strings.Join(row, ","),
							Error:     err.Error(),
							ImportId:  importID,
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

	failedCount := 0
	for err := range errorChan {
		if err != nil {
			failedCount++
		}
	}

	notificationType := "SUCCESS"
	if failedCount > 0 {
		notificationType = "ERROR"
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Supplier Import Status",
		Message:          fmt.Sprintf("Import completed. Success: %d, Failed: %d out of %d records.", totalRows-failedCount, failedCount, totalRows),
		NotificationType: notificationType,
		ErrorLink:        fmt.Sprintf("/suppliers/import-errors/%d", importID),
	})
}

// GetImportErrors retrieves the list of errors encountered during a specific import.
func (s *SupplierService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}
