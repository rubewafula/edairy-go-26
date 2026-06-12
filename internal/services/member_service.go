package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"runtime"

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	models "github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type MemberService struct {
	notificationService *UINotificationService
}

func NewMemberService() *MemberService {
	return &MemberService{
		notificationService: NewUINotificationService(),
	}
}

// Create user
func (s *MemberService) CreateMember(
	ctx context.Context,
	req dtos.CreateMemberRequest,
	userID uint64) (*models.Member, error) {

	// Validate ID number uniqueness before processing files or starting transaction
	var count int64
	db.DB.Model(&models.Member{}).Where("id_no = ?", req.IDNo).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("a member with ID number %s is already registered", req.IDNo)
	}

	memberNo := req.MemberNo
	if memberNo == "" {
		memberNo = utils.GenerateMemberNo()
	}

	idFrontPath, err := utils.SaveFile(req.IDFrontPhoto, "members")
	if err != nil {
		return nil, err
	}

	idBackPath, err := utils.SaveFile(req.IDBackPhoto, "members")
	if err != nil {
		return nil, err
	}

	passportPath, err := utils.SaveFile(req.PassportPhoto, "members")
	if err != nil {
		return nil, err
	}

	newMember := &models.Member{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		MemberNo:       memberNo,
		IDNo:           req.IDNo,
		MemberTypeID:   req.MemberTypeID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		OtherNames:     req.OtherNames,
		RouteID:        req.RouteID,
		DateOfBirth:    utils.ParseDate(req.DOB),
		Gender:         req.Gender,
		BirthCity:      req.BirthCity,
		PrimaryPhone:   utils.NormalizePhone(req.PrimaryPhone),
		SecondaryPhone: utils.NormalizePhone(req.SecondaryPhone),
		Email:          req.Email,
		NumberOfCows:   req.NumberOfCows,
		IdFrontPhoto:   idFrontPath,
		IdBackPhoto:    idBackPath,
		PassportPhoto:  passportPath,
		TaxNumber:      req.TaxNumber,
		MaritalStatus:  req.MaritalStatus, //
		Title:          req.Title,         //
		Status:         "PENDING",         //
		DateRegistered: utils.Now(),       //
	}

	if req.IDDateOfIssue != "" {
		t := utils.ParseDate(req.IDDateOfIssue)
		newMember.IdDateOfIssue = &t
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error { // Start of transaction
		if err := tx.Create(newMember).Error; err != nil {
			return err // Return error to transaction, not nil, err
		}

		// Handle NextOfKins
		var primaryNOK *dtos.CreateMemberNextOfKinRequest
		if len(req.NextOfKins) > 0 {
			for i := range req.NextOfKins { // Correctly iterate over slice to avoid referencing loop variable
				nokReq := &req.NextOfKins[i]
				nok := &models.MemberNextOfKin{
					BaseModel:              models.BaseModel{CreatedBy: userID},
					MemberID:               newMember.ID,
					FullName:               nokReq.FullName,
					Relationship:           nokReq.Relationship,
					PhoneNumber:            nokReq.PhoneNumber,
					AlternativePhoneNumber: nokReq.AlternativePhoneNumber,
					EmailAddress:           nokReq.EmailAddress,
					NationalIDNo:           nokReq.NationalIDNo,
					PostalAddress:          nokReq.PostalAddress,
					PhysicalAddress:        nokReq.PhysicalAddress,
					Occupation:             nokReq.Occupation,
					IsPrimary:              nokReq.IsPrimary,
					Status:                 nokReq.Status,
					Remarks:                nokReq.Remarks,
				}
				if err := tx.Create(nok).Error; err != nil {
					return err
				}
				if nokReq.IsPrimary && primaryNOK == nil { // Take the first primary NOK
					primaryNOK = nokReq // Assign value, not address of loop variable
				}
			}
		}

		// Set primary next of kin on the member model
		if primaryNOK != nil {
			newMember.NextOfKinFullName = primaryNOK.FullName
			if primaryNOK.PhoneNumber != nil {
				newMember.NextOfKinPhone = *primaryNOK.PhoneNumber
			}
		} else {
			// Fallback to direct fields if no primary NOK in list
			newMember.NextOfKinFullName = req.NextOfKinFullName
			newMember.NextOfKinPhone = utils.NormalizePhone(req.NextOfKinPhone)
		}
		if err := tx.Save(newMember).Error; err != nil { // Save updated member with NOK details
			return err
		}

		// Handle Member Bank Account
		return s.handleMemberBankAccount(tx, newMember.ID, req.BankID, req.BankBranch, req.AccountNo, req.AccountName, userID)
	}) // End of transaction
	return newMember, err // Return newMember and the error from the transaction
}

