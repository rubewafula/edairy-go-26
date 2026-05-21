package apperrors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrPermissionExists = &AppError{
		Code:    "PERMISSION_EXISTS",
		Message: "Permission already exists",
	}

	ErrInternal = &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Something went wrong",
	}
	ErrRoleExists = &AppError{
		Code:    "ROLE_EXISTS",
		Message: "Role already exists",
	}
)
