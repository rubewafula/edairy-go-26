package services

import (
	"encoding/csv"
	"fmt"
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

type CustomerService struct {
	notificationService *UINotificationService
}

func NewCustomerService() *CustomerService {
	return &CustomerService{
		notificationService: NewUINotificationService(),
	}
}

func (s *CustomerService) CreateCustomer(req dtos.CreateCustomerRequest) (*dtos.CustomerResponse, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	customer := &models.Customer{
		CustomerTypeID: req.CustomerTypeID,
		FullNames:      req.FullNames,
		Phone:          req.Phone,
		EmailAddress:   req.EmailAddress,
		CustomerNo:     req.CustomerNo,
		KraPin:         req.KraPin,
		Status:         "ACTIVE",
		CreditLimit:    req.CreditLimit,
		PostalAddress:  req.PostalAddress,
		PostalCode:     req.PostalCode,
		PostalTown:     req.PostalTown,
		Terms:          req.Terms,
		Rate:           req.Rate,
	}

	if err := db.DB.Create(customer).Error; err != nil {
		return nil, err
	}
	return s.GetCustomer(fmt.Sprintf("%d", customer.ID))
}

func (s *CustomerService) GetCustomers(page, limit int) ([]dtos.CustomerResponse, int64, error) {
	var results []dtos.CustomerResponse
	var total int64
	db.DB.Model(&models.Customer{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			c.*, ct.description as customer_type_name
		FROM customers c
		LEFT JOIN customer_types ct ON c.customer_type_id = ct.id
		WHERE c.deleted_at IS NULL
		ORDER BY c.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CustomerService) GetCustomer(id string) (*dtos.CustomerResponse, error) {
	var result dtos.CustomerResponse
	query := `
		SELECT 
			c.*, ct.description as customer_type_name
		FROM customers c
		LEFT JOIN customer_types ct ON c.customer_type_id = ct.id
		WHERE c.id = ? AND c.deleted_at IS NULL
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

func (s *CustomerService) UpdateCustomer(id string, req dtos.UpdateCustomerRequest) error {
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		return err
	}

	customer.CustomerTypeID = req.CustomerTypeID
	customer.FullNames = req.FullNames
	customer.Phone = req.Phone
	customer.EmailAddress = req.EmailAddress
	customer.KraPin = req.KraPin
	if req.Status != "" {
		customer.Status = req.Status
	}
	customer.CreditLimit = req.CreditLimit
	customer.PostalAddress = req.PostalAddress
	customer.PostalCode = req.PostalCode
	customer.PostalTown = req.PostalTown
	customer.Terms = req.Terms
	customer.Rate = req.Rate

	return db.DB.Save(&customer).Error
}

func (s *CustomerService) DeleteCustomer(id string) error {
	return db.DB.Delete(&models.Customer{}, id).Error
}

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
			return fmt.Errorf("no sheets found")
		}
		data, err = f.GetRows(sheets[0])
	} else {
		return fmt.Errorf("unsupported format")
	}

	if err != nil {
		return err
	}

	go s.processCustomerRowsInBackground(data, userID)
	return nil
}

func (s *CustomerService) processCustomerRowsInBackground(data [][]string, userID uint64) {
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
						if r := recover(); r != nil {
							db.DB.Create(&models.CustomerImportError{
								BaseModel: models.BaseModel{CreatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic: %v", r),
							})
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						if len(row) < 5 {
							return fmt.Errorf("insufficient columns")
						}
						custNo := strings.TrimSpace(row[0])
						var count int64
						tx.Model(&models.Customer{}).Where("customer_no = ?", custNo).Count(&count)
						if count > 0 {
							return fmt.Errorf("customer %s already exists", custNo)
						}

						rate, _ := utils.ParseFloat(row[5])
						customer := models.Customer{
							BaseModel:    models.BaseModel{CreatedBy: userID},
							CustomerNo:   custNo,
							FullNames:    row[1],
							Phone:        utils.NormalizePhone(row[2]),
							EmailAddress: row[3],
							KraPin:       row[4],
							Rate:         rate,
							Status:       "ACTIVE",
						}
						return tx.Create(&customer).Error
					})

					if err != nil {
						db.DB.Create(&models.CustomerImportError{
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

	failedCount := len(errorChan)
	msg := "Customer import completed successfully."
	if failedCount > 0 {
		msg = fmt.Sprintf("Customer import finished with %d failures. Check logs.", failedCount)
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title: "Customer Import Status", Message: msg, NotificationType: "IMPORT_STATUS",
	})
}