// Get all users
func (s *MemberService) GetMembers(page, limit int, memberNo, primaryPhone, memberTypeID, routeID, q string) ([]dtos.MemberResponse, int64, error) {
	var results []dtos.MemberResponse
	var total int64

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Base query for data
	baseQuery := `
		SELECT 
			m.id, m.member_no, 
			m.id_no, 
			m.member_type_id, mt.name AS member_type_name,
			m.first_name, m.last_name, m.other_names, m.route_id, r.route_name,
			m.date_of_birth, m.id_no, m.gender, m.birth_city,
			m.primary_phone, m.secondary_phone, m.email, m.number_of_cows,
			m.id_front_photo, m.id_back_photo, m.passport_photo, m.id_date_of_issue,
			m.tax_number, m.marital_status, m.title,
			m.next_of_kin_full_name, m.next_of_kin_phone, m.status, m.date_registered,
			m.created_at, m.updated_at
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
	`
	// Base query for count
	baseCountQuery := `
		SELECT COUNT(*)
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
	`

	var args []interface{}

	whereClauses := []string{"m.deleted_at IS NULL"}

	if memberNo != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+memberNo+"%")
	}
	if primaryPhone != "" {
		primaryPhone = utils.NormalizePhone(primaryPhone)
		whereClauses = append(whereClauses, "m.primary_phone LIKE ?")
		args = append(args, "%"+primaryPhone+"%")
	}
	if memberTypeID != "" {
		whereClauses = append(whereClauses, "m.member_type_id = ?")
		args = append(args, memberTypeID)
	}
	if routeID != "" {
		whereClauses = append(whereClauses, "m.route_id = ?")
		args = append(args, routeID)
	}
	if q != "" {
		whereClauses = append(whereClauses, "(m.first_name LIKE ? OR m.last_name LIKE ? OR m.id_no LIKE ?)")
		search := "%" + q + "%"
		args = append(args, search, search, search)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	// Execute count query first to get filtered total
	if err := db.DB.Raw(baseCountQuery+whereSql, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Execute main data query
	fullQuery := baseQuery + whereSql + " ORDER BY m.id DESC LIMIT ? OFFSET ?"
	if err := db.DB.Raw(fullQuery, append(args, limit, offset)...).Scan(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

// Get single user
func (s *MemberService) GetMember(id string) (*dtos.MemberResponse, error) {
	var result dtos.MemberResponse
	query := `
		SELECT 
			m.id, m.member_no, m.member_type_id, mt.name AS member_type_name,
			m.first_name, m.last_name, m.other_names, m.route_id, r.route_name,
			m.dob, m.id_no, m.gender, m.birth_city,
			m.primary_phone, m.secondary_phone, m.email, m.number_of_cows,
			m.id_front_photo, m.id_back_photo, m.passport_photo, m.id_date_of_issue,
			m.tax_number, m.marital_status, m.title,
			m.next_of_kin_full_name, m.next_of_kin_phone, m.status, m.date_registered,
			m.created_at, m.updated_at
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
		WHERE m.id = ? AND m.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Fetch associated Next of Kins
	db.DB.Raw(`
		SELECT 
			mnok.*, 
			m.member_no, 
			CONCAT(m.first_name, ' ', m.last_name, ' ', COALESCE(m.other_names, '')) AS member_full_name
		FROM member_next_of_kins mnok
		JOIN member_registrations m ON mnok.member_id = m.id
		WHERE mnok.member_id = ? AND mnok.deleted_at IS NULL
	`, result.ID).Scan(&result.NextOfKins)

	// Fetch associated Bank Accounts
	db.DB.Raw(`
		SELECT 
			mba.id, mba.member_id, m.member_no, m.first_name, m.last_name,
			mba.bank_id, b.bank_name, mba.bank_branch_id,
			mba.account_number, mba.account_name, mba.status,
			mba.created_at, mba.updated_at
		FROM member_bank_accounts mba
		LEFT JOIN member_registrations m ON mba.member_id = m.id
		LEFT JOIN banks b ON mba.bank_id = b.id
		WHERE mba.member_id = ? AND mba.deleted_at IS NULL
	`, result.ID).Scan(&result.BankAccounts)

	return &result, nil
}

func (s *MemberService) UpdateMember(
	id string,
	req dtos.UpdateMemberRequest,
	userID uint64,
	idFront *multipart.FileHeader,
	idBack *multipart.FileHeader,
	passport *multipart.FileHeader,
) error {
	var member models.Member

	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}

	// Handle File Updates (only if provided)
	if idFront != nil {
		path, err := utils.SaveFile(idFront, "members")
		if err == nil {
			member.IdFrontPhoto = path
		}
	}
	if idBack != nil {
		path, err := utils.SaveFile(idBack, "members")
		if err == nil {
			member.IdBackPhoto = path
		}
	}
	if passport != nil {
		path, err := utils.SaveFile(passport, "members")
		if err == nil {
			member.PassportPhoto = path
		}
	}

	// Update Fields mapping according to CreateMember
	member.MemberTypeID = req.MemberTypeID //
	member.FirstName = req.FirstName
	member.LastName = req.LastName
	member.OtherNames = req.OtherNames
	member.RouteID = req.RouteID
	member.DateOfBirth = utils.ParseDate(req.DOB)
	member.IDNo = req.IDNo
	member.Gender = req.Gender
	member.BirthCity = req.BirthCity
	member.PrimaryPhone = utils.NormalizePhone(req.PrimaryPhone)
	member.SecondaryPhone = utils.NormalizePhone(req.SecondaryPhone)
	member.Email = req.Email
	member.NumberOfCows = req.NumberOfCows
	if req.IDDateOfIssue != "" {
		t := utils.ParseDate(req.IDDateOfIssue)
		member.IdDateOfIssue = &t
	} else {
		member.IdDateOfIssue = nil
	}
	member.TaxNumber = req.TaxNumber
	member.MaritalStatus = req.MaritalStatus
	member.Title = req.Title

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&member).Error; err != nil {
			return err
		}

		// Handle NextOfKins: Delete existing and re-create
		if err := tx.Where("member_id = ?", member.ID).Delete(&models.MemberNextOfKin{}).Error; err != nil {
			return err
		}

		var primaryNOK *dtos.CreateMemberNextOfKinRequest
		if len(req.NextOfKins) > 0 {
			for i := range req.NextOfKins { // Correctly iterate over slice to avoid referencing loop variable
				nokReq := &req.NextOfKins[i]
				nok := &models.MemberNextOfKin{
					BaseModel:              models.BaseModel{CreatedBy: userID}, // Assuming userID for updates too
					MemberID:               member.ID,
					FullName:               nokReq.FullName,
					Relationship:           nokReq.Relationship,
					PhoneNumber:            nokReq.PhoneNumber,
					AlternativePhoneNumber: nokReq.AlternativePhoneNumber,
					EmailAddress:           nokReq.EmailAddress,
					NationalIDNo:           nokReq.NationalIDNo,
					PostalAddress:          nokReq.PostalAddress,
					PhysicalAddress:        nokReq.PhysicalAddress,
					Occupation:             nokReq.Occupation,
					IsPrimary:              nokReq.IsPrimary,
					Status:                 nokReq.Status,
					Remarks:                nokReq.Remarks,
				}
				if err := tx.Create(nok).Error; err != nil {
					return err
				}
				if nokReq.IsPrimary && primaryNOK == nil {
					primaryNOK = nokReq // Assign value, not address of loop variable
				}
			}
		}

		// Update primary next of kin on the member model
		if primaryNOK != nil {
			member.NextOfKinFullName = primaryNOK.FullName
			if primaryNOK.PhoneNumber != nil {
				member.NextOfKinPhone = *primaryNOK.PhoneNumber
			}
		} else {
			// If no primary NOK in list, clear direct fields
			member.NextOfKinFullName = ""
			member.NextOfKinPhone = ""
		}
		if err := tx.Save(&member).Error; err != nil {
			return err
		}

		// Handle Member Bank Account (assuming one primary bank account per member for simplicity)
		return s.handleMemberBankAccount(tx, member.ID, req.BankID, req.BankBranch, req.AccountNo, req.AccountName, userID)
	})
}

func (s *MemberService) DeleteMember(id string) error {
	var member models.Member

	// ensure record exists
	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}

	return db.DB.Delete(&member).Error
}

func (s *MemberService) SuspendMember(id string) error {
	var member models.Member
	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&member).Update("status", "SUSPENDED").Error
}

// ImportMembers bulk imports members from CSV, XLS, or XLSX files.
func (s *MemberService) ImportMembers(file *multipart.FileHeader, userID uint64) error {
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

	// The actual processing will happen in a goroutine
	go s.processMemberRowsInBackground(data, userID)

	// Return immediately to the API caller
	return nil
}

func (s *MemberService) processMemberRowsInBackground(data [][]string, userID uint64) {
	totalRows := len(data) - 1
	if totalRows < 0 {
		return
	}

	log.Printf("[MemberService.processMemberRowsInBackground] Starting background member import for %d rows.", totalRows)

	importID := uint64(utils.Now().UnixNano())

	var wg sync.WaitGroup
	jobs := make(chan []string, totalRows)
	errorChan := make(chan error, totalRows)

	numWorkers := runtime.NumCPU() * 2 // Use a reasonable number of workers
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Start worker goroutines
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[MemberService] Worker panicked during member import for row: %v, error: %v", row, r)
							db.DB.Create(&models.ImportError{
								BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic during import: %v", r),
								ImportId:  importID,
							})
							errorChan <- fmt.Errorf("panic during import for row: %v", r)
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						// Expected columns based on: Member No, First Name, Last Name, Other Names, ID No, Gender, Date of Birth, Date Registered, Member Type, Route, Primary Phone, Secondary Phone, Email, Bank, Bank Branch, Account No, Account Name, Number of Cows
						if len(row) < 18 {
							return fmt.Errorf("row has insufficient columns (%d < 18)", len(row))
						}

						idNo := strings.TrimSpace(row[4])
						if idNo == "" {
							return fmt.Errorf("ID number is empty for row")
						}

						// Idempotency: skip if member already exists
						var count int64
						tx.Model(&models.Member{}).Where("id_no = ?", idNo).Count(&count)
						if count > 0 {
							return fmt.Errorf("member with ID number %s already exists", idNo)
						}

						// Lookups for Route, Bank, and Member Type
						routeVal := strings.TrimSpace(row[9])
						var route models.Route
						tx.Where("route_code = ? OR name = ?", routeVal, routeVal).First(&route)

						bankName := strings.TrimSpace(row[13])
						var bank models.Bank
						tx.Where("bank_name LIKE ?", "%"+bankName+"%").First(&bank)

						memberTypeName := strings.TrimSpace(row[8])
						var memberType models.MemberType
						if err := tx.Where("name = ?", memberTypeName).First(&memberType).Error; err != nil {
							// Fallback to first member type if none matches by name
							tx.First(&memberType)
						}

						dob := utils.ParseFlexibleDate(row[6])
						joiningDate := utils.ParseFlexibleDate(row[7])
						numberOfCows, _ := strconv.Atoi(strings.TrimSpace(row[17]))

						member := models.Member{
							BaseModel:      models.BaseModel{CreatedBy: userID},
							MemberNo:       strings.TrimSpace(row[0]),
							FirstName:      strings.TrimSpace(row[1]),
							LastName:       strings.TrimSpace(row[2]),
							OtherNames:     strings.TrimSpace(row[3]),
							IDNo:           idNo,
							DateOfBirth:    dob,
							PrimaryPhone:   utils.NormalizePhone(row[10]),
							SecondaryPhone: utils.NormalizePhone(row[11]),
							Email:          strings.TrimSpace(row[12]),
							Gender:         strings.ToUpper(strings.TrimSpace(row[5])),
							RouteID:        route.ID,
							MemberTypeID:   memberType.ID,
							DateRegistered: joiningDate,
							NumberOfCows:   numberOfCows,
							Status:         "PENDING",
						}

						if err := tx.Create(&member).Error; err != nil {
							return err
						}

						// Attachment of Bank Account if provided
						accNo := strings.TrimSpace(row[15])
						fullName := strings.TrimSpace(row[16])
						if accNo != "" && bank.ID != 0 {
							bankAcc := models.MemberBankAccount{
								BaseModel:     models.BaseModel{CreatedBy: userID},
								MemberID:      member.ID,
								BankID:        bank.ID,
								AccountNumber: accNo,
								AccountName:   fullName,
								Status:        "ACTIVE",
							}
							if err := tx.Create(&bankAcc).Error; err != nil {
								return err
							}
						}
						return nil
					})

					if err != nil {
						log.Printf("[MemberService.processMemberRowsInBackground] Error processing row: %v, error: %v", row, err)
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

	// Send rows to job channel, skipping header
	for i := 1; i < len(data); i++ {
		jobs <- data[i]
	}
	close(jobs)

	wg.Wait() // Wait for all workers to finish
	close(errorChan)

	failedCount := 0
	for err := range errorChan {
		if err != nil {
			failedCount++
		}
	}

	message := fmt.Sprintf("Member import completed. Success: %d, Failed: %d out of %d records.", totalRows-failedCount, failedCount, totalRows)
	notificationType := "SUCCESS"
	errorLink := ""

	if failedCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/members/import-errors/%d", importID)
		log.Printf("[MemberService.processMemberRowsInBackground] Member import completed with %d failures.", failedCount)
	} else if totalRows == 0 {
		message = "Member import completed. No records were processed."
		notificationType = "SUCCESS"
		log.Printf("[MemberService.processMemberRowsInBackground] Member import completed. No records were processed.")
	} else {
		log.Printf("[MemberService.processMemberRowsInBackground] Member import completed successfully.")
	}

	_, err := s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Member Import Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
	})
	if err != nil {
		log.Printf("[MemberService.processMemberRowsInBackground] Failed to create UI notification: %v", err)
	}
}

func (s *MemberService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}

// handleMemberBankAccount creates or updates a member's bank account.
func (s *MemberService) handleMemberBankAccount(tx *gorm.DB, memberID, bankID uint64, bankBranch, accountNo, accountName string, userID uint64) error {
	if bankID == 0 || accountNo == "" || accountName == "" {
		// No bank details provided, or incomplete. Do nothing.
		return nil
	}

	var existingAccount models.MemberBankAccount
	result := tx.Where("member_id = ?", memberID).First(&existingAccount)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Create new account
		newAccount := &models.MemberBankAccount{
			BaseModel:     models.BaseModel{CreatedBy: userID},
			MemberID:      memberID,
			BankID:        bankID,
			AccountNumber: accountNo,
			AccountName:   accountName,
			Status:        "ACTIVE", // Default status
		}
		return tx.Create(newAccount).Error
	} else {
		// Update existing account
		existingAccount.BankID = bankID
		existingAccount.AccountNumber = accountNo
		existingAccount.AccountName = accountName
		existingAccount.UpdatedBy = userID
		return tx.Save(&existingAccount).Error
	}
}

// ExportMembers initiates a background process to export member data to a CSV file.
// It takes filtering parameters similar to GetMembers and the userID for notification.
func (s *MemberService) ExportMembers(userID uint64, memberNo, primaryPhone, memberTypeID, routeID, gender, status, reportType string) error {
	go s.processMemberExportInBackground(userID, memberNo, primaryPhone, memberTypeID, routeID, gender, status, reportType)
	return nil
}

// memberExportQueryResult is a helper struct to hold the results of the member export query,
// including joined fields from related tables like member types, routes, banks, and bank branches.
type memberExportQueryResult struct {
	models.Member
	MemberTypeName string `gorm:"column:member_type_name"`
	RouteName      string `gorm:"column:route_name"`
	BankName       string `gorm:"column:bank_name"`
	BankBranchName string `gorm:"column:bank_branch_name"`
	AccountNumber  string `gorm:"column:account_number"`
	AccountName    string `gorm:"column:account_name"`
}

// processMemberExportInBackground performs the actual data export in a separate goroutine.
// It queries the database, generates a CSV file, saves it, and sends a UI notification.
func (s *MemberService) processMemberExportInBackground(userID uint64, memberNo, primaryPhone, memberTypeID, routeID, gender, status, reportType string) {
	var results []memberExportQueryResult

	// Base SQL query to fetch member details along with their member type, route, and primary bank account information.
	baseQuery := `
		SELECT
			m.id, m.member_no, m.member_type_id, mt.name AS member_type_name,
			m.first_name, m.last_name, m.other_names, m.route_id, r.route_name,
			m.date_of_birth, m.id_no, m.gender, m.birth_city,
			m.primary_phone, m.secondary_phone, m.email, m.number_of_cows,
			m.id_front_photo, m.id_back_photo, m.passport_photo, m.id_date_of_issue,
			m.tax_number, m.marital_status, m.title,
			m.next_of_kin_full_name, m.next_of_kin_phone, m.status, m.date_registered,
			m.created_at, m.updated_at,
			b.bank_name, bb.branch_name AS bank_branch_name, mba.account_number, mba.account_name
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
		LEFT JOIN member_bank_accounts mba ON m.id = mba.member_id AND mba.deleted_at IS NULL
		LEFT JOIN banks b ON mba.bank_id = b.id
		LEFT JOIN bank_branches bb ON mba.bank_branch_id = bb.id
	`

	var args []interface{}
	whereClauses := []string{"m.deleted_at IS NULL"}

	// Apply filters based on the request parameters.
	if memberNo != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+memberNo+"%")
	}
	if primaryPhone != "" {
		primaryPhone = utils.NormalizePhone(primaryPhone)
		whereClauses = append(whereClauses, "m.primary_phone LIKE ?")
		args = append(args, "%"+primaryPhone+"%")
	}
	if memberTypeID != "" {
		whereClauses = append(whereClauses, "m.member_type_id = ?")
		args = append(args, memberTypeID)
	}
	if routeID != "" {
		whereClauses = append(whereClauses, "m.route_id = ?")
		args = append(args, routeID)
	}
	if gender != "" {
		whereClauses = append(whereClauses, "m.gender = ?")
		args = append(args, gender)
	}
	if status != "" {
		whereClauses = append(whereClauses, "m.status = ?")
		args = append(args, status)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	// Execute the query to retrieve all matching member records.
	err := db.DB.Raw(baseQuery+whereSql+" ORDER BY m.id DESC", args...).Scan(&results).Error
	if err != nil {
		log.Printf("[MemberService.processMemberExportInBackground] Error querying members for export: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Member Export Failed",
			Message:          fmt.Sprintf("Failed to export member data: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	// Generate file data based on requested type
	var fileData []byte
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateMemberPDF(results, routeID, gender, status)
	} else {
		fileData, err = s.generateMemberCSV(results)
	}

	if err != nil {
		log.Printf("[MemberService.processMemberExportInBackground] Error generating %s: %v", ext, err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Member Export Failed",
			Message:          fmt.Sprintf("Failed to generate export file: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	// Create the directory for exports if it doesn't exist.
	exportDir := "./storage/exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Printf("[MemberService.processMemberExportInBackground] Failed to create export directory: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Member Export Failed",
			Message:          fmt.Sprintf("Failed to create export directory: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	// Generate a unique filename for the exported CSV.
	filename := fmt.Sprintf("members_export_%d.%s", utils.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[MemberService.processMemberExportInBackground] Failed to save exported CSV: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Member Export Failed",
			Message:          fmt.Sprintf("Failed to save exported file: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	// Send a UI notification to the user with a link to download the generated file.
	downloadLink := fmt.Sprintf("/api/members/export/download/%s", filename)
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Member Export Complete",
		Message:          "Your member data export is ready for download.",
		NotificationType: "SUCCESS",
		DownloadLink:     downloadLink,
	})
	log.Printf("[MemberService.processMemberExportInBackground] Member export completed successfully. File: %s", filename)
}

func (s *MemberService) generateMemberCSV(results []memberExportQueryResult) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Define the CSV header, matching the import columns for consistency.
	headers := []string{
		"Member No", "First Name", "Last Name", "Other Names", "ID No", "Gender",
		"Date of Birth", "Date Registered", "Member Type", "Route", "Primary Phone",
		"Secondary Phone", "Email", "Bank", "Bank Branch", "Account No", "Account Name",
		"Number of Cows",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Iterate through the query results and write each member's data to the CSV.
	for _, member := range results {
		row := []string{
			member.MemberNo,
			member.FirstName,
			member.LastName,
			member.OtherNames,
			member.IDNo,
			member.Gender,
			utils.FormatDate(member.DateOfBirth),
			utils.FormatDate(member.DateRegistered),
			member.MemberTypeName,
			member.RouteName,
			member.PrimaryPhone,
			member.SecondaryPhone,
			member.Email,
			member.BankName,
			member.BankBranchName,
			member.AccountNumber,
			member.AccountName,
			strconv.Itoa(member.NumberOfCows),
		}
		if err := writer.Write(row); err != nil {
			log.Printf("[MemberService.generateMemberCSV] Error writing CSV row for member %s: %v", member.MemberNo, err)
		}
	}
	writer.Flush()

	return buf.Bytes(), writer.Error()
}

func (s *MemberService) generateMemberPDF(results []memberExportQueryResult, routeID, gender, status string) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
		Phone          string `gorm:"column:phone"`
		Email          string `gorm:"column:email"`
	}
	db.DB.Table("organization_details").First(&org)

	routeName := "ALL"
	if routeID != "" {
		db.DB.Table("routes").Where("id = ?", routeID).Select("route_name").Scan(&routeName)
	}

	genderVal := "ALL"
	if gender != "" {
		genderVal = gender
	}

	statusVal := "ALL"
	if status != "" {
		statusVal = status
	}

	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape orientation
	pdf.AddPage()

	// Standard Professional Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Phone: %s | Email: %s", org.Phone, org.Email), "", 1, "C", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "MEMBER REGISTER", "", 1, "C", false, 0, "")

	// Summary Line
	pdf.SetFont("Arial", "I", 9)
	filterSummary := fmt.Sprintf("Showing (%d) Records for; Route %s, Gender %s, Status %s [filters applied]",
		len(results), routeName, genderVal, statusVal)
	pdf.CellFormat(0, 8, filterSummary, "", 1, "L", false, 0, "")

	pdf.Ln(2)

	// Define column headers and widths
	pdf.SetFont("Arial", "B", 8)
	headers := []string{"Member No", "Full Name", "ID No", "Phone", "Member Type", "Route", "Status"}
	widths := []float64{25, 65, 25, 30, 45, 45, 25}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Add data rows
	pdf.SetFont("Arial", "", 8)
	for _, m := range results {
		fullName := strings.TrimSpace(fmt.Sprintf("%s %s %s", m.FirstName, m.LastName, m.OtherNames))
		pdf.CellFormat(25, 8, m.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(65, 8, fullName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, m.IDNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, m.PrimaryPhone, "1", 0, "L", false, 0, "")
		pdf.CellFormat(45, 8, m.MemberTypeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(45, 8, m.RouteName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, m.Status, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

// ExportAGMReport initiates a background process to export a printable AGM attendance register.
func (s *MemberService) ExportAGMReport(userID uint64, memberNo, primaryPhone, memberTypeID, routeID, gender, status string) error {
	go s.processAGMReportInBackground(userID, memberNo, primaryPhone, memberTypeID, routeID, gender, status)
	return nil
}

func (s *MemberService) processAGMReportInBackground(userID uint64, memberNo, primaryPhone, memberTypeID, routeID, gender, status string) {
	var results []memberExportQueryResult

	// Query optimized for attendance (Alphabetical sorting)
	baseQuery := `
		SELECT
			m.id, m.member_no, m.id_no, m.first_name, m.last_name, m.other_names,
			m.primary_phone, m.status, r.route_name
		FROM member_registrations m
		LEFT JOIN routes r ON m.route_id = r.id
	`

	var args []interface{}
	whereClauses := []string{"m.deleted_at IS NULL"}

	if memberNo != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+memberNo+"%")
	}
	if primaryPhone != "" {
		primaryPhone = utils.NormalizePhone(primaryPhone)
		whereClauses = append(whereClauses, "m.primary_phone LIKE ?")
		args = append(args, "%"+primaryPhone+"%")
	}
	if memberTypeID != "" {
		whereClauses = append(whereClauses, "m.member_type_id = ?")
		args = append(args, memberTypeID)
	}
	if routeID != "" {
		whereClauses = append(whereClauses, "m.route_id = ?")
		args = append(args, routeID)
	}
	if gender != "" {
		whereClauses = append(whereClauses, "m.gender = ?")
		args = append(args, gender)
	}
	if status != "" {
		whereClauses = append(whereClauses, "m.status = ?")
		args = append(args, status)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	// Sort by Name for easier lookup on a printed list
	err := db.DB.Raw(baseQuery+whereSql+" ORDER BY m.first_name ASC, m.last_name ASC", args...).Scan(&results).Error
	if err != nil {
		log.Printf("[MemberService.processAGMReportInBackground] Error querying members: %v", err)
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "AGM Report Failed",
			Message:          fmt.Sprintf("Failed to export AGM report: %v", err),
			NotificationType: "ERROR",
		})
		return
	}

	fileData, err := s.generateAGMReportPDF(results)
	if err != nil {
		log.Printf("[MemberService.processAGMReportInBackground] Error generating PDF: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("agm_attendance_%d.pdf", utils.Now().UnixNano())
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[MemberService.processAGMReportInBackground] Failed to save file: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "AGM Printable Report Ready",
		Message:          "Your AGM attendance register is ready for download.",
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/members/agm-report/download/%s", filename),
	})
}

func (s *MemberService) generateAGMReportPDF(results []memberExportQueryResult) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
	}
	db.DB.Table("organization_details").First(&org)

	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape gives more room for signatures
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "AGM ATTENDANCE REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"Member No", "ID No", "Full Names", "Phone", "____________"}
	widths := []float64{30, 35, 90, 45, 65}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 10, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, m := range results {
		fullName := strings.TrimSpace(fmt.Sprintf("%s %s %s", m.FirstName, m.LastName, m.OtherNames))
		pdf.CellFormat(widths[0], 12, m.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 12, m.IDNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 12, fullName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 12, m.PrimaryPhone, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 12, "", "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
