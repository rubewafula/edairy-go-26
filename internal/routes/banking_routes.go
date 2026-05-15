package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerBankingRoutes(api *gin.RouterGroup) {
	bankController := controllers.NewBankController()
	bankBranchController := controllers.NewBankBranchController()

	// Bank Routes
	api.POST("/banks", bankController.CreateBank)
	api.GET("/banks", bankController.GetBanks)
	api.GET("/banks/:id", bankController.GetBank)
	api.PUT("/banks/:id", bankController.UpdateBank)
	api.DELETE("/banks/:id", bankController.DeleteBank)

	// Bank Branch Routes
	api.POST("/bank-branches", bankBranchController.CreateBankBranch)
	api.GET("/bank-branches", bankBranchController.GetBankBranches)
	api.GET("/bank-branches/:id", bankBranchController.GetBankBranch)
	api.PUT("/bank-branches/:id", bankBranchController.UpdateBankBranch)
	api.DELETE("/bank-branches/:id", bankBranchController.DeleteBankBranch)
}
