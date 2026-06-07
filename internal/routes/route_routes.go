package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerRouteRoutes(api *gin.RouterGroup) {
	routeController := controllers.NewRouteController()

	api.POST("/routes", routeController.CreateRoute)
	api.GET("/routes", routeController.GetRoutes)
	api.GET("/routes/:id", routeController.GetRoute)
	api.PUT("/routes/:id", routeController.UpdateRoute)
	api.DELETE("/routes/:id", routeController.DeleteRoute)
}
