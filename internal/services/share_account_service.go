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

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ShareAccountService struct {
	notificationService *UINotificationService
}

func NewShareAccountService() *ShareAccountService {
	return &ShareAccountService{
		notificationService: NewUINotificationService(),
	}
}

func (s *ShareAccountService) CreateAccount(req dtos.CreateShareAccountRequest) (*models.ShareAccount, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	openedAt := time.Now()
	if req.OpenedAt != "" {
		openedAt = utils.ParseDate(req.OpenedAt)
	}

	account := &models.ShareAccount{
		MemberID:    req.MemberID,
		ShareTypeID: req.ShareTypeID,
		Status:      status,
		OpenedAt:    openedAt,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *ShareAccountService) GetShareAccounts(memberID, shareTypeID string) ([]dtos.ShareAccountResponse, int64, error) {
	var results []dtos.ShareAccountResponse
	var total int64

	countQuery := db.DB.Model(&models.ShareAccount{})
	if memberID != "" {
		countQuery = countQuery.Where("member_id = ?", memberID)
	}
	if shareTypeID != "" {
		countQuery = countQuery.Where("share_type_id = ?", shareTypeID)
	}
	countQuery.Count(&total)

	whereClause := "sa.deleted_at IS NULL"
	var args []interface{}
	if memberID != "" {
		whereClause += " AND sa.member_id = ?"
		args = append(args, memberID)
	}
	if shareTypeID != "" {
		whereClause += " AND sa.share_type_id = ?"
		args = append(args, shareTypeID)
	}

	query := fmt.Sprintf(`
		SELECT 
			sa.id, sa.member_id, m.member_no, m.first_name, m.last_name,
			sa.share_type_id, st.share_code, st.share_type AS share_type_name, st.description,
			sa.status, sa.share_units, sa.share_amount, sa.opened_at, sa.created_at, sa.updated_at
		FROM share_accounts sa
		LEFT JOIN member_registrations m ON sa.member_id = m.id
		LEFT JOIN share_types st ON sa.share_type_id = st.id
		WHERE %s
	`, whereClause)

	err := db.DB.Raw(query, args...).Scan(&results).Error
	return results, total, err
}

func (s *ShareAccountService) GetShareAccount(id string) (*dtos.ShareAccountResponse, error) {
	var result dtos.ShareAccountResponse
	query := `
		SELECT 
			sa.id, sa.member_id, m.member_no, m.first_name, m.last_name,
			sa.share_type_id, st.share_code, st.share_type AS share_type_name, st.description,
			sa.status, sa.share_units, sa.share_amount, sa.opened_at, sa.created_at, sa.updated_at
		FROM share_accounts sa
		LEFT JOIN member_registrations m ON sa.member_id = m.id
		LEFT JOIN share_types st ON sa.share_type_id = st.id
		WHERE sa.id = ? AND sa.deleted_at IS NULL
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

func (s *ShareAccountService) UpdateAccount(id string, req dtos.UpdateShareAccountRequest) error {
	var account models.ShareAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	account.MemberID = req.MemberID
	account.ShareTypeID = req.ShareTypeID
	account.Status = req.Status
	if req.OpenedAt != "" {
		account.OpenedAt = utils.ParseDate(req.OpenedAt)
	}

	return db.DB.Save(&account).Error
}

func (s *ShareAccountService) DeleteAccount(id string) error {
	return db.DB.Delete(&models.ShareAccount{}, id).Error
}

// ImportAccounts bulk imports share accounts from CSV, XLS, or XLSX files.
func (s *ShareAccountService) ImportAccounts(file *multipart.FileHeader, userID uint64) error {
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
		if err != nil {
			return err
		}
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
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	go s.processShareAccountRowsInBackground(data, userID)

	return nil
}

func (s *ShareAccountService) processShareAccountRowsInBackground(data [][]string, userID uint64) {
	totalRows := len(data) - 1
	if totalRows < 0 {
		return
	}

	log.Printf("[ShareAccountService] Starting background import for %d rows.", totalRows)
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
						// Column mapping: [Member Number(0), share_type(1), share_code(2), share_units(3), share_amount(4), status(5)]
						if len(row) < 5 {
							return fmt.Errorf("row has insufficient columns (found %d, need at least 5)", len(row))
						}

						memberNo := strings.TrimSpace(row[0])
						var member models.Member
						if err := tx.Where("member_no = ?", memberNo).First(&member).Error; err != nil {
							return fmt.Errorf("member with number %s not found", memberNo)
						}

						sTypeName := strings.TrimSpace(row[1])
						sCode := strings.TrimSpace(row[2])
						units, _ := utils.ParseFloat(row[3])
						amount, _ := utils.ParseFloat(row[4])

						var st models.ShareType
						if err := tx.Where("share_code = ? OR share_type = ?", sCode, sTypeName).First(&st).Error; err != nil {
							// If share type does not exist, create one based on share_code and share_type
							shareValue := 0.0
							if units > 0 {
								shareValue = amount / units
							}

							st = models.ShareType{
								ShareCode:     sCode,
								ShareType:     sTypeName,
								Description:   "Auto-created during import",
								ShareValue:    shareValue,
								HasShareValue: 1,
								Mandatory:     0,
							}

							if err := tx.Create(&st).Error; err != nil {
								return fmt.Errorf("failed to auto-create share type '%s': %w", sTypeName, err)
							}
						}

						status := "ACTIVE"
						if len(row) > 5 && strings.TrimSpace(row[5]) != "" {
							status = strings.ToUpper(strings.TrimSpace(row[5]))
						}

						var account models.ShareAccount
						if err := tx.Where("member_id = ? AND share_type_id = ?", member.ID, st.ID).First(&account).Error; err == nil {
							return fmt.Errorf("share account for member %s and share type %s already exists", memberNo, st.ShareType)
						}

						account = models.ShareAccount{
							MemberID:    member.ID,
							ShareTypeID: st.ID,
							ShareUnits:  units,
							ShareAmount: amount,
							Status:      status,
							OpenedAt:    time.Now(),
						}

						if err := tx.Create(&account).Error; err != nil {
							return err
						}

						// 2. Create Master Transaction Record
						transaction := models.Transaction{
							BaseModel:       models.BaseModel{CreatedBy: userID},
							Reference:       fmt.Sprintf("SHR-IMP-%d-%s", importID, member.MemberNo),
							TransactionName: "SHARE ACCOUNT IMPORT",
							TransactionType: "SHARE",
							TransactionDate: time.Now(),
							Description:     fmt.Sprintf("Bulk share account import for member %s", member.MemberNo),
							Status:          "POSTED",
						}
						if err := tx.Create(&transaction).Error; err != nil {
							return err
						}

						// 3. Record Share Movement (IMPORT)
						shareTx := models.ShareTransaction{
							TransactionID:   transaction.ID,
							ShareAccountID:  account.ID,
							MemberID:        member.ID,
							TransactionType: "TRANSFER_IN",
							ShareUnits:      units,
							UnitPrice:       st.ShareValue,
							Debit:           amount,
							Credit:          0,
							BalanceAfter:    amount,
							TransactionDate: time.Now(),
						}
						if err := tx.Create(&shareTx).Error; err != nil {
							return err
						}

						// 4. Post GL Entries (Accounting)
						if amount > 0 {
							rule := "SHARES_IMPORT"
							desc := fmt.Sprintf("Opening balance for share account: %s", st.ShareType)
							if err := s.postGLEntry(tx, transaction.ID, rule, true, amount, desc, time.Now(), userID); err != nil {
								return err
							}
							if err := s.postGLEntry(tx, transaction.ID, rule, false, amount, desc, time.Now(), userID); err != nil {
								return err
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
		Title:            "Share Account Import Status",
		Message:          fmt.Sprintf("Import completed. Success: %d, Failed: %d out of %d records.", totalRows-failedCount, failedCount, totalRows),
		NotificationType: notificationType,
		ErrorLink:        fmt.Sprintf("/share-accounts/import-errors/%d", importID),
	})
}

func (s *ShareAccountService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

// postGLEntry creates a single ledger entry based on a rule type found in transaction_posting_rules.
func (s *ShareAccountService) postGLEntry(
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

// ExportAccounts initiates a background process to export share accounts.
func (s *ShareAccountService) ExportAccounts(userID uint64, memberID, shareTypeID, status, reportType string) error {
	go s.processShareAccountExportInBackground(userID, memberID, shareTypeID, status, reportType)
	return nil
}

type shareAccountExportQueryResult struct {
	MemberNo      string  `gorm:"column:member_no"`
	FirstName     string  `gorm:"column:first_name"`
	LastName      string  `gorm:"column:last_name"`
	ShareTypeName string  `gorm:"column:share_type_name"`
	ShareCode     string  `gorm:"column:share_code"`
	ShareUnits    float64 `gorm:"column:share_units"`
	ShareAmount   float64 `gorm:"column:share_amount"`
	Status        string  `gorm:"column:status"`
}

func (s *ShareAccountService) processShareAccountExportInBackground(userID uint64, memberID, shareTypeID, status, reportType string) {
	var results []shareAccountExportQueryResult

	whereClause := "sa.deleted_at IS NULL"
	var args []interface{}
	if memberID != "" {
		whereClause += " AND sa.member_id = ?"
		args = append(args, memberID)
	}
	if shareTypeID != "" {
		whereClause += " AND sa.share_type_id = ?"
		args = append(args, shareTypeID)
	}
	if status != "" {
		whereClause += " AND sa.status = ?"
		args = append(args, status)
	}

	query := fmt.Sprintf(`
		SELECT 
			m.member_no, m.first_name, m.last_name,
			st.share_type AS share_type_name, st.share_code,
			sa.share_units, sa.share_amount, sa.status
		FROM share_accounts sa
		LEFT JOIN member_registrations m ON sa.member_id = m.id
		LEFT JOIN share_types st ON sa.share_type_id = st.id
		WHERE %s
		ORDER BY sa.id DESC
	`, whereClause)

	if err := db.DB.Raw(query, args...).Scan(&results).Error; err != nil {
		log.Printf("[ShareAccountService] Export query error: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title: "Export Failed", Message: "Failed to query share accounts.", NotificationType: "ERROR",
		})
		return
	}

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateShareAccountPDF(results, memberID, shareTypeID, status)
	} else {
		fileData, err = s.generateShareAccountCSV(results)
	}

	if err != nil {
		log.Printf("[ShareAccountService] Export generation error: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title: "Export Failed", Message: "Failed to generate file.", NotificationType: "ERROR",
		})
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("share_accounts_export_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[ShareAccountService] File write error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Share Account Export Ready",
		Message:          "Your share accounts export has been generated successfully.",
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/share-accounts/export/download/%s", filename),
	})
}

func (s *ShareAccountService) generateShareAccountCSV(results []shareAccountExportQueryResult) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"Member Number", "share_type", "share_code", "share_units", "share_amount", "status"})
	for _, r := range results {
		writer.Write([]string{r.MemberNo, r.ShareTypeName, r.ShareCode, fmt.Sprintf("%.2f", r.ShareUnits), fmt.Sprintf("%.2f", r.ShareAmount), r.Status})
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *ShareAccountService) generateShareAccountPDF(results []shareAccountExportQueryResult, memberID, shareTypeID, status string) ([]byte, error) {
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
	pdf.CellFormat(0, 10, "SHARE ACCOUNT REGISTER", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	memberVal := "ALL"
	if memberID != "" {
		db.DB.Table("member_registrations").Where("id = ?", memberID).Pluck("member_no", &memberVal)
	}
	statusVal := "ALL"
	if status != "" {
		statusVal = status
	}

	pdf.CellFormat(0, 8, fmt.Sprintf("Showing (%d) Records for; Member %s, Status %s [filters applied]", len(results), memberVal, statusVal), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	headers := []string{"Member No", "Full Name", "Share Type", "Code", "Units", "Amount", "Status"}
	widths := []float64{30, 70, 50, 25, 30, 40, 30}
	pdf.SetFont("Arial", "B", 8)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, r := range results {
		pdf.CellFormat(30, 8, r.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 8, fmt.Sprintf("%s %s", r.FirstName, r.LastName), "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 8, r.ShareTypeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, r.ShareCode, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", r.ShareUnits), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%.2f", r.ShareAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, r.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
