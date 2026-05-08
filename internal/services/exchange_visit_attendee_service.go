package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type ExchangeVisitAttendeeService struct{}

func NewExchangeVisitAttendeeService() *ExchangeVisitAttendeeService {
	return &ExchangeVisitAttendeeService{}
}

func (s *ExchangeVisitAttendeeService) CreateAttendee(req dtos.CreateExchangeVisitAttendeeRequest) (*models.ExchangeVisitAttendee, error) {
	attendee := &models.ExchangeVisitAttendee{
		ExchangeVisitID:      req.ExchangeVisitID,
		Attendee:             req.Attendee,
		AttendeeOrganization: req.AttendeeOrganization,
		AttendeeDesignation:  req.AttendeeDesignation,
		Attended:             req.Attended,
		Comments:             req.Comments,
		AttendanceEmployeeID: req.AttendanceEmployeeID,
	}

	if err := db.DB.Create(attendee).Error; err != nil {
		return nil, err
	}
	return attendee, nil
}

func (s *ExchangeVisitAttendeeService) GetAttendees() ([]dtos.ExchangeVisitAttendeeResponse, int64, error) {
	var results []dtos.ExchangeVisitAttendeeResponse
	var total int64
	db.DB.Model(&models.ExchangeVisitAttendee{}).Count(&total)

	query := `
		SELECT 
			eva.id, eva.exchange_visit_id, ev.exchange_visit_partner AS partner, 
			eva.attendee, eva.attendee_organization, eva.attendee_designation, 
			eva.attended, eva.comments, eva.attendance_employee_id, 
			e.first_name AS employee_first_name, e.surname AS employee_surname,
			eva.created_at, eva.updated_at
		FROM exchange_visit_attendees eva
		LEFT JOIN exchange_visits ev ON eva.exchange_visit_id = ev.id
		LEFT JOIN employees e ON eva.attendance_employee_id = e.id
		WHERE eva.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ExchangeVisitAttendeeService) GetAttendee(id string) (*dtos.ExchangeVisitAttendeeResponse, error) {
	var result dtos.ExchangeVisitAttendeeResponse
	query := `
		SELECT 
			eva.id, eva.exchange_visit_id, ev.exchange_visit_partner AS partner, 
			eva.attendee, eva.attendee_organization, eva.attendee_designation, 
			eva.attended, eva.comments, eva.attendance_employee_id, 
			e.first_name AS employee_first_name, e.surname AS employee_surname,
			eva.created_at, eva.updated_at
		FROM exchange_visit_attendees eva
		LEFT JOIN exchange_visits ev ON eva.exchange_visit_id = ev.id
		LEFT JOIN employees e ON eva.attendance_employee_id = e.id
		WHERE eva.id = ? AND eva.deleted_at IS NULL
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

func (s *ExchangeVisitAttendeeService) UpdateAttendee(id string, req dtos.UpdateExchangeVisitAttendeeRequest) error {
	var attendee models.ExchangeVisitAttendee
	if err := db.DB.First(&attendee, id).Error; err != nil {
		return err
	}

	attendee.ExchangeVisitID = req.ExchangeVisitID
	attendee.Attendee = req.Attendee
	attendee.AttendeeOrganization = req.AttendeeOrganization
	attendee.AttendeeDesignation = req.AttendeeDesignation
	attendee.Attended = req.Attended
	attendee.Comments = req.Comments
	attendee.AttendanceEmployeeID = req.AttendanceEmployeeID

	return db.DB.Save(&attendee).Error
}

func (s *ExchangeVisitAttendeeService) DeleteAttendee(id string) error {
	return db.DB.Delete(&models.ExchangeVisitAttendee{}, id).Error
}
