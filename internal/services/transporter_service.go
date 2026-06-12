package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
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

type TransporterService struct {
	notificationService *UINotificationService
}

func NewTransporterService() *TransporterService {
	return &TransporterService{
		notificationService: NewUINotificationService(),
	}
}

// sanitizeNumericString removes all non-digit characters from a string.
func sanitizeNumericString(s string) string {
	reg := regexp.MustCompile(`[^0-9]+`)
	return reg.ReplaceAllString(s, "")
}

func (s *TransporterService) CreateTransporter(req dtos.CreateTransporterRequest) (*dtos.TransporterResponse, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	transporter := &models.Transporter{
		TransporterNo: req.TransporterNo,
		Category:      req.Category,
		PrimaryPhone:  req.PrimaryPhone,
		EmailAddress:  req.EmailAddress,
		Status:        status,
		Restricted:    req.Restricted,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transporter).Error; err != nil {
			return err
		}

		if req.Category == "INDIVIDUAL" {
			passportPath, err := utils.SaveFile(req.PassportPhoto, "transporters")
			if err != nil {
				return fmt.Errorf("failed to save passport photo: %w", err)
			}
			idFrontPath, err := utils.SaveFile(req.IDFrontPhoto, "transporters")
			if err != nil {
				return fmt.Errorf("failed to save ID front photo: %w", err)
			}
			idBackPath, err := utils.SaveFile(req.IDBackPhoto, "transporters")
			if err != nil {
				return fmt.Errorf("failed to save ID back photo: %w", err)
			}

			nationalIDNoSanitized := sanitizeNumericString(req.NationalIDNo)
			kraPinSanitized := sanitizeNumericString(req.KraPin)
			individual := &models.IndividualTransporter{
				TransporterID:     transporter.ID,
				FirstName:         req.FirstName,
				LastName:          req.LastName,
				OtherNames:        utils.StringPtr(req.OtherNames),
				Gender:            utils.StringPtr(req.Gender),
				NationalIDNo:      utils.StringPtr(nationalIDNoSanitized),
				MaritalStatus:     utils.StringPtr(req.MaritalStatus),
				KraPin:            utils.StringPtr(kraPinSanitized), // KRA PIN is common to both individual and company
				NextOfKinFullName: utils.StringPtr(req.NextOfKinFullName),
				NextOfKinPhone:    utils.StringPtr(req.NextOfKinPhone),
				PassportPhoto:     utils.StringPtr(passportPath),
				IDFrontPhoto:      utils.StringPtr(idFrontPath),
				IDBackPhoto:       utils.StringPtr(idBackPath),
			}

			if req.DateOfBirth != "" {
				t := utils.ParseDate(req.DateOfBirth)
				if !t.IsZero() {
					individual.DateOfBirth = &t
				}
			}

			return tx.Create(individual).Error
		} else {
			certificatePath := ""
			if req.CertificateOfIncorporation != nil {
				certificatePath, _ = utils.SaveFile(req.CertificateOfIncorporation, "transporters")
			}

			company := &models.CompanyTransporter{
				TransporterID:              transporter.ID,
				CompanyName:                req.CompanyName,
				RegistrationNo:             utils.StringPtr(req.RegistrationNo),
				KraPin:                     utils.StringPtr(req.KraPin),
				ContactPersonName:          utils.StringPtr(req.ContactPersonName),
				ContactPersonPhone:         utils.StringPtr(req.ContactPersonPhone),
				PostalAddress:              utils.StringPtr(req.PostalAddress),
				PostalCode:                 utils.StringPtr(req.PostalCode),
				Town:                       utils.StringPtr(req.Town),
				CertificateOfIncorporation: utils.StringPtr(certificatePath),
			}
			return tx.Create(company).Error
		}
	})

	if err != nil {
		log.Printf("Transporter: Error creating: %s", err.Error())
		return nil, err
	}

	var created models.Transporter
	if err := db.DB.Preload("Individual").Preload("Company").First(&created, transporter.ID).Error; err != nil {
		return nil, err
	}
	res := s.toTransporterResponse(created)
	return &res, nil
}

