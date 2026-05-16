package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeLeaveTypeService struct{}

func NewEmployeeLeaveTypeService() *EmployeeLeaveTypeService {
	return &EmployeeLeaveTypeService{}
}

func (s *EmployeeLeaveTypeService) CreateEmployeeLeaveType(req dtos.CreateEmployeeLeaveTypeRequest, userID uint64) (*models.EmployeeLeaveType, error) {
	leaveType := &models.EmployeeLeaveType{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		Code:        req.Code,
		Description: req.Description,
		Days:        req.Days,
		Gender:      req.Gender,
	}
	if err := db.DB.Create(leaveType).Error; err != nil {
		return nil, err
	}
	return leaveType, nil
}

func (s *EmployeeLeaveTypeService) GetEmployeeLeaveTypes(page, limit int) ([]dtos.EmployeeLeaveTypeResponse, int64, error) {
	var results []dtos.EmployeeLeaveTypeResponse
	var total int64
	if err := db.DB.Model(&models.EmployeeLeaveType{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := db.DB.Model(&models.EmployeeLeaveType{}).Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *EmployeeLeaveTypeService) GetEmployeeLeaveType(id string) (*dtos.EmployeeLeaveTypeResponse, error) {
	var result dtos.EmployeeLeaveTypeResponse
	if err := db.DB.Model(&models.EmployeeLeaveType{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error; err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *EmployeeLeaveTypeService) UpdateEmployeeLeaveType(id string, req dtos.UpdateEmployeeLeaveTypeRequest, userID uint64) error {
	var leaveType models.EmployeeLeaveType
	if err := db.DB.First(&leaveType, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"code":        req.Code,
		"description": req.Description,
		"days":        req.Days,
		"gender":      req.Gender,
		"updated_by":  userID,
	}
	return db.DB.Model(&leaveType).Updates(updates).Error
}

func (s *EmployeeLeaveTypeService) DeleteEmployeeLeaveType(id string, userID uint64) error {
	var leaveType models.EmployeeLeaveType
	if err := db.DB.First(&leaveType, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&leaveType).Update("updated_by", userID).Delete(&leaveType).Error
}
