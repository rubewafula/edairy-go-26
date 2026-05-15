package dtos

import "time"

type SuppliedItemRequest struct {
	ItemID    uint64  `json:"item_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required"`
	UnitPrice float64 `json:"unit_price" validate:"required"`
}

type CreateSupplyRequest struct {
	VendorID        uint64                `json:"vendor_id" validate:"required"`
	PaymentTypeID   *uint64               `json:"payment_type_id"`
	PurchaseOrderID *uint64               `json:"purchase_order_id"`
	Reference       string                `json:"reference"`
	Activity        string                `json:"activity"`
	SuppliedDate    string                `json:"supplied_date" validate:"required,datetime"`
	StoreID         *uint64               `json:"store_id"`
	PaymentTermID   *uint64               `json:"payment_term_id"`
	Items           []SuppliedItemRequest `json:"items" validate:"required,min=1"`
}

type UpdateSuppliedItemRequest struct {
	ItemID    uint64  `json:"item_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required"`
	UnitPrice float64 `json:"unit_price" validate:"required"`
}

type SupplyResponse struct {
	ID           uint64    `json:"ID"`
	VendorName   string    `json:"VendorName"`
	TotalAmount  float64   `json:"TotalAmount"`
	ItemCount    uint64    `json:"ItemCount"`
	Reference    string    `json:"Reference"`
	SuppliedDate time.Time `json:"SuppliedDate"`
	Settled      bool      `json:"Settled"`
	StoreName    string    `json:"StoreName"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

type CreateSupplyRejectRequest struct {
	ItemID   uint64 `json:"item_id" validate:"required"`
	SupplyID uint64 `json:"supply_id" validate:"required"`
	Quantity string `json:"quantity" validate:"required"`
	Reason   string `json:"reason" validate:"required"`
}

type UpdateSupplyRejectRequest struct {
	Quantity string `json:"quantity" validate:"required"`
	Reason   string `json:"reason" validate:"required"`
}

type SupplyRejectResponse struct {
	ID         uint64    `json:"ID"`
	ItemID     uint64    `json:"ItemID"`
	ItemName   string    `json:"ItemName"`
	SupplyID   uint64    `json:"SupplyID"`
	VendorName string    `json:"VendorName"`
	Quantity   string    `json:"Quantity"`
	Reason     string    `json:"Reason"`
	CreatedBy  uint64    `json:"CreatedBy"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
}

type SuppliedItemResponse struct {
	ID         uint64  `json:"ID"`
	SupplyID   uint64  `json:"SupplyID"`
	ItemID     uint64  `json:"ItemID"`
	ItemName   string  `json:"ItemName"`
	Quantity   int     `json:"Quantity"`
	UnitPrice  float64 `json:"UnitPrice"`
	TotalPrice float64 `json:"TotalPrice"`
}