func (s *TransporterService) GetTransporters(page, limit int) ([]dtos.TransporterResponse, int64, error) {
	var transporters []models.Transporter
	var total int64

	db.DB.Model(&models.Transporter{}).Count(&total)

	offset := (page - 1) * limit
	err := db.DB.Preload("Individual").Preload("Company").
		Limit(limit).Offset(offset).Order("id DESC").Find(&transporters).Error
	if err != nil {
		return nil, 0, err
	}

	var responses []dtos.TransporterResponse
	for _, t := range transporters {
		responses = append(responses, s.toTransporterResponse(t))
	}
	return responses, total, nil
}

func (s *TransporterService) GetTransporter(id string) (*dtos.TransporterResponse, error) {
	var transporter models.Transporter
	if err := db.DB.Preload("Individual").Preload("Company").First(&transporter, id).Error; err != nil {
		return nil, err
	}
	response := s.toTransporterResponse(transporter)
	return &response, nil
}

func (s *TransporterService) toTransporterResponse(t models.Transporter) dtos.TransporterResponse {
	response := dtos.TransporterResponse{
		ID:            t.ID,
		TransporterNo: t.TransporterNo,
		Category:      t.Category,
		PrimaryPhone:  t.PrimaryPhone,
		EmailAddress:  t.EmailAddress,
		Status:        t.Status,
		Restricted:    t.Restricted,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	if t.Individual != nil {
		response.Individual = &dtos.IndividualTransporterResponse{
			ID:                t.Individual.ID,
			TransporterID:     t.Individual.TransporterID,
			FirstName:         t.Individual.FirstName,
			LastName:          t.Individual.LastName,
			OtherNames:        utils.StringValue(t.Individual.OtherNames),
			Gender:            utils.StringValue(t.Individual.Gender),
			DateOfBirth:       t.Individual.DateOfBirth,
			NationalIDNo:      utils.StringValue(t.Individual.NationalIDNo),
			KraPin:            utils.StringValue(t.Individual.KraPin),
			MaritalStatus:     utils.StringValue(t.Individual.MaritalStatus),
			NextOfKinFullName: utils.StringValue(t.Individual.NextOfKinFullName),
			NextOfKinPhone:    utils.StringValue(t.Individual.NextOfKinPhone),
			PassportPhoto:     utils.StringValue(t.Individual.PassportPhoto),
			IDFrontPhoto:      utils.StringValue(t.Individual.IDFrontPhoto),
			IDBackPhoto:       utils.StringValue(t.Individual.IDBackPhoto),
		}
	}
	if t.Company != nil {
		response.Company = &dtos.CompanyTransporterResponse{
			ID:                         t.Company.ID,
			TransporterID:              t.Company.TransporterID,
			CompanyName:                t.Company.CompanyName,
			RegistrationNo:             utils.StringValue(t.Company.RegistrationNo),
			KraPin:                     utils.StringValue(t.Company.KraPin),
			ContactPersonName:          utils.StringValue(t.Company.ContactPersonName),
			ContactPersonPhone:         utils.StringValue(t.Company.ContactPersonPhone),
			PostalAddress:              utils.StringValue(t.Company.PostalAddress),
			PostalCode:                 utils.StringValue(t.Company.PostalCode),
			Town:                       utils.StringValue(t.Company.Town),
			CertificateOfIncorporation: utils.StringValue(t.Company.CertificateOfIncorporation),
		}
	}
	return response
}

func (s *TransporterService) UpdateTransporter(id string, req dtos.UpdateTransporterRequest) error {
	var transporter models.Transporter
	if err := db.DB.First(&transporter, id).Error; err != nil {
		return err
	}

	transporter.PrimaryPhone = req.PrimaryPhone
	transporter.EmailAddress = req.EmailAddress
	transporter.Status = req.Status
	transporter.Restricted = req.Restricted

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&transporter).Error; err != nil {
			return err
		}

		// Update individual or company details based on category
		if transporter.Category == "INDIVIDUAL" {
			var individual models.IndividualTransporter
			if err := tx.Where("transporter_id = ?", transporter.ID).First(&individual).Error; err != nil {
				return fmt.Errorf("individual transporter details not found: %w", err)
			}

			individual.FirstName = req.FirstName
			individual.LastName = req.LastName
			individual.OtherNames = utils.StringPtr(req.OtherNames)
			individual.Gender = utils.StringPtr(req.Gender)
			individual.MaritalStatus = utils.StringPtr(req.MaritalStatus)
			individual.NextOfKinFullName = utils.StringPtr(req.NextOfKinFullName)
			individual.NextOfKinPhone = utils.StringPtr(req.NextOfKinPhone)

			// Sanitize and update NationalIDNo and KraPin
			nationalIDNoSanitized := sanitizeNumericString(req.NationalIDNo)
			kraPinSanitized := sanitizeNumericString(req.KraPin)
			individual.NationalIDNo = utils.StringPtr(nationalIDNoSanitized)
			individual.KraPin = utils.StringPtr(kraPinSanitized)

			if req.DateOfBirth != "" {
				t := utils.ParseDate(req.DateOfBirth)
				individual.DateOfBirth = &t
			}

			// Handle file uploads for individual
			if req.PassportPhoto != nil {
				passportPath, err := utils.SaveFile(req.PassportPhoto, "transporters")
				if err != nil {
					return fmt.Errorf("failed to save passport photo: %w", err)
				}
				if passportPath != "" { // Only update if a new file was successfully saved
					individual.PassportPhoto = utils.StringPtr(passportPath)
				}
			}
			if req.IDFrontPhoto != nil {
				idFrontPath, err := utils.SaveFile(req.IDFrontPhoto, "transporters")
				if err != nil {
					return fmt.Errorf("failed to save ID front photo: %w", err)
				}
				if idFrontPath != "" { // Only update if a new file was successfully saved
					individual.IDFrontPhoto = utils.StringPtr(idFrontPath)
				}
			}
			if req.IDBackPhoto != nil {
				idBackPath, err := utils.SaveFile(req.IDBackPhoto, "transporters")
				if err != nil {
					return fmt.Errorf("failed to save ID back photo: %w", err)
				}
				if idBackPath != "" { // Only update if a new file was successfully saved
					individual.IDBackPhoto = utils.StringPtr(idBackPath)
				}
			}

			return tx.Save(&individual).Error
		} else if transporter.Category == "COMPANY" {
			var company models.CompanyTransporter
			if err := tx.Where("transporter_id = ?", transporter.ID).First(&company).Error; err != nil {
				return fmt.Errorf("company transporter details not found: %w", err)
			}

			company.CompanyName = req.CompanyName
			company.ContactPersonName = utils.StringPtr(req.ContactPersonName)
			company.ContactPersonPhone = utils.StringPtr(req.ContactPersonPhone)
			company.PostalAddress = utils.StringPtr(req.PostalAddress)
			company.PostalCode = utils.StringPtr(req.PostalCode)
			company.Town = utils.StringPtr(req.Town)

			// Sanitize and update RegistrationNo and KraPin
			registrationNoSanitized := sanitizeNumericString(req.RegistrationNo)
			kraPinCompanySanitized := sanitizeNumericString(req.KraPin)
			company.RegistrationNo = utils.StringPtr(registrationNoSanitized)
			company.KraPin = utils.StringPtr(kraPinCompanySanitized)

			// Handle file upload for company
			if req.CertificateOfIncorporation != nil {
				certificatePath, err := utils.SaveFile(req.CertificateOfIncorporation, "transporters")
				if err != nil {
					return fmt.Errorf("failed to save certificate of incorporation: %w", err)
				}
				if certificatePath != "" { // Only update if a new file was successfully saved
					company.CertificateOfIncorporation = utils.StringPtr(certificatePath)
				}
			}
			return tx.Save(&company).Error
		}
		return nil // No specific individual/company update needed or category not matched
	})
}

