package models

import "time"

// SMSInAppConfiguration represents the configuration for in-app SMS activities.
type SMSInAppConfiguration struct {
	BaseModel
	ActivityCode        string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	ActivityDescription string     `gorm:"type:varchar(255);not null"`
	IsEnabled           bool       `gorm:"type:boolean;default:true"`
	CreatedAt           time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt           time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	DeletedAt           *time.Time `gorm:"index"`
}

func (SMSInAppConfiguration) TableName() string {
	return "sms_in_app_configurations"
}
