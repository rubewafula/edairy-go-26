package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type EmployeeQualificationService struct{}

func NewEmployeeQualificationService() *EmployeeQualificationService {
	return &EmployeeQualificationService{}
}

func (s *EmployeeQualificationService) CreateQualification(req dtos.CreateEmployeeQualificationRequest, userID uint64) (*models.EmployeeQualification, error) {
	qualification := &models.EmployeeQualification{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		EmployeeID:    req.EmployeeID,
		Qualification: req.Qualification,
		Institution:   req.Institution,
		StartDate:     utils.ParseDate(req.StartDate),
		EndDate:       utils.ParseDate(req.EndDate),
		Score:         req.Score,
	}

	if err := db.DB.Create(qualification).Error; err != nil {
		return nil, err
	}
	return qualification, nil
}

func (s *EmployeeQualificationService) GetEmployeeQualifications(employeeID string, page, limit int) ([]dtos.EmployeeQualificationResponse, int64, error) {
	var results []dtos.EmployeeQualificationResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeQualification{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			eq.id, eq.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eq.qualification, eq.institution,
			eq.start_date, eq.end_date, eq.score, eq.created_at, eq.updated_at
		FROM employee_qualifications eq
		LEFT JOIN employees e ON eq.employee_id = e.id
		WHERE eq.deleted_at IS NULL AND (? = '' OR eq.employee_id = ?)
		ORDER BY eq.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeQualificationService) GetQualification(id string) (*dtos.EmployeeQualificationResponse, error) {
	var result dtos.EmployeeQualificationResponse
	query := `
		SELECT 
			eq.id, eq.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eq.qualification, eq.institution,
			eq.start_date, eq.end_date, eq.score, eq.created_at, eq.updated_at
		FROM employee_qualifications eq
		LEFT JOIN employees e ON eq.employee_id = e.id
		WHERE eq.id = ? AND eq.deleted_at IS NULL
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

func (s *EmployeeQualificationService) UpdateQualification(id string, req dtos.UpdateEmployeeQualificationRequest, userID uint64) error {
	var qualification models.EmployeeQualification
	if err := db.DB.First(&qualification, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"qualification": req.Qualification,
		"institution":   req.Institution,
		"start_date":    utils.ParseDate(req.StartDate),
		"end_date":      utils.ParseDate(req.EndDate),
		"score":         req.Score,
		"updated_by":    userID,
	}

	return db.DB.Model(&qualification).Updates(updates).Error
}

func (s *EmployeeQualificationService) DeleteQualification(id string, userID uint64) error {
	var qualification models.EmployeeQualification
	if err := db.DB.First(&qualification, id).Error; err != nil {
		return err
	}

	// Audit the update before soft delete
	return db.DB.Model(&qualification).Update("updated_by", userID).Delete(&qualification).Error
}

func (s *EmployeeQualificationService) GetQualificationsByEmployeeID(employeeID string) ([]dtos.EmployeeQualificationResponse, error) {
	var results []dtos.EmployeeQualificationResponse
	err := db.DB.Model(&models.EmployeeQualification{}).Where("employee_id = ?", employeeID).Find(&results).Error
	return results, err
}
