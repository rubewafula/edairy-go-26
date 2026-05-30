package dtos

import "time"

type UINotificationResponse struct {
	ID               uint64    `json:"id"`
	UserID           uint64    `json:"user_id"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	NotificationType string    `json:"notification_type"`
	ReferenceID      *uint64   `json:"reference_id,omitempty"`
	ReferenceType    *string   `json:"reference_type,omitempty"`
	IsRead           bool      `json:"is_read"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateUINotificationRequest struct {
	Title            string  `json:"title" validate:"required,max=255"`
	Message          string  `json:"message" validate:"required"`
	NotificationType string  `json:"notification_type" validate:"required,max=100"`
	ReferenceID      *uint64 `json:"reference_id"`
	ReferenceType    *string `json:"reference_type"`
}
