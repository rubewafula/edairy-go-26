package dtos

type CreateTransporterBenefitRequest struct {
	Name        string  `json:"name" validate:"required"`
	MinQuantity uint64  `json:"min_quantity" validate:"required"`
	Rate        float64 `json:"rate" validate:"required"`
	RouteID     uint64  `json:"route_id" validate:"required"`
	Status      string  `json:"status"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}

type UpdateTransporterBenefitRequest struct {
	Name        string  `json:"name" validate:"required"`
	MinQuantity uint64  `json:"min_quantity" validate:"required"`
	Rate        float64 `json:"rate" validate:"required"`
	RouteID     uint64  `json:"route_id" validate:"required"`
	Status      string  `json:"status"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}
