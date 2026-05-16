package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerLoanRoutes(api *gin.RouterGroup) {
	loanController := controllers.NewLoanController()
	loanManagementController := controllers.NewLoanManagementController()
	loanOrganizationProfileController := controllers.NewLoanOrganizationProfileController()

	// Loan Routes (General Loan)
	api.POST("/loans", loanController.CreateLoan)
	api.GET("/loans", loanController.GetLoans)
	api.GET("/loans/:id", loanController.GetLoan)
	api.PUT("/loans/:id", loanController.UpdateLoan)
	api.DELETE("/loans/:id", loanController.DeleteLoan)

	// Loan Account Routes
	api.POST("/loan-accounts", loanManagementController.CreateLoanAccount)
	api.GET("/loan-accounts", loanManagementController.GetLoanAccounts)
	api.GET("/loan-accounts/:id", loanManagementController.GetLoanAccount)
	api.PUT("/loan-accounts/:id", loanManagementController.UpdateLoanAccount)
	api.DELETE("/loan-accounts/:id", loanManagementController.DeleteLoanAccount)

	// Loan Callback Routes
	api.POST("/loan-callbacks", loanManagementController.CreateLoanCallback)
	api.GET("/loan-callbacks", loanManagementController.GetLoanCallbacks)
	api.GET("/loan-callbacks/:id", loanManagementController.GetLoanCallback)
	api.PUT("/loan-callbacks/:id", loanManagementController.UpdateLoanCallback)
	api.DELETE("/loan-callbacks/:id", loanManagementController.DeleteLoanCallback)

	// Loan Transaction Routes
	api.POST("/loan-transactions", loanManagementController.CreateLoanTransaction)
	api.GET("/loan-transactions", loanManagementController.GetLoanTransactions)
	api.GET("/loan-transactions/:id", loanManagementController.GetLoanTransaction)
	api.PUT("/loan-transactions/:id", loanManagementController.UpdateLoanTransaction)
	api.DELETE("/loan-transactions/:id", loanManagementController.DeleteLoanTransaction)
	api.GET("/loan-transactions/loan/:loan_id", loanManagementController.GetLoanTransactionsByLoanID) // Assuming this is GetTransactionsByLoanID

	// Loan Origination Callback Log Routes
	api.POST("/loan-origination-logs", loanManagementController.CreateLoanOriginationCallbackLog)
	api.GET("/loan-origination-logs", loanManagementController.GetLoanOriginationCallbackLogs)
	api.GET("/loan-origination-logs/:id", loanManagementController.GetLoanOriginationCallbackLog)
	api.PUT("/loan-origination-logs/:id", loanManagementController.UpdateLoanOriginationCallbackLog)
	api.DELETE("/loan-origination-logs/:id", loanManagementController.DeleteLoanOriginationCallbackLog)

	// Loan Organization Profile Routes
	api.POST("/loan-organization-profiles", loanOrganizationProfileController.CreateLoanOrganizationProfile)
	api.GET("/loan-organization-profiles", loanOrganizationProfileController.GetLoanOrganizationProfiles)
	api.GET("/loan-organization-profiles/:id", loanOrganizationProfileController.GetLoanOrganizationProfile)
	api.PUT("/loan-organization-profiles/:id", loanOrganizationProfileController.UpdateLoanOrganizationProfile)
	api.DELETE("/loan-organization-profiles/:id", loanOrganizationProfileController.DeleteLoanOrganizationProfile)
}
