package dtos

import "time"

// SMS Group & Contact
type CreateSMSGroupRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type CreateSMSContactRequest struct {
	SMSGroupID  uint64 `json:"sms_group_id" validate:"required"`
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
	ID        uint64    `json:"ID"`
	Recipient string    `json:"Recipient"`
	Message   string    `json:"Message"`
	Status    string    `json:"Status"`
	CreatedAt time.Time `json:"CreatedAt"`
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
	ID           uint64    `json:"ID"`
	ProviderCode string    `json:"ProviderCode"`
	ProviderName string    `json:"ProviderName"`
	IsDefault    string    `json:"IsDefault"`
	Status       string    `json:"Status"`
	CreatedAt    time.Time `json:"CreatedAt"`
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
	ID           uint64    `json:"ID"`
	TemplateCode string    `json:"TemplateCode"`
	TemplateName string    `json:"TemplateName"`
	ModuleName   string    `json:"ModuleName"`
	Message      string    `json:"Message"`
	Status       string    `json:"Status"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

// SMS Campaign
type CreateSMSCampaignRequest struct {
	CampaignCode string `json:"campaign_code" validate:"required"`
	CampaignName string `json:"campaign_name" validate:"required"`
	Description  string `json:"description"`
	Message      string `json:"message" validate:"required"`
	SMSGroupID   uint64 `json:"sms_group_id" validate:"required"`
	ScheduledAt  string `json:"scheduled_at" validate:"omitempty,datetime"`
	SiteID       uint64 `json:"site_id"`
}

type SMSCampaignResponse struct {
	ID              uint64     `json:"ID"`
	CampaignCode    string     `json:"CampaignCode"`
	CampaignName    string     `json:"CampaignName"`
	GroupName       string     `json:"GroupName"`
	TotalRecipients int        `json:"TotalRecipients"`
	TotalSent       int        `json:"TotalSent"`
	Status          string     `json:"Status"`
	ScheduledAt     *time.Time `json:"ScheduledAt"`
	CreatedAt       time.Time  `json:"CreatedAt"`
}

type SMSCampaignRecipientResponse struct {
	ID              uint64     `json:"ID"`
	RecipientName   string     `json:"RecipientName"`
	PhoneNo         string     `json:"PhoneNo"`
	Status          string     `json:"Status"`
	SentAt          *time.Time `json:"SentAt"`
	DeliveredAt     *time.Time `json:"DeliveredAt"`
	ResponseMessage string     `json:"ResponseMessage"`
}

type SMSQueueResponse struct {
	ID                 uint64     `json:"ID"`
	SMSMessageID       uint64     `json:"SMSMessageID"`
	Processed          string     `json:"Processed"`
	ProcessingAttempts int        `json:"ProcessingAttempts"`
	LastAttemptAt      *time.Time `json:"LastAttemptAt"`
}
