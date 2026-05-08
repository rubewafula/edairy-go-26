package dtos

type CreateSubRouteRequest struct {
	RouteID     uint64 `json:"route_id" validate:"required"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

type UpdateSubRouteRequest struct {
	RouteID     uint64 `json:"route_id" validate:"required"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}
