package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerMilkRoutes(api *gin.RouterGroup) {
	produceDashboardController := controllers.NewProduceDashboardController()
	milkDeliveryShiftController := controllers.NewMilkDeliveryShiftController()
	coolerMilkCollectionController := controllers.NewCoolerMilkCollectionController()
	milkRejectController := controllers.NewMilkRejectController()
	milkJournalController := controllers.NewMilkJournalController()
	milkJournalEntryController := controllers.NewMilkJournalEntryController()
	milkDeliveryController := controllers.NewMilkDeliveryController()
	milkLocalSaleController := controllers.NewMilkLocalSaleController()
	dailyMilkVarianceController := controllers.NewDailyMilkVarianceController()
	defaultMilkRateController := controllers.NewDefaultMilkRateController()
	milkCanController := controllers.NewMilkCanController()
	canMovementController := controllers.NewCanMovementController()
	specialRateController := controllers.NewMilkSpecialRateController()

	api.GET("/produce-dashboard", produceDashboardController.GetDashboard)

	api.POST("/milk-delivery-shifts", milkDeliveryShiftController.CreateShift)
	api.GET("/milk-delivery-shifts", milkDeliveryShiftController.GetShifts)
	api.GET("/milk-delivery-shifts/:id", milkDeliveryShiftController.GetShift)
	api.PUT("/milk-delivery-shifts/:id", milkDeliveryShiftController.UpdateShift)
	api.DELETE("/milk-delivery-shifts/:id", milkDeliveryShiftController.DeleteShift)

	api.POST("/cooler-milk-collections", coolerMilkCollectionController.CreateCollection)
	api.GET("/cooler-milk-collections", coolerMilkCollectionController.GetCollections)
	api.GET("/cooler-milk-collections/:id", coolerMilkCollectionController.GetCollection)
	api.PUT("/cooler-milk-collections/:id", coolerMilkCollectionController.UpdateCollection)
	api.DELETE("/cooler-milk-collections/:id", coolerMilkCollectionController.DeleteCollection)

	api.POST("/milk-journals", milkJournalController.CreateMilkJournal)
	api.GET("/milk-journals", milkJournalController.GetMilkJournals)
	api.GET("/milk-journals/:id", milkJournalController.GetMilkJournal)
	api.PUT("/milk-journals/:id", milkJournalController.UpdateMilkJournal)
	api.GET("/milk-journals-today", milkJournalController.GetDailyJournals)
	api.DELETE("/milk-journals/:id", milkJournalController.DeleteMilkJournal)
	api.POST("/milk-journals/import", milkJournalController.ImportJournals)
	api.GET("/milk-journals/import-errors/:importid", milkJournalController.GetMilkJournalImportErrors)

	api.POST("/milk-journal-entries", milkJournalEntryController.CreateEntry)
	api.GET("/milk-journal-entries", milkJournalEntryController.GetEntries)
	api.GET("/milk-journal-entries/:id", milkJournalEntryController.GetEntry)
	api.PUT("/milk-journal-entries/:id", milkJournalEntryController.UpdateEntry)
	api.DELETE("/milk-journal-entries/:id", milkJournalEntryController.DeleteEntry)
	api.POST("/milk-journal-entries/upload", milkJournalEntryController.UploadEntries)
	api.GET("/stray-milk-collections", milkJournalEntryController.GetStrayEntries)

	api.POST("/milk-rejects", milkRejectController.CreateReject)
	api.GET("/milk-rejects", milkRejectController.GetRejects)
	api.GET("/milk-rejects/:id", milkRejectController.GetReject)
	api.DELETE("/milk-rejects/:id", milkRejectController.DeleteReject)

	api.POST("/milk-deliveries", milkDeliveryController.CreateDelivery)
	api.GET("/milk-deliveries", milkDeliveryController.GetDeliveries)
	api.GET("/milk-deliveries/:id", milkDeliveryController.GetDelivery)
	api.PUT("/milk-deliveries/:id", milkDeliveryController.UpdateDelivery)
	api.DELETE("/milk-deliveries/:id", milkDeliveryController.DeleteDelivery)

	api.POST("/milk-local-sales", milkLocalSaleController.CreateMilkLocalSale)
	api.GET("/milk-local-sales", milkLocalSaleController.GetMilkLocalSales)
	api.GET("/milk-local-sales/:id", milkLocalSaleController.GetMilkLocalSale)
	api.PUT("/milk-local-sales/:id", milkLocalSaleController.UpdateMilkLocalSale)
	api.DELETE("/milk-local-sales/:id", milkLocalSaleController.DeleteMilkLocalSale)

	api.GET("/daily-milk-variances", dailyMilkVarianceController.GetDailyVariances)

	api.POST("/default-milk-rates", defaultMilkRateController.CreateRate)
	api.GET("/default-milk-rates", defaultMilkRateController.GetRates)
	api.GET("/default-milk-rates/:id", defaultMilkRateController.GetRate)
	api.PUT("/default-milk-rates/:id", defaultMilkRateController.UpdateRate)
	api.DELETE("/default-milk-rates/:id", defaultMilkRateController.DeleteRate)

	api.POST("/milk-special-rates", specialRateController.Create)
	api.GET("/milk-special-rates", specialRateController.List)
	api.GET("/milk-special-rates/:id", specialRateController.Get)
	api.PUT("/milk-special-rates/:id", specialRateController.Update)
	api.DELETE("/milk-special-rates/:id", specialRateController.Delete)

	api.POST("/milk-cans", milkCanController.CreateMilkCan)
	api.GET("/milk-cans", milkCanController.GetMilkCans)
	api.GET("/milk-cans/:id", milkCanController.GetMilkCan)
	api.PUT("/milk-cans/:id", milkCanController.UpdateMilkCan)
	api.DELETE("/milk-cans/:id", milkCanController.DeleteMilkCan)

	api.POST("/can-movements", canMovementController.CreateMovement)
	api.GET("/can-movements", canMovementController.GetMovements)
	api.GET("/can-movements/:id", canMovementController.GetMovement)
	api.PUT("/can-movements/:id", canMovementController.UpdateMovement)
	api.DELETE("/can-movements/:id", canMovementController.DeleteMovement)
}
