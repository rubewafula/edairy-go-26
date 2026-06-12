package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeLeaveDetailService struct{}

func NewEmployeeLeaveDetailService() *EmployeeLeaveDetailService {
	return &EmployeeLeaveDetailService{}
}

func (s *EmployeeLeaveDetailService) CreateLeaveDetail(req dtos.CreateEmployeeLeaveDetailRequest, userID uint64) (*models.EmployeeLeaveDetail, error) {
	detail := &models.EmployeeLeaveDetail{
		BaseModel:           models.BaseModel{CreatedBy: userID},
		EmployeeID:          req.EmployeeID,
		BalanceBF:           req.BalanceBF,
		AllocatedDays:       req.AllocatedDays,
		EmployeeLeaveTypeID: req.LeaveTypeID,
	}

	if err := db.DB.Create(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func (s *EmployeeLeaveDetailService) GetLeaveDetails(employeeID string, page, limit int) ([]dtos.EmployeeLeaveDetailResponse, int64, error) {
	var results []dtos.EmployeeLeaveDetailResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeLeaveDetail{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			eld.*, 
			elt.code as leave_type_code, 
			elt.description as leave_type_name,
			CONCAT(e.first_name, ' ', COALESCE(e.middle_name, ''), ' ', e.surname) as employee_name
		FROM employee_leave_details eld
		LEFT JOIN employee_leave_types elt ON eld.employee_leave_type_id = elt.id
		LEFT JOIN employees e ON eld.employee_id = e.id
		WHERE eld.deleted_at IS NULL 
		AND (? = '' OR eld.employee_id = ?)
		ORDER BY eld.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeLeaveDetailService) GetLeaveDetail(id string) (*dtos.EmployeeLeaveDetailResponse, error) {
	var result dtos.EmployeeLeaveDetailResponse
	query := `
		SELECT 
			eld.*, 
			elt.code as leave_type_code, 
			elt.description as leave_type_name,
			CONCAT(e.first_name, ' ', COALESCE(e.middle_name, ''), ' ', e.surname) as employee_name
		FROM employee_leave_details eld
		LEFT JOIN employee_leave_types elt ON eld.employee_leave_type_id = elt.id
		LEFT JOIN employees e ON eld.employee_id = e.id
		WHERE eld.id = ? AND eld.deleted_at IS NULL
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

func (s *EmployeeLeaveDetailService) UpdateLeaveDetail(id string, req dtos.UpdateEmployeeLeaveDetailRequest, userID uint64) error {
	var detail models.EmployeeLeaveDetail
	if err := db.DB.First(&detail, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"balance_bf":             req.BalanceBF,
		"allocated_days":         req.AllocatedDays,
		"employee_leave_type_id": req.LeaveTypeID,
		"updated_by":             userID,
	}

	return db.DB.Model(&detail).Updates(updates).Error
}

func (s *EmployeeLeaveDetailService) DeleteLeaveDetail(id string, userID uint64) error {
	var detail models.EmployeeLeaveDetail
	if err := db.DB.First(&detail, id).Error; err != nil {
		return err
	}
	// Standard audit before soft delete
	return db.DB.Model(&detail).Update("updated_by", userID).Delete(&detail).Error
}
