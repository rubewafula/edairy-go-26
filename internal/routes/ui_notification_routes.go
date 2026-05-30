package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerUINotificationRoutes(api *gin.RouterGroup) {
	c := controllers.NewUINotificationController()

	api.GET("/notifications", c.GetMyNotifications)
	api.GET("/notifications/unread-count", c.GetUnreadCount)
	api.PUT("/notifications/:id/read", c.MarkAsRead)
	api.PUT("/notifications/read-all", c.MarkAllAsRead)
}
