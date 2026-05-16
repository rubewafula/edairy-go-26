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

	err := queryBuilder.Limit(limit).Offset(offset).Order("id DESC").Find(&results).Error
	return results, total, err
}

func (s *EmployeeDependantService) GetEmployeeDependant(id string) (*dtos.EmployeeDependantResponse, error) {
	var result dtos.EmployeeDependantResponse
	err := db.DB.Model(&models.EmployeeDependant{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
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
