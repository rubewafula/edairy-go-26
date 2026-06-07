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

type AssetService struct {
	notificationService *UINotificationService
}

func NewAssetService() *AssetService {
	return &AssetService{
		notificationService: NewUINotificationService(),
	}
}

func (s *AssetService) CreateAsset(req dtos.CreateAssetRequest) (*models.Asset, error) {
	asset := &models.Asset{

		AssetCode:               req.AssetCode,
		AssetName:               req.AssetName,
		CategoryID:              req.CategoryID,
		SerialNo:                req.SerialNo,
		Barcode:                 req.Barcode,
		Manufacturer:            req.Manufacturer,
		VendorID:                req.VendorID,
		PurchaseCost:            req.PurchaseCost,
		SalvageValue:            req.SalvageValue,
		AcquisitionDate:         utils.ParseDate(req.AcquisitionDate),
		UsefulLifeYears:         req.UsefulLifeYears,
		DepreciationMethod:      req.DepreciationMethod,
		DepreciationRate:        req.DepreciationRate,
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		BookValue:               req.BookValue,
		WarrantyEndDate:         utils.ParseDate(req.WarrantyEndDate),
		Status:                  req.Status,
		Loanable:                req.Loanable,
		Comments:                req.Comments,
		LocationID:              req.LocationID, // Assuming LocationID is in DTO
	}

	if err := db.DB.Create(asset).Error; err != nil {
		log.Println("Found error trying to created asset: %s", err.Error())
		log.Printf("Found error trying to created asset: %s", err.Error()) // Use Printf
		return nil, err
	}
	return asset, nil
}

