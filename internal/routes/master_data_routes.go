package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerMasterDataRoutes(api *gin.RouterGroup) {
	cattleBreedController := controllers.NewCattleBreedController()
	deductionTypeController := controllers.NewDeductionTypeController()
	deductionPricingRuleController := controllers.NewDeductionPricingRuleController()
	paymentModeController := controllers.NewPaymentModeController()

	activityLogController := controllers.NewActivityLogController()

	// Cattle Breed Routes
	api.POST("/cattle-breeds", cattleBreedController.CreateCattleBreed)
	api.GET("/cattle-breeds", cattleBreedController.GetCattleBreeds)
	api.GET("/cattle-breeds/:id", cattleBreedController.GetCattleBreed)
	api.PUT("/cattle-breeds/:id", cattleBreedController.UpdateCattleBreed)
	api.DELETE("/cattle-breeds/:id", cattleBreedController.DeleteCattleBreed)

	// Deduction Type Routes
	api.POST("/deduction-types", deductionTypeController.CreateDeductionType)
	api.GET("/deduction-types", deductionTypeController.GetDeductionTypes)
	api.GET("/deduction-types/:id", deductionTypeController.GetDeductionType)
	api.PUT("/deduction-types/:id", deductionTypeController.UpdateDeductionType)
	api.DELETE("/deduction-types/:id", deductionTypeController.DeleteDeductionType)

	// Deduction Pricing Rule Routes
	api.POST("/deduction-pricing-rules", deductionPricingRuleController.CreateRule)
	api.GET("/deduction-pricing-rules", deductionPricingRuleController.GetRules)
	api.GET("/deduction-pricing-rules/:id", deductionPricingRuleController.GetRule)
	api.PUT("/deduction-pricing-rules/:id", deductionPricingRuleController.UpdateRule)
	api.DELETE("/deduction-pricing-rules/:id", deductionPricingRuleController.DeleteRule)

	// Payment Mode Routes
	api.POST("/payment-modes", paymentModeController.CreatePaymentMode)
	api.GET("/payment-modes", paymentModeController.GetPaymentModes)
	api.GET("/payment-modes/:id", paymentModeController.GetPaymentMode)
	api.PUT("/payment-modes/:id", paymentModeController.UpdatePaymentMode)
	api.DELETE("/payment-modes/:id", paymentModeController.DeletePaymentMode)

	// Activity Log Routes
	api.GET("/activity-logs", activityLogController.GetLogs)
	api.GET("/activity-logs/:id", activityLogController.GetLog)

}
