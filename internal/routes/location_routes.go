package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerLocationRoutes(api *gin.RouterGroup) {
	routeController := controllers.NewRouteController()
	routeCenterController := controllers.NewRouteCenterController()
	locationController := controllers.NewLocationController()

	// Route APIs
	api.POST("/routes", routeController.CreateRoute)
	api.GET("/routes", routeController.GetRoutes)
	api.GET("/routes/:id", routeController.GetRoute)
	api.PUT("/routes/:id", routeController.UpdateRoute)
	api.DELETE("/routes/:id", routeController.DeleteRoute)

	// Route Center APIs
	api.POST("/route-centers", routeCenterController.CreateCenter)
	api.GET("/route-centers", routeCenterController.GetCenters)
	api.GET("/route-centers/:id", routeCenterController.GetCenter)
	api.PUT("/route-centers/:id", routeCenterController.UpdateCenter)
	api.DELETE("/route-centers/:id", routeCenterController.DeleteCenter)

	// Administrative Location APIs
	api.GET("/counties", locationController.GetCounties)
	api.POST("/counties", locationController.CreateCounty)
	api.GET("/sub-counties", locationController.GetSubCounties)
	api.POST("/sub-counties", locationController.CreateSubCounty)
	api.GET("/wards", locationController.GetWards)
	api.POST("/wards", locationController.CreateWard)
}
