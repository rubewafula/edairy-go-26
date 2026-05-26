package dtos

type CreateCustomerTypeRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}

type UpdateCustomerTypeRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}
