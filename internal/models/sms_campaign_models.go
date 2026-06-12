package models

import (
	"time"
)

// SMSOutbox represents an individual SMS message sent or to be sent as part of a campaign.
type SMSOutbox struct {
	BaseModel
	SMSCampaignID *uint64    `gorm:"column:sms_campaign_id" json:"sms_campaign_id"`
	PhoneNo       string     `gorm:"column:phone_no;size:45;not null" json:"phone_no"`
	Status        string     `gorm:"column:status;type:enum('pending','sent','delivered','failed');default:'pending'" json:"status"`
	SMSMessageID  *uint64    `gorm:"column:sms_message_id" json:"sms_message_id"`
	SentAt        *time.Time `gorm:"column:sent_at" json:"sent_at"`
	DeliveredAt   *time.Time `gorm:"column:delivered_at" json:"delivered_at"`
	SMSContactID  *uint64    `gorm:"column:sms_contact_id" json:"sms_contact_id"`
}

func (SMSOutbox) TableName() string {
	return "sms_outboxes"
}
