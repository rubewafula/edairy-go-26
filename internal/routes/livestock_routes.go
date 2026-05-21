package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerLivestockRoutes(api *gin.RouterGroup) {
	livestockController := controllers.NewLivestockController()

	api.POST("/livestocks", livestockController.CreateLivestocks)
	api.GET("/livestocks", livestockController.GetLivestocks) // Paginated list
	api.GET("/livestocks/:id", livestockController.GetLivestock)
	api.PUT("/livestocks/:id", livestockController.UpdateLivestocks)
	api.DELETE("/livestocks/:id", livestockController.DeleteLivestocks)

	api.POST("/livestock-categories", livestockController.CreateCategory)
	api.GET("/livestock-categories", livestockController.GetCategories) // Paginated list
	api.GET("/livestock-categories/:id", livestockController.GetCategory)
	api.PUT("/livestock-categories/:id", livestockController.UpdateCategory)
	api.DELETE("/livestock-categories/:id", livestockController.DeleteCategory)

	api.POST("/livestock-breeds", livestockController.CreateBreed)
	api.GET("/livestock-breeds", livestockController.GetBreeds) // Paginated list with category name
	api.GET("/livestock-breeds/:id", livestockController.GetBreed)
	api.PUT("/livestock-breeds/:id", livestockController.UpdateBreed)
	api.DELETE("/livestock-breeds/:id", livestockController.DeleteBreed)

	api.POST("/livestock-deaths", livestockController.CreateDeath)
	api.GET("/livestock-deaths", livestockController.GetDeaths) // Paginated list
	api.GET("/livestock-deaths/:id", livestockController.GetDeath)
	api.PUT("/livestock-deaths/:id", livestockController.UpdateDeath)
	api.DELETE("/livestock-deaths/:id", livestockController.DeleteDeath)

	api.POST("/livestock-feedings", livestockController.CreateFeeding)
	api.GET("/livestock-feedings", livestockController.GetFeedings) // Paginated list
	api.GET("/livestock-feedings/:id", livestockController.GetFeeding)
	api.PUT("/livestock-feedings/:id", livestockController.UpdateFeeding)
	api.DELETE("/livestock-feedings/:id", livestockController.DeleteFeeding)

	api.POST("/livestock-health-records", livestockController.CreateHealth)
	api.GET("/livestock-health-records", livestockController.GetHealthRecords) // Paginated list
	api.GET("/livestock-health-records/:id", livestockController.GetHealthRecord)
	api.PUT("/livestock-health-records/:id", livestockController.UpdateHealthRecord)
	api.DELETE("/livestock-health-records/:id", livestockController.DeleteHealthRecord)

	api.POST("/livestock-movements", livestockController.CreateMovement)
	api.GET("/livestock-movements", livestockController.GetMovements) // Paginated list
	api.GET("/livestock-movements/:id", livestockController.GetMovement)
	api.PUT("/livestock-movements/:id", livestockController.UpdateMovement)
	api.DELETE("/livestock-movements/:id", livestockController.DeleteMovement)

	api.POST("/livestock-photos", livestockController.CreatePhoto)
	api.GET("/livestock-photos", livestockController.GetPhotos) // Paginated list
	api.GET("/livestock-photos/:id", livestockController.GetPhoto)
	api.PUT("/livestock-photos/:id", livestockController.UpdatePhoto)
	api.DELETE("/livestock-photos/:id", livestockController.DeletePhoto)

	api.POST("/livestock-production", livestockController.CreateProduction)
	api.GET("/livestock-production", livestockController.GetProductionRecords) // Paginated list
	api.GET("/livestock-production/:id", livestockController.GetProductionRecord)
	api.PUT("/livestock-production/:id", livestockController.UpdateProductionRecord)
	api.DELETE("/livestock-production/:id", livestockController.DeleteProductionRecord)

	api.POST("/livestock-sales", livestockController.CreateSale)
	api.GET("/livestock-sales", livestockController.GetSales) // Paginated list
	api.GET("/livestock-sales/:id", livestockController.GetSale)
	api.PUT("/livestock-sales/:id", livestockController.UpdateSale)
	api.DELETE("/livestock-sales/:id", livestockController.DeleteSale)

	api.POST("/livestock-weight-records", livestockController.CreateWeight)
	api.GET("/livestock-weight-records", livestockController.GetWeightRecords) // Paginated list
	api.GET("/livestock-weight-records/:id", livestockController.GetWeightRecord)
	api.PUT("/livestock-weight-records/:id", livestockController.UpdateWeightRecord)
	api.DELETE("/livestock-weight-records/:id", livestockController.DeleteWeightRecord)
}
