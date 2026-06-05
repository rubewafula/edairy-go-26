package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerUINotificationRoutes(api *gin.RouterGroup) {
	c := controllers.NewUINotificationController()

	api.GET("/ui-notifications", c.GetMyNotifications)
	api.GET("/ui-notifications/unread", c.GetUnread)
	api.PUT("/ui-notifications/:id/read", c.MarkAsRead)
	api.PUT("/ui-notifications/read-all", c.MarkAllAsRead)
}
