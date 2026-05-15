package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeBenefitService struct{}

func NewEmployeeBenefitService() *EmployeeBenefitService {
	return &EmployeeBenefitService{}
}

func (s *EmployeeBenefitService) CreateEmployeeBenefit(req dtos.CreateEmployeeBenefitRequest, userID uint64) (*models.EmployeeBenefit, error) {
	benefit := &models.EmployeeBenefit{
		BaseModel:  models.BaseModel{CreatedBy: userID},
		EmployeeID: req.EmployeeID,
		BenefitID:  req.BenefitID,
		Amount:     req.Amount,
		Status:     req.Status,
	}
	if benefit.Status == "" {
		benefit.Status = "ACTIVE"
	}
	if err := db.DB.Create(benefit).Error; err != nil {
		return nil, err
	}
	return benefit, nil
}

func (s *EmployeeBenefitService) GetEmployeeBenefits(page, limit int) ([]dtos.EmployeeBenefitResponse, int64, error) {
	var results []dtos.EmployeeBenefitResponse
	var total int64
	db.DB.Model(&models.EmployeeBenefit{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			eb.id, eb.employee_id, eb.benefit_id, b.name AS benefit_name,
			eb.amount, eb.status, eb.created_at, eb.updated_at
		FROM employee_benefits eb
		LEFT JOIN benefits b ON eb.benefit_id = b.id
		WHERE eb.deleted_at IS NULL
		ORDER BY eb.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeBenefitService) GetEmployeeBenefit(id string) (*dtos.EmployeeBenefitResponse, error) {
	var result dtos.EmployeeBenefitResponse
	query := `
		SELECT 
			eb.id, eb.employee_id, eb.benefit_id, b.name AS benefit_name,
			eb.amount, eb.status, eb.created_at, eb.updated_at
		FROM employee_benefits eb
		LEFT JOIN benefits b ON eb.benefit_id = b.id
		WHERE eb.id = ? AND eb.deleted_at IS NULL
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

func (s *EmployeeBenefitService) UpdateEmployeeBenefit(id string, req dtos.UpdateEmployeeBenefitRequest, userID uint64) error {
	var benefit models.EmployeeBenefit
	if err := db.DB.First(&benefit, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"benefit_id": req.BenefitID,
		"amount":     req.Amount,
		"status":     req.Status,
		"updated_by": userID,
	}
	return db.DB.Model(&benefit).Updates(updates).Error
}

func (s *EmployeeBenefitService) DeleteEmployeeBenefit(id string, userID uint64) error {
	var benefit models.EmployeeBenefit
	if err := db.DB.First(&benefit, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&benefit).Update("updated_by", userID).Delete(&benefit).Error
}
