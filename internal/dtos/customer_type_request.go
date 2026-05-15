package dtos

type CreateCustomerTypeRequest struct {
	Name        string `json:"name" validate:"required,max=125"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateCustomerTypeRequest struct {
	Name        string `json:"name" validate:"required,max=125"`
	Description string `json:"description" validate:"max=255"`
}
