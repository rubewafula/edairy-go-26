package models

import "time"

// SMS Module
type SMSGroup struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type SMSContact struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	SMSGroupID  uint64    `gorm:"column:sms_group_id"`
	Name        string    `gorm:"column:name"`
	PhoneNumber string    `gorm:"column:phone_number"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type SMSMessage struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	SMSGroupID *uint64   `gorm:"column:sms_group_id"`
	Recipient  string    `gorm:"column:recipient"`
	Message    string    `gorm:"column:message"`
	Status     string    `gorm:"column:status;default:'pending'"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

type SMSQueue struct {
	ID                 uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	SMSMessageID       uint64     `gorm:"column:sms_message_id"`
	PriorityLevel      int        `gorm:"column:priority_level;default:1"`
	ScheduledAt        *time.Time `gorm:"column:scheduled_at"`
	Processed          string     `gorm:"type:enum('yes','no');default:'no';column:processed"`
	ProcessingAttempts int        `gorm:"column:processing_attempts;default:0"`
	LastAttemptAt      *time.Time `gorm:"column:last_attempt_at"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
}

func (SMSQueue) TableName() string {
	return "sms_queue"
}

type SMSProvider struct {
	BaseModel
	ProviderCode string `gorm:"uniqueIndex;column:provider_code"`
	ProviderName string `gorm:"column:provider_name"`
	ApiUrl       string `gorm:"column:api_url"`
	ApiKey       string `gorm:"column:api_key"`
	ApiSecret    string `gorm:"column:api_secret"`
	SenderID     string `gorm:"column:sender_id"`
	Username     string `gorm:"column:username"`
	Password     string `gorm:"column:password"`
	IsDefault    string `gorm:"type:enum('yes','no');default:'no';column:is_default"`
	Status       string `gorm:"type:enum('active','inactive');default:'active';column:status"`
	Notes        string `gorm:"column:notes"`
}

type SMSTemplate struct {
	BaseModel
	TemplateCode string `gorm:"uniqueIndex;column:template_code"`
	TemplateName string `gorm:"column:template_name"`
	ModuleName   string `gorm:"column:module_name"`
	Message      string `gorm:"column:message"`
	Variables    string `gorm:"column:variables"`
	Status       string `gorm:"type:enum('active','inactive');default:'active';column:status"`
}

type SMSCampaign struct {
	BaseModel
	CampaignCode    string     `gorm:"uniqueIndex;column:campaign_code"`
	CampaignName    string     `gorm:"column:campaign_name"`
	Description     string     `gorm:"column:description"`
	Message         string     `gorm:"column:message"`
	SMSGroupID      *uint64    `gorm:"column:sms_group_id"`
	TotalRecipients int        `gorm:"column:total_recipients;default:0"`
	TotalSent       int        `gorm:"column:total_sent;default:0"`
	TotalFailed     int        `gorm:"column:total_failed;default:0"`
	Status          string     `gorm:"type:enum('draft','scheduled','running','completed','cancelled');default:'draft';column:status"`
	ScheduledAt     *time.Time `gorm:"column:scheduled_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at"`
	SiteID          *uint64    `gorm:"column:site_id"`
}

type SMSCampaignRecipient struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	SMSCampaignID   uint64     `gorm:"column:sms_campaign_id"`
	RecipientName   string     `gorm:"column:recipient_name"`
	PhoneNo         string     `gorm:"column:phone_no"`
	Status          string     `gorm:"type:enum('pending','sent','delivered','failed');default:'pending';column:status"`
	SMSMessageID    *uint64    `gorm:"column:sms_message_id"`
	ResponseMessage string     `gorm:"column:response_message"`
	SentAt          *time.Time `gorm:"column:sent_at"`
	DeliveredAt     *time.Time `gorm:"column:delivered_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (AssetDepreciationEntry) TableName() string {
	return "asset_depreciation_entries"
}
