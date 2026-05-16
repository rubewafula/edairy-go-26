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
	ID           uint64    `json:"id"`
	VendorName   string    `json:"vendor_name"`
	TotalAmount  float64   `json:"total_amount"`
	ItemCount    uint64    `json:"item_count"`
	Reference    string    `json:"reference"`
	SuppliedDate time.Time `json:"supplied_date"`
	Settled      bool      `json:"settled"`
	StoreName    string    `json:"store_name"`
	CreatedAt    time.Time `json:"created_at"`
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
	ID         uint64    `json:"id"`
	ItemID     uint64    `json:"item_id"`
	ItemName   string    `json:"item_name"`
	SupplyID   uint64    `json:"supply_id"`
	VendorName string    `json:"vendor_name"`
	Quantity   string    `json:"quantity"`
	Reason     string    `json:"reason"`
	CreatedBy  uint64    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SuppliedItemResponse struct {
	ID         uint64  `json:"id"`
	SupplyID   uint64  `json:"supply_id"`
	ItemID     uint64  `json:"item_id"`
	ItemName   string  `json:"item_name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}
