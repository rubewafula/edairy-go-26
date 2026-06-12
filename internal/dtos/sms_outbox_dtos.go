package dtos

import "time"

// CreateSMSOutboxRequest defines the request body for creating an SMS outbox entry.
type CreateSMSOutboxRequest struct {
	SMSCampaignID uint64 `json:"sms_campaign_id" binding:"required"`
	PhoneNo       string `json:"phone_no" binding:"required"`
	Message       string `json:"message" binding:"required"`
}

// UpdateSMSOutboxRequest defines the request body for updating an SMS outbox entry.
type UpdateSMSOutboxRequest struct {
	Status       *string    `json:"status"` // enum('pending','sent','delivered','failed')
	SentAt       *time.Time `json:"sent_at"`
	DeliveredAt  *time.Time `json:"delivered_at"`
	SMSMessageID *uint64    `json:"sms_message_id"`
	ContactID    *uint64    `json:"contact_id"`
	SMSContactID *uint64    `json:"sms_contact_id"`
}

// SMSOutboxResponse defines the response structure for an SMS outbox entry.
type SMSOutboxResponse struct {
	ID            uint64     `json:"id"`
	SMSCampaignID uint64     `json:"sms_campaign_id"`
	CampaignName  string     `json:"campaign_name"`
	PhoneNo       string     `json:"phone_no"`
	Status        string     `json:"status"`
	SMSMessageID  *uint64    `json:"sms_message_id"`
	MessageText   string     `json:"message_text"`
	SentAt        *time.Time `json:"sent_at"`
	DeliveredAt   *time.Time `json:"delivered_at"`
	ContactID     *uint64    `json:"contact_id"`
	SMSContactID  *uint64    `json:"sms_contact_id"`
	ContactName   string     `json:"contact_name"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
