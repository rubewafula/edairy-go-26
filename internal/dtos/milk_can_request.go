package dtos

import "time"

type CreateMilkCanRequest struct {
	CanID      string  `json:"CanID" validate:"required"`
	CanType    string  `json:"CanType" validate:"required"`
	CanSize    float64 `json:"CanSize" validate:"required"`
	Units      string  `json:"Units"`
	TareWeight float64 `json:"TareWeight"`
	RouteID    uint64  `json:"RouteID" validate:"required"`
}

type UpdateMilkCanRequest struct {
	CanID      string  `json:"CanID" validate:"required"`
	CanType    string  `json:"CanType" validate:"required"`
	CanSize    float64 `json:"CanSize" validate:"required"`
	Units      string  `json:"Units"`
	TareWeight float64 `json:"TareWeight"`
	RouteID    uint64  `json:"RouteID" validate:"required"`
}

type MilkCanResponse struct {
	ID         uint64    `json:"ID"`
	CanID      string    `json:"CanID"`
	CanType    string    `json:"CanType"`
	CanSize    float64   `json:"CanSize"`
	Units      string    `json:"Units"`
	TareWeight float64   `json:"TareWeight"`
	RouteID    uint64    `json:"RouteID"`
	RouteName  string    `json:"RouteName"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
}
