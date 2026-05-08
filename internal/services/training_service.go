package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type TrainingService struct{}

func NewTrainingService() *TrainingService {
	return &TrainingService{}
}

func (s *TrainingService) CreateTraining(req dtos.CreateTrainingRequest) (*models.Training, error) {
	status := req.Status
	if status == "" {
		status = "SCHEDULED"
	}
	training := &models.Training{
		Topic:          req.Topic,
		Description:    req.Description,
		Venue:          req.Venue,
		TrainingUserID: req.TrainingUserID,
		TrainingDate:   utils.ParseDate(req.TrainingDate),
		Status:         status,
	}

	if err := db.DB.Create(training).Error; err != nil {
		return nil, err
	}
	return training, nil
}

func (s *TrainingService) GetTrainings() ([]dtos.TrainingResponse, int64, error) {
	var results []dtos.TrainingResponse
	var total int64
	db.DB.Model(&models.Training{}).Count(&total)

	query := `
		SELECT 
			t.id, t.topic, t.description, t.venue, u.name AS facilitator, 
			t.training_date, t.status, t.created_at, t.updated_at
		FROM trainings t
		LEFT JOIN users u ON t.training_user_id = u.id
		WHERE t.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *TrainingService) GetTraining(id string) (*dtos.TrainingResponse, error) {
	var result dtos.TrainingResponse
	query := `
		SELECT 
			t.id, t.topic, t.description, t.venue, u.name AS facilitator, 
			t.training_date, t.status, t.created_at, t.updated_at
		FROM trainings t
		LEFT JOIN users u ON t.training_user_id = u.id
		WHERE t.id = ? AND t.deleted_at IS NULL
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

func (s *TrainingService) UpdateTraining(id string, req dtos.UpdateTrainingRequest) error {
	var training models.Training
	if err := db.DB.First(&training, id).Error; err != nil {
		return err
	}

	training.Topic = req.Topic
	training.Description = req.Description
	training.Venue = req.Venue
	training.TrainingUserID = req.TrainingUserID
	training.TrainingDate = utils.ParseDate(req.TrainingDate)
	training.Status = req.Status

	return db.DB.Save(&training).Error
}

func (s *TrainingService) DeleteTraining(id string) error {
	return db.DB.Delete(&models.Training{}, id).Error
}
