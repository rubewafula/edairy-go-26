package dtos

import "time"

type DailyMilkVarianceResponse struct {
	ID               uint64    `json:"id"`
	Transporter      string    `json:"transporter"`
	Day              time.Time `json:"day"`
	Month            string    `json:"month"`
	FieldCollections float64   `json:"field_collections"`
	MCC              float64   `json:"mcc"`
	CashSales        float64   `json:"cash_sales"`
	CreditSales      float64   `json:"credit_sales"`
	Rejects          float64   `json:"rejects"`
	Balance          float64   `json:"balance"`
}
