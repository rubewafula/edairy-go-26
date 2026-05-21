package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerAssetRoutes(api *gin.RouterGroup) {
	assetCategoryController := controllers.NewAssetCategoryController()
	assetController := controllers.NewAssetController()
	assetAssignmentController := controllers.NewAssetAssignmentController()
	assetDepreciationController := controllers.NewAssetDepreciationController()

	// Asset Category Routes
	api.POST("/asset-categories", assetCategoryController.CreateCategory)
	api.GET("/asset-categories", assetCategoryController.GetCategories)
	api.GET("/asset-categories/:id", assetCategoryController.GetCategory)
	api.PUT("/asset-categories/:id", assetCategoryController.UpdateCategory)
	api.DELETE("/asset-categories/:id", assetCategoryController.DeleteCategory)

	// Fixed Asset Routes
	api.POST("/fixed-assets", assetController.CreateAsset)
	api.GET("/fixed-assets", assetController.GetAssets)
	api.GET("/fixed-assets/:id", assetController.GetAsset)
	api.PUT("/fixed-assets/:id", assetController.UpdateAsset)
	api.DELETE("/fixed-assets/:id", assetController.DeleteAsset)

	// Asset Assignment Routes
	api.POST("/asset-assignments", assetAssignmentController.CreateAssignment)
	api.GET("/asset-assignments", assetAssignmentController.GetAssignments)
	api.GET("/asset-assignments/:id", assetAssignmentController.GetAssignment)
	api.PUT("/asset-assignments/:id", assetAssignmentController.UpdateAssignment)
	api.DELETE("/asset-assignments/:id", assetAssignmentController.DeleteAssignment)

	// Asset Depreciation Routes
	api.POST("/asset-depreciation-entries", assetDepreciationController.CreateEntry)
	api.PUT("/asset-depreciation-entries/:id", assetDepreciationController.UpdateEntry)
	api.GET("/asset-depreciation-entries", assetDepreciationController.GetEntries)
	api.GET("/asset-depreciation-entries/:id", assetDepreciationController.GetEntry)
	api.DELETE("/asset-depreciation-entries/:id", assetDepreciationController.DeleteEntry)
}
