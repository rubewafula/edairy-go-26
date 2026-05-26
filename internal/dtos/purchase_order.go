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
	PoDate          string                     `json:"po_date" validate:"required"`
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

type PurchaseRequisitionItemRequest struct {
	ItemID   uint64  `json:"item_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required"`
}

type CreatePurchaseRequisitionRequest struct {
	RequisitionNo   string                           `json:"requisition_no" validate:"required"`
	RequisitionDate string                           `json:"requisition_date" validate:"required"`
	Description     string                           `json:"description"`
	Status          string                           `json:"status"`
	Items           []PurchaseRequisitionItemRequest `json:"items" validate:"required,min=1"`
}

type UpdatePurchaseRequisitionRequest struct {
	RequisitionNo   string `json:"requisition_no" validate:"required"`
	RequisitionDate string `json:"requisition_date" validate:"required,datetime"`
	Description     string `json:"description"`
	Status          string `json:"status" validate:"required,oneof=draft pending approved rejected cancelled"`
}

type PurchaseOrderResponse struct {
	ID           uint64    `json:"id"`
	PoNumber     string    `json:"po_number"`
	PoDate       time.Time `json:"po_date"`
	SupplierName string    `json:"supplier_name"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
}

type PurchaseRequisitionResponse struct {
	ID              uint64    `json:"id"`
	RequisitionNo   string    `json:"requisition_no"`
	RequisitionDate time.Time `json:"requisition_date"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	CreatedBy       uint64    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreatePurchaseRequisitionItemRequest struct {
	PurchaseRequisitionID uint64  `json:"purchase_requisition_id" validate:"required"`
	ItemID                uint64  `json:"item_id" validate:"required"`
	Quantity              float64 `json:"quantity" validate:"required"`
}

type PurchaseRequisitionItemResponse struct {
	ID        uint64    `json:"id"`
	ItemName  string    `json:"item_name"`
	Quantity  float64   `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
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
	ID              uint64    `json:"id"`
	PurchaseOrderID uint64    `json:"purchase_order_id"`
	ItemID          uint64    `json:"item_id"`
	ItemName        string    `json:"item_name"`
	Description     string    `json:"description"`
	Quantity        float64   `json:"quantity"`
	UnitPrice       float64   `json:"unit_price"`
	TotalPrice      float64   `json:"total_price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
