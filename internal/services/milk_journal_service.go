package services

import (
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MilkJournalService struct{}

func NewMilkJournalService() *MilkJournalService {
	return &MilkJournalService{}
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

	query := `sss
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
		jd := utils.ParseDate(req.JournalDate)
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
