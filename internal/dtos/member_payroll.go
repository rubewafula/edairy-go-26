package dtos

import "time"

type CreateMemberPayrollRequest struct {
	DateOpened     string `json:"date_opened" validate:"required"`
	Description    string `json:"description"`
	PayDateRangeID uint64 `json:"pay_date_range_id" validate:"required"`
	FiscalPeriod   string `json:"fiscal_period" validate:"required"`
}

type UpdateMemberPayrollRequest struct {
	Description     string  `json:"description"`
	Status          string  `json:"status" validate:"omitempty,oneof=draft confirmed approved closed cancelled"`
	PhysicalPeriod  string  `json:"physical_period"`
	GrossKilos      float64 `json:"gross_kilos"`
	RejectKilos     float64 `json:"reject_kilos"`
	NetKilos        float64 `json:"net_kilos"`
	TotalDeductions float64 `json:"total_deductions"`
	GrossPay        float64 `json:"gross_pay"`
	NetPay          float64 `json:"net_pay"`
	TransportCost   float64 `json:"transport_cost"`
}

type MemberPayrollResponse struct {
	ID               uint64     `json:"id"`
	DateOpened       *time.Time `json:"date_opened"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	PayDateRangeID   *uint64    `json:"pay_date_range_id"`
	PayDateRangeName string     `json:"pay_date_range_name"`
	FiscalPeriod     string     `json:"fiscal_period"`
	GrossKilos       float64    `json:"gross_kilos"`
	RejectKilos      float64    `json:"reject_kilos"`
	NetKilos         float64    `json:"net_kilos"`
	TotalDeductions  float64    `json:"total_deductions"`
	GrossPay         float64    `json:"gross_pay"`
	NetPay           float64    `json:"net_pay"`
	TransportCost    float64    `json:"transport_cost"`
	PostedAt         *time.Time `json:"posted_at"`
	PostedByName     string     `json:"posted_by_name"`
	ConfirmedAt      *time.Time `json:"confirmed_at"`
	ConfirmedByName  string     `json:"confirmed_by_name"`
	ApprovedAt       *time.Time `json:"approved_at"`
	ApprovedByName   string     `json:"approved_by_name"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type MemberPayrollStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=confirmed approved closed cancelled"`
}
