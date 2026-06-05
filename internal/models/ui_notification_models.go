package models

import "time"

type UINotification struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID           uint64    `gorm:"column:user_id;not null;index:idx_ui_notifications_user" json:"user_id"`
	Title            string    `gorm:"column:title;size:255;not null" json:"title"`
	Message          string    `gorm:"column:message;type:text;not null" json:"message"`
	NotificationType string    `gorm:"column:notification_type;size:100;not null" json:"notification_type"`
	ReferenceID      *uint64   `gorm:"column:reference_id" json:"reference_id,omitempty"`
	ReferenceType    *string   `gorm:"column:reference_type;size:100" json:"reference_type,omitempty"`
	IsRead           bool      `gorm:"column:is_read;default:0;index:idx_ui_notifications_read" json:"is_read"`
	CreatedAt        time.Time `gorm:"column:created_at;index:idx_ui_notifications_created;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	ErrorLink        string    `gorm:"column:error_link;default:null" json:"error_link"`
	DownloadLink     string    `gorm:"column:download_link;default:null" json:"download_link"`
}

func (UINotification) TableName() string {
	return "ui_notifications"
}
