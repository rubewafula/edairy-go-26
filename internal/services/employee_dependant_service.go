package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeDependantService struct{}

func NewEmployeeDependantService() *EmployeeDependantService {
	return &EmployeeDependantService{}
}

func (s *EmployeeDependantService) CreateEmployeeDependant(req dtos.CreateEmployeeDependantRequest, userID uint64) (*models.EmployeeDependant, error) {
	dependant := &models.EmployeeDependant{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		EmployeeID:   req.EmployeeID,
		Name:         req.Name,
		Relationship: req.Relationship,
	}
	if err := db.DB.Create(dependant).Error; err != nil {
		return nil, err
	}
	return dependant, nil
}

func (s *EmployeeDependantService) GetEmployeeDependants(employeeID string, page, limit int) ([]dtos.EmployeeDependantResponse, int64, error) {
	var results []dtos.EmployeeDependantResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeDependant{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ed.id, ed.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			ed.name, ed.relationship, ed.created_at, ed.updated_at
		FROM employee_dependants ed
		LEFT JOIN employees e ON ed.employee_id = e.id
		WHERE ed.deleted_at IS NULL AND (? = '' OR ed.employee_id = ?)
		ORDER BY ed.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeDependantService) GetEmployeeDependant(id string) (*dtos.EmployeeDependantResponse, error) {
	var result dtos.EmployeeDependantResponse
	query := `
		SELECT 
			ed.id, ed.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			ed.name, ed.relationship, ed.created_at, ed.updated_at
		FROM employee_dependants ed
		LEFT JOIN employees e ON ed.employee_id = e.id
		WHERE ed.id = ? AND ed.deleted_at IS NULL
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

func (s *EmployeeDependantService) UpdateEmployeeDependant(id string, req dtos.UpdateEmployeeDependantRequest, userID uint64) error {
	var dependant models.EmployeeDependant
	if err := db.DB.First(&dependant, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":         req.Name,
		"relationship": req.Relationship,
		"updated_by":   userID,
	}

	return db.DB.Model(&dependant).Updates(updates).Error
}

func (s *EmployeeDependantService) DeleteEmployeeDependant(id string, userID uint64) error {
	var dependant models.EmployeeDependant
	if err := db.DB.First(&dependant, id).Error; err != nil {
		return err
	}

	// Audit update before soft delete
	return db.DB.Model(&dependant).Update("updated_by", userID).Delete(&dependant).Error
}
