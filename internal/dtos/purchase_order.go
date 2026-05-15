package dtos

import "time"

type PurchaseOrderItemRequest struct {
	ItemID      uint64  `json:"item_id" validate:"required"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity" validate:"required"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
}

type CreatePurchaseOrderRequest struct {
	SupplierID      *uint64                    `json:"supplier_id"`
	SupplierQuoteID *uint64                    `json:"supplier_quote_id"`
	PoNumber        string                     `json:"po_number" validate:"required"`
	PoDate          string                     `json:"po_date" validate:"required,datetime"`
	Items           []PurchaseOrderItemRequest `json:"items" validate:"required,min=1"`
}

type UpdatePurchaseOrderRequest struct {
	SupplierID      *uint64 `json:"supplier_id"`
	SupplierQuoteID *uint64 `json:"supplier_quote_id"`
	PoNumber        string  `json:"po_number" validate:"required"`
	PoDate          string  `json:"po_date" validate:"required,datetime"`
	Status          string  `json:"status" validate:"required,oneof=draft pending approved rejected cancelled"`
	TotalAmount     float64 `json:"total_amount"` // This might be recalculated by the system, but included for completeness
}

type CreatePurchaseRequisitionRequest struct {
	RequisitionNo   string `json:"requisition_no" validate:"required"`
	RequisitionDate string `json:"requisition_date" validate:"required,datetime"`
	Description     string `json:"description"`
}

type UpdatePurchaseRequisitionRequest struct {
	RequisitionNo   string `json:"requisition_no" validate:"required"`
	RequisitionDate string `json:"requisition_date" validate:"required,datetime"`
	Description     string `json:"description"`
	Status          string `json:"status" validate:"required,oneof=draft pending approved rejected cancelled"`
}

type PurchaseOrderResponse struct {
	ID           uint64    `json:"ID"`
	PoNumber     string    `json:"PoNumber"`
	PoDate       time.Time `json:"PoDate"`
	SupplierName string    `json:"SupplierName"`
	Status       string    `json:"Status"`
	TotalAmount  float64   `json:"TotalAmount"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

type PurchaseRequisitionResponse struct {
	ID              uint64    `json:"ID"`
	RequisitionNo   string    `json:"RequisitionNo"`
	RequisitionDate time.Time `json:"RequisitionDate"`
	Description     string    `json:"Description"`
	Status          string    `json:"Status"`
	CreatedBy       uint64    `json:"CreatedBy"`
	CreatedAt       time.Time `json:"CreatedAt"`
}

type CreatePurchaseRequisitionItemRequest struct {
	PurchaseRequisitionID uint64  `json:"purchase_requisition_id" validate:"required"`
	ItemID                uint64  `json:"item_id" validate:"required"`
	Quantity              float64 `json:"quantity" validate:"required"`
}

type PurchaseRequisitionItemResponse struct {
	ID        uint64    `json:"ID"`
	ItemName  string    `json:"ItemName"`
	Quantity  float64   `json:"Quantity"`
	Status    string    `json:"Status"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type UpdatePurchaseRequisitionItemRequest struct {
	ItemID   uint64  `json:"item_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required"`
	Status   string  `json:"status"`
}

type UpdatePurchaseOrderItemRequest struct {
	ItemID      uint64  `json:"item_id" validate:"required"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity" validate:"required"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
}

type PurchaseOrderItemResponse struct {
	ID              uint64    `json:"ID"`
	PurchaseOrderID uint64    `json:"PurchaseOrderID"`
	ItemID          uint64    `json:"ItemID"`
	ItemName        string    `json:"ItemName"`
	Description     string    `json:"Description"`
	Quantity        float64   `json:"Quantity"`
	UnitPrice       float64   `json:"UnitPrice"`
	TotalPrice      float64   `json:"TotalPrice"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
