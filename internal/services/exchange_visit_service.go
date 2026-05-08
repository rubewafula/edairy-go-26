package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type ExchangeVisitService struct{}

func NewExchangeVisitService() *ExchangeVisitService {
	return &ExchangeVisitService{}
}

func (s *ExchangeVisitService) CreateVisit(req dtos.CreateExchangeVisitRequest) (*models.ExchangeVisit, error) {
	visit := &models.ExchangeVisit{
		Partner:    req.Partner,
		VisitDate:  utils.ParseDate(req.VisitDate),
		Purpose:    req.Purpose,
		Venue:      req.Venue,
		EmployeeID: req.EmployeeID,
		VisitNotes: req.VisitNotes,
	}

	if err := db.DB.Create(visit).Error; err != nil {
		return nil, err
	}
	return visit, nil
}

func (s *ExchangeVisitService) GetVisits() ([]dtos.ExchangeVisitResponse, int64, error) {
	var results []dtos.ExchangeVisitResponse
	var total int64
	db.DB.Model(&models.ExchangeVisit{}).Count(&total)

	query := `
		SELECT 
			ev.id, ev.exchange_visit_partner AS partner, ev.exchange_visit_date AS visit_date, 
			ev.purpose, ev.venue, ev.exchange_visit_employee_id AS employee_id, 
			e.first_name AS employee_first_name, e.surname AS employee_surname,
			ev.visit_notes, ev.created_at, ev.updated_at
		FROM exchange_visits ev
		LEFT JOIN employees e ON ev.exchange_visit_employee_id = e.id
		WHERE ev.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ExchangeVisitService) GetVisit(id string) (*dtos.ExchangeVisitResponse, error) {
	var result dtos.ExchangeVisitResponse
	query := `
		SELECT 
			ev.id, ev.exchange_visit_partner AS partner, ev.exchange_visit_date AS visit_date, 
			ev.purpose, ev.venue, ev.exchange_visit_employee_id AS employee_id, 
			e.first_name AS employee_first_name, e.surname AS employee_surname,
			ev.visit_notes, ev.created_at, ev.updated_at
		FROM exchange_visits ev
		LEFT JOIN employees e ON ev.exchange_visit_employee_id = e.id
		WHERE ev.id = ? AND ev.deleted_at IS NULL
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

func (s *ExchangeVisitService) UpdateVisit(id string, req dtos.UpdateExchangeVisitRequest) error {
	var visit models.ExchangeVisit
	if err := db.DB.First(&visit, id).Error; err != nil {
		return err
	}

	visit.Partner = req.Partner
	visit.VisitDate = utils.ParseDate(req.VisitDate)
	visit.Purpose = req.Purpose
	visit.Venue = req.Venue
	visit.EmployeeID = req.EmployeeID
	visit.VisitNotes = req.VisitNotes

	return db.DB.Save(&visit).Error
}

func (s *ExchangeVisitService) DeleteVisit(id string) error {
	return db.DB.Delete(&models.ExchangeVisit{}, id).Error
}
