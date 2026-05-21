package services

import (
	"log"

	"github.com/rubewafula/edairy-go-26/internal/apperrors"
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

	var count int64

	err := db.DB.Model(&models.Permission{}).
		Where("name = ? AND guard_name = ? AND deleted_at IS NULL", req.Name, req.GuardName).
		Count(&count).Error

	if err != nil {
		log.Printf("[PermissionService] Found error trying to query duplicate: %s", err.Error())
		return nil, err
	}

	log.Printf("[PermissionService] Found duplicates : %d", count)

	if count > 0 {
		return nil, apperrors.ErrPermissionExists
	}

	permission := &models.Permission{
		Name:      req.Name,
		GuardName: req.GuardName,
	}

	if err := db.DB.Create(permission).Error; err != nil {
		return nil, err
	}

	log.Printf("[PermissionService] Created permission: %s (%s)", permission.Name, permission.GuardName)
	return permission, nil
}

func (s *PermissionService) GetPermissions(page int, perPage int) ([]dtos.PermissionResponse, int64, error) {
	var results []dtos.PermissionResponse
	var total int64

	// defaults
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// total count
	err := db.DB.Model(&models.Permission{}).
		Where("deleted_at IS NULL").
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// paginated query
	query := `
		SELECT 
			id, 
			name, 
			guard_name, 
			created_at, 
			updated_at
		FROM permissions
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`

	err = db.DB.Raw(query, perPage, offset).Scan(&results).Error

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

	if err := db.DB.Save(&permission).Error; err != nil {
		log.Printf("[PermissionService] Failed to update permission ID %s: %v", id, err)
		return err
	}

	log.Printf("[PermissionService] Updated permission ID %s", id)
	return nil
}

func (s *PermissionService) DeletePermission(id string) error {
	if err := db.DB.Delete(&models.Permission{}, id).Error; err != nil {
		log.Printf("[PermissionService] Failed to delete permission ID %s: %v", id, err)
		return err
	}

	log.Printf("[PermissionService] Deleted permission ID %s", id)
	return nil
}
