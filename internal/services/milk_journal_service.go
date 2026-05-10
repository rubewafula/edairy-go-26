package services

import (
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
	journal := &models.MilkJournal{
		Journal:             req.Journal,
		JournalDate:         utils.ParseDate(req.JournalDate),
		MilkDeliveryShiftID: req.MilkDeliveryShiftID,
		RouteID:             req.RouteID,
		UserID:              req.UserID,
		TransporterID:       req.TransporterID,
		Confirmed:           req.Confirmed,
	}

	if err := db.DB.Create(journal).Error; err != nil {
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
			mj.id, mj.journal, mj.journal_date, mj.milk_delivery_shift_id, mds.name AS milk_delivery_shift,
			mj.route_id, r.route_name, mj.user_id, mj.transporter_id, mj.confirmed,
			mj.created_at, mj.updated_at,
			(SELECT COUNT(*) FROM milk_journal_entries mje WHERE mje.milk_journal_id = mj.id AND mje.deleted_at IS NULL) AS entries_count
		FROM milk_journals mj
		LEFT JOIN milk_delivery_shifts mds ON mj.milk_delivery_shift_id = mds.id
		LEFT JOIN routes r ON mj.route_id = r.id
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
			mj.id, mj.journal, mj.journal_date, mj.milk_delivery_shift_id, mds.name AS milk_delivery_shift,
			mj.route_id, r.route_name, mj.user_id, mj.transporter_id, mj.confirmed,
			mj.created_at, mj.updated_at,
			(SELECT COUNT(*) FROM milk_journal_entries mje WHERE mje.milk_journal_id = mj.id AND mje.deleted_at IS NULL) AS entries_count
		FROM milk_journals mj
		LEFT JOIN milk_delivery_shifts mds ON mj.milk_delivery_shift_id = mds.id
		LEFT JOIN routes r ON mj.route_id = r.id
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

	journal.Journal = req.Journal
	journal.JournalDate = utils.ParseDate(req.JournalDate)
	journal.MilkDeliveryShiftID = req.MilkDeliveryShiftID
	journal.RouteID = req.RouteID
	journal.UserID = req.UserID
	journal.TransporterID = req.TransporterID
	journal.Confirmed = req.Confirmed

	return db.DB.Save(&journal).Error
}

func (s *MilkJournalService) DeleteMilkJournal(id string) error {
	return db.DB.Delete(&models.MilkJournal{}, id).Error
}
