package dtos

type CreateTransportRateRequest struct {
	RouteID       uint64  `json:"route_id" validate:"required"`
	TransporterID uint64  `json:"transporter_id" validate:"required"`
	TransportRate float64 `json:"rate" validate:"required,min=0"`
	MemberID      uint64  `json:"member_id" validate:"required"`
	Status        string  `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateTransportRateRequest struct {
	RouteID       uint64  `json:"route_id" validate:"required"`
	TransporterID uint64  `json:"transporter_id" validate:"required"`
	TransportRate float64 `json:"rate" validate:"required,min=0"`
	MemberID      uint64  `json:"member_id" validate:"required"`
	Status        string  `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}
