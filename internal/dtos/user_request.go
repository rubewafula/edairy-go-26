package dtos

type CreateUserRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	IsVerified bool   `json:"is_verified"`
}

type UpdateUserRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Email      string `json:"email" validate:"required,email"`
	IsVerified bool   `json:"is_verified"`
}
