package dtos

type CreateRouteCenterRequest struct {
	RouteID uint64 `json:"route_id" validate:"required"`
	Center  string `json:"center" validate:"required,max=96"`
}

type UpdateRouteCenterRequest struct {
	RouteID uint64 `json:"route_id" validate:"required"`
	Center  string `json:"center" validate:"required,max=96"`
}
