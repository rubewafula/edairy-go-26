package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
	"github.com/rubewafula/edairy-go-26/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rubewafula/edairy-go-26/docs"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware to allow cross-origin requests
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	authController := controllers.NewAuthController()
	memberController := controllers.NewMemberController()
	bankController := controllers.NewBankController()
	bankBranchController := controllers.NewBankBranchController()
	memberBankAccountController := controllers.NewMemberBankAccountController()
	milkDeliveryShiftController := controllers.NewMilkDeliveryShiftController()
	transporterController := controllers.NewTransporterController()
	subRouteController := controllers.NewSubRouteController()
	milkJournalController := controllers.NewMilkJournalController()
	milkCanController := controllers.NewMilkCanController()
	storeController := controllers.NewStoreController()
	memberTypeController := controllers.NewMemberTypeController()
	routeController := controllers.NewRouteController()
	customerController := controllers.NewCustomerController()
	customerClassController := controllers.NewCustomerClassController()
	customerMilkRateController := controllers.NewCustomerMilkRateController()
	customerPayDateRangeController := controllers.NewCustomerPayDateRangeController()
	memberDependantController := controllers.NewMemberDependantController()
	loanController := controllers.NewLoanController()
	transportRateController := controllers.NewTransportRateController()
	trainingController := controllers.NewTrainingController()
	trainingSessionController := controllers.NewTrainingSessionController()
	trainingAttendeeController := controllers.NewTrainingAttendeeController()
	exchangeVisitController := controllers.NewExchangeVisitController()
	exchangeVisitAttendeeController := controllers.NewExchangeVisitAttendeeController()
	paymentModeController := controllers.NewPaymentModeController()
	routeCenterController := controllers.NewRouteCenterController()

	adminDashboardController := controllers.NewAdminDashboardController()

	jwtSecret := os.Getenv("JWT_SECRET")

	authMiddleware := middleware.AuthMiddleware([]byte(jwtSecret))

	r.POST("/api/signup", authController.Signup)
	r.GET("/api/verify", authController.Verify)
	r.POST("/api/login", authController.Login)

	r.POST("/api/forgot-password", authController.ForgotPassword)
	r.POST("/api/reset-password", authController.ResetPassword)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	api.Use(authMiddleware)
	api.Use(middleware.RequireAutoPermission())
	{
		//auth routes

		api.POST("/change-password", authController.ChangePassword)
		//Dashboard stats
		api.GET("/admin-dashboard", adminDashboardController.GetDashboard)
		//Member routes
		api.POST("/member/create", memberController.CreateMember)
		api.GET("/members", memberController.GetMembers)
		api.PUT("/member/:id", memberController.UpdateMember)
		api.DELETE("/member/:id", memberController.DeleteMember)
		api.GET("/member/:id", memberController.GetMember)

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

		// Bank Account Routes
		api.POST("/member-bank-accounts", memberBankAccountController.CreateAccount)
		api.GET("/member-bank-accounts", memberBankAccountController.GetAccounts)
		api.GET("/member-bank-accounts/:id", memberBankAccountController.GetAccount)
		api.PUT("/member-bank-accounts/:id", memberBankAccountController.UpdateAccount)
		api.DELETE("/member-bank-accounts/:id", memberBankAccountController.DeleteAccount)

		// Milk Delivery Shift Routes
		api.POST("/milk-delivery-shifts", milkDeliveryShiftController.CreateShift)
		api.GET("/milk-delivery-shifts", milkDeliveryShiftController.GetShifts)
		api.GET("/milk-delivery-shifts/:id", milkDeliveryShiftController.GetShift)
		api.PUT("/milk-delivery-shifts/:id", milkDeliveryShiftController.UpdateShift)
		api.DELETE("/milk-delivery-shifts/:id", milkDeliveryShiftController.DeleteShift)

		// Transporter Routes
		api.POST("/transporters", transporterController.CreateTransporter)
		api.GET("/transporters", transporterController.GetTransporters)
		api.GET("/transporters/:id", transporterController.GetTransporter)
		api.PUT("/transporters/:id", transporterController.UpdateTransporter)
		api.DELETE("/transporters/:id", transporterController.DeleteTransporter)

		// SubRoute Routes
		api.POST("/sub-routes", subRouteController.CreateSubRoute)
		api.GET("/sub-routes", subRouteController.GetSubRoutes)
		api.GET("/sub-routes/:id", subRouteController.GetSubRoute)
		api.PUT("/sub-routes/:id", subRouteController.UpdateSubRoute)
		api.DELETE("/sub-routes/:id", subRouteController.DeleteSubRoute)

		// Milk Journal Routes
		api.POST("/milk-journals", milkJournalController.CreateMilkJournal)
		api.GET("/milk-journals", milkJournalController.GetMilkJournals)
		api.GET("/milk-journals/:id", milkJournalController.GetMilkJournal)
		api.PUT("/milk-journals/:id", milkJournalController.UpdateMilkJournal)
		api.DELETE("/milk-journals/:id", milkJournalController.DeleteMilkJournal)

		// Milk Can Routes
		api.POST("/milk-cans", milkCanController.CreateMilkCan)
		api.GET("/milk-cans", milkCanController.GetMilkCans)
		api.GET("/milk-cans/:id", milkCanController.GetMilkCan)
		api.PUT("/milk-cans/:id", milkCanController.UpdateMilkCan)
		api.DELETE("/milk-cans/:id", milkCanController.DeleteMilkCan)

		// Store Routes
		api.POST("/stores", storeController.CreateStore)
		api.GET("/stores", storeController.GetStores)
		api.GET("/stores/:id", storeController.GetStore)
		api.PUT("/stores/:id", storeController.UpdateStore)
		api.DELETE("/stores/:id", storeController.DeleteStore)

		// Member Type Routes
		api.POST("/member-types", memberTypeController.CreateMemberType)
		api.GET("/member-types", memberTypeController.GetMemberTypes)
		api.GET("/member-types/:id", memberTypeController.GetMemberType)
		api.PUT("/member-types/:id", memberTypeController.UpdateMemberType)
		api.DELETE("/member-types/:id", memberTypeController.DeleteMemberType)

		// Route Routes
		api.POST("/routes", routeController.CreateRoute)
		api.GET("/routes", routeController.GetRoutes)
		api.GET("/routes/:id", routeController.GetRoute)
		api.PUT("/routes/:id", routeController.UpdateRoute)
		api.DELETE("/routes/:id", routeController.DeleteRoute)

		// Wallet Type Routes
		api.POST("/wallet-types", controllers.NewWalletTypeController().CreateWalletType)
		api.GET("/wallet-types", controllers.NewWalletTypeController().GetWalletTypes)
		api.GET("/wallet-types/:id", controllers.NewWalletTypeController().GetWalletType)
		api.PUT("/wallet-types/:id", controllers.NewWalletTypeController().UpdateWalletType)
		api.DELETE("/wallet-types/:id", controllers.NewWalletTypeController().DeleteWalletType)

		// Cattle Breed Routes
		api.POST("/cattle-breeds", controllers.NewCattleBreedController().CreateCattleBreed)
		api.GET("/cattle-breeds", controllers.NewCattleBreedController().GetCattleBreeds)
		api.GET("/cattle-breeds/:id", controllers.NewCattleBreedController().GetCattleBreed)
		api.PUT("/cattle-breeds/:id", controllers.NewCattleBreedController().UpdateCattleBreed)
		api.DELETE("/cattle-breeds/:id", controllers.NewCattleBreedController().DeleteCattleBreed)

		// Customer Routes
		api.POST("/customers", customerController.CreateCustomer)
		api.GET("/customers", customerController.GetCustomers)
		api.GET("/customers/:id", customerController.GetCustomer)
		api.PUT("/customers/:id", customerController.UpdateCustomer)
		api.DELETE("/customers/:id", customerController.DeleteCustomer)

		// Customer Class Routes
		api.POST("/customer-classes", customerClassController.CreateClass)
		api.GET("/customer-classes", customerClassController.GetClasses)
		api.GET("/customer-classes/:id", customerClassController.GetClass)
		api.PUT("/customer-classes/:id", customerClassController.UpdateClass)
		api.DELETE("/customer-classes/:id", customerClassController.DeleteClass)

		// Customer Milk Rate Routes
		api.POST("/customer-milk-rates", customerMilkRateController.CreateRate)
		api.GET("/customer-milk-rates", customerMilkRateController.GetRates)
		api.GET("/customer-milk-rates/:id", customerMilkRateController.GetRate)
		api.PUT("/customer-milk-rates/:id", customerMilkRateController.UpdateRate)
		api.DELETE("/customer-milk-rates/:id", customerMilkRateController.DeleteRate)

		// Customer Pay Date Range Routes
		api.POST("/customer-pay-date-ranges", customerPayDateRangeController.CreateRange)
		api.GET("/customer-pay-date-ranges", customerPayDateRangeController.GetRanges)
		api.GET("/customer-pay-date-ranges/:id", customerPayDateRangeController.GetRange)
		api.PUT("/customer-pay-date-ranges/:id", customerPayDateRangeController.UpdateRange)
		api.DELETE("/customer-pay-date-ranges/:id", customerPayDateRangeController.DeleteRange)

		// Member Dependant Routes
		api.POST("/member-dependants", memberDependantController.CreateDependant)
		api.GET("/member-dependants", memberDependantController.GetDependants)
		api.GET("/member-dependants/:id", memberDependantController.GetDependant)
		api.PUT("/member-dependants/:id", memberDependantController.UpdateDependant)
		api.DELETE("/member-dependants/:id", memberDependantController.DeleteDependant)

		// Loan Routes
		api.POST("/loans", loanController.CreateLoan)
		api.GET("/loans", loanController.GetLoans)
		api.GET("/loans/:id", loanController.GetLoan)
		api.PUT("/loans/:id", loanController.UpdateLoan)
		api.DELETE("/loans/:id", loanController.DeleteLoan)

		// Transport Rate Routes
		api.POST("/transport-rates", transportRateController.CreateRate)
		api.GET("/transport-rates", transportRateController.GetRates)
		api.GET("/transport-rates/:id", transportRateController.GetRate)
		api.PUT("/transport-rates/:id", transportRateController.UpdateRate)
		api.DELETE("/transport-rates/:id", transportRateController.DeleteRate)

		// Training Routes
		api.POST("/trainings", trainingController.CreateTraining)
		api.GET("/trainings", trainingController.GetTrainings)
		api.GET("/trainings/:id", trainingController.GetTraining)
		api.PUT("/trainings/:id", trainingController.UpdateTraining)
		api.DELETE("/trainings/:id", trainingController.DeleteTraining)

		// Training Session Routes
		api.POST("/training-sessions", trainingSessionController.CreateSession)
		api.GET("/training-sessions", trainingSessionController.GetSessions)
		api.GET("/training-sessions/:id", trainingSessionController.GetSession)
		api.PUT("/training-sessions/:id", trainingSessionController.UpdateSession)
		api.DELETE("/training-sessions/:id", trainingSessionController.DeleteSession)

		// Training Session Attendee Routes
		api.POST("/training-session-attendees", trainingAttendeeController.CreateAttendee)
		api.GET("/training-session-attendees", trainingAttendeeController.GetAttendees)
		api.GET("/training-session-attendees/:id", trainingAttendeeController.GetAttendee)
		api.PUT("/training-session-attendees/:id", trainingAttendeeController.UpdateAttendee)
		api.DELETE("/training-session-attendees/:id", trainingAttendeeController.DeleteAttendee)

		// Exchange Visit Routes
		api.POST("/exchange-visits", exchangeVisitController.CreateVisit)
		api.GET("/exchange-visits", exchangeVisitController.GetVisits)
		api.GET("/exchange-visits/:id", exchangeVisitController.GetVisit)
		api.PUT("/exchange-visits/:id", exchangeVisitController.UpdateVisit)
		api.DELETE("/exchange-visits/:id", exchangeVisitController.DeleteVisit)

		// Exchange Visit Attendee Routes
		api.POST("/exchange-visit-attendees", exchangeVisitAttendeeController.CreateAttendee)
		api.GET("/exchange-visit-attendees", exchangeVisitAttendeeController.GetAttendees)
		api.GET("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.GetAttendee)
		api.PUT("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.UpdateAttendee)
		api.DELETE("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.DeleteAttendee)

		// Payment Mode Routes
		api.POST("/payment-modes", paymentModeController.CreatePaymentMode)
		api.GET("/payment-modes", paymentModeController.GetPaymentModes)
		api.GET("/payment-modes/:id", paymentModeController.GetPaymentMode)
		api.PUT("/payment-modes/:id", paymentModeController.UpdatePaymentMode)
		api.DELETE("/payment-modes/:id", paymentModeController.DeletePaymentMode)

		// Route Center Routes
		api.POST("/route-centers", routeCenterController.CreateCenter)
		api.GET("/route-centers", routeCenterController.GetCenters)
		api.GET("/route-centers/:id", routeCenterController.GetCenter)
		api.PUT("/route-centers/:id", routeCenterController.UpdateCenter)
		api.DELETE("/route-centers/:id", routeCenterController.DeleteCenter)
	}

	return r
}
