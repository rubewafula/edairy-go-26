package dtos

type CreateAssetCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=45"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateAssetCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=45"`
	Description string `json:"description" validate:"max=255"`
}
