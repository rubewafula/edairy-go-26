package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerFinanceRoutes(api *gin.RouterGroup) {
	financeDashboardController := controllers.NewFinanceDashboardController()
	cashTransactionController := controllers.NewCashTransactionController()

	accountController := controllers.NewAccountController()
	accountSubAccountController := controllers.NewAccountSubAccountController()
	transactionPostingRuleController := controllers.NewTransactionPostingRuleController()
	accountTypeController := controllers.NewAccountTypeController()
	accountCategoryController := controllers.NewAccountCategoryController()

	api.GET("/finance-dashboard", financeDashboardController.GetDashboard)

	api.POST("/cash-transactions", cashTransactionController.Create)
	api.GET("/cash-transactions", cashTransactionController.List)
	api.GET("/cash-transactions/:id", cashTransactionController.Get)
	api.PUT("/cash-transactions/:id", cashTransactionController.Update)
	api.DELETE("/cash-transactions/:id", cashTransactionController.Delete)

	// Account routes
	api.GET("/trial-balance", accountController.GetTrialBalance)
	api.GET("/profit-loss", accountController.GetProfitLoss)
	api.GET("/balance-sheet", accountController.GetBalanceSheet)
	api.POST("/accounts", accountController.CreateAccount)
	api.GET("/accounts", accountController.GetAccounts)
	api.GET("/accounts/:id", accountController.GetAccount)
	api.PUT("/accounts/:id", accountController.UpdateAccount)
	api.DELETE("/accounts/:id", accountController.DeleteAccount)

	// Account Sub-Account routes
	api.POST("/sub-accounts", accountSubAccountController.CreateAccountSubAccount)
	api.GET("/sub-accounts", accountSubAccountController.GetAccountSubAccounts)
	api.GET("/sub-accounts/:id", accountSubAccountController.GetAccountSubAccount)
	api.PUT("/sub-accounts/:id", accountSubAccountController.UpdateAccountSubAccount)
	api.DELETE("/sub-accounts/:id", accountSubAccountController.DeleteAccountSubAccount)

	// Transaction Posting Rule routes
	api.POST("/transaction-posting-rules", transactionPostingRuleController.CreateTransactionPostingRule)
	api.GET("/transaction-posting-rules", transactionPostingRuleController.GetTransactionPostingRules)
	api.GET("/transaction-posting-rules/:id", transactionPostingRuleController.GetTransactionPostingRule)
	api.PUT("/transaction-posting-rules/:id", transactionPostingRuleController.UpdateTransactionPostingRule)
	api.DELETE("/transaction-posting-rules/:id", transactionPostingRuleController.DeleteTransactionPostingRule)

	// Account Type routes
	api.POST("/account-types", accountTypeController.Create)
	api.GET("/account-types", accountTypeController.List)
	api.GET("/account-types/:id", accountTypeController.Get)
	api.PUT("/account-types/:id", accountTypeController.Update)
	api.DELETE("/account-types/:id", accountTypeController.Delete)

	// Account Category routes
	api.POST("/account-categories", accountCategoryController.Create)
	api.GET("/account-categories", accountCategoryController.List)
	api.GET("/account-categories/:id", accountCategoryController.Get)
	api.PUT("/account-categories/:id", accountCategoryController.Update)
	api.DELETE("/account-categories/:id", accountCategoryController.Delete)
}