func (s *TransporterService) DeleteTransporter(id string) error {
	var transporter models.Transporter
	if err := db.DB.First(&transporter, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&transporter).Error
}

// ImportTransporters bulk imports transporters from CSV or Excel files.
func (s *TransporterService) ImportTransporters(file *multipart.FileHeader, userID uint64) error {
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

	go s.processTransporterRowsInBackground(data, userID)

	return nil
}

func (s *TransporterService) processTransporterRowsInBackground(data [][]string, userID uint64) {
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
				// Pad the row to 14 columns with empty strings to prevent index out of range errors
				for len(row) < 14 {
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
						// Columns: [transporter_no(0), category(1), primary_phone(2), email_address(3), national_id_no(4), first_name(5), last_name(6), other_names(7), gender(8), date_of_birth(9), company_name(10), registration_no(11), kra_pin(12), certificate_of_incorporation(13)]
						tNo := strings.TrimSpace(row[0])
						category := strings.ToUpper(strings.TrimSpace(row[1]))

						var count int64
						tx.Model(&models.Transporter{}).Where("transporter_no = ?", tNo).Count(&count)
						if count > 0 {
							return fmt.Errorf("transporter with number %s already exists", tNo)
						}

						transporter := models.Transporter{
							BaseModel:     models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
							TransporterNo: tNo,
							Category:      category,
							PrimaryPhone:  utils.NormalizePhone(row[2]),
							EmailAddress:  strings.TrimSpace(row[3]),
							Status:        "ACTIVE",
						}
						if err := tx.Create(&transporter).Error; err != nil {
							return err
						}

						if category == "INDIVIDUAL" {
							dob := utils.ParseFlexibleDate(row[9])
							nationalIDNoInput := strings.TrimSpace(row[4])
							kraPinInput := strings.TrimSpace(row[12]) // KRA PIN can be in row 12 for both types

							individual := models.IndividualTransporter{
								BaseModel:     models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
								TransporterID: transporter.ID,
								FirstName:     strings.TrimSpace(row[5]),
								LastName:      strings.TrimSpace(row[6]),
								// Logic: Even if row 10 (Company Name) has data, we ignore it because Category is INDIVIDUAL
								OtherNames:    utils.StringPtr(strings.TrimSpace(row[7])), // Use StringPtr for nullable fields
								Gender:        utils.StringPtr(strings.TrimSpace(row[8])),
								NationalIDNo:  utils.StringPtr(nationalIDNoInput), // Use StringPtr for nullable fields
								KraPin:        utils.StringPtr(kraPinInput),       // Use StringPtr for nullable fields
								MaritalStatus: nil,                                // Not in import, so nil by default if *string
							}

							if !dob.IsZero() {
								individual.DateOfBirth = &dob
							}

							return tx.Create(&individual).Error
						} else {
							// Logic: If Category is COMPANY, we read specific company columns (10, 11, 13)
							companyNameInput := strings.TrimSpace(row[10])
							regNoInput := strings.TrimSpace(row[11])
							kraPinInput := strings.TrimSpace(row[12])
							certInput := strings.TrimSpace(row[13])

							company := models.CompanyTransporter{
								BaseModel:                  models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
								TransporterID:              transporter.ID,
								CompanyName:                companyNameInput, // Assuming CompanyName is a required field (string)
								RegistrationNo:             utils.StringPtr(regNoInput),
								KraPin:                     utils.StringPtr(kraPinInput),
								CertificateOfIncorporation: utils.StringPtr(certInput),
							}

							return tx.Create(&company).Error
						}
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
		Title:            "Transporter Import Status",
		Message:          fmt.Sprintf("Import completed. Success: %d, Failed: %d out of %d records.", totalRows-failedCount, failedCount, totalRows),
		NotificationType: notificationType,
		ErrorLink:        fmt.Sprintf("/transporters/import-errors/%d", importID),
	})
}

