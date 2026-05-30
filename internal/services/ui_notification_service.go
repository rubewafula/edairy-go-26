package services

import (
	"strconv"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	socket "github.com/rubewafula/edairy-go-26/internal/socket-io"
)

type UINotificationService struct{}

func NewUINotificationService() *UINotificationService {
	return &UINotificationService{}
}

// CreateNotification creates a new UI notification.
func (s *UINotificationService) CreateNotification(userID uint64, req dtos.CreateUINotificationRequest) (*models.UINotification, error) {
	notification := &models.UINotification{
		UserID:           userID,
		Title:            req.Title,
		Message:          req.Message,
		NotificationType: req.NotificationType,
		ReferenceID:      req.ReferenceID,
		ReferenceType:    req.ReferenceType,
		IsRead:           false,
	}

	if err := db.DB.Create(notification).Error; err != nil {
		return nil, err
	}

	// Emit real-time notification to the user's private room
	socket.EmitNotification(strconv.FormatUint(userID, 10), notification)

	return notification, nil
}

// GetUserNotifications retrieves paginated notifications for a user.
func (s *UINotificationService) GetUserNotifications(userID uint64, page, limit int) ([]models.UINotification, int64, error) {
	var notifications []models.UINotification
	var total int64

	offset := (page - 1) * limit
	query := db.DB.Model(&models.UINotification{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, total, err
}

// MarkAsRead marks a specific notification as read.
func (s *UINotificationService) MarkAsRead(notificationID uint64, userID uint64) error {
	return db.DB.Model(&models.UINotification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead marks all unread notifications for a user as read.
func (s *UINotificationService) MarkAllAsRead(userID uint64) error {
	return db.DB.Model(&models.UINotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// GetUnreadCount returns the count of unread notifications for a user.
func (s *UINotificationService) GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.UINotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
