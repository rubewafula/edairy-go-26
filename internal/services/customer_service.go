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

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type CustomerService struct {
	notificationService *UINotificationService
}

func NewCustomerService() *CustomerService {
	return &CustomerService{
		notificationService: NewUINotificationService(),
	}
}

func (s *CustomerService) CreateCustomer(req dtos.CreateCustomerRequest, userID uint64) (*models.Customer, error) {
	customer := &models.Customer{
		BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		CustomerTypeID: req.CustomerTypeID,
		FullNames:      req.FullNames,
		CustomerNo:     req.CustomerNo,
		Phone:          utils.NormalizePhone(req.Phone),
		KraPin:         req.KraPin,
		PostalAddress:  req.PostalAddress,
		PostalCode:     req.PostalCode,
		PostalTown:     req.PostalTown,
		Rate:           req.Rate,
		Status:         req.Status,
	}

	if customer.Status == "" {
		customer.Status = "ACTIVE"
	}

	if err := db.DB.Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetCustomers(page, limit int) ([]dtos.CustomerResponse, int64, error) {
	var results []dtos.CustomerResponse
	var total int64
	db.DB.Model(&models.Customer{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT c.*, ct.name as type_name
		FROM customers c
		LEFT JOIN customer_types ct ON c.customer_type_id = ct.id
		WHERE c.deleted_at IS NULL
		ORDER BY c.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerService) GetCustomer(id string) (*models.Customer, error) {
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) UpdateCustomer(id string, req dtos.UpdateCustomerRequest, userID uint64) error {
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&customer).Updates(map[string]interface{}{
		"customer_type_id": req.CustomerTypeID,
		"full_names":       req.FullNames,
		"phone":            utils.NormalizePhone(req.Phone),
		"kra_pin":          req.KraPin,
		"postal_address":   req.PostalAddress,
		"postal_code":      req.PostalCode,
		"postal_town":      req.PostalTown,
		"rate":             req.Rate,
		"status":           req.Status,
		"updated_by":       userID,
	}).Error
}

func (s *CustomerService) DeleteCustomer(id string) error {
	return db.DB.Delete(&models.Customer{}, id).Error
}

// ImportCustomers bulk imports customers from CSV, XLS, or XLSX files.
func (s *CustomerService) ImportCustomers(file *multipart.FileHeader, userID uint64) error {
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

	go s.processCustomerRowsInBackground(data, userID)

	return nil
}

func (s *CustomerService) processCustomerRowsInBackground(data [][]string, userID uint64) {
	totalRows := len(data) - 1
	if totalRows < 0 {
		return
	}

	importID := uint64(utils.Now().UnixNano())
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
				// Pad the row to 9 columns to prevent index out of range errors
				for len(row) < 9 {
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
						// Column mapping: [full_names(0), customer_type(1), customer_no(2), phone(3), kra_pin(4), postal_address(5), postal_code(6), postal_town(7), rate(8)]

						cNo := strings.TrimSpace(row[2])
						if cNo == "" {
							return fmt.Errorf("customer number is required")
						}

						var count int64
						tx.Model(&models.Customer{}).Where("customer_no = ?", cNo).Count(&count)
						if count > 0 {
							return fmt.Errorf("customer with number %s already exists", cNo)
						}

						// 1. Resolve/Create Customer Type
						typeName := strings.TrimSpace(row[1])
						var cType models.CustomerType
						if err := tx.Where("name = ?", typeName).First(&cType).Error; err != nil {
							cType = models.CustomerType{
								BaseModel: models.BaseModel{CreatedBy: userID},
								Name:      typeName,
							}
							if err := tx.Create(&cType).Error; err != nil {
								return fmt.Errorf("failed to create customer type: %w", err)
							}
						}

						rate, _ := utils.ParseFloat(row[8])

						customer := models.Customer{
							BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							CustomerTypeID: cType.ID,
							FullNames:      strings.TrimSpace(row[0]),
							CustomerNo:     cNo,
							Phone:          utils.NormalizePhone(row[3]),
							KraPin:         strings.TrimSpace(row[4]),
							PostalAddress:  strings.TrimSpace(row[5]),
							PostalCode:     strings.TrimSpace(row[6]),
							PostalTown:     strings.TrimSpace(row[7]),
							Rate:           rate,
							Status:         "ACTIVE",
						}

						return tx.Create(&customer).Error
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
		Title:            "Customer Import Status",
		Message:          fmt.Sprintf("Import completed. Success: %d, Failed: %d out of %d records.", totalRows-failedCount, failedCount, totalRows),
		NotificationType: notificationType,
		ErrorLink:        fmt.Sprintf("/customers/import-errors/%d", importID),
	})
}

// ExportCustomers initiates a background process to export customer data.
func (s *CustomerService) ExportCustomers(userID uint64, status, format string) error {
	go s.processCustomerExportInBackground(userID, status, format)
	return nil
}

func (s *CustomerService) processCustomerExportInBackground(userID uint64, status, format string) {
	var results []struct {
		models.Customer
		TypeName string `gorm:"column:type_name"`
	}

	query := `
		SELECT c.*, ct.name as type_name
		FROM customers c
		LEFT JOIN customer_types ct ON c.customer_type_id = ct.id
		WHERE c.deleted_at IS NULL`

	if status != "" {
		query += fmt.Sprintf(" AND c.status = '%s'", status)
	}

	if err := db.DB.Raw(query).Scan(&results).Error; err != nil {
		log.Printf("[CustomerService] Export query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"

	if strings.ToLower(format) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateCustomerPDF(results, status)
	} else {
		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)
		writer.Write([]string{"full_names", "customer_type", "customer_no", "phone", "kra_pin", "postal_address", "postal_code", "postal_town", "rate"})

		for _, c := range results {
			writer.Write([]string{
				c.FullNames, c.TypeName, c.CustomerNo, c.Phone, c.KraPin,
				c.PostalAddress, c.PostalCode, c.PostalTown, fmt.Sprintf("%.2f", c.Rate),
			})
		}
		writer.Flush()
		fileData = buf.Bytes()
		err = writer.Error()
	}

	if err != nil {
		log.Printf("[CustomerService] Error generating export: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Export Failed",
			Message:          fmt.Sprintf("Failed to generate %s export: %v", format, err),
			NotificationType: "ERROR",
		})
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("customers_export_%d.%s", utils.Now().UnixNano(), ext)
	filepath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filepath, fileData, 0644); err != nil {
		log.Printf("[CustomerService] Error writing file: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            fmt.Sprintf("Customer Export (%s) Ready", strings.ToUpper(ext)),
		Message:          fmt.Sprintf("Your customer data %s export is ready for download.", ext),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/customers/export/download/%s", filename),
	})
}

func (s *CustomerService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

func (s *CustomerService) generateCustomerPDF(results []struct {
	models.Customer
	TypeName string `gorm:"column:type_name"`
}, status string) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
		Phone          string `gorm:"column:phone"`
		Email          string `gorm:"column:email"`
	}
	db.DB.Table("organization_details").First(&org)

	statusVal := "ALL"
	if status != "" {
		statusVal = status
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Phone: %s | Email: %s", org.Phone, org.Email), "", 1, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "CUSTOMER REGISTER", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(0, 8, fmt.Sprintf("Showing (%d) Records for Status: %s", len(results), statusVal), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table Headers
	pdf.SetFont("Arial", "B", 8)
	headers := []string{"C-No", "Name", "Type", "Phone", "Town", "Rate", "Status"}
	widths := []float64{30, 80, 40, 40, 30, 30, 30}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 8)
	for _, c := range results {
		pdf.CellFormat(widths[0], 8, c.CustomerNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, c.FullNames, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, c.TypeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, c.Phone, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 8, c.PostalTown, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, fmt.Sprintf("%.2f", c.Rate), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[6], 8, c.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
