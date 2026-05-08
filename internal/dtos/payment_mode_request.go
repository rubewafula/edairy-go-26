package dtos

type CreatePaymentModeRequest struct {
	Code string `json:"code" validate:"omitempty,max=255"`
	Name string `json:"name" validate:"required,max=255"`
}

type UpdatePaymentModeRequest struct {
	Code string `json:"code" validate:"omitempty,max=255"`
	Name string `json:"name" validate:"required,max=255"`
}
