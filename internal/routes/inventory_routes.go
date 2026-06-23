package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerInventoryRoutes(api *gin.RouterGroup) {
	storeDashboardController := controllers.NewStoreDashboardController()
	productGradeController := controllers.NewProductGradeController()
	itemCategoryController := controllers.NewItemCategoryController()
	storeController := controllers.NewStoreController()
	storeItemController := controllers.NewStoreItemController()
	storeSaleController := controllers.NewStoreSaleController()
	storeSaleItemController := controllers.NewStoreSaleItemController()
	storeInventoryController := controllers.NewStoreInventoryController()
	storeStockController := controllers.NewStoreStockController()
	storeStockTakingController := controllers.NewStoreStockTakingController()
	interStoreTransferController := controllers.NewInterStoreTransferController()
	interStoreTransferItemController := controllers.NewInterStoreTransferItemController()
	storeStockMovementController := controllers.NewStoreStockMovementController()
	storeStockMovementTypeController := controllers.NewStoreStockMovementTypeController()
	storeItemUnitController := controllers.NewStoreItemUnitController()

	api.GET("/store-dashboard", storeDashboardController.GetDashboard)

	// Product Grade Routes
	api.POST("/product-grades", productGradeController.CreateGrade)
	api.GET("/product-grades", productGradeController.GetGrades)
	api.GET("/product-grades/:id", productGradeController.GetGrade)
	api.PUT("/product-grades/:id", productGradeController.UpdateGrade)
	api.DELETE("/product-grades/:id", productGradeController.DeleteGrade)

	// Item Category Routes
	api.POST("/item-categories", itemCategoryController.CreateCategory)
	api.GET("/item-categories", itemCategoryController.GetCategories)
	api.GET("/item-categories/:id", itemCategoryController.GetCategory)
	api.PUT("/item-categories/:id", itemCategoryController.UpdateCategory)
	api.DELETE("/item-categories/:id", itemCategoryController.DeleteCategory)

	// Store Routes
	api.POST("/stores", storeController.CreateStore)
	api.GET("/stores", storeController.GetStores)
	api.GET("/stores/:id", storeController.GetStore)
	api.PUT("/stores/:id", storeController.UpdateStore)
	api.DELETE("/stores/:id", storeController.DeleteStore)

	// Store Item Routes
	api.POST("/store-items", storeItemController.CreateItem)
	api.GET("/store-items", storeItemController.GetItems)
	api.GET("/store-items/:id", storeItemController.GetItem)
	api.PUT("/store-items/:id", storeItemController.UpdateItem)
	api.DELETE("/store-items/:id", storeItemController.DeleteItem)

	// Store Sale Routes
	api.POST("/store-sales", storeSaleController.CreateSale)
	api.GET("/store-sales", storeSaleController.GetSales)
	api.GET("/store-sales/:id", storeSaleController.GetSale)
	api.PUT("/store-sales/:id", storeSaleController.UpdateSale)
	api.DELETE("/store-sales/:id", storeSaleController.DeleteSale)

	// Store Sale Item Routes
	api.POST("/store-sale-items", storeSaleItemController.CreateSaleItem)
	api.GET("/store-sale-items", storeSaleItemController.GetSaleItems)
	api.GET("/store-sale-items/:id", storeSaleItemController.GetSaleItem)
	api.PUT("/store-sale-items/:id", storeSaleItemController.UpdateSaleItem)
	api.DELETE("/store-sale-items/:id", storeSaleItemController.DeleteSaleItem)

	// Store Inventory Routes
	api.POST("/store-inventories", storeInventoryController.CreateInventory)
	api.GET("/store-inventories", storeInventoryController.GetInventories)
	api.GET("/store-inventories/:id", storeInventoryController.GetInventory)
	api.PUT("/store-inventories/:id", storeInventoryController.UpdateInventory)
	api.DELETE("/store-inventories/:id", storeInventoryController.DeleteInventory)

	// Store Stock Routes
	api.POST("/store-stocks", storeStockController.CreateStock)
	api.GET("/store-stocks", storeStockController.GetStocks)
	api.GET("/store-stocks/:id", storeStockController.GetStock)
	api.PUT("/store-stocks/:id", storeStockController.UpdateStock)
	api.DELETE("/store-stocks/:id", storeStockController.DeleteStock)
	api.POST("/store-stocks/import", storeController.ImportStoreStock)
	api.GET("/store-stocks/import-errors/:importid", storeController.GetImportErrors)

	// Store Stock Taking Routes
	api.POST("/store-stock-takings", storeStockTakingController.CreateStockTaking)
	api.GET("/store-stock-takings", storeStockTakingController.GetStockTakings)
	api.GET("/store-stock-takings/:id", storeStockTakingController.GetStockTaking)

	// Inter Store Transfer Routes
	api.POST("/inter-store-transfers", interStoreTransferController.CreateTransfer)
	api.GET("/inter-store-transfers", interStoreTransferController.GetTransfers)
	api.GET("/inter-store-transfers/:id", interStoreTransferController.GetTransfer)
	api.PUT("/inter-store-transfers/:id", interStoreTransferController.UpdateTransfer)
	api.DELETE("/inter-store-transfers/:id", interStoreTransferController.DeleteTransfer)

	// Inter Store Transfer Item Routes
	api.POST("/inter-store-transfer-items", interStoreTransferItemController.CreateTransferItem)
	api.GET("/inter-store-transfer-items", interStoreTransferItemController.GetTransferItems)
	api.GET("/inter-store-transfer-items/:id", interStoreTransferItemController.GetTransferItem)
	api.PUT("/inter-store-transfer-items/:id", interStoreTransferItemController.UpdateTransferItem)
	api.DELETE("/inter-store-transfer-items/:id", interStoreTransferItemController.DeleteTransferItem)

	// Store Stock Movement Routes
	api.POST("/store-stock-movements", storeStockMovementController.CreateMovement)
	api.GET("/store-stock-movements", storeStockMovementController.GetMovements)
	api.GET("/store-stock-movements/:id", storeStockMovementController.GetMovement)

	// Store Stock Movement Type Routes
	api.POST("/store-stock-movement-types", storeStockMovementTypeController.CreateMovementType)
	api.GET("/store-stock-movement-types", storeStockMovementTypeController.GetMovementTypes)
	api.GET("/store-stock-movement-types/:id", storeStockMovementTypeController.GetMovementType)
	api.PUT("/store-stock-movement-types/:id", storeStockMovementTypeController.UpdateMovementType)
	api.DELETE("/store-stock-movement-types/:id", storeStockMovementTypeController.DeleteMovementType)

	// Store Item Unit Routes
	api.POST("/store-item-units", storeItemUnitController.CreateUnit)
	api.GET("/store-item-units", storeItemUnitController.GetUnits)
	api.GET("/store-item-units/:id", storeItemUnitController.GetUnit)
	api.PUT("/store-item-units/:id", storeItemUnitController.UpdateUnit)
	api.DELETE("/store-item-units/:id", storeItemUnitController.DeleteUnit)
}
