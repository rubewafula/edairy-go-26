package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"runtime"

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
	req dtos.CreateMemberRequest) (*models.Member, error) {

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
		MemberNo:       memberNo,
		MemberTypeID:   req.MemberTypeID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		OtherNames:     req.OtherNames,
		RouteID:        req.RouteID,
		DateOfBirth:    utils.ParseDate(req.DOB),
		IDNumber:       req.IDNo,
		Gender:         req.Gender,
		BirthCity:      req.BirthCity,
		PrimaryPhone:   utils.NormalizePhone(req.PrimaryPhone),
		SecondaryPhone: utils.NormalizePhone(req.SecondaryPhone),
		Email:          req.Email,
		NumberOfCows:   req.NumberOfCows,
		IdFrontPhoto:   idFrontPath,
		IdBackPhoto:    idBackPath,
		PassportPhoto:  passportPath,
		IdDateOfIssue:  utils.ParseDate(req.IDDateOfIssue),
		TaxNumber:      req.TaxNumber,
		MaritalStatus:  req.MaritalStatus,
		Title:          req.Title,

		NextOfKinFullName: req.NextOfKinFullName,
		NextOfKinPhone:    utils.NormalizePhone(req.NextOfKinPhone),
		Status:            "PENDING",
		DateRegistered:    time.Now(),
	}

	if err := db.DB.Create(newMember).Error; err != nil {
		return nil, err
	}

	return newMember, nil
}

// Get all users
func (s *MemberService) GetMembers(page, limit int, memberNo, primaryPhone, memberTypeID, routeID string) ([]dtos.MemberResponse, int64, error) {
	var results []dtos.MemberResponse
	var total int64

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Base query for data
	baseQuery := `
		SELECT 
			m.id, m.member_no, m.member_type_id, mt.name AS member_type_name,
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
			m.date_of_birth, m.id_no, m.gender, m.birth_city,
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
	return &result, nil
}

func (s *MemberService) UpdateMember(
	id string,
	req dtos.UpdateMemberRequest,
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
	member.MemberTypeID = req.MemberTypeID
	member.FirstName = req.FirstName
	member.LastName = req.LastName
	member.OtherNames = req.OtherNames
	member.RouteID = req.RouteID
	member.DateOfBirth = utils.ParseDate(req.DOB)
	member.IDNumber = req.IDNo
	member.Gender = req.Gender
	member.BirthCity = req.BirthCity
	member.PrimaryPhone = utils.NormalizePhone(req.PrimaryPhone)
	member.SecondaryPhone = utils.NormalizePhone(req.SecondaryPhone)
	member.Email = req.Email
	member.NumberOfCows = req.NumberOfCows
	member.IdDateOfIssue = utils.ParseDate(req.IDDateOfIssue)
	member.TaxNumber = req.TaxNumber
	member.MaritalStatus = req.MaritalStatus
	member.Title = req.Title
	member.NextOfKinFullName = req.NextOfKinFullName
	member.NextOfKinPhone = utils.NormalizePhone(req.NextOfKinPhone)

	return db.DB.Save(&member).Error
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
	log.Printf("[MemberService.processMemberRowsInBackground] Starting background member import for %d rows.", len(data)-1)

	var wg sync.WaitGroup
	jobs := make(chan []string, len(data)-1) // Buffer jobs channel
	errorChan := make(chan error, len(data)-1)

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
							db.DB.Create(&models.MemberImportError{
								BaseModel: models.BaseModel{CreatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic during import: %v", r),
							})
							errorChan <- fmt.Errorf("panic during import for row: %v", r)
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						// Skip header and rows that don't have enough columns
						if len(row) < 19 {
							return fmt.Errorf("row has insufficient columns (%d < 19)", len(row))
						}

						idNo := strings.TrimSpace(row[3])
						if idNo == "" {
							return fmt.Errorf("ID number is empty for row")
						}

						// Idempotency: skip if member already exists
						var count int64
						tx.Model(&models.Member{}).Where("id_no = ?", idNo).Count(&count)
						if count > 0 {
							return fmt.Errorf("member with ID number %s already exists", idNo)
						}

						fullName := strings.TrimSpace(row[2])
						nameParts := strings.Fields(fullName)
						firstName, lastName, otherNames := "", "", ""
						if len(nameParts) > 0 {
							firstName = nameParts[0]
							if len(nameParts) > 1 {
								lastName = nameParts[len(nameParts)-1]
								if len(nameParts) > 2 {
									otherNames = strings.Join(nameParts[1:len(nameParts)-1], " ")
								}
							}
						}

						// Lookups for Route and Bank
						routeVal := strings.TrimSpace(row[18])
						var route models.Route
						tx.Where("code = ? OR name = ?", routeVal, routeVal).First(&route)

						bankName := strings.TrimSpace(row[16])
						var bank models.Bank
						tx.Where("bank_name LIKE ?", "%"+bankName+"%").First(&bank)

						var memberType models.MemberType
						// Fallback to first member type if none matches
						tx.First(&memberType)

						dob := utils.ParseFlexibleDate(row[5])
						joiningDate := utils.ParseFlexibleDate(row[4])

						member := models.Member{
							BaseModel:      models.BaseModel{CreatedBy: userID},
							MemberNo:       strings.TrimSpace(row[0]),
							Title:          strings.TrimSpace(row[1]),
							FirstName:      firstName,
							LastName:       lastName,
							OtherNames:     otherNames,
							IDNumber:       idNo,
							DateOfBirth:    dob,
							PrimaryPhone:   utils.NormalizePhone(row[8]),
							Email:          strings.TrimSpace(row[9]),
							Status:         strings.ToUpper(strings.TrimSpace(row[10])),
							Gender:         strings.ToUpper(strings.TrimSpace(row[11])),
							BirthCity:      strings.TrimSpace(row[14]),
							RouteID:        route.ID,
							MemberTypeID:   memberType.ID,
							DateRegistered: joiningDate,
						}

						if err := tx.Create(&member).Error; err != nil {
							return err
						}

						// Attachment of Bank Account if provided
						accNo := strings.TrimSpace(row[17])
						if accNo != "" && bank.ID != 0 {
							bankAcc := models.MemberBankAccount{
								BaseModel:     models.BaseModel{CreatedBy: userID},
								MemberID:      member.ID,
								BankID:        bank.ID,
								AccountNumber: accNo,
								AccountName:   fullName,
								Status:        "ACTIVE",
							}
							tx.Create(&bankAcc)
						}
						return nil
					})

					if err != nil {
						log.Printf("[MemberService.processMemberRowsInBackground] Error processing row: %v, error: %v", row, err)
						db.DB.Create(&models.MemberImportError{
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

	if failedCount > 0 {
		log.Printf("[MemberService.processMemberRowsInBackground] Member import completed with %d failures.", failedCount)
	} else {
		log.Printf("[MemberService.processMemberRowsInBackground] Member import completed successfully.")
	}

	// Create and emit UI notification
	var notificationTitle string
	var notificationMessage string
	if failedCount == 0 {
		notificationTitle = "Member Import Successful"
		notificationMessage = "All members from your file were imported successfully."
	} else {
		notificationTitle = "Member Import Completed with Errors"
		notificationMessage = fmt.Sprintf("Member import finished with %d failures. Please check the import error logs for details.", failedCount)
	}

	_, err := s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            notificationTitle,
		Message:          notificationMessage,
		NotificationType: "MEMBER_IMPORT_STATUS",
	})
	if err != nil {
		log.Printf("[MemberService.processMemberRowsInBackground] Failed to create UI notification: %v", err)
	}
}
