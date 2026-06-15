package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type UINotificationController struct {
	service *services.UINotificationService
}

func NewUINotificationController() *UINotificationController {
	return &UINotificationController{
		service: services.NewUINotificationService(),
	}
}

func (c *UINotificationController) GetMyNotifications(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	notifications, total, err := c.service.GetUserNotifications(userID, page, limit)
	if err != nil {
		log.Printf("[UINotificationController.GetMyNotifications] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  notifications,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (c *UINotificationController) MarkAsRead(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err := c.service.MarkAsRead(id, userID); err != nil {
		log.Printf("[UINotificationController.MarkAsRead] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (c *UINotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.MarkAllAsRead(userID); err != nil {
		log.Printf("[UINotificationController.MarkAllAsRead] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (c *UINotificationController) GetUnread(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	notifications, total, err := c.service.GetUserUnreadNotifications(userID, page, limit)
	if err != nil {
		log.Printf("[UINotificationController.GetUnread] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve unread notifications"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  notifications,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
