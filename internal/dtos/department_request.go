package dtos

type CreateDepartmentRequest struct {
	DepartmentCode string `json:"department_code" validate:"required,max=255"`
	DepartmentName string `json:"department_name" validate:"required,max=255"`
	Description    string `json:"description" validate:"max=255"`
}

type UpdateDepartmentRequest struct {
	DepartmentCode string `json:"department_code" validate:"required,max=255"`
	DepartmentName string `json:"department_name" validate:"required,max=255"`
	Description    string `json:"description" validate:"max=255"`
}
