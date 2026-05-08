package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
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
		Confirmed:           req.Confirmed,
	}

	if err := db.DB.Create(journal).Error; err != nil {
		return nil, err
	}
	return journal, nil
}

func (s *MilkJournalService) GetMilkJournals() ([]models.MilkJournal, int64, error) {
	var journals []models.MilkJournal
	var total int64
	db.DB.Model(&models.MilkJournal{}).Count(&total)
	err := db.DB.Find(&journals).Error
	return journals, total, err
}

func (s *MilkJournalService) GetMilkJournal(id string) (*models.MilkJournal, error) {
	var journal models.MilkJournal
	if err := db.DB.First(&journal, id).Error; err != nil {
		return nil, err
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
	journal.Confirmed = req.Confirmed

	return db.DB.Save(&journal).Error
}

func (s *MilkJournalService) DeleteMilkJournal(id string) error {
	return db.DB.Delete(&models.MilkJournal{}, id).Error
}
