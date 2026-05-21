package services

import (
	"log"

	"github.com/rubewafula/edairy-go-26/internal/apperrors"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (s *RoleService) CreateRole(req dtos.CreateRoleRequest) (*models.Role, error) {
	var count int64
	err := db.DB.Model(&models.Role{}).
		Where("name = ? AND guard_name = ? AND deleted_at IS NULL", req.Name, req.GuardName).
		Count(&count).Error

	if err != nil {
		log.Printf("[RoleService] Found error trying to query duplicate role: %s", err.Error())
		return nil, err
	}

	if count > 0 {
		return nil, apperrors.ErrRoleExists
	}

	role := &models.Role{
		Name:      req.Name,
		GuardName: req.GuardName,
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(req.PermissionIDs) > 0 {
			var permissions []models.Permission
			if err := tx.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
				return err
			}

			if err := tx.Model(role).Association("Permissions").Append(permissions); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetRoles() ([]dtos.RoleResponse, int64, error) {
	var results []dtos.RoleResponse
	var total int64
	db.DB.Model(&models.Role{}).Count(&total)

	query := `
		SELECT id, name, guard_name, created_at, updated_at
		FROM roles
		WHERE deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// Efficiently fetch permissions for all roles in one query
	if len(results) > 0 {
		roleIDs := make([]uint64, len(results))
		for i, r := range results {
			roleIDs[i] = r.ID
		}

		type RolePerm struct {
			dtos.PermissionResponse
			RoleID uint64
		}
		var rolePerms []RolePerm
		permQuery := `
			SELECT p.id, p.name, p.guard_name, p.created_at, p.updated_at, rp.role_id
			FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id IN ? AND p.deleted_at IS NULL AND rp.deleted_at IS NULL
		`
		db.DB.Raw(permQuery, roleIDs).Scan(&rolePerms)

		permMap := make(map[uint64][]dtos.PermissionResponse)
		for _, rp := range rolePerms {
			permMap[rp.RoleID] = append(permMap[rp.RoleID], rp.PermissionResponse)
		}

		for i := range results {
			if perms, ok := permMap[results[i].ID]; ok {
				results[i].Permissions = perms
			} else {
				results[i].Permissions = []dtos.PermissionResponse{}
			}
		}
	}

	return results, total, nil
}

func (s *RoleService) GetRole(id string) (*dtos.RoleResponse, error) {
	var result dtos.RoleResponse
	query := `
		SELECT id, name, guard_name, created_at, updated_at
		FROM roles
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

	// Fetch permissions using the pivot table
	var permissions []dtos.PermissionResponse
	permQuery := `
		SELECT p.id, p.name, p.guard_name, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ? AND p.deleted_at IS NULL AND rp.deleted_at IS NULL
	`
	db.DB.Raw(permQuery, result.ID).Scan(&permissions)
	result.Permissions = permissions

	return &result, nil
}

func (s *RoleService) UpdateRole(id string, req dtos.UpdateRoleRequest) error {
	var role models.Role
	if err := db.DB.First(&role, id).Error; err != nil {
		return err
	}

	// Check for duplicate role name/guard (excluding the current role)
	var count int64
	err := db.DB.Model(&models.Role{}).
		Where("name = ? AND guard_name = ? AND id != ? AND deleted_at IS NULL", req.Name, req.GuardName, id).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return apperrors.ErrRoleExists
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		role.Name = req.Name
		role.GuardName = req.GuardName

		if err := tx.Save(&role).Error; err != nil {
			return err
		}

		var permissions []models.Permission
		if len(req.PermissionIDs) > 0 {
			if err := tx.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
				return err
			}
		}

		// Replace clears existing permission associations and assigns the new list
		return tx.Model(&role).Association("Permissions").Replace(permissions)
	})
}

func (s *RoleService) DeleteRole(id string) error {
	return db.DB.Delete(&models.Role{}, id).Error
}

func (s *RoleService) AppendAllPermissionsToRole(id string) error {
	var role models.Role
	if err := db.DB.First(&role, id).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := db.DB.Find(&permissions).Error; err != nil {
		return err
	}

	return db.DB.Model(&role).Association("Permissions").Append(permissions)
}
