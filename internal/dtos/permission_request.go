package dtos

type CreatePermissionRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	GuardName string `json:"guard_name" validate:"required,max=255"`
}

type UpdatePermissionRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	GuardName string `json:"guard_name" validate:"required,max=255"`
}
