package models

import (
	"time"

	"gorm.io/gorm"
)

// SMSContactGroup represents the many-to-many relationship between SMS groups and contacts.
type SMSContactGroup struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	SMSGroupID   uint64         `gorm:"column:sms_group_id;not null;index" json:"sms_group_id"`
	SMSContactID uint64         `gorm:"column:sms_contact_id;not null;index" json:"sms_contact_id"`
	Status       string         `gorm:"column:status;type:enum('active','inactive');default:'active'" json:"status"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Associations are often useful for preloading related data
	SMSGroup   *SMSGroup   `gorm:"foreignKey:SMSGroupID" json:"sms_group,omitempty"`
	SMSContact *SMSContact `gorm:"foreignKey:SMSContactID" json:"sms_contact,omitempty"`
}

func (SMSContactGroup) TableName() string {
	return "sms_contact_groups"
}
