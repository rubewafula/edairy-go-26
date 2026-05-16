package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type EmployeeSalaryService struct{}

func NewEmployeeSalaryService() *EmployeeSalaryService {
	return &EmployeeSalaryService{}
}

func (s *EmployeeSalaryService) Create(req dtos.CreateEmployeeSalaryRequest, userID uint64) (*models.EmployeeSalary, error) {
	salary := &models.EmployeeSalary{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		EmployeeID:  req.EmployeeID,
		BasicSalary: req.BasicSalary,
		Status:      req.Status,
	}
	if err := db.DB.Create(salary).Error; err != nil {
		return nil, err
	}
	return salary, nil
}

func (s *EmployeeSalaryService) List(employeeID string, page, limit int) ([]dtos.EmployeeSalaryResponse, int64, error) {
	var results []dtos.EmployeeSalaryResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeSalary{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	err := queryBuilder.Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *EmployeeSalaryService) Get(id string) (*dtos.EmployeeSalaryResponse, error) {
	var result dtos.EmployeeSalaryResponse
	err := db.DB.Model(&models.EmployeeSalary{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmployeeSalaryService) Update(id string, req dtos.UpdateEmployeeSalaryRequest, userID uint64) error {
	var salary models.EmployeeSalary
	if err := db.DB.First(&salary, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"basic_salary": req.BasicSalary,
		"status":       req.Status,
		"updated_by":   userID,
	}
	return db.DB.Model(&salary).Updates(updates).Error
}

func (s *EmployeeSalaryService) Delete(id string, userID uint64) error {
	var salary models.EmployeeSalary
	if err := db.DB.First(&salary, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&salary).Update("updated_by", userID).Delete(&salary).Error
}
