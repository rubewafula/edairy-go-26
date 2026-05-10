package dtos

import "time"

type DailyMilkVarianceResponse struct {
	ID               uint64    `json:"ID"`
	Transporter      string    `json:"Transporter"`
	Day              time.Time `json:"Day"`
	Month            string    `json:"Month"`
	FieldCollections float64   `json:"FieldCollections"`
	MCC              float64   `json:"MCC"`
	CashSales        float64   `json:"CashSales"`
	CreditSales      float64   `json:"CreditSales"`
	Rejects          float64   `json:"Rejects"`
	Balance          float64   `json:"Balance"`
}
