package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type DepartmentService struct{}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{}
}

func (s *DepartmentService) CreateDepartment(req dtos.CreateDepartmentRequest, userID uint64) (*models.Department, error) {
	department := &models.Department{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		DepartmentCode: req.DepartmentCode,
		DepartmentName: req.DepartmentName,
		Description:    req.Description,
	}

	if err := db.DB.Create(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (s *DepartmentService) GetDepartments(page, limit int) ([]dtos.DepartmentResponse, int64, error) {
	var results []dtos.DepartmentResponse
	var total int64
	db.DB.Model(&models.Department{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.Department{}).
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Scan(&results).Error

	return results, total, err
}

func (s *DepartmentService) GetDepartment(id string) (*dtos.DepartmentResponse, error) {
	var result dtos.DepartmentResponse
	if err := db.DB.Model(&models.Department{}).First(&result, id).Error; err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *DepartmentService) UpdateDepartment(id string, req dtos.UpdateDepartmentRequest, userID uint64) error {
	var department models.Department
	if err := db.DB.First(&department, id).Error; err != nil {
		return err
	}

	department.DepartmentCode = req.DepartmentCode
	department.DepartmentName = req.DepartmentName
	department.Description = req.Description
	department.UpdatedBy = userID

	return db.DB.Save(&department).Error
}

func (s *DepartmentService) DeleteDepartment(id string, userID uint64) error {
	return db.DB.Model(&models.Department{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.Department{}).Error
}
