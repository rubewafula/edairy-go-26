package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type PermissionService struct{}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

func (s *PermissionService) CreatePermission(req dtos.CreatePermissionRequest) (*models.Permission, error) {
	permission := &models.Permission{
		Name:      req.Name,
		GuardName: req.GuardName,
	}

	if err := db.DB.Create(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *PermissionService) GetPermissions() ([]dtos.PermissionResponse, int64, error) {
	var results []dtos.PermissionResponse
	var total int64
	db.DB.Model(&models.Permission{}).Count(&total)

	query := `
		SELECT id, name, guard_name, created_at, updated_at
		FROM permissions
		WHERE deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *PermissionService) GetPermission(id string) (*dtos.PermissionResponse, error) {
	var result dtos.PermissionResponse
	query := `
		SELECT id, name, guard_name, created_at, updated_at
		FROM permissions
		WHERE id = ? AND deleted_at IS NULL
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

func (s *PermissionService) UpdatePermission(id string, req dtos.UpdatePermissionRequest) error {
	var permission models.Permission
	if err := db.DB.First(&permission, id).Error; err != nil {
		return err
	}

	permission.Name = req.Name
	permission.GuardName = req.GuardName

	return db.DB.Save(&permission).Error
}

func (s *PermissionService) DeletePermission(id string) error {
	return db.DB.Delete(&models.Permission{}, id).Error
}
