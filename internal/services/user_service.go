package services

import (
	"github.com/google/uuid"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(req dtos.CreateUserRequest) (*models.User, error) {
	// dtos.CreateUserRequest is correctly defined and used here.
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hashedPassword,
		IsVerified:        req.IsVerified,
		VerificationToken: uuid.NewString(),
	}

	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUsers() ([]dtos.UserResponse, int64, error) {
	var results []dtos.UserResponse
	var total int64
	db.DB.Model(&models.User{}).Count(&total)

	query := `
		SELECT id, name, email, email_verified_at, is_verified,
		       created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	if len(results) > 0 {
		userIDs := make([]uint64, len(results))
		for i, u := range results {
			userIDs[i] = u.ID
		}

		type UserRoleRes struct {
			dtos.RoleResponse
			UserID uint64
		}
		var userRoles []UserRoleRes
		roleQuery := `
			SELECT r.id, r.name, r.guard_name, r.created_at, r.updated_at, ur.user_id
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id IN ? AND r.deleted_at IS NULL
		`
		db.DB.Raw(roleQuery, userIDs).Scan(&userRoles)

		roleMap := make(map[uint64][]dtos.RoleResponse)
		for _, ur := range userRoles {
			roleMap[ur.UserID] = append(roleMap[ur.UserID], ur.RoleResponse)
		}

		for i := range results {
			results[i].Roles = roleMap[results[i].ID]
		}

		// Fetch direct permissions
		type UserPermRes struct {
			dtos.PermissionResponse
			UserID uint64
		}
		var userPerms []UserPermRes
		permQuery := `
			SELECT p.id, p.name, p.guard_name, p.created_at, p.updated_at, up.user_id
			FROM permissions p
			INNER JOIN user_permissions up ON p.id = up.permission_id
			WHERE up.user_id IN ? AND p.deleted_at IS NULL
		`
		db.DB.Raw(permQuery, userIDs).Scan(&userPerms)

		permMap := make(map[uint64][]dtos.PermissionResponse)
		for _, up := range userPerms {
			permMap[up.UserID] = append(permMap[up.UserID], up.PermissionResponse)
		}

		for i := range results {
			results[i].Permissions = permMap[results[i].ID]
		}
	}

	return results, total, nil
}

func (s *UserService) GetUser(id string) (*dtos.UserResponse, error) {
	var result dtos.UserResponse
	query := `SELECT id, name, email, email_verified_at, is_verified,
		       created_at, updated_at 
		FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var roles []dtos.RoleResponse
	roleQuery := `SELECT r.id, r.name, r.guard_name, r.created_at, r.updated_at FROM roles r INNER JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = ? AND r.deleted_at IS NULL`
	db.DB.Raw(roleQuery, result.ID).Scan(&roles)
	result.Roles = roles

	var permissions []dtos.PermissionResponse
	permQuery := `SELECT p.id, p.name, p.guard_name, p.created_at, p.updated_at FROM permissions p INNER JOIN user_permissions up ON p.id = up.permission_id WHERE up.user_id = ? AND p.deleted_at IS NULL`
	db.DB.Raw(permQuery, result.ID).Scan(&permissions)
	result.Permissions = permissions

	return &result, nil
}

func (s *UserService) UpdateUser(id string, req dtos.UpdateUserRequest) error {
	// dtos.UpdateUserRequest is correctly defined and used here.
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return err
	}
	user.Name = req.Name
	user.Email = req.Email
	user.IsVerified = req.IsVerified
	return db.DB.Save(&user).Error
}

func (s *UserService) DeleteUser(id string) error {
	return db.DB.Delete(&models.User{}, id).Error
}
