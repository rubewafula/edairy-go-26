package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerShareRoutes(api *gin.RouterGroup) {
	walletTypeController := controllers.NewWalletTypeController()
	shareTypeController := controllers.NewShareTypeController()
	shareAccountController := controllers.NewShareAccountController()
	shareDividendController := controllers.NewShareDividendController()
	dividendDeclarationController := controllers.NewDividendDeclarationController()
	sharePaymentController := controllers.NewSharePaymentController()
	shareTransferController := controllers.NewShareTransferController()
	shareTransactionController := controllers.NewShareTransactionController()

	// Wallet Type Routes
	api.POST("/wallet-types", walletTypeController.CreateWalletType)
	api.GET("/wallet-types", walletTypeController.GetWalletTypes)
	api.GET("/wallet-types/:id", walletTypeController.GetWalletType)
	api.PUT("/wallet-types/:id", walletTypeController.UpdateWalletType)
	api.DELETE("/wallet-types/:id", walletTypeController.DeleteWalletType)

	// Share Type Routes
	api.POST("/share-types", shareTypeController.CreateShareType)
	api.GET("/share-types", shareTypeController.GetShareTypes)
	api.GET("/share-types/:id", shareTypeController.GetShareType)
	api.PUT("/share-types/:id", shareTypeController.UpdateShareType)
	api.DELETE("/share-types/:id", shareTypeController.DeleteShareType)

	// Share Account Routes
	api.POST("/share-accounts", shareAccountController.CreateAccount)
	api.GET("/share-accounts", shareAccountController.GetAccounts)
	api.GET("/share-accounts/:id", shareAccountController.GetAccount)
	api.PUT("/share-accounts/:id", shareAccountController.UpdateAccount)
	api.DELETE("/share-accounts/:id", shareAccountController.DeleteAccount)

	// Share Dividend Routes
	api.POST("/share-dividends", shareDividendController.CreateDividend)
	api.GET("/share-dividends", shareDividendController.GetDividends)
	api.GET("/share-dividends/:id", shareDividendController.GetDividend)
	api.PUT("/share-dividends/:id", shareDividendController.UpdateDividend)
	api.DELETE("/share-dividends/:id", shareDividendController.DeleteDividend)

	// Dividend Declaration Routes
	api.POST("/dividend-declarations", dividendDeclarationController.CreateDeclaration)
	api.GET("/dividend-declarations", dividendDeclarationController.GetDeclarations)
	api.GET("/dividend-declarations/:id", dividendDeclarationController.GetDeclaration)
	api.PUT("/dividend-declarations/:id", dividendDeclarationController.UpdateDeclaration)
	//api.DELETE("/dividend-declarations/:id", dividendDeclarationController.DeleteDividend)

	// Share Payment Routes
	api.POST("/share-payments", sharePaymentController.CreateSharePayment)
	api.GET("/share-payments", sharePaymentController.GetSharePayments)
	api.GET("/share-payments/:id", sharePaymentController.GetSharePayment)
	api.PUT("/share-payments/:id", sharePaymentController.UpdateSharePayment)
	api.DELETE("/share-payments/:id", sharePaymentController.DeleteSharePayment)

	// Share Transaction Routes
	api.POST("/share-transactions", shareTransactionController.CreateShareTransaction)
	api.GET("/share-transactions", shareTransactionController.GetShareTransactions)
	api.GET("/share-transactions/:id", shareTransactionController.GetShareTransaction)
	api.PUT("/share-transactions/:id", shareTransactionController.UpdateShareTransaction)
	api.DELETE("/share-transactions/:id", shareTransactionController.DeleteShareTransaction)

	// Share Transfer Routes
	api.POST("/share-transfers", shareTransferController.CreateShareTransfer)
	api.GET("/share-transfers", shareTransferController.GetShareTransfers)
	api.GET("/share-transfers/:id", shareTransferController.GetShareTransfer)
	api.PUT("/share-transfers/:id", shareTransferController.UpdateShareTransfer)
	api.DELETE("/share-transfers/:id", shareTransferController.DeleteShareTransfer)
}
