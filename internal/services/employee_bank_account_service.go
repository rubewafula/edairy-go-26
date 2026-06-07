package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type EmployeeBankAccountService struct {
	notificationService *UINotificationService
}

func NewEmployeeBankAccountService() *EmployeeBankAccountService {
	return &EmployeeBankAccountService{
		notificationService: NewUINotificationService(),
	}
}

func (s *EmployeeBankAccountService) CreateAccount(req dtos.CreateEmployeeBankAccountRequest, userID uint64) (*models.EmployeeBankAccount, error) {
	account := &models.EmployeeBankAccount{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		EmployeeID:    req.EmployeeID,
		BankID:        req.BankID,
		BranchID:      req.BranchID,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *EmployeeBankAccountService) GetAccounts(employeeID string, page, limit int) ([]dtos.EmployeeBankAccountResponse, int64, error) {
	var results []dtos.EmployeeBankAccountResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeBankAccount{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			eba.id, eba.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eba.bank_id, b.bank_name,
			eba.account_number, eba.account_name, eba.created_at, eba.updated_at,
			eba.created_by, eba.updated_by
		FROM employee_bank_accounts eba
		LEFT JOIN banks b ON eba.bank_id = b.id
		LEFT JOIN employees e ON eba.employee_id = e.id
		WHERE eba.deleted_at IS NULL AND (? = '' OR eba.employee_id = ?)
		ORDER BY eba.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeBankAccountService) GetAccount(id string) (*dtos.EmployeeBankAccountResponse, error) {
	var result dtos.EmployeeBankAccountResponse
	query := `
		SELECT 
			eba.id, eba.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eba.bank_id, b.bank_name,
			eba.account_number, eba.account_name, eba.created_at, eba.updated_at
		FROM employee_bank_accounts eba
		LEFT JOIN banks b ON eba.bank_id = b.id
		LEFT JOIN employees e ON eba.employee_id = e.id
		WHERE eba.id = ? AND eba.deleted_at IS NULL
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

func (s *EmployeeBankAccountService) UpdateAccount(id string, req dtos.UpdateEmployeeBankAccountRequest, userID uint64) error {
	var account models.EmployeeBankAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"bank_id":        req.BankID,
		"branch_id":      req.BranchID,
		"account_number": req.AccountNumber,
		"account_name":   req.AccountName,
		"updated_by":     userID,
	}

	return db.DB.Model(&account).Updates(updates).Error
}

func (s *EmployeeBankAccountService) DeleteAccount(id string, userID uint64) error {
	// Audit the update before soft delete
	return db.DB.Model(&models.EmployeeBankAccount{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.EmployeeBankAccount{}).Error
}

// ImportAccounts bulk imports employee bank accounts from CSV, XLS, or XLSX files.
func (s *EmployeeBankAccountService) ImportAccounts(file *multipart.FileHeader, userID uint64) error {
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

	go s.processImportRowsInBackground(data, userID)

	return nil
}

func (s *EmployeeBankAccountService) processImportRowsInBackground(data [][]string, userID uint64) {
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
				// Column mapping: [Employee NO(0), Employee Names(1), Bank(2), Branch(3), Account Name(4), Account Number(5)]
				for len(row) < 6 {
					row = append(row, "")
				}

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
						empNo := strings.TrimSpace(row[0])
						if empNo == "" {
							return fmt.Errorf("employee number is required")
						}

						var employee models.Employee
						if err := tx.Where("employee_no = ?", empNo).First(&employee).Error; err != nil {
							return fmt.Errorf("employee with number %s not found", empNo)
						}

						bankNameInput := strings.TrimSpace(row[2])
						branchNameInput := strings.TrimSpace(row[3])
						accountNameInput := strings.TrimSpace(row[4])
						accNo := strings.TrimSpace(row[5])

						if bankNameInput == "" || accNo == "" {
							return fmt.Errorf("bank name and account number are required")
						}

						var bank models.Bank
						if err := tx.Where("bank_name = ?", bankNameInput).First(&bank).Error; err != nil {
							bank = models.Bank{BankName: bankNameInput}
							if err := tx.Create(&bank).Error; err != nil {
								return fmt.Errorf("failed to create bank: %w", err)
							}
						}

						var branch models.BankBranch
						if branchNameInput != "" {
							if err := tx.Where("bank_id = ? AND name = ?", bank.ID, branchNameInput).First(&branch).Error; err != nil {
								branch = models.BankBranch{BankID: bank.ID, Name: branchNameInput}
								if err := tx.Create(&branch).Error; err != nil {
									return fmt.Errorf("failed to create bank branch: %w", err)
								}
							}
						}

						// Idempotency: Create or Update existing primary account
						var bankAcc models.EmployeeBankAccount
						res := tx.Where("employee_id = ?", employee.ID).First(&bankAcc)

						bankAcc.EmployeeID = employee.ID
						bankAcc.BankID = bank.ID
						bankAcc.BranchID = branch.ID
						bankAcc.AccountNumber = accNo
						bankAcc.AccountName = accountNameInput
						bankAcc.UpdatedBy = userID

						if errors.Is(res.Error, gorm.ErrRecordNotFound) {
							bankAcc.BaseModel = models.BaseModel{CreatedBy: userID}
							return tx.Create(&bankAcc).Error
						}
						return tx.Save(&bankAcc).Error
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
	for range errorChan {
		failedCount++
	}

	notificationType := "SUCCESS"
	if failedCount > 0 {
		notificationType = "ERROR"
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Employee Bank Accounts Import Status",
		Message:          fmt.Sprintf("Import completed. Success: %d, Failed: %d out of %d records.", totalRows-failedCount, failedCount, totalRows),
		NotificationType: notificationType,
		ErrorLink:        fmt.Sprintf("/employee-bank-accounts/import-errors/%d", importID),
	})
}

// ExportAccounts initiates a background process to export employee bank accounts.
func (s *EmployeeBankAccountService) ExportAccounts(userID uint64, format string) error {
	go s.processExportInBackground(userID, format)
	return nil
}

// GetImportErrors retrieves the list of errors encountered during a specific import session.
func (s *EmployeeBankAccountService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

type employeeBankAccountExportResult struct {
	EmployeeNo    string `gorm:"column:employee_no"`
	EmployeeNames string `gorm:"column:employee_names"`
	Bank          string `gorm:"column:bank"`
	Branch        string `gorm:"column:branch"`
	AccountName   string `gorm:"column:account_name"`
	AccountNo     string `gorm:"column:account_no"`
}

func (s *EmployeeBankAccountService) processExportInBackground(userID uint64, format string) {
	var results []employeeBankAccountExportResult

	// Unified query to fetch all requested columns
	query := `
		SELECT 
			e.employee_no, 
			CONCAT(e.first_name, ' ', COALESCE(e.middle_name, ''), ' ', e.surname) as employee_names,
			b.bank_name as bank,
			COALESCE(bb.branch_name, '') as branch,
			eba.account_name,
			eba.account_number as account_no
		FROM employee_bank_accounts eba
		LEFT JOIN employees e ON eba.employee_id = e.id
		LEFT JOIN banks b ON eba.bank_id = b.id
		LEFT JOIN bank_branches bb ON eba.branch_id = bb.id
		WHERE eba.deleted_at IS NULL
		ORDER BY eba.id DESC
	`

	if err := db.DB.Raw(query).Scan(&results).Error; err != nil {
		log.Printf("[EmployeeBankAccountService] Export query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"

	if strings.ToLower(format) == "pdf" {
		ext = "pdf"
		fileData, err = s.generatePDF(results)
	} else {
		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)
		writer.Write([]string{"Employee NO", "Employee Names", "Bank", "Branch", "account_name", "account_no"})

		for _, r := range results {
			writer.Write([]string{r.EmployeeNo, r.EmployeeNames, r.Bank, r.Branch, r.AccountName, r.AccountNo})
		}
		writer.Flush()
		fileData = buf.Bytes()
		err = writer.Error()
	}

	if err != nil {
		log.Printf("[EmployeeBankAccountService] Error generating export: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("employee_bank_accounts_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[EmployeeBankAccountService] Error saving export file: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            fmt.Sprintf("Employee Bank Accounts Export (%s) Ready", strings.ToUpper(ext)),
		Message:          fmt.Sprintf("The employee bank accounts %s export is ready for download.", ext),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/employee-bank-accounts/export/download/%s", filename),
	})
}

func (s *EmployeeBankAccountService) generatePDF(results []employeeBankAccountExportResult) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
		Phone          string `gorm:"column:phone"`
		Email          string `gorm:"column:email"`
	}
	db.DB.Table("organization_details").First(&org)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Phone: %s | Email: %s", org.Phone, org.Email), "", 1, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "EMPLOYEE BANK ACCOUNTS REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"Emp-No", "Employee Names", "Bank", "Branch", "Account Name", "Account No"}
	widths := []float64{25, 70, 45, 45, 50, 40}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, r := range results {
		pdf.CellFormat(widths[0], 8, r.EmployeeNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, r.EmployeeNames, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, r.Bank, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, r.Branch, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 8, r.AccountName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, r.AccountNo, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
