package dtos

type CreateBankRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	SwiftCode   string `json:"swift_code" validate:"omitempty,max=20"`
	Description string `json:"description"`
}

type UpdateBankRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	SwiftCode   string `json:"swift_code" validate:"omitempty,max=20"`
	Description string `json:"description"`
}
