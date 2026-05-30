package controllers

import (
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (c *UINotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.MarkAllAsRead(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (c *UINotificationController) GetUnreadCount(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	count, err := c.service.GetUnreadCount(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"count": count})
}
