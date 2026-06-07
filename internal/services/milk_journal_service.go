package services

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
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

type MilkJournalService struct {
	notificationService *UINotificationService
}

func NewMilkJournalService() *MilkJournalService {
	return &MilkJournalService{
		notificationService: NewUINotificationService(),
	}
}

func (s *MilkJournalService) CreateMilkJournal(req dtos.CreateMilkJournalRequest) (*models.MilkJournal, error) {

	var journalDate *time.Time

	if req.JournalDate != "" {
		jd := utils.ParseDate(req.JournalDate)
		journalDate = &jd
	} else {
		journalDate = nil
	}

	journal := &models.MilkJournal{
		Journal:             req.Journal,
		JournalDate:         journalDate,
		MilkDeliveryShiftID: req.MilkDeliveryShiftID,
		RouteID:             req.RouteID,
		UserID:              req.UserID,
		TransporterID:       req.TransporterID,
		Confirmed:           false,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the Booklet Header (Journal)
		if err := tx.Create(journal).Error; err != nil {
			return err
		}

		// 2. Process each Page (Batch)
		for _, bReq := range req.Batches {
			batch := &models.MilkJournalBatch{
				MilkJournalID: journal.ID,
				BatchNo:       bReq.BatchNo,
			}
			if err := tx.Create(batch).Error; err != nil {
				return err
			}

			// 3. Process each line item (Entry) on the page
			for _, eReq := range bReq.Entries {
				entry := &models.MilkJournalEntry{
					MemberID:           eReq.MemberID,
					MilkJournalID:      journal.ID,
					MilkJournalBatchID: batch.ID,
					Status:             "PENDING",
					Quantity:           eReq.Quantity,
					RouteCenterID:      eReq.RouteCenterID,
					CanID:              eReq.CanID,
				}
				if err := tx.Create(entry).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return journal, nil
}

func (s *MilkJournalService) GetMilkJournals(page, limit int) ([]dtos.MilkJournalResponse, int64, error) {
	var journals []dtos.MilkJournalResponse
	var total int64
	db.DB.Model(&models.MilkJournal{}).Count(&total)
	offset := (page - 1) * limit

	query := `	
		SELECT 
					mj.id,
					mj.journal,
					mj.journal_date,
					mj.milk_delivery_shift_id,
					mds.name AS milk_delivery_shift,
					mj.route_id,
					r.route_name,
					mj.user_id,
					mj.transporter_id,
					mj.confirmed,
					mj.created_at,
					mj.updated_at,
					mjb.batch_no,

					COALESCE(e.entries_count, 0) AS entries_count,
					COALESCE(e.total_litres, 0) AS collections

				FROM milk_journals mj

				LEFT JOIN milk_journal_batches mjb 
					ON mj.id = mjb.milk_journal_id

				LEFT JOIN milk_delivery_shifts mds 
					ON mj.milk_delivery_shift_id = mds.id

				LEFT JOIN routes r 
					ON mj.route_id = r.id

				LEFT JOIN (
					SELECT 
						milk_journal_id,
						COUNT(*) AS entries_count,
						SUM(quantity) AS total_litres
					FROM milk_journal_entries
					WHERE deleted_at IS NULL
					GROUP BY milk_journal_id
				) e ON e.milk_journal_id = mj.id

				WHERE mj.deleted_at IS NULL
				LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&journals).Error
	return journals, total, err
}

func (s *MilkJournalService) GetMilkJournal(id string) (*dtos.MilkJournalResponse, error) {
	var journal dtos.MilkJournalResponse
	query := `
		SELECT 
			mj.id,
			mj.journal,
			mj.journal_date,
			mj.milk_delivery_shift_id,
			mds.name AS milk_delivery_shift,
			mj.route_id,
			r.route_name,
			mj.user_id,
			mj.transporter_id,
			mj.confirmed,
			mj.created_at,
			mj.updated_at,
			mjb.batch_no,

			COALESCE(e.entries_count, 0) AS entries_count,
			COALESCE(e.total_litres, 0) AS collections

		FROM milk_journals mj

		LEFT JOIN milk_journal_batches mjb 
			ON mj.id = mjb.milk_journal_id

		LEFT JOIN milk_delivery_shifts mds
			ON mj.milk_delivery_shift_id = mds.id

		LEFT JOIN routes r 
			ON mj.route_id = r.id

		LEFT JOIN (
			SELECT 
				milk_journal_id,
				COUNT(*) AS entries_count,
				SUM(quantity) AS total_litres
			FROM milk_journal_entries
			WHERE deleted_at IS NULL
			GROUP BY milk_journal_id
		) e ON e.milk_journal_id = mj.id

		WHERE mj.id = ? AND mj.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&journal).Error
	if err != nil {
		return nil, err
	}

	if journal.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &journal, nil
}

func (s *MilkJournalService) UpdateMilkJournal(id string, req dtos.UpdateMilkJournalRequest) error {
	var journal models.MilkJournal
	if err := db.DB.First(&journal, id).Error; err != nil {
		return err
	}

	var journalDate *time.Time

	if req.JournalDate != "" {
		jd := utils.ParseFlexibleDate(req.JournalDate)
		journalDate = &jd
	} else {
		journalDate = nil
	}
	journal.Journal = req.Journal
	journal.JournalDate = journalDate
	journal.MilkDeliveryShiftID = req.MilkDeliveryShiftID
	journal.RouteID = req.RouteID
	journal.UserID = req.UserID
	journal.TransporterID = req.TransporterID
	journal.Confirmed = req.Confirmed

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&journal).Error; err != nil {
			return err
		}

		// Delete old entries and batches to replace with updated payload
		if err := tx.Where("milk_journal_id = ?", journal.ID).Delete(&models.MilkJournalEntry{}).Error; err != nil {
			return err
		}
		if err := tx.Where("milk_journal_id = ?", journal.ID).Delete(&models.MilkJournalBatch{}).Error; err != nil {
			return err
		}

		// Process updated batches and entries
		for _, bReq := range req.Batches {
			batch := &models.MilkJournalBatch{
				MilkJournalID: journal.ID,
				BatchNo:       bReq.BatchNo,
			}
			if err := tx.Create(batch).Error; err != nil {
				return err
			}

			for _, eReq := range bReq.Entries {
				entry := &models.MilkJournalEntry{
					MemberID:           eReq.MemberID,
					MilkJournalID:      journal.ID,
					MilkJournalBatchID: batch.ID,
					Status:             "PENDING",
					Quantity:           eReq.Quantity,
					RouteCenterID:      eReq.RouteCenterID,
					CanID:              eReq.CanID,
				}
				if err := tx.Create(entry).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *MilkJournalService) DeleteMilkJournal(id string) error {
	return db.DB.Delete(&models.MilkJournal{}, id).Error
}

func (s *MilkJournalService) GetDailyJournals() ([]dtos.DailyJournalSummaryResponse, error) {
	var results []struct {
		Journal    string `gorm:"column:journal"`
		Route      string `gorm:"column:route"`
		Shift      string `gorm:"column:shift"`
		BatchesStr string `gorm:"column:batches_str"`
	}

	query := `
		SELECT 
			mj.journal,
			r.route_name AS route,
			mds.name AS shift,
			GROUP_CONCAT(mjb.batch_no) AS batches_str
		FROM milk_journals mj
		LEFT JOIN routes r ON mj.route_id = r.id
		LEFT JOIN milk_delivery_shifts mds ON mj.milk_delivery_shift_id = mds.id
		LEFT JOIN milk_journal_batches mjb ON mj.id = mjb.milk_journal_id
		WHERE DATE(mj.journal_date) = CURDATE() 
		  AND mj.deleted_at IS NULL
		GROUP BY mj.id, mj.journal, r.route_name, mds.name
	`
	err := db.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	journals := make([]dtos.DailyJournalSummaryResponse, len(results))
	for i, res := range results {
		batches := []string{}
		if res.BatchesStr != "" {
			batches = strings.Split(res.BatchesStr, ",")
		}
		journals[i] = dtos.DailyJournalSummaryResponse{
			Journal: res.Journal,
			Route:   res.Route,
			Shift:   res.Shift,
			Batches: batches,
		}
	}

	return journals, nil
}

// ImportJournals bulk imports milk journals from CSV, XLS, or XLSX files.
func (s *MilkJournalService) ImportJournals(file *multipart.FileHeader, userID uint64) error {
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

	go s.processJournalRowsInBackground(data, userID)

	return nil
}

func (s *MilkJournalService) processJournalRowsInBackground(data [][]string, userID uint64) {
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
				func() {
					defer func() {
						if r := recover(); r != nil {
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
						// Expected format: Member No(0), Batch No(1), Route(2), Can ID(3), Quantity(4), Journal(5), Journal Date(6), Shift(7)
						if len(row) < 8 {
							return fmt.Errorf("insufficient columns (found %d, need at least 8)", len(row))
						}

						memberNo := strings.TrimSpace(row[0])
						var member models.Member
						if err := tx.Where("member_no = ?", memberNo).First(&member).Error; err != nil {
							return fmt.Errorf("member '%s' not found", memberNo)
						}

						routeVal := strings.TrimSpace(row[2])
						var route models.Route
						if err := tx.Where("route_name = ? OR route_code = ?", routeVal, routeVal).First(&route).Error; err != nil {
							return fmt.Errorf("route '%s' not found", routeVal)
						}

						canVal := strings.TrimSpace(row[3])
						var canID uint64
						if canVal != "" {
							var can models.MilkCan
							if err := tx.Where("can_id = ?", canVal).First(&can).Error; err == nil {
								canID = can.ID
							}
						}

						journalName := strings.TrimSpace(row[5])
						journalDate := utils.ParseFlexibleDate(row[6])

						// Resolve Shift
						var shift models.MilkDeliveryShift
						shiftFound := false
						shiftVal := strings.TrimSpace(row[7])
						if shiftVal != "" {
							if err := tx.Where("name = ?", shiftVal).First(&shift).Error; err == nil {
								shiftFound = true
							}
						}

						if !shiftFound {
							tx.First(&shift) // Resolve a default shift if not specified or not found
						}

						var journal models.MilkJournal
						if err := tx.Where("journal = ? AND journal_date = ? AND route_id = ? AND milk_delivery_shift_id = ?", journalName, journalDate, route.ID, shift.ID).First(&journal).Error; err != nil {
							if err == gorm.ErrRecordNotFound {
								journal = models.MilkJournal{
									Journal:             journalName,
									JournalDate:         &journalDate,
									MilkDeliveryShiftID: shift.ID,
									RouteID:             route.ID,
									UserID:              userID,
									Confirmed:           false, //
								}
								if err := tx.Create(&journal).Error; err != nil {
									return err
								}
							} else {
								return err
							}
						}

						batchNo := strings.TrimSpace(row[1])
						var batch models.MilkJournalBatch
						if err := tx.Where("milk_journal_id = ? AND batch_no = ?", journal.ID, batchNo).First(&batch).Error; err != nil {
							if err == gorm.ErrRecordNotFound {
								batch = models.MilkJournalBatch{
									MilkJournalID: journal.ID,
									BatchNo:       batchNo,
								}
								if err := tx.Create(&batch).Error; err != nil {
									return err
								}
							} else {
								return err
							}
						}

						qty, _ := utils.ParseFloat(row[4])
						entry := models.MilkJournalEntry{
							MemberID:           member.ID,
							MilkJournalID:      journal.ID,
							MilkJournalBatchID: batch.ID,
							Status:             "PENDING",
							Quantity:           qty,
							CanID:              canID,
						}
						return tx.Create(&entry).Error
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

	message := fmt.Sprintf("Milk journal import completed. Success: %d, Failed: %d out of %d records.", totalRows-failCount, failCount, totalRows)
	notificationType := "SUCCESS"
	errorLink := ""
	if failCount > 0 {
		notificationType = "ERROR"
		errorLink = fmt.Sprintf("/milk-journals/import-errors/%d", importID)
	} else if totalRows == 0 {
		message = "Milk journal import completed. No records were processed."
		notificationType = "SUCCESS"
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Milk Journal Import Status",
		Message:          message,
		NotificationType: notificationType,
		ErrorLink:        errorLink,
	})
}

func (s *MilkJournalService) GetImportErrors(importID uint64) ([]models.ImportError, error) {
	var importErrors []models.ImportError
	err := db.DB.Where("import_id = ?", importID).Order("id DESC").Find(&importErrors).Error
	return importErrors, err
}
