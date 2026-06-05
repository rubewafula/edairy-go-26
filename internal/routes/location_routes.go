package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerLocationRoutes(api *gin.RouterGroup) {
	locationController := controllers.NewLocationController()

	// Administrative Location APIs
	api.POST("/locations", locationController.CreateLocation)
	api.GET("/locations", locationController.GetLocations)
	api.GET("/locations/:id", locationController.GetLocation)
	api.PUT("/locations/:id", locationController.UpdateLocation)
	api.DELETE("/locations/:id", locationController.DeleteLocation)

	api.GET("/counties", locationController.GetCounties)
	api.POST("/counties", locationController.CreateCounty)
	api.GET("/sub-counties", locationController.GetSubCounties)
	api.POST("/sub-counties", locationController.CreateSubCounty)
	api.GET("/wards", locationController.GetWards)
	api.POST("/wards", locationController.CreateWard)
}
