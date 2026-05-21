package dtos

type CreatePermissionRequest struct {
	Name      string `json:"name" validate:"required,max=128"`
	GuardName string `json:"guard_name" validate:"required,max=25"`
}

type UpdatePermissionRequest struct {
	Name      string `json:"name" validate:"required,max=128"`
	GuardName string `json:"guard_name" validate:"required,max=25"`
}
