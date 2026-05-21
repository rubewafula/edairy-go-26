package dtos

type CreateRoleRequest struct {
	Name          string   `json:"name" validate:"required,max=255"`
	GuardName     string   `json:"guard_name" validate:"required,max=255"`
	PermissionIDs []uint64 `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string   `json:"name" validate:"required,max=255"`
	GuardName     string   `json:"guard_name" validate:"required,max=255"`
	PermissionIDs []uint64 `json:"permission_ids"`
}
