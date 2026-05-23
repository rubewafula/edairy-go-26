package dtos

import "time"

type CreateMilkDeliveryRequest struct {
	DeliveryNoteNumber string  `json:"delivery_note_number" validate:"required"`
	CustomerID         uint64  `json:"customer_id" validate:"required"`
	QuantityAccepted   float64 `json:"quantity_accepted" validate:"required"`
	Cooler             string  `json:"cooler"`
	Invoiced           int     `json:"invoiced"`
	TransactionDate    string  `json:"transaction_date" validate:"required"`
	Amount             float64 `json:"amount"`
	AmountPaid         float64 `json:"amount_paid"`
	RouteID            uint64  `json:"route_id"`
	Confirmed          int     `json:"confirmed"`
	Processed          string  `json:"processed"`
	TransporterID      uint64  `json:"transporter_id"`
}

type UpdateMilkDeliveryRequest struct {
	DeliveryNoteNumber string  `json:"delivery_note_number" validate:"required"`
	CustomerID         uint64  `json:"customer_id" validate:"required"`
	QuantityAccepted   float64 `json:"quantity_accepted" validate:"required"`
	Cooler             string  `json:"cooler"`
	Invoiced           int     `json:"invoiced"`
	TransactionDate    string  `json:"transaction_date" validate:"required"`
	Amount             float64 `json:"amount"`
	AmountPaid         float64 `json:"amount_paid"`
	RouteID            uint64  `json:"route_id"`
	Confirmed          int     `json:"confirmed"`
	Processed          string  `json:"processed"`
	TransporterID      uint64  `json:"transporter_id"`
}

type MilkDeliveryResponse struct {
	ID                 uint64    `json:"id"`
	DeliveryNoteNumber string    `json:"delivery_note_number"`
	CustomerID         uint64    `json:"customer_id"`
	CustomerName       string    `json:"customer_name"`
	QuantityAccepted   float64   `json:"quantity_accepted"`
	Cooler             string    `json:"cooler"`
	Invoiced           int       `json:"invoiced"`
	TransactionDate    time.Time `json:"transaction_date"`
	Amount             float64   `json:"amount"`
	AmountPaid         float64   `json:"amount_paid"`
	RouteID            uint64    `json:"route_id"`
	RouteName          string    `json:"route_name"`
	Confirmed          int       `json:"confirmed"`
	Processed          string    `json:"processed"`
	TransporterID      uint64    `json:"transporter_id"`
	TransporterName    string    `json:"transporter_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
