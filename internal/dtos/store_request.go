package dtos

type CreateStoreRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

type UpdateStoreRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}