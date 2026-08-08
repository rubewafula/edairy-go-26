package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerLoanModuleRoutes(api *gin.RouterGroup) {
	h := controllers.NewLoanModuleController()
	g := api.Group("/loan-module")

	// Products
	g.GET("/products", h.ListProducts)
	g.POST("/products", h.CreateProduct)
	g.GET("/products/:id", h.GetProduct)
	g.PUT("/products/:id", h.UpdateProduct)

	// Applications
	g.GET("/applications", h.ListApplications)
	g.POST("/applications", h.CreateApplication)
	g.GET("/applications/:id", h.GetApplication)
	g.POST("/applications/:id/submit", h.SubmitApplication)
	g.POST("/applications/:id/approve", h.ApproveApplication)
	g.POST("/applications/:id/reject", h.RejectApplication)

	// Contracts
	g.GET("/contracts", h.ListContracts)
	g.GET("/contracts/:id", h.GetContract)
	g.GET("/contracts/:id/disbursement-quote", h.GetDisbursementQuote)
	g.GET("/disbursement-channels", h.ListDisbursementChannels)
	g.POST("/disbursement-channels", h.CreateDisbursementChannel)
	g.GET("/disbursement-channels/:id", h.GetDisbursementChannel)
	g.PUT("/disbursement-channels/:id", h.UpdateDisbursementChannel)
	g.DELETE("/disbursement-channels/:id", h.DeleteDisbursementChannel)
	g.GET("/disbursement-channels/:id/config", h.GetDisbursementChannelConfig)
	g.PUT("/disbursement-channels/:id/config", h.UpdateDisbursementChannelConfig)
	g.POST("/contracts/:id/disburse", h.DisburseContract)
	g.POST("/disbursements/:id/confirm", h.ConfirmDisbursement)
	g.GET("/disbursements/:id/provider-status", h.GetDisbursementProviderStatus)
	g.GET("/contracts/:id/schedule", h.GetSchedule)
	g.POST("/contracts/:id/repayments", h.RecordRepayment)
	g.GET("/contracts/:id/statement", h.GetStatement)
	g.POST("/contracts/:id/settle", h.SettleContract)
	g.POST("/contracts/:id/restructure", h.RestructureContract)
	g.POST("/contracts/:id/write-off", h.WriteOffContract)

	// Interest & penalties
	g.POST("/interest/accrue", h.AccrueInterest)
	g.POST("/penalties/run", h.RunPenalties)

	// Reports
	g.GET("/reports/portfolio", h.PortfolioReport)
	g.GET("/reports/aging", h.AgingReport)
	g.GET("/reports/par", h.PARReport)
	g.GET("/reports/register", h.RegisterReport)
	g.GET("/reports/applications", h.ApplicationsPipelineReport)
	g.GET("/reports/disbursements", h.DisbursementsReport)
}
