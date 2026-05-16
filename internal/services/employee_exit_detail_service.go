package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type EmployeeExitDetailService struct{}

func NewEmployeeExitDetailService() *EmployeeExitDetailService {
	return &EmployeeExitDetailService{}
}

func (s *EmployeeExitDetailService) CreateExitDetail(req dtos.CreateEmployeeExitDetailRequest, userID uint64) (*models.EmployeeExitDetail, error) {
	exitDetail := &models.EmployeeExitDetail{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		EmployeeID:      req.EmployeeID,
		ContractType:    req.ContractType,
		ContractEndDate: utils.ParseDate(req.ContractEndDate),
		DateOfLeaving:   utils.ParseDate(req.DateOfLeaving),
		ExitCategory:    req.ExitCategory,
		Reasons:         req.Reasons,
	}

	if err := db.DB.Create(exitDetail).Error; err != nil {
		return nil, err
	}
	return exitDetail, nil
}

func (s *EmployeeExitDetailService) GetExitDetails(employeeID string, page, limit int) ([]dtos.EmployeeExitDetailResponse, int64, error) {
	var results []dtos.EmployeeExitDetailResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeExitDetail{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			eed.id, eed.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eed.contract_type, eed.contract_end_date,
			eed.date_of_leaving, eed.exit_category, eed.reasons, eed.created_at, eed.updated_at
		FROM employee_exit_details eed
		LEFT JOIN employees e ON eed.employee_id = e.id
		WHERE eed.deleted_at IS NULL AND (? = '' OR eed.employee_id = ?)
		ORDER BY eed.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeExitDetailService) GetExitDetail(id string) (*dtos.EmployeeExitDetailResponse, error) {
	var result dtos.EmployeeExitDetailResponse
	query := `
		SELECT 
			eed.id, eed.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eed.contract_type, eed.contract_end_date,
			eed.date_of_leaving, eed.exit_category, eed.reasons, eed.created_at, eed.updated_at
		FROM employee_exit_details eed
		LEFT JOIN employees e ON eed.employee_id = e.id
		WHERE eed.id = ? AND eed.deleted_at IS NULL
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

func (s *EmployeeExitDetailService) UpdateExitDetail(id string, req dtos.UpdateEmployeeExitDetailRequest, userID uint64) error {
	var exitDetail models.EmployeeExitDetail
	if err := db.DB.First(&exitDetail, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"contract_type":     req.ContractType,
		"contract_end_date": utils.ParseDate(req.ContractEndDate),
		"date_of_leaving":   utils.ParseDate(req.DateOfLeaving),
		"exit_category":     req.ExitCategory,
		"reasons":           req.Reasons,
		"updated_by":        userID,
	}

	return db.DB.Model(&exitDetail).Updates(updates).Error
}

func (s *EmployeeExitDetailService) DeleteExitDetail(id string, userID uint64) error {
	var exitDetail models.EmployeeExitDetail
	if err := db.DB.First(&exitDetail, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&exitDetail).Update("updated_by", userID).Delete(&exitDetail).Error
}
