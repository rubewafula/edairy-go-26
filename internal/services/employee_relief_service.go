package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeReliefService struct{}

func NewEmployeeReliefService() *EmployeeReliefService {
	return &EmployeeReliefService{}
}

func (s *EmployeeReliefService) Create(req dtos.CreateEmployeeReliefRequest, userID uint64) (*models.EmployeeRelief, error) {
	relief := &models.EmployeeRelief{
		BaseModel:  models.BaseModel{CreatedBy: userID},
		EmployeeID: req.EmployeeID,
		ReliefID:   req.ReliefID,
		Status:     req.Status,
	}

	if err := db.DB.Create(relief).Error; err != nil {
		return nil, err
	}
	return relief, nil
}

func (s *EmployeeReliefService) List(employeeID string, page, limit int) ([]dtos.EmployeeReliefResponse, int64, error) {
	var results []dtos.EmployeeReliefResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeRelief{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			er.id, er.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			er.relief_id, tr.relief as relief_name,
			er.status, er.created_at, er.updated_at
		FROM employee_reliefs er
		LEFT JOIN employees e ON er.employee_id = e.id
		LEFT JOIN payroll_reliefs tr ON er.relief_id = tr.id
		WHERE er.deleted_at IS NULL AND (? = '' OR er.employee_id = ?)
		ORDER BY er.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeReliefService) Get(id string) (*dtos.EmployeeReliefResponse, error) {
	var result dtos.EmployeeReliefResponse
	query := `
		SELECT 
			er.id, er.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			er.relief_id, tr.relief as relief_name,
			er.status, er.created_at, er.updated_at
		FROM employee_reliefs er
		LEFT JOIN employees e ON er.employee_id = e.id
		LEFT JOIN payroll_reliefs tr ON er.relief_id = tr.id
		WHERE er.id = ? AND er.deleted_at IS NULL
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

func (s *EmployeeReliefService) Update(id string, req dtos.UpdateEmployeeReliefRequest, userID uint64) error {
	var relief models.EmployeeRelief
	if err := db.DB.First(&relief, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"employee_id": req.EmployeeID,
		"relief_id":   req.ReliefID,
		"status":      req.Status,
		"updated_by":  userID,
	}

	return db.DB.Model(&relief).Updates(updates).Error
}

func (s *EmployeeReliefService) Delete(id string, userID uint64) error {
	var relief models.EmployeeRelief
	if err := db.DB.First(&relief, id).Error; err != nil {
		return err
	}

	// Audit before soft delete
	return db.DB.Model(&relief).Update("updated_by", userID).Delete(&relief).Error
}
