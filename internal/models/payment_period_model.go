package models

type PaymentPeriod struct {
	BaseModel
	Name          string `gorm:"column:name;type:enum('WEEKLY','BI-WEEKLY','MONTHLY')"`
	Description   string `gorm:"column:description;not null"`
	DefaultPeriod int    `gorm:"column:default_period;default:0"`
}

func (PaymentPeriod) TableName() string {
	return "payment_periods"
}
