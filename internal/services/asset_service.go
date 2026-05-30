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
		CurrentLocation:         req.CurrentLocation,
		Status:                  req.Status,
		Loanable:                req.Loanable,
		Comments:                req.Comments,
	}

	if err := db.DB.Create(asset).Error; err != nil {
		log.Println("Found error trying to created asset: %s", err.Error())
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
			a.accumulated_depreciation, a.book_value, a.warranty_end_date, a.current_location,
			a.status, a.loanable, a.comments, a.created_by, a.updated_by, a.created_at, a.updated_at
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.ud
	`

	baseCountQuery := `
		SELECT COUNT(*)
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.ud
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
			a.accumulated_depreciation, a.book_value, a.warranty_end_date, a.current_location,
			a.status, a.loanable, a.comments, a.created_by, a.updated_by, a.created_at, a.updated_at
		FROM fixed_assets a
		LEFT JOIN asset_categories ac ON a.asset_category_id = ac.ud
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
	asset.CurrentLocation = req.CurrentLocation
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

	go s.processAssetRowsInBackground(data, userID)
	return nil
}

func (s *AssetService) processAssetRowsInBackground(data [][]string, userID uint64) {
	var wg sync.WaitGroup
	jobs := make(chan []string, len(data)-1)
	errorChan := make(chan error, len(data)-1)
	numWorkers := runtime.NumCPU() * 2

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil { //
							db.DB.Create(&models.AssetImportError{
								BaseModel: models.BaseModel{CreatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic during import: %v", r),
							})
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						if len(row) < 5 {
							return fmt.Errorf("insufficient columns (found %d, need at least 5)", len(row))
						}

						assetCode := strings.TrimSpace(row[1])
						if assetCode == "" {
							return fmt.Errorf("asset code is required")
						}

						var count int64
						tx.Model(&models.Asset{}).Where("asset_code = ?", assetCode).Count(&count)
						if count > 0 {
							return fmt.Errorf("asset with code %s already exists", assetCode)
						}

						var category models.AssetCategory
						tx.Where("name = ?", strings.TrimSpace(row[2])).First(&category)

						cost, _ := utils.ParseFloat(row[4])
						fixedAsset := models.Asset{ //
							BaseModel:       models.BaseModel{CreatedBy: userID}, //
							AssetName:       strings.TrimSpace(row[0]),
							AssetCode:       assetCode,
							CategoryID:      category.ID,
							SerialNo:        strings.TrimSpace(row[3]),
							PurchaseCost:    cost,
							AcquisitionDate: utils.ParseFlexibleDate(row[5]),
							Status:          "ACTIVE",
						} //
						return tx.Create(&fixedAsset).Error
					})

					if err != nil {
						db.DB.Create(&models.AssetImportError{
							BaseModel: models.BaseModel{CreatedBy: userID},
							RowData:   strings.Join(row, ","),
							Error:     err.Error(),
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

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title: "Asset Import Status", Message: "Fixed asset import process has completed.", NotificationType: "ASSET_IMPORT",
	})
}