// ExportTransporters initiates a background process to export transporter data to CSV.
func (s *TransporterService) ExportTransporters(userID uint64, status, format string) error {
	go s.processTransporterExportInBackground(userID, status, format)
	return nil
}

func (s *TransporterService) processTransporterExportInBackground(userID uint64, status, format string) {
	var results []struct {
		models.Transporter
		NationalIDNo               string     `gorm:"column:national_id_no"`
		FirstName                  string     `gorm:"column:first_name"`
		LastName                   string     `gorm:"column:last_name"`
		OtherNames                 string     `gorm:"column:other_names"`
		Gender                     string     `gorm:"column:gender"`
		DateOfBirth                *time.Time `gorm:"column:date_of_birth"`
		CompanyName                string     `gorm:"column:company_name"`
		RegistrationNo             string     `gorm:"column:registration_no"`
		KraPin                     string     `gorm:"column:kra_pin"`
		CertificateOfIncorporation string     `gorm:"column:certificate_of_incorporation"`
	}

	query := `
		SELECT 
			t.*, i.national_id_no, i.first_name, i.last_name, i.other_names, i.gender, i.date_of_birth,
			c.company_name, c.registration_no, c.certificate_of_incorporation,
			COALESCE(i.kra_pin, c.kra_pin) as kra_pin
		FROM transporters t
		LEFT JOIN individual_transporters i ON t.id = i.transporter_id
		LEFT JOIN company_transporters c ON t.id = c.transporter_id
		WHERE t.deleted_at IS NULL`

	if status != "" {
		query += fmt.Sprintf(" AND t.status = '%s'", status)
	}

	if err := db.DB.Raw(query).Scan(&results).Error; err != nil {
		log.Printf("[TransporterService] Export query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"

	if strings.ToLower(format) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateTransporterPDF(results, status)
	} else {
		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)
		writer.Write([]string{"transporter_no", "category", "primary_phone", "email_address", "national_id_no", "first_name", "last_name", "other_names", "gender", "date_of_birth", "company_name", "registration_no", "kra_pin", "certificate_of_incorporation"})

		for _, t := range results {
			dob := ""
			if t.DateOfBirth != nil {
				dob = t.DateOfBirth.Format("2006-01-02")
			}
			writer.Write([]string{t.TransporterNo, t.Category, t.PrimaryPhone, t.EmailAddress, t.NationalIDNo, t.FirstName, t.LastName, t.OtherNames, t.Gender, dob, t.CompanyName, t.RegistrationNo, t.KraPin, t.CertificateOfIncorporation})
		}
		writer.Flush()
		fileData = buf.Bytes()
		err = writer.Error()
	}

	if err != nil {
		log.Printf("[TransporterService] Error generating export: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Export Failed",
			Message:          fmt.Sprintf("Failed to generate %s export: %v", format, err),
			NotificationType: "ERROR",
		})
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("transporters_export_%d.%s", utils.Now().UnixNano(), ext)
	filepath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filepath, fileData, 0644); err != nil {
		log.Printf("[TransporterService] Error writing file: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            fmt.Sprintf("Transporter Export (%s) Ready", strings.ToUpper(ext)),
		Message:          fmt.Sprintf("Your transporter data %s export is ready for download.", ext),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/transporters/export/download/%s", filename),
	})
}

