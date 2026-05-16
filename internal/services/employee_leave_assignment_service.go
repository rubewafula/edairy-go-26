package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeLeaveAssignmentService struct{}

func NewEmployeeLeaveAssignmentService() *EmployeeLeaveAssignmentService {
	return &EmployeeLeaveAssignmentService{}
}

func (s *EmployeeLeaveAssignmentService) CreateAssignment(req dtos.CreateEmployeeLeaveAssignmentRequest, userID uint64) (*models.EmployeeLeaveAssignment, error) {
	assignment := &models.EmployeeLeaveAssignment{
		BaseModel:          models.BaseModel{CreatedBy: userID},
		EmployeeID:         req.EmployeeID,
		LeaveApplicationID: req.LeaveApplicationID,
		RelieverID:         req.RelieverID,
	}

	if err := db.DB.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *EmployeeLeaveAssignmentService) GetAssignments(employeeID string, page, limit int) ([]dtos.EmployeeLeaveAssignmentResponse, int64, error) {
	var results []dtos.EmployeeLeaveAssignmentResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeLeaveAssignment{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ela.id, ela.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			ela.leave_application_id, lapp.application_no,
			ela.reliever_id, CONCAT(r.first_name, ' ', r.surname) as reliever_name,
			ela.created_at, ela.updated_at
		FROM employee_leave_assignments ela
		LEFT JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN employee_leave_applications lapp ON ela.leave_application_id = lapp.id
		LEFT JOIN employees r ON ela.reliever_id = r.id
		WHERE ela.deleted_at IS NULL AND (? = '' OR ela.employee_id = ?)
		ORDER BY ela.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeLeaveAssignmentService) GetAssignment(id string) (*dtos.EmployeeLeaveAssignmentResponse, error) {
	var result dtos.EmployeeLeaveAssignmentResponse
	query := `
		SELECT 
			ela.id, ela.employee_id, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			ela.leave_application_id, lapp.application_no,
			ela.reliever_id, CONCAT(r.first_name, ' ', r.surname) as reliever_name,
			ela.created_at, ela.updated_at
		FROM employee_leave_assignments ela
		LEFT JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN employee_leave_applications lapp ON ela.leave_application_id = lapp.id
		LEFT JOIN employees r ON ela.reliever_id = r.id
		WHERE ela.id = ? AND ela.deleted_at IS NULL
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

func (s *EmployeeLeaveAssignmentService) UpdateAssignment(id string, req dtos.UpdateEmployeeLeaveAssignmentRequest, userID uint64) error {
	var assignment models.EmployeeLeaveAssignment
	if err := db.DB.First(&assignment, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"employee_id":          req.EmployeeID,
		"leave_application_id": req.LeaveApplicationID,
		"reliever_id":          req.RelieverID,
		"updated_by":           userID,
	}
	return db.DB.Model(&assignment).Updates(updates).Error
}

func (s *EmployeeLeaveAssignmentService) DeleteAssignment(id string, userID uint64) error {
	return db.DB.Model(&models.EmployeeLeaveAssignment{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.EmployeeLeaveAssignment{}).Error
}
