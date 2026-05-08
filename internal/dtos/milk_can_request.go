package dtos

type CreateMilkCanRequest struct {
	CanID      string  `json:"can_id" validate:"required,max=255"`
	CanType    string  `json:"can_type" validate:"required,max=255"`
	CanSize    float64 `json:"can_size" validate:"required,min=0"`
	TareWeight float64 `json:"tare_weight" validate:"required,min=0"`
	RouteID    uint64  `json:"route_id" validate:"required"`
}

type UpdateMilkCanRequest struct {
	CanID      string  `json:"can_id" validate:"required,max=255"`
	CanType    string  `json:"can_type" validate:"required,max=255"`
	CanSize    float64 `json:"can_size" validate:"required,min=0"`
	TareWeight float64 `json:"tare_weight" validate:"required,min=0"`
	RouteID    uint64  `json:"route_id" validate:"required"`
}