package dtos

type CreateTransporterRouteAssignmentRequest struct {
	TransporterID uint64 `json:"transporter_id" validate:"required"`
	RouteID       uint64 `json:"route_id" validate:"required"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Active        bool   `json:"active"`
}

type UpdateTransporterRouteAssignmentRequest struct {
	TransporterID uint64 `json:"transporter_id" validate:"required"`
	RouteID       uint64 `json:"route_id" validate:"required"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Active        bool   `json:"active"`
}
