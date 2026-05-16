package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type EmployeeLeaveApplicationService struct{}

func NewEmployeeLeaveApplicationService() *EmployeeLeaveApplicationService {
	return &EmployeeLeaveApplicationService{}
}

func (s *EmployeeLeaveApplicationService) CreateApplication(req dtos.CreateEmployeeLeaveApplicationRequest, userID uint64) (*models.EmployeeLeaveApplication, error) {
	application := &models.EmployeeLeaveApplication{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		ApplicationNo: fmt.Sprintf("LV-%d-%d", req.EmployeeID, time.Now().Unix()),
		EmployeeID:    req.EmployeeID,
		LeaveTypeID:   req.LeaveTypeID,
		DaysApplied:   req.DaysApplied,
		StartDate:     utils.ParseDate(req.StartDate),
		EndDate:       utils.ParseDate(req.EndDate),
		ReturnDate:    utils.ParseDate(req.ReturnDate),
		Status:        "PENDING",
	}

	if err := db.DB.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (s *EmployeeLeaveApplicationService) GetApplications(employeeID string, page, limit int) ([]dtos.EmployeeLeaveApplicationResponse, int64, error) {
	var results []dtos.EmployeeLeaveApplicationResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeLeaveApplication{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ela.*, 
			CONCAT(e.first_name, ' ', e.surname) as employee_name,
			elt.description as leave_type,
			CONCAT(app.first_name, ' ', app.surname) as approver_name
		FROM employee_leave_applications ela
		LEFT JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN employee_leave_types elt ON ela.leave_type_id = elt.id
		LEFT JOIN employees app ON ela.approver_id = app.id
		WHERE ela.deleted_at IS NULL AND (? = '' OR ela.employee_id = ?)
		ORDER BY ela.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeLeaveApplicationService) GetApplication(id string) (*dtos.EmployeeLeaveApplicationResponse, error) {
	var result dtos.EmployeeLeaveApplicationResponse
	query := `
		SELECT 
			ela.*, 
			CONCAT(e.first_name, ' ', e.surname) as employee_name,
			elt.description as leave_type,
			CONCAT(app.first_name, ' ', app.surname) as approver_name
		FROM employee_leave_applications ela
		LEFT JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN employee_leave_types elt ON ela.leave_type_id = elt.id
		LEFT JOIN employees app ON ela.approver_id = app.id
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

func (s *EmployeeLeaveApplicationService) UpdateApplication(id string, req dtos.UpdateEmployeeLeaveApplicationRequest, userID uint64) error {
	var application models.EmployeeLeaveApplication
	if err := db.DB.First(&application, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"approver_id":   req.ApproverID,
		"days_approved": req.DaysApproved,
		"status":        req.Status,
		"approved":      req.Approved,
		"updated_by":    userID,
	}

	return db.DB.Model(&application).Updates(updates).Error
}

func (s *EmployeeLeaveApplicationService) DeleteApplication(id string, userID uint64) error {
	var application models.EmployeeLeaveApplication
	if err := db.DB.First(&application, id).Error; err != nil {
		return err
	}
	// Audit before soft delete
	return db.DB.Model(&application).Update("updated_by", userID).Delete(&application).Error
}
