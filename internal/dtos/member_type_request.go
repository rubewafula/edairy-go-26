package dtos

type CreateMemberTypeRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

type UpdateMemberTypeRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}