func (s *AssetService) GetAssets(page, limit int, assetCode, assetName, categoryID, fromDate, toDate string) ([]dtos.AssetResponse, int64, error) {
	var results []dtos.AssetResponse
	var total int64

	// Calculate offset for pagination
	offset := (page - 1) * limit

	query := `
		SELECT 
			a.id, a.asset_code, a.asset_name, a.asset_category_id as category_id, ac.name as category_name,
			a.serial_no, a.barcode, a.manufacturer, a.vendor_id, a.purchase_cost, a.salvage_value,
			a.acquisition_date, a.useful_life_years, a.depreciation_method, a.depreciation_rate,
			a.accumulated_depreciation, a.book_value, a.warranty_end_date, l.location_name as current_location,
			a.status, a.loanable, a.comments, a.created_by, a.updated_by, a.created_at, a.updated_at
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.id
		LEFT JOIN locations  l ON a.location_id = l.id
	`

	baseCountQuery := `
		SELECT COUNT(*)
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.id
	`

	var args []interface{}
	whereClauses := []string{"a.deleted_at IS NULL"}

	if assetCode != "" {
		whereClauses = append(whereClauses, "a.asset_code LIKE ?")
		args = append(args, "%"+assetCode+"%")
	}
	if assetName != "" {
		whereClauses = append(whereClauses, "a.asset_name LIKE ?")
		args = append(args, "%"+assetName+"%")
	}
	if categoryID != "" {
		whereClauses = append(whereClauses, "a.asset_category_id = ?")
		args = append(args, categoryID)
	}
	if fromDate != "" && toDate != "" {
		whereClauses = append(whereClauses, "a.acquisition_date BETWEEN ? AND ?")
		args = append(args, fromDate, toDate)
	} else if fromDate != "" {
		whereClauses = append(whereClauses, "a.acquisition_date >= ?")
		args = append(args, fromDate)
	} else if toDate != "" {
		whereClauses = append(whereClauses, "a.acquisition_date <= ?")
		args = append(args, toDate)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	if err := db.DB.Raw(baseCountQuery+whereSql, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.DB.Raw(query+whereSql+" ORDER BY a.id DESC LIMIT ? OFFSET ?", append(args, limit, offset)...).Scan(&results).Error
	return results, total, err
}

func (s *AssetService) GetAsset(id string) (*dtos.AssetResponse, error) {
	var result dtos.AssetResponse
	query := `
		SELECT 
			a.id, a.asset_code, a.asset_name, a.asset_category_id as category_id, ac.name as category_name,
			a.serial_no, a.barcode, a.manufacturer, a.vendor_id, a.purchase_cost, a.salvage_value,
			a.acquisition_date, a.useful_life_years, a.depreciation_method, a.depreciation_rate,
			a.accumulated_depreciation, a.book_value, a.warranty_end_date, l.location_name as current_location,
			a.status, a.loanable, a.comments, a.created_by, a.updated_by, a.created_at, a.updated_at
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.id
		LEFT JOIN locations  l ON a.location_id = l.id
		WHERE a.id = ? AND a.deleted_at IS NULL
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

func (s *AssetService) UpdateAsset(id string, req dtos.UpdateAssetRequest) error {
	var asset models.Asset
	if err := db.DB.First(&asset, id).Error; err != nil {
		return err
	}

	asset.AssetCode = req.AssetCode
	asset.AssetName = req.AssetName
	asset.CategoryID = req.CategoryID
	asset.SerialNo = req.SerialNo
	asset.Barcode = req.Barcode
	asset.Manufacturer = req.Manufacturer
	asset.VendorID = req.VendorID
	asset.PurchaseCost = req.PurchaseCost
	asset.SalvageValue = req.SalvageValue
	asset.AcquisitionDate = utils.ParseDate(req.AcquisitionDate)
	asset.UsefulLifeYears = req.UsefulLifeYears
	asset.DepreciationMethod = req.DepreciationMethod
	asset.DepreciationRate = req.DepreciationRate
	asset.AccumulatedDepreciation = req.AccumulatedDepreciation
	asset.BookValue = req.BookValue
	asset.WarrantyEndDate = utils.ParseDate(req.WarrantyEndDate)
	asset.LocationID = req.LocationID
	asset.Status = req.Status
	asset.Loanable = req.Loanable
	asset.Comments = req.Comments

	return db.DB.Save(&asset).Error

}

func (s *AssetService) DeleteAsset(id string) error {
	return db.DB.Delete(&models.Asset{}, id).Error
}

func (s *AssetService) ImportAssets(file *multipart.FileHeader, userID uint64) error {
	src, err := file.Open() //
	if err != nil {
		return err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var data [][]string

	if ext == ".csv" {
		reader := csv.NewReader(src) //
		data, err = reader.ReadAll()
	} else if ext == ".xlsx" || ext == ".xls" {
		f, err := excelize.OpenReader(src) //
		if err != nil {
			return err
		}
		sheets := f.GetSheetList() //
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

	log.Printf("Asset Import found asset file will proceed in background :%s", file.Filename)

	go s.processAssetRowsInBackground(data, userID)
	return nil
}

func (s *AssetService) processAssetRowsInBackground(data [][]string, userID uint64) {
	totalRows := len(data) - 1
	if totalRows < 0 {
		return
	}

	importID := uint64(utils.Now().UnixNano())

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
						if r := recover(); r != nil { //
							db.DB.Create(&models.ImportError{
								BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic during import: %v", r),
								ImportId:  importID,
							})
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						// Expected columns: [Code, Asset Name, Asset Category, Purchase Cost, Acquisition Date, Depreciation Method, Location, Warranty End Date]
						if len(row) < 8 {
							return fmt.Errorf("insufficient columns (found %d, need at least 8)", len(row))
						}

						assetCode := strings.TrimSpace(row[0])
						if assetCode == "" {
							return fmt.Errorf("asset code is required")
						}

						var count int64
						tx.Model(&models.Asset{}).Where("asset_code = ?", assetCode).Count(&count)
						if count > 0 {
							return fmt.Errorf("asset with code %s already exists", assetCode)
						}

						// Resolve/Create Asset Category
						categoryName := strings.TrimSpace(row[2])
						var category models.AssetCategory
						if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
							if err == gorm.ErrRecordNotFound {
								category = models.AssetCategory{
									BaseModel: models.BaseModel{CreatedBy: userID},
									Name:      categoryName,
								}
								if err := tx.Create(&category).Error; err != nil {
									return fmt.Errorf("failed to create asset category %s: %w", categoryName, err)
								}
							} else {
								return fmt.Errorf("failed to query asset category %s: %w", categoryName, err)
							}
						}

						// Resolve Location
						locationName := strings.TrimSpace(row[6])
						var location models.Location
						var locationID uint64
						if err := tx.Where("location_name = ? OR code = ?", locationName, locationName).First(&location).Error; err != nil {
							if err == gorm.ErrRecordNotFound {
								return fmt.Errorf("location '%s' not found in system", locationName)
							} else {
								return fmt.Errorf("failed to query location %s: %w", locationName, err)
							}
						}
						locationID = location.ID

						cost, _ := utils.ParseFloat(row[3])
						fixedAsset := models.Asset{ //
							BaseModel:          models.BaseModel{CreatedBy: userID}, //
							AssetName:          strings.TrimSpace(row[1]),
							AssetCode:          assetCode,
							CategoryID:         category.ID,
							LocationID:         locationID,
							PurchaseCost:       cost,
							AcquisitionDate:    utils.ParseFlexibleDate(row[4]),
							DepreciationMethod: strings.TrimSpace(row[5]),
							WarrantyEndDate:    utils.ParseFlexibleDate(row[7]),
							Status:             "ACTIVE",
						} //
						return tx.Create(&fixedAsset).Error
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

	failCount := 0
	for range errorChan {
		failCount++
	}

	message := fmt.Sprintf("Fixed asset import completed. Success: %d, Failed: %d out of %d records.", totalRows-failCount, failCount, totalRows)
	notificationType := "SUCCESS"
	errorLink := ""
	if failCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/fixed-assets/import-errors/%d", importID)
	} else if totalRows == 0 {
		// If no rows were processed (e.g., empty file or only header), it's an info or warning
		message = "Fixed asset import completed. No records were processed."
		notificationType = "INFO"
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Asset Import Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
	})
}

func (s *AssetService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	query := db.DB.Where("import_id = ?", importID)
	err := query.Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

func (s *AssetService) ProcessReportInBackground(userID uint64, format, assetCode, assetName, categoryID, fromDate, toDate string) {
	go func() {
		content, filename, err := s.GenerateReport(format, assetCode, assetName, categoryID, fromDate, toDate)
		if err != nil {
			log.Printf("[AssetService.ProcessReportInBackground] Error generating report: %v", err)
			s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
				Title:            "Report Generation Failed",
				Message:          "There was an error generating your fixed assets report.",
				NotificationType: "ERROR",
			})
			return
		}

		// Ensure directory exists
		reportDir := "./storage/reports"
		if err := os.MkdirAll(reportDir, 0755); err != nil {
			log.Printf("[AssetService.ProcessReportInBackground] Failed to create reports directory: %v", err)
			return
		}

		filePath := filepath.Join(reportDir, filename)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			log.Printf("[AssetService.ProcessReportInBackground] Failed to save report: %v", err)
			return
		}

		downloadLink := fmt.Sprintf("/api/fixed-assets/report/download/%s", filename)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Report Generated",
			Message:          fmt.Sprintf("Fixed assets report (%s) is ready for download.", strings.ToUpper(format)),
			NotificationType: "SUCCESS",
			ErrorLink:        downloadLink,
		})
	}()
}

func (s *AssetService) GenerateReport(format, assetCode, assetName, categoryID, fromDate, toDate string) ([]byte, string, error) {
	var results []dtos.AssetResponse
	var args []interface{}
	whereClauses := []string{"a.deleted_at IS NULL"}

	if assetCode != "" {
		whereClauses = append(whereClauses, "a.asset_code LIKE ?")
		args = append(args, "%"+assetCode+"%")
	}
	if assetName != "" {
		whereClauses = append(whereClauses, "a.asset_name LIKE ?")
		args = append(args, "%"+assetName+"%")
	}
	if categoryID != "" {
		whereClauses = append(whereClauses, "a.asset_category_id = ?")
		args = append(args, categoryID)
	}
	if fromDate != "" && toDate != "" {
		whereClauses = append(whereClauses, "a.acquisition_date BETWEEN ? AND ?")
		args = append(args, fromDate, toDate)
	} else if fromDate != "" {
		whereClauses = append(whereClauses, "a.acquisition_date >= ?")
		args = append(args, fromDate)
	} else if toDate != "" {
		whereClauses = append(whereClauses, "a.acquisition_date <= ?")
		args = append(args, toDate)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")
	query := `
		SELECT 
			a.id, a.asset_code, a.asset_name, a.asset_category_id as category_id, ac.name as category_name,
			a.serial_no, a.barcode, a.manufacturer, a.purchase_cost, a.acquisition_date,
			l.location_name as current_location, a.status
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.id
		LEFT JOIN locations l ON a.location_id = l.id
	` + whereSql + " ORDER BY a.id DESC"

	if err := db.DB.Raw(query, args...).Scan(&results).Error; err != nil {
		return nil, "", err
	}

	org := s.getOrgInfo()
	filterSummary := s.getFilterSummary(assetCode, assetName, categoryID, fromDate, toDate)
	filename := fmt.Sprintf("assets_report_%d.%s", utils.Now().Unix(), format)

	if format == "csv" {
		content, err := s.generateCSV(results, org, filterSummary)
		return content, filename, err
	}

	content, err := s.generatePDF(results, org, filterSummary)
	return content, filename, err
}

func (s *AssetService) getOrgInfo() map[string]string {
	var org struct {
		Name          string `gorm:"column:name"`
		PostalAddress string `gorm:"column:postal_address"`
		PostalCode    string `gorm:"column:postal_code"`
		Town          string `gorm:"column:town"`
		Phone         string `gorm:"column:phone"`
		Email         string `gorm:"column:email"`
	}
	if err := db.DB.Table("organizations").First(&org).Error; err != nil {
		return map[string]string{
			"name":    "EDAIRY ERP SYSTEM",
			"address": "2772",
			"code":    "00606",
			"town":    "NAIROBI",
			"phone":   "0726986944",
			"email":   "technology@edairy.africa",
		}
	}
	return map[string]string{
		"name":    org.Name,
		"address": org.PostalAddress,
		"code":    org.PostalCode,
		"town":    org.Town,
		"phone":   org.Phone,
		"email":   org.Email,
	}
}

func (s *AssetService) getFilterSummary(code, name, cat, from, to string) string {
	var filters []string
	if code != "" {
		filters = append(filters, "Code: "+code)
	}
	if name != "" {
		filters = append(filters, "Name: "+name)
	}
	if cat != "" {
		filters = append(filters, "Category ID: "+cat)
	}
	if from != "" || to != "" {
		filters = append(filters, fmt.Sprintf("Period: %s to %s", from, to))
	}
	if len(filters) == 0 {
		return "No filters applied"
	}
	return strings.Join(filters, " | ")
}

func (s *AssetService) generateCSV(data []dtos.AssetResponse, org map[string]string, filters string) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{org["name"]})
	writer.Write([]string{fmt.Sprintf("P.O Box %s, %s, %s", org["address"], org["code"], org["town"])})
	writer.Write([]string{fmt.Sprintf("Phone: %s | Email: %s", org["phone"], org["email"])})
	writer.Write([]string{"Fixed Assets Report"})
	writer.Write([]string{"Report generated on " + utils.Now().Format("02/01/2006 15:04:05")})
	writer.Write([]string{fmt.Sprintf("Showing (%d) asset records", len(data))})
	writer.Write([]string{filters})
	writer.Write([]string{})
	writer.Write([]string{"ID", "Asset Code", "Asset Name", "Category", "Acquisition Date", "Cost", "Location", "Status"})
	for _, a := range data {
		writer.Write([]string{
			fmt.Sprintf("%d", a.ID), a.AssetCode, a.AssetName, a.CategoryName,
			a.AcquisitionDate.Format("02/01/2006"), fmt.Sprintf("%.2f", a.PurchaseCost),
			a.CurrentLocation, a.Status,
		})
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *AssetService) generatePDF(data []dtos.AssetResponse, org map[string]string, filters string) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, org["name"])
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("P.O Box %s, %s, %s", org["address"], org["code"], org["town"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Phone: %s | Email: %s", org["phone"], org["email"]))
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Fixed Assets Report")
	pdf.Ln(6)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 10, "Report generated on "+utils.Now().Format("02/01/2006 15:04:05"))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Showing (%d) asset records", len(data)))
	pdf.Ln(5)
	pdf.Cell(0, 10, filters)
	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 10)
	headers := []string{"ID", "Code", "Asset Name", "Category", "Acq. Date", "Cost", "Location", "Status"}
	w := []float64{15, 30, 60, 40, 30, 30, 45, 25}
	for i, h := range headers {
		pdf.CellFormat(w[i], 7, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	for _, a := range data {
		pdf.CellFormat(w[0], 6, fmt.Sprintf("%d", a.ID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(w[1], 6, a.AssetCode, "1", 0, "L", false, 0, "")
		pdf.CellFormat(w[2], 6, a.AssetName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(w[3], 6, a.CategoryName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(w[4], 6, a.AcquisitionDate.Format("02/01/2006"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(w[5], 6, fmt.Sprintf("%.2f", a.PurchaseCost), "1", 0, "R", false, 0, "")
		pdf.CellFormat(w[6], 6, a.CurrentLocation, "1", 0, "L", false, 0, "")
		pdf.CellFormat(w[7], 6, a.Status, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
