package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type TrainingAttendeeService struct{}

func NewTrainingAttendeeService() *TrainingAttendeeService {
	return &TrainingAttendeeService{}
}

func (s *TrainingAttendeeService) CreateAttendee(req dtos.CreateTrainingAttendeeRequest) (*models.TrainingAttendee, error) {
	attendee := &models.TrainingAttendee{
		TrainingSessionID: req.TrainingSessionID,
		Names:             req.Names,
		IDNumber:          req.IDNumber,
		PhoneNumber:       req.PhoneNumber,
		MembershipNumber:  req.MembershipNumber,
		Comments:          req.Comments,
		MemberID:          req.MemberID,
	}

	if err := db.DB.Create(attendee).Error; err != nil {
		return nil, err
	}
	return attendee, nil
}

func (s *TrainingAttendeeService) GetAttendees() ([]dtos.TrainingAttendeeResponse, int64, error) {
	var results []dtos.TrainingAttendeeResponse
	var total int64
	db.DB.Model(&models.TrainingAttendee{}).Count(&total)

	query := `
		SELECT 
			ta.id, ta.training_session_id, t.topic, 
			ta.names, ta.id_number, ta.phone_number, ta.membership_number, 
			ta.comments, ta.member_id, ta.created_at, ta.updated_at
		FROM training_session_attendees ta
		LEFT JOIN training_sessions ts ON ta.training_session_id = ts.id
		LEFT JOIN trainings t ON ts.training_id = t.id
		WHERE ta.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *TrainingAttendeeService) GetAttendee(id string) (*dtos.TrainingAttendeeResponse, error) {
	var result dtos.TrainingAttendeeResponse
	query := `
		SELECT 
			ta.id, ta.training_session_id, t.topic, 
			ta.names, ta.id_number, ta.phone_number, ta.membership_number, 
			ta.comments, ta.member_id, ta.created_at, ta.updated_at
		FROM training_session_attendees ta
		LEFT JOIN training_sessions ts ON ta.training_session_id = ts.id
		LEFT JOIN trainings t ON ts.training_id = t.id
		WHERE ta.id = ? AND ta.deleted_at IS NULL
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

func (s *TrainingAttendeeService) UpdateAttendee(id string, req dtos.UpdateTrainingAttendeeRequest) error {
	var attendee models.TrainingAttendee
	if err := db.DB.First(&attendee, id).Error; err != nil {
		return err
	}

	attendee.TrainingSessionID = req.TrainingSessionID
	attendee.Names = req.Names
	attendee.IDNumber = req.IDNumber
	attendee.PhoneNumber = req.PhoneNumber
	attendee.MembershipNumber = req.MembershipNumber
	attendee.Comments = req.Comments
	attendee.MemberID = req.MemberID

	return db.DB.Save(&attendee).Error
}

func (s *TrainingAttendeeService) DeleteAttendee(id string) error {
	return db.DB.Delete(&models.TrainingAttendee{}, id).Error
}
