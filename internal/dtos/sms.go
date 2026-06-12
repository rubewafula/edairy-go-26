package dtos

import "time"

// SMS Group & Contact
type CreateSMSGroupRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type CreateSMSContactRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// SMS Messaging
type SendSMSRequest struct {
	Recipient string `json:"recipient" validate:"required"`
	Message   string `json:"message" validate:"required"`
	GroupID   uint64 `json:"sms_group_id"`
}

type SMSMessageResponse struct {
	ID        uint64    `json:"id"`
	Recipient string    `json:"recipient"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// SMS Provider
type CreateSMSProviderRequest struct {
	ProviderCode string `json:"provider_code" validate:"required"`
	ProviderName string `json:"provider_name" validate:"required"`
	ApiUrl       string `json:"api_url"`
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	SenderID     string `json:"sender_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	IsDefault    string `json:"is_default" validate:"oneof=yes no"`
	Status       string `json:"status" validate:"oneof=active inactive"`
	Notes        string `json:"notes"`
}

type SMSProviderResponse struct {
	ID           uint64    `json:"id"`
	ProviderCode string    `json:"provider_code"`
	ProviderName string    `json:"provider_name"`
	IsDefault    string    `json:"is_default"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// SMS Template
type CreateSMSTemplateRequest struct {
	TemplateCode string `json:"template_code" validate:"required"`
	TemplateName string `json:"template_name" validate:"required"`
	ModuleName   string `json:"module_name"`
	Message      string `json:"message" validate:"required"`
	Variables    string `json:"variables"`
	Status       string `json:"status" validate:"oneof=active inactive"`
}

type SMSTemplateResponse struct {
	ID           uint64    `json:"id"`
	TemplateCode string    `json:"template_code"`
	TemplateName string    `json:"template_name"`
	ModuleName   string    `json:"module_name"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// SMS Campaign
type CreateSMSCampaignRequest struct {
	CampaignCode string   `json:"campaign_code" validate:"required"`
	CampaignName string   `json:"campaign_name" validate:"required"`
	Description  string   `json:"description"`
	Message      string   `json:"message" validate:"required"`
	SMSGroupIDs  []uint64 `json:"sms_group_id" validate:"required"`
	ScheduledAt  string   `json:"scheduled_at" validate:"omitempty,datetime"`
	ExpiryDate   uint64   `json:"expiry_date"`
}

type SMSCampaignResponse struct {
	ID              uint64     `json:"id"`
	CampaignCode    string     `json:"campaign_code"`
	CampaignName    string     `json:"campaign_name"`
	GroupName       string     `json:"group_name"`
	TotalRecipients int        `json:"total_recipients"`
	TotalSent       int        `json:"total_sent"`
	Status          string     `json:"status"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type UpdateSMSCampaignRequest struct {
	CampaignCode string `json:"campaign_code"`
	CampaignName string `json:"campaign_name"`
	Description  string `json:"description"`
	Message      string `json:"message"`
	SMSGroupID   uint64 `json:"sms_group_id"`
	ScheduledAt  string `json:"scheduled_at" validate:"omitempty,datetime"`
	Status       string `json:"status" validate:"omitempty,oneof=draft pending approved sent failed"`
}

type SMSCampaignRecipientResponse struct {
	ID              uint64     `json:"id"`
	RecipientName   string     `json:"recipient_name"`
	PhoneNo         string     `json:"phone_no"`
	Status          string     `json:"status"`
	SentAt          *time.Time `json:"sent_at"`
	DeliveredAt     *time.Time `json:"delivered_at"`
	ResponseMessage string     `json:"response_message"`
}

type SMSQueueResponse struct {
	ID                 uint64     `json:"id"`
	SMSMessageID       uint64     `json:"sms_message_id"`
	Processed          string     `json:"processed"`
	ProcessingAttempts int        `json:"processing_attempts"`
	LastAttemptAt      *time.Time `json:"last_attempt_at"`
}

type CreateSMSCampaignRecipientRequest struct {
	SMSCampaignID uint64 `json:"sms_campaign_id" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNo       string `json:"phone_no" binding:"required"`
	Status        string `json:"status"`
}

type UpdateSMSCampaignRecipientRequest struct {
	RecipientName string `json:"recipient_name"`
	PhoneNo       string `json:"phone_no"`
	Status        string `json:"status"`
}

// SMS In-App Configuration
type CreateSMSInAppConfigurationRequest struct {
	ActivityCode        string `json:"activity_code" validate:"required"`
	ActivityDescription string `json:"activity_description"`
	IsEnabled           bool   `json:"is_enabled"`
}

type UpdateSMSInAppConfigurationRequest struct {
	ActivityDescription string `json:"activity_description"`
	IsEnabled           bool   `json:"is_enabled"`
}

type SMSInAppConfigurationResponse struct {
	ID                  uint64    `json:"id"`
	ActivityCode        string    `json:"activity_code"`
	ActivityDescription string    `json:"activity_description"`
	IsEnabled           bool      `json:"is_enabled"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
