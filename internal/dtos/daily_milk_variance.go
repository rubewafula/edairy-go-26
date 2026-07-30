package dtos

import "time"

type DailyMilkVarianceResponse struct {
	ID               uint64    `json:"id" gorm:"column:id"`
	Transporter      string    `json:"transporter" gorm:"column:transporter"`
	TransporterID    *uint64   `json:"transporter_id,omitempty" gorm:"column:transporter_id"`
	Day              time.Time `json:"day" gorm:"column:day"`
	Month            string    `json:"month" gorm:"column:month"`
	FieldCollections float64   `json:"field_collections" gorm:"column:field_collections"`
	MCC              float64   `json:"mcc" gorm:"column:mcc"`
	CashSales        float64   `json:"cash_sales" gorm:"column:cash_sales"`
	CreditSales      float64   `json:"credit_sales" gorm:"column:credit_sales"`
	Rejects          float64   `json:"rejects" gorm:"column:rejects"`
	Balance          float64   `json:"balance" gorm:"column:balance"`
}
