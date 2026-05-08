package dtos

type CreateCustomerClassRequest struct {
	ClassCode   string `json:"class_code" validate:"required,max=50"`
	Description string `json:"description"`
}

type UpdateCustomerClassRequest struct {
	ClassCode   string `json:"class_code" validate:"required,max=50"`
	Description string `json:"description"`
}
