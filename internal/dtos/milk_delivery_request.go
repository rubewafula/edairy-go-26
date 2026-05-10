package dtos

import "time"

type CreateMilkDeliveryRequest struct {
	DeliveryNoteNumber string  `json:"DeliveryNoteNumber" validate:"required"`
	CustomerID         uint64  `json:"CustomerID" validate:"required"`
	QuantityAccepted   float64 `json:"QuantityAccepted" validate:"required"`
	Cooler             string  `json:"Cooler"`
	Invoiced           int     `json:"Invoiced"`
	TransactionDate    string  `json:"TransactionDate" validate:"required,datetime"`
	Amount             float64 `json:"Amount"`
	AmountPaid         float64 `json:"AmountPaid"`
	RouteID            uint64  `json:"RouteID"`
	Confirmed          int     `json:"Confirmed"`
	Processed          string  `json:"Processed"`
	TransporterID      uint64  `json:"TransporterID"`
}

type UpdateMilkDeliveryRequest struct {
	DeliveryNoteNumber string  `json:"DeliveryNoteNumber" validate:"required"`
	CustomerID         uint64  `json:"CustomerID" validate:"required"`
	QuantityAccepted   float64 `json:"QuantityAccepted" validate:"required"`
	Cooler             string  `json:"Cooler"`
	Invoiced           int     `json:"Invoiced"`
	TransactionDate    string  `json:"TransactionDate" validate:"required,datetime"`
	Amount             float64 `json:"Amount"`
	AmountPaid         float64 `json:"AmountPaid"`
	RouteID            uint64  `json:"RouteID"`
	Confirmed          int     `json:"Confirmed"`
	Processed          string  `json:"Processed"`
	TransporterID      uint64  `json:"TransporterID"`
}

type MilkDeliveryResponse struct {
	ID                 uint64    `json:"ID"`
	DeliveryNoteNumber string    `json:"DeliveryNoteNumber"`
	CustomerID         uint64    `json:"CustomerID"`
	CustomerName       string    `json:"CustomerName"`
	QuantityAccepted   float64   `json:"QuantityAccepted"`
	Cooler             string    `json:"Cooler"`
	Invoiced           int       `json:"Invoiced"`
	TransactionDate    time.Time `json:"TransactionDate"`
	Amount             float64   `json:"Amount"`
	AmountPaid         float64   `json:"AmountPaid"`
	RouteID            uint64    `json:"RouteID"`
	RouteName          string    `json:"RouteName"`
	Confirmed          int       `json:"Confirmed"`
	Processed          string    `json:"Processed"`
	TransporterID      uint64    `json:"TransporterID"`
	TransporterName    string    `json:"TransporterName"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
}
