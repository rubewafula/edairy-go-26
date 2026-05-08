package dtos

type CreateCattleBreedRequest struct {
	Code string `json:"code" validate:"required,max=255"`
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateCattleBreedRequest struct {
	Code string `json:"code" validate:"required,max=255"`
	Name string `json:"name" validate:"required,max=255"`
}
