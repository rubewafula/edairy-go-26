package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeProfessionalTitleService struct{}

func NewEmployeeProfessionalTitleService() *EmployeeProfessionalTitleService {
	return &EmployeeProfessionalTitleService{}
}

func (s *EmployeeProfessionalTitleService) Create(req dtos.CreateEmployeeProfessionalTitleRequest, userID uint64) (*models.EmployeeProfessionalTitle, error) {
	title := &models.EmployeeProfessionalTitle{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Code:      req.Code,
		Name:      req.Title,
	}

	if err := db.DB.Create(title).Error; err != nil {
		return nil, err
	}
	return title, nil
}

func (s *EmployeeProfessionalTitleService) List(employeeID string, page, limit int) ([]dtos.EmployeeProfessionalTitleResponse, int64, error) {
	var results []dtos.EmployeeProfessionalTitleResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeProfessionalTitle{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	err := queryBuilder.Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *EmployeeProfessionalTitleService) Get(id string) (*dtos.EmployeeProfessionalTitleResponse, error) {
	var result dtos.EmployeeProfessionalTitleResponse
	err := db.DB.Model(&models.EmployeeProfessionalTitle{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *EmployeeProfessionalTitleService) Update(id string, req dtos.UpdateEmployeeProfessionalTitleRequest, userID uint64) error {
	var title models.EmployeeProfessionalTitle
	if err := db.DB.First(&title, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"title":      req.Title,
		"updated_by": userID,
	}

	return db.DB.Model(&title).Updates(updates).Error
}

func (s *EmployeeProfessionalTitleService) Delete(id string, userID uint64) error {
	var title models.EmployeeProfessionalTitle
	if err := db.DB.First(&title, id).Error; err != nil {
		return err
	}
	// Audit update before soft delete
	return db.DB.Model(&title).Update("updated_by", userID).Delete(&title).Error
}
