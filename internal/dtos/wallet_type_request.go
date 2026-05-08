package dtos

type CreateWalletTypeRequest struct {
	Code        string `json:"code" validate:"required,max=255"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

type UpdateWalletTypeRequest struct {
	Code        string `json:"code" validate:"required,max=255"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}
