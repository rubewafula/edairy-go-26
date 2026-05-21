package dtos

type CreateUserRequest struct {
	Name          string   `json:"name" validate:"required,max=255"`
	Email         string   `json:"email" validate:"required,email"`
	Password      string   `json:"password" validate:"required,min=8"`
	IsVerified    bool     `json:"is_verified"`
	RoleIDs       []uint64 `json:"role_ids"`
	PermissionIDs []uint64 `json:"permission_ids"`
}

type UpdateUserRequest struct {
	Name          string   `json:"name" validate:"required,max=255"`
	Email         string   `json:"email" validate:"required,email"`
	IsVerified    bool     `json:"is_verified"`
	RoleIDs       []uint64 `json:"role_ids"`
	PermissionIDs []uint64 `json:"permission_ids"`
}