func (s *TransporterService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

func (s *TransporterService) generateTransporterPDF(results []struct {
	models.Transporter
	NationalIDNo               string     `gorm:"column:national_id_no"`
	FirstName                  string     `gorm:"column:first_name"`
	LastName                   string     `gorm:"column:last_name"`
	OtherNames                 string     `gorm:"column:other_names"`
	Gender                     string     `gorm:"column:gender"`
	DateOfBirth                *time.Time `gorm:"column:date_of_birth"`
	CompanyName                string     `gorm:"column:company_name"`
	RegistrationNo             string     `gorm:"column:registration_no"`
	KraPin                     string     `gorm:"column:kra_pin"`
	CertificateOfIncorporation string     `gorm:"column:certificate_of_incorporation"`
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
	pdf.CellFormat(0, 10, "TRANSPORTER REGISTER", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(0, 8, fmt.Sprintf("Showing (%d) Records for Status: %s", len(results), statusVal), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table Headers
	pdf.SetFont("Arial", "B", 8)
	headers := []string{"T-No", "Name / Company", "ID / Reg No", "Phone", "Category", "Status"}
	widths := []float64{30, 80, 40, 40, 30, 30}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 8)
	for _, t := range results {
		displayName := t.CompanyName
		displayID := t.RegistrationNo
		if t.Category == "INDIVIDUAL" {
			displayName = strings.TrimSpace(fmt.Sprintf("%s %s %s", t.FirstName, t.LastName, t.OtherNames))
			displayID = t.NationalIDNo
		}

		// Truncate name if too long for cell
		if len(displayName) > 40 {
			displayName = displayName[:37] + "..."
		}

		pdf.CellFormat(widths[0], 8, t.TransporterNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, displayName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, displayID, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, t.PrimaryPhone, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 8, t.Category, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 8, t.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
