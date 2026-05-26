package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type JobPositionService struct{}

func NewJobPositionService() *JobPositionService {
	return &JobPositionService{}
}

func (s *JobPositionService) CreateJobPosition(req dtos.CreateJobPositionRequest, userID uint64) (*models.JobPosition, error) {
	jobPosition := &models.JobPosition{
		BaseModel:         models.BaseModel{CreatedBy: userID},
		Code:              req.Code,
		Name:              req.Name,
		JobDescription:    req.JobDescription,
		DepartmentID:      req.DepartmentID,
		GradeID:           req.GradeID,
		NoOfPosts:         req.NoOfPosts,
		VaccantPositions:  req.NoOfPosts, // Initially, all posts are vacant
		OccupiedPositions: 0,
	}

	if err := db.DB.Create(jobPosition).Error; err != nil {
		return nil, err
	}
	return jobPosition, nil
}

func (s *JobPositionService) GetJobPositions(page, limit int) ([]dtos.JobPositionResponse, int64, error) {
	var results []dtos.JobPositionResponse
	var total int64
	db.DB.Model(&models.JobPosition{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			jp.*, d.department_name
		FROM job_positions jp
		LEFT JOIN departments d ON jp.department_id = d.id
		WHERE jp.deleted_at IS NULL
		ORDER BY jp.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error

	return results, total, err
}

func (s *JobPositionService) GetJobPosition(id string) (*dtos.JobPositionResponse, error) {
	var result dtos.JobPositionResponse
	query := `
		SELECT 
			jp.*, d.department_name
		FROM job_positions jp
		LEFT JOIN departments d ON jp.department_id = d.id
		WHERE jp.id = ? AND jp.deleted_at IS NULL
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

func (s *JobPositionService) UpdateJobPosition(id string, req dtos.UpdateJobPositionRequest, userID uint64) error {
	var jobPosition models.JobPosition
	if err := db.DB.First(&jobPosition, id).Error; err != nil {
		return err
	}

	jobPosition.Code = req.Code
	jobPosition.Name = req.Name
	jobPosition.JobDescription = req.JobDescription
	jobPosition.DepartmentID = req.DepartmentID
	jobPosition.GradeID = req.GradeID
	jobPosition.NoOfPosts = req.NoOfPosts
	jobPosition.OccupiedPositions = req.OccupiedPositions
	jobPosition.VaccantPositions = req.NoOfPosts - req.OccupiedPositions // Recalculate vacant
	jobPosition.UpdatedBy = userID

	return db.DB.Save(&jobPosition).Error
}

func (s *JobPositionService) DeleteJobPosition(id string, userID uint64) error {
	return db.DB.Model(&models.JobPosition{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.JobPosition{}).Error
}
