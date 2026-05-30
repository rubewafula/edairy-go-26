package dtos

type CreateLocationRequest struct {
	Code         string `json:"code" validate:"required"`
	LocationName string `json:"location_name" validate:"required"`
}

type UpdateLocationRequest struct {
	Code         string `json:"code" validate:"required"`
	LocationName string `json:"location_name" validate:"required"`
}
