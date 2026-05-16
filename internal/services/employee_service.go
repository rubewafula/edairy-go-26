package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type EmployeeService struct{}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{}
}

func (s *EmployeeService) CreateEmployee(req dtos.CreateEmployeeRequest, userID uint64) (*models.Employee, error) {
	employee := &models.Employee{
		BaseModel:         models.BaseModel{CreatedBy: userID},
		UserID:            req.UserID,
		Surname:           req.Surname,
		FirstName:         req.FirstName,
		MiddleName:        req.MiddleName,
		EmployeeNo:        req.EmployeeNo,
		IDNo:              req.IDNo,
		KraPin:            req.KraPin,
		NssfNo:            req.NssfNo,
		NhifNo:            req.NhifNo,
		Gender:            req.Gender,
		DateOfBirth:       utils.ParseDate(req.DateOfBirth),
		Phone:             req.Phone,
		Email:             req.Email,
		JobPositionID:     req.JobPositionID,
		Status:            req.Status,
		Title:             req.Title,
		Town:              req.Town,
		SiteID:            req.SiteID,
		MaritalStatus:     req.MaritalStatus,
		Religion:          req.Religion,
		Disabled:          req.Disabled,
		StoreID:           req.StoreID,
		PostalAddress:     req.PostalAddress,
		PostalCode:        req.PostalCode,
		BirthCity:         req.BirthCity,
		NextOfKinFullName: req.NextOfKinFullName,
		NextOfKinPhone:    req.NextOfKinPhone,
	}

	if err := db.DB.Create(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (s *EmployeeService) GetEmployees(page, limit int) ([]dtos.EmployeeResponse, int64, error) {
	var employees []dtos.EmployeeResponse
	var total int64
	db.DB.Model(&models.Employee{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT e.*, jp.name as job_position_name, d.department_name as department_name
		FROM employees e
		LEFT JOIN job_positions jp ON e.job_position_id = jp.id
		LEFT JOIN departments d ON jp.department_id = d.id
		WHERE e.deleted_at IS NULL
		ORDER BY e.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&employees).Error
	return employees, total, err
}

func (s *EmployeeService) GetEmployee(id string) (*dtos.EmployeeResponse, error) {
	var employee dtos.EmployeeResponse
	query := `
		SELECT e.*, jp.name as job_position_name, d.department_name as department_name
		FROM employees e
		LEFT JOIN job_positions jp ON e.job_position_id = jp.id
		LEFT JOIN departments d ON jp.department_id = d.id
		WHERE e.id = ? AND e.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&employee).Error
	if err != nil {
		return nil, err
	}
	if employee.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &employee, nil
}

func (s *EmployeeService) UpdateEmployee(id string, req dtos.UpdateEmployeeRequest, userID uint64) error {
	var employee models.Employee
	if err := db.DB.First(&employee, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"surname":       req.Surname,
		"first_name":    req.FirstName,
		"middle_name":   req.MiddleName,
		"phone_number":  req.Phone,
		"email_address": req.Email,
		"updated_by":    userID,
	}

	return db.DB.Model(&employee).Updates(updates).Error
}

func (s *EmployeeService) DeleteEmployee(id string, userID uint64) error {
	var employee models.Employee
	if err := db.DB.First(&employee, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&employee).Update("updated_by", userID).Delete(&employee).Error
}

func (s *EmployeeService) GetSalaries(employeeID string) ([]models.EmployeeSalary, error) {
	var salaries []models.EmployeeSalary
	err := db.DB.Where("employee_id = ?", employeeID).Find(&salaries).Error
	return salaries, err
}

func (s *EmployeeService) GetSalary(id string) (*models.EmployeeSalary, error) {
	var salary models.EmployeeSalary
	if err := db.DB.First(&salary, id).Error; err != nil {
		return nil, err
	}
	return &salary, nil
}

func (s *EmployeeService) UpdateSalary(id string, req dtos.UpdateEmployeeSalaryRequest, userID uint64) error {
	return db.DB.Model(&models.EmployeeSalary{}).Where("id = ?", id).
		Updates(map[string]interface{}{"basic_salary": req.BasicSalary, "status": req.Status, "updated_by": userID}).Error
}

func (s *EmployeeService) CreateSalary(req dtos.CreateEmployeeSalaryRequest, userID uint64) (*models.EmployeeSalary, error) {
	salary := &models.EmployeeSalary{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		EmployeeID:  req.EmployeeID,
		BasicSalary: req.BasicSalary,
		Status:      req.Status,
	}
	if salary.Status == "" {
		salary.Status = "ACTIVE"
	}
	err := db.DB.Create(salary).Error
	return salary, err
}
