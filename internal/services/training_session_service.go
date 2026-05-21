package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type TrainingSessionService struct{}

func NewTrainingSessionService() *TrainingSessionService {
	return &TrainingSessionService{}
}

func (s *TrainingSessionService) CreateSession(req dtos.CreateTrainingSessionRequest) (*models.TrainingSession, error) {

	session := &models.TrainingSession{
		TrainingID:       req.TrainingID,
		Partner:          req.Partner,
		SessionStartTime: utils.ParseDate(req.SessionStartTime),
		SessionEndTime:   utils.ParseDate(req.SessionEndTime),
		Topic:            req.Topic,
		Description:      req.Description,
		Trainers:         req.Trainers,
		Status:           "ACTIVE",
	}

	if err := db.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *TrainingSessionService) GetSessions() ([]dtos.TrainingSessionResponse, int64, error) {
	var results []dtos.TrainingSessionResponse
	var total int64
	db.DB.Model(&models.TrainingSession{}).Count(&total)

	query := `
		SELECT 
			ts.id, 
			ts.training_id, 
			t.topic, 
			ts.partner, 
			ts.session_start_time, 
			ts.session_end_time, 
			ts.trainers,
			ts.status, ts.description, ts.created_at, ts.updated_at
		FROM training_sessions ts
		LEFT JOIN trainings t ON ts.training_id = t.id
		WHERE ts.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *TrainingSessionService) GetSession(id string) (*dtos.TrainingSessionResponse, error) {
	var result dtos.TrainingSessionResponse
	query := `
		SELECT 
			ts.id, 
			ts.training_id, 
			t.topic, 
			ts.partner, 
			ts.session_start_time, 
			ts.session_end_time,
			ts.trainers,
			ts.status, ts.description, ts.created_at, ts.updated_at
		FROM training_sessions ts
		LEFT JOIN trainings t ON ts.training_id = t.id
		WHERE ts.id = ? AND ts.deleted_at IS NULL
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

func (s *TrainingSessionService) UpdateSession(id string, req dtos.UpdateTrainingSessionRequest) error {
	var session models.TrainingSession
	if err := db.DB.First(&session, id).Error; err != nil {
		return err
	}

	session.TrainingID = req.TrainingID
	session.Partner = req.Partner
	session.SessionStartTime = utils.ParseDate(req.SessionStartTime)
	session.SessionEndTime = utils.ParseDate(req.SessionEndTime)
	session.Topic = req.Topic
	session.Description = req.Description
	session.Trainers = req.Trainers
	session.Status = "ACTIVE"

	return db.DB.Save(&session).Error
}

func (s *TrainingSessionService) DeleteSession(id string) error {
	return db.DB.Delete(&models.TrainingSession{}, id).Error
}
