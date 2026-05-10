package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MilkJournalEntryService struct{}

func NewMilkJournalEntryService() *MilkJournalEntryService {
	return &MilkJournalEntryService{}
}

func (s *MilkJournalEntryService) CreateEntry(req dtos.CreateMilkJournalEntryRequest) (*models.MilkJournalEntry, error) {
	entry := &models.MilkJournalEntry{
		MemberID:            req.MemberID,
		MilkJournalID:       req.MilkJournalID,
		MilkJournalBatchID:  req.MilkJournalBatchID,
		RouteID:             req.RouteID,
		MilkDeliveryShiftID: req.MilkDeliveryShiftID,
		Status:              req.Status,
		JournalDate:         utils.ParseDate(req.JournalDate),
		Quantity:            req.Quantity,
		TransporterID:       req.TransporterID,
		RouteCenterID:       req.RouteCenterID,
		CanID:               req.CanID,
	}

	if err := db.DB.Create(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *MilkJournalEntryService) GetEntries(page, limit int) ([]dtos.MilkJournalEntryResponse, int64, error) {
	var entries []dtos.MilkJournalEntryResponse
	var total int64
	db.DB.Model(&models.MilkJournalEntry{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			mje.*, 
			m.member_no,
			CONCAT(m.first_name, ' ', m.last_name) AS member_name,
			r.route_name, 
			mds.name AS milk_delivery_shift
		FROM milk_journal_entries mje
		LEFT JOIN member_registrations m ON mje.member_id = m.id
		LEFT JOIN routes r ON mje.route_id = r.id
		LEFT JOIN milk_delivery_shifts mds ON mje.milk_delivery_shift_id = mds.id
		WHERE mje.deleted_at IS NULL
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&entries).Error
	return entries, total, err
}

func (s *MilkJournalEntryService) GetEntry(id string) (*dtos.MilkJournalEntryResponse, error) {
	var entry dtos.MilkJournalEntryResponse
	query := `
		SELECT 
			mje.*, 
			m.member_no,
			CONCAT(m.first_name, ' ', m.last_name) AS member_name,
			r.route_name, 
			mds.name AS milk_delivery_shift
		FROM milk_journal_entries mje
		LEFT JOIN member_registrations m ON mje.member_id = m.id
		LEFT JOIN routes r ON mje.route_id = r.id
		LEFT JOIN milk_delivery_shifts mds ON mje.milk_delivery_shift_id = mds.id
		WHERE mje.id = ? AND mje.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&entry).Error
	if err != nil {
		return nil, err
	}
	if entry.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &entry, nil
}

func (s *MilkJournalEntryService) UpdateEntry(id string, req dtos.UpdateMilkJournalEntryRequest) error {
	var entry models.MilkJournalEntry
	if err := db.DB.First(&entry, id).Error; err != nil {
		return err
	}

	entry.MemberID = req.MemberID
	entry.MilkJournalID = req.MilkJournalID
	entry.MilkJournalBatchID = req.MilkJournalBatchID
	entry.RouteID = req.RouteID
	entry.MilkDeliveryShiftID = req.MilkDeliveryShiftID
	entry.Status = req.Status
	entry.JournalDate = utils.ParseDate(req.JournalDate)
	entry.Quantity = req.Quantity
	entry.TransporterID = req.TransporterID
	entry.RouteCenterID = req.RouteCenterID
	entry.CanID = req.CanID

	return db.DB.Save(&entry).Error
}

func (s *MilkJournalEntryService) DeleteEntry(id string) error {
	return db.DB.Delete(&models.MilkJournalEntry{}, id).Error
}

func (s *MilkJournalEntryService) GetStrayEntries(page, limit int) ([]dtos.StrayMilkCollectionResponse, int64, error) {
	var entries []dtos.StrayMilkCollectionResponse
	var total int64
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM milk_journal_entries mje
		INNER JOIN member_registrations m ON mje.member_id = m.id
		WHERE mje.deleted_at IS NULL AND mje.route_id != m.route_id
	`
	db.DB.Raw(countQuery).Scan(&total)

	query := `
		SELECT 
			mje.id, 
			mje.member_id, 
			m.member_no, 
			CONCAT(m.first_name, ' ', m.last_name) AS member_name,
			m.route_id AS member_route_id,
			mr.route_name AS member_route,
			mje.route_id AS journal_route_id,
			jr.route_name AS stray_route,
			mje.quantity,
			mje.journal_date,
			mds.name AS milk_delivery_shift,
			mje.created_at
		FROM milk_journal_entries mje
		INNER JOIN member_registrations m ON mje.member_id = m.id
		LEFT JOIN routes mr ON m.route_id = mr.id
		LEFT JOIN routes jr ON mje.route_id = jr.id
		LEFT JOIN milk_delivery_shifts mds ON mje.milk_delivery_shift_id = mds.id
		WHERE mje.deleted_at IS NULL AND mje.route_id != m.route_id
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&entries).Error
	return entries, total, err
}
