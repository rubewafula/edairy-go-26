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
		BaseModel:     models.BaseModel{CreatedBy: userID},
		EmployeeID:    req.EmployeeID,
		BalanceBF:     req.BalanceBF,
		AllocatedDays: req.AllocatedDays,
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

	err := queryBuilder.Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *EmployeeLeaveDetailService) GetLeaveDetail(id string) (*dtos.EmployeeLeaveDetailResponse, error) {
	var result dtos.EmployeeLeaveDetailResponse
	err := db.DB.Model(&models.EmployeeLeaveDetail{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
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
		"balance_bf":     req.BalanceBF,
		"allocated_days": req.AllocatedDays,
		"updated_by":     userID,
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
