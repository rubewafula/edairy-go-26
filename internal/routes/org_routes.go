package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerOrgRoutes(api *gin.RouterGroup) {
	organizationAddressController := controllers.NewOrganizationAddressController()
	organizationDocumentController := controllers.NewOrganizationDocumentController()
	organizationBankController := controllers.NewOrganizationBankController()
	organizationLeadershipController := controllers.NewOrganizationLeadershipController()
	organizationWalletController := controllers.NewOrganizationWalletController()
	organizationKybCommentController := controllers.NewOrganizationKybCommentController()
	userStoreAssignmentController := controllers.NewUserStoreAssignmentController()
	paymentPeriodController := controllers.NewPaymentPeriodController()

	// Organization Address Routes
	api.POST("/organization-addresses", organizationAddressController.CreateAddress)
	api.GET("/organization-addresses", organizationAddressController.GetAddresses)
	api.GET("/organization-addresses/:id", organizationAddressController.GetAddress)
	api.PUT("/organization-addresses/:id", organizationAddressController.UpdateAddress)
	api.DELETE("/organization-addresses/:id", organizationAddressController.DeleteAddress)

	// Organization Bank Routes
	api.POST("/organization-banks", organizationBankController.CreateBank)
	api.GET("/organization-banks", organizationBankController.GetBanks)
	api.GET("/organization-banks/:id", organizationBankController.GetBank)
	api.PUT("/organization-banks/:id", organizationBankController.UpdateBank)
	api.DELETE("/organization-banks/:id", organizationBankController.DeleteBank)

	// Organization Wallet Routes
	api.POST("/organization-wallets", organizationWalletController.CreateWallet)
	api.GET("/organization-wallets", organizationWalletController.GetWallets)
	api.GET("/organization-wallets/:id", organizationWalletController.GetWallet)
	api.PUT("/organization-wallets/:id", organizationWalletController.UpdateWallet)
	api.DELETE("/organization-wallets/:id", organizationWalletController.DeleteWallet)

	// Organization Leadership Routes
	api.POST("/organization-leaderships", organizationLeadershipController.CreateLeadership)
	api.GET("/organization-leaderships", organizationLeadershipController.GetLeaderships)
	api.GET("/organization-leaderships/:id", organizationLeadershipController.GetLeadership)
	api.PUT("/organization-leaderships/:id", organizationLeadershipController.UpdateLeadership)
	api.DELETE("/organization-leaderships/:id", organizationLeadershipController.DeleteLeadership)
	api.GET("/organization-leaderships/national-id/:id_no", organizationLeadershipController.GetLeadershipByNationalID)

	// Organization Document Routes
	api.POST("/organization-documents", organizationDocumentController.CreateDocument)
	api.GET("/organization-documents", organizationDocumentController.GetDocuments)
	api.GET("/organization-documents/:id", organizationDocumentController.GetDocument)
	api.PUT("/organization-documents/:id", organizationDocumentController.UpdateDocument)
	api.DELETE("/organization-documents/:id", organizationDocumentController.DeleteDocument)
	api.GET("/organization-documents/astra/:astra_id", organizationDocumentController.GetDocumentsByAstraID)

	// Organization KYB Comment Routes
	api.POST("/organization-kyb-comments", organizationKybCommentController.CreateComment)
	api.GET("/organization-kyb-comments", organizationKybCommentController.GetComments)
	api.GET("/organization-kyb-comments/:id", organizationKybCommentController.GetComment)
	api.PUT("/organization-kyb-comments/:id", organizationKybCommentController.UpdateComment)
	api.DELETE("/organization-kyb-comments/:id", organizationKybCommentController.DeleteComment)
	api.GET("/organization-kyb-comments/iteration/:iteration", organizationKybCommentController.GetCommentsByIteration)

	// User Store Assignment Routes
	api.POST("/user-store-assignments", userStoreAssignmentController.CreateAssignment)
	api.GET("/user-store-assignments", userStoreAssignmentController.GetAssignments)
	api.GET("/user-store-assignments/:id", userStoreAssignmentController.GetAssignment)
	api.PUT("/user-store-assignments/:id", userStoreAssignmentController.UpdateAssignment)
	api.DELETE("/user-store-assignments/:id", userStoreAssignmentController.DeleteAssignment)

	// Payment Period Routes
	api.POST("/payment-periods", paymentPeriodController.Create)
	api.GET("/payment-periods", paymentPeriodController.List)
	api.GET("/payment-periods/:id", paymentPeriodController.Get)
	api.PUT("/payment-periods/:id", paymentPeriodController.Update)
	api.DELETE("/payment-periods/:id", paymentPeriodController.Delete)
}
