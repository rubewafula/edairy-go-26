package dtos

type CreateRouteRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Code        string `json:"code" validate:"required,max=50"`
	LocationID  uint64 `json:"location_id" validate:"required"`
}

type UpdateRouteRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Code        string `json:"code" validate:"required,max=50"`
	LocationID  uint64 `json:"location_id" validate:"required"`
}
