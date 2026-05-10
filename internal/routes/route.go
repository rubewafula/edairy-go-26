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
	coolerMilkCollectionController := controllers.NewCoolerMilkCollectionController()
	milkRejectController := controllers.NewMilkRejectController()
	transporterController := controllers.NewTransporterController()
	individualTransporterController := controllers.NewIndividualTransporterController()
	companyTransporterController := controllers.NewCompanyTransporterController()
	driverController := controllers.NewTransporterDriverController()
	vehicleController := controllers.NewTransporterVehicleController()
	routeAssignmentController := controllers.NewTransporterRouteAssignmentController()
	driverAssignmentController := controllers.NewTransporterDriverAssignmentController()
	transporterBankAccountController := controllers.NewTransporterBankAccountController()
	transporterBenefitController := controllers.NewTransporterBenefitController()
	subRouteController := controllers.NewSubRouteController()
	milkJournalController := controllers.NewMilkJournalController()
	milkJournalEntryController := controllers.NewMilkJournalEntryController()
	milkDeliveryController := controllers.NewMilkDeliveryController()
	milkLocalSaleController := controllers.NewMilkLocalSaleController()
	milkCanController := controllers.NewMilkCanController()
	canMovementController := controllers.NewCanMovementController()
	productGradeController := controllers.NewProductGradeController()
	deductionTypeController := controllers.NewDeductionTypeController()
	itemCategoryController := controllers.NewItemCategoryController()
	deductionPricingRuleController := controllers.NewDeductionPricingRuleController()
	storeController := controllers.NewStoreController()
	storeItemController := controllers.NewStoreItemController()
	storeSaleController := controllers.NewStoreSaleController()
	storeSaleItemController := controllers.NewStoreSaleItemController()
	storeInventoryController := controllers.NewStoreInventoryController()
	storeStockController := controllers.NewStoreStockController()
	storeStockTakingController := controllers.NewStoreStockTakingController()
	storeStockMovementController := controllers.NewStoreStockMovementController()
	storeStockMovementTypeController := controllers.NewStoreStockMovementTypeController()
	memberTypeController := controllers.NewMemberTypeController()
	routeController := controllers.NewRouteController()
	customerController := controllers.NewCustomerController()
	customerClassController := controllers.NewCustomerClassController()
	customerMilkRateController := controllers.NewCustomerMilkRateController()
	customerPayDateRangeController := controllers.NewCustomerPayDateRangeController()
	memberDependantController := controllers.NewMemberDependantController()
	dailyMilkVarianceController := controllers.NewDailyMilkVarianceController()
	loanController := controllers.NewLoanController()
	transportRateController := controllers.NewTransportRateController()
	trainingController := controllers.NewTrainingController()
	trainingSessionController := controllers.NewTrainingSessionController()
	trainingAttendeeController := controllers.NewTrainingAttendeeController()
	exchangeVisitController := controllers.NewExchangeVisitController()
	exchangeVisitAttendeeController := controllers.NewExchangeVisitAttendeeController()
	paymentModeController := controllers.NewPaymentModeController()
	routeCenterController := controllers.NewRouteCenterController()
	shareTypeController := controllers.NewShareTypeController()
	shareAccountController := controllers.NewShareAccountController()
	shareDividendController := controllers.NewShareDividendController()
	dividendDeclarationController := controllers.NewDividendDeclarationController()
	sharePaymentController := controllers.NewSharePaymentController()
	shareTransferController := controllers.NewShareTransferController()
	shareTransactionController := controllers.NewShareTransactionController()
	assetCategoryController := controllers.NewAssetCategoryController()
	assetController := controllers.NewAssetController()
	assetAssignmentController := controllers.NewAssetAssignmentController()
	assetDepreciationController := controllers.NewAssetDepreciationController()
	permissionController := controllers.NewPermissionController()
	roleController := controllers.NewRoleController()
	userController := controllers.NewUserController()

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

		// Individual Transporter Routes
		api.GET("/individual-transporters", individualTransporterController.GetIndividualTransporters)
		api.GET("/individual-transporters/:id", individualTransporterController.GetIndividualTransporter)
		api.PUT("/individual-transporters/:id", individualTransporterController.UpdateIndividualTransporter)

		// Company Transporter Routes
		api.GET("/company-transporters", companyTransporterController.GetCompanyTransporters)
		api.GET("/company-transporters/:id", companyTransporterController.GetCompanyTransporter)
		api.PUT("/company-transporters/:id", companyTransporterController.UpdateCompanyTransporter)

		// Transporter Driver Routes
		api.POST("/transporter-drivers", driverController.CreateDriver)
		api.GET("/transporter-drivers", driverController.GetDrivers)
		api.GET("/transporter-drivers/:id", driverController.GetDriver)
		api.PUT("/transporter-drivers/:id", driverController.UpdateDriver)
		api.DELETE("/transporter-drivers/:id", driverController.DeleteDriver)

		// Transporter Vehicle Routes
		api.POST("/transporter-vehicles", vehicleController.CreateVehicle)
		api.GET("/transporter-vehicles", vehicleController.GetVehicles)
		api.GET("/transporter-vehicles/:id", vehicleController.GetVehicle)
		api.PUT("/transporter-vehicles/:id", vehicleController.UpdateVehicle)
		api.DELETE("/transporter-vehicles/:id", vehicleController.DeleteVehicle)

		// Transporter Route Assignment Routes
		api.POST("/transporter-route-assignments", routeAssignmentController.CreateAssignment)
		api.GET("/transporter-route-assignments", routeAssignmentController.GetAssignments)
		api.GET("/transporter-route-assignments/:id", routeAssignmentController.GetAssignment)
		api.PUT("/transporter-route-assignments/:id", routeAssignmentController.UpdateAssignment)
		api.DELETE("/transporter-route-assignments/:id", routeAssignmentController.DeleteAssignment)

		// Transporter Driver Assignment Routes
		api.POST("/transporter-driver-assignments", driverAssignmentController.CreateAssignment)
		api.GET("/transporter-driver-assignments", driverAssignmentController.GetAssignments)
		api.GET("/transporter-driver-assignments/:id", driverAssignmentController.GetAssignment)
		api.PUT("/transporter-driver-assignments/:id", driverAssignmentController.UpdateAssignment)
		api.DELETE("/transporter-driver-assignments/:id", driverAssignmentController.DeleteAssignment)

		// Transporter Bank Account Routes
		api.POST("/transporter-bank-accounts", transporterBankAccountController.CreateAccount)
		api.GET("/transporter-bank-accounts", transporterBankAccountController.GetAccounts)
		api.GET("/transporter-bank-accounts/:id", transporterBankAccountController.GetAccount)
		api.PUT("/transporter-bank-accounts/:id", transporterBankAccountController.UpdateAccount)
		api.DELETE("/transporter-bank-accounts/:id", transporterBankAccountController.DeleteAccount)

		// Transporter Benefit Routes
		api.POST("/transporter-benefits", transporterBenefitController.CreateBenefit)
		api.GET("/transporter-benefits", transporterBenefitController.GetBenefits)
		api.GET("/transporter-benefits/:id", transporterBenefitController.GetBenefit)
		api.PUT("/transporter-benefits/:id", transporterBenefitController.UpdateBenefit)
		api.DELETE("/transporter-benefits/:id", transporterBenefitController.DeleteBenefit)

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

		// Milk Journal Entry Routes
		api.POST("/milk-journal-entries", milkJournalEntryController.CreateEntry)
		api.GET("/milk-journal-entries", milkJournalEntryController.GetEntries)
		api.GET("/milk-journal-entries/:id", milkJournalEntryController.GetEntry)
		api.PUT("/milk-journal-entries/:id", milkJournalEntryController.UpdateEntry)
		api.DELETE("/milk-journal-entries/:id", milkJournalEntryController.DeleteEntry)
		api.POST("/milk-journal-entries/upload", milkJournalEntryController.UploadEntries)

		// Stray Milk Collection Routes
		api.GET("/stray-milk-collections", milkJournalEntryController.GetStrayEntries)

		// Milk Reject Routes
		api.POST("/milk-rejects", milkRejectController.CreateReject)
		api.GET("/milk-rejects", milkRejectController.GetRejects)
		api.GET("/milk-rejects/:id", milkRejectController.GetReject)
		api.DELETE("/milk-rejects/:id", milkRejectController.DeleteReject)

		// Cooler Milk Collection Routes
		api.POST("/cooler-milk-collections", coolerMilkCollectionController.CreateCollection)
		api.GET("/cooler-milk-collections", coolerMilkCollectionController.GetCollections)
		api.GET("/cooler-milk-collections/:id", coolerMilkCollectionController.GetCollection)
		api.PUT("/cooler-milk-collections/:id", coolerMilkCollectionController.UpdateCollection)
		api.DELETE("/cooler-milk-collections/:id", coolerMilkCollectionController.DeleteCollection)

		// Milk Delivery Routes
		api.POST("/milk-deliveries", milkDeliveryController.CreateDelivery)
		api.GET("/milk-deliveries", milkDeliveryController.GetDeliveries)
		api.GET("/milk-deliveries/:id", milkDeliveryController.GetDelivery)
		api.PUT("/milk-deliveries/:id", milkDeliveryController.UpdateDelivery)
		api.DELETE("/milk-deliveries/:id", milkDeliveryController.DeleteDelivery)

		// Milk Local Sale Routes
		api.POST("/milk-local-sales", milkLocalSaleController.CreateMilkLocalSale)
		api.GET("/milk-local-sales", milkLocalSaleController.GetMilkLocalSales)
		api.GET("/milk-local-sales/:id", milkLocalSaleController.GetMilkLocalSale)
		api.PUT("/milk-local-sales/:id", milkLocalSaleController.UpdateMilkLocalSale)
		api.DELETE("/milk-local-sales/:id", milkLocalSaleController.DeleteMilkLocalSale)

		// Daily Milk Variance Routes
		api.GET("/daily-milk-variances", dailyMilkVarianceController.GetDailyVariances)

		// Milk Can Routes
		api.POST("/milk-cans", milkCanController.CreateMilkCan)
		api.GET("/milk-cans", milkCanController.GetMilkCans)
		api.GET("/milk-cans/:id", milkCanController.GetMilkCan)
		api.PUT("/milk-cans/:id", milkCanController.UpdateMilkCan)
		api.DELETE("/milk-cans/:id", milkCanController.DeleteMilkCan)

		// Can Movement Routes
		api.POST("/can-movements", canMovementController.CreateMovement)
		api.GET("/can-movements", canMovementController.GetMovements)
		api.GET("/can-movements/:id", canMovementController.GetMovement)
		api.PUT("/can-movements/:id", canMovementController.UpdateMovement)
		api.DELETE("/can-movements/:id", canMovementController.DeleteMovement)

		// Product Grade Routes
		api.POST("/product-grades", productGradeController.CreateGrade)
		api.GET("/product-grades", productGradeController.GetGrades)
		api.GET("/product-grades/:id", productGradeController.GetGrade)
		api.PUT("/product-grades/:id", productGradeController.UpdateGrade)
		api.DELETE("/product-grades/:id", productGradeController.DeleteGrade)

		// Deduction Type Routes
		api.POST("/deduction-types", deductionTypeController.CreateDeductionType)
		api.GET("/deduction-types", deductionTypeController.GetDeductionTypes)
		api.GET("/deduction-types/:id", deductionTypeController.GetDeductionType)
		api.PUT("/deduction-types/:id", deductionTypeController.UpdateDeductionType)
		api.DELETE("/deduction-types/:id", deductionTypeController.DeleteDeductionType)

		// Item Category Routes
		api.POST("/item-categories", itemCategoryController.CreateCategory)
		api.GET("/item-categories", itemCategoryController.GetCategories)
		api.GET("/item-categories/:id", itemCategoryController.GetCategory)
		api.PUT("/item-categories/:id", itemCategoryController.UpdateCategory)
		api.DELETE("/item-categories/:id", itemCategoryController.DeleteCategory)

		// Deduction Pricing Rule Routes
		api.POST("/deduction-pricing-rules", deductionPricingRuleController.CreateRule)
		api.GET("/deduction-pricing-rules", deductionPricingRuleController.GetRules)
		api.GET("/deduction-pricing-rules/:id", deductionPricingRuleController.GetRule)
		api.PUT("/deduction-pricing-rules/:id", deductionPricingRuleController.UpdateRule)
		api.DELETE("/deduction-pricing-rules/:id", deductionPricingRuleController.DeleteRule)

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

		// Store Stock Taking Routes
		api.POST("/store-stock-takings", storeStockTakingController.CreateStockTaking)
		api.GET("/store-stock-takings", storeStockTakingController.GetStockTakings)
		api.GET("/store-stock-takings/:id", storeStockTakingController.GetStockTaking)
		api.PUT("/store-stock-takings/:id", storeStockTakingController.UpdateStockTaking)
		api.DELETE("/store-stock-takings/:id", storeStockTakingController.DeleteStockTaking)

		// Store Stock Movement Routes
		api.POST("/store-stock-movements", storeStockMovementController.CreateMovement)
		api.GET("/store-stock-movements", storeStockMovementController.GetMovements)
		api.GET("/store-stock-movements/:id", storeStockMovementController.GetMovement)
		api.PUT("/store-stock-movements/:id", storeStockMovementController.UpdateMovement)
		api.DELETE("/store-stock-movements/:id", storeStockMovementController.DeleteMovement)

		// Store Stock Movement Type Routes
		api.POST("/store-stock-movement-types", storeStockMovementTypeController.CreateMovementType)
		api.GET("/store-stock-movement-types", storeStockMovementTypeController.GetMovementTypes)
		api.GET("/store-stock-movement-types/:id", storeStockMovementTypeController.GetMovementType)
		api.PUT("/store-stock-movement-types/:id", storeStockMovementTypeController.UpdateMovementType)
		api.DELETE("/store-stock-movement-types/:id", storeStockMovementTypeController.DeleteMovementType)

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
		api.DELETE("/dividend-declarations/:id", dividendDeclarationController.DeleteDeclaration)

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

		api.POST("/asset-categories", assetCategoryController.CreateCategory)
		api.GET("/asset-categories", assetCategoryController.GetCategories)
		api.GET("/asset-categories/:id", assetCategoryController.GetCategory)
		api.PUT("/asset-categories/:id", assetCategoryController.UpdateCategory)
		api.DELETE("/asset-categories/:id", assetCategoryController.DeleteCategory)

		// Asset Routes
		api.POST("/fixed-assets", assetController.CreateAsset)
		api.GET("/fixed-assets", assetController.GetAssets)
		api.GET("/fixed-assets/:id", assetController.GetAsset)
		api.PUT("/fixed-assets/:id", assetController.UpdateAsset)
		api.DELETE("/fixed-assets/:id", assetController.DeleteAsset)

		// Asset Assignment Routes
		api.POST("/asset-assignments", assetAssignmentController.CreateAssignment)
		api.GET("/asset-assignments", assetAssignmentController.GetAssignments)
		api.GET("/asset-assignments/:id", assetAssignmentController.GetAssignment)
		api.PUT("/asset-assignments/:id", assetAssignmentController.UpdateAssignment)
		api.DELETE("/asset-assignments/:id", assetAssignmentController.DeleteAssignment)

		// Asset Depreciation Routes
		api.POST("/asset-depreciation-entries", assetDepreciationController.CreateEntry)
		api.GET("/asset-depreciation-entries", assetDepreciationController.GetEntries)
		api.GET("/asset-depreciation-entries/:id", assetDepreciationController.GetEntry)
		api.DELETE("/asset-depreciation-entries/:id", assetDepreciationController.DeleteEntry)

		// Permission Routes
		api.POST("/permissions", permissionController.CreatePermission)
		api.GET("/permissions", permissionController.GetPermissions)
		api.GET("/permissions/:id", permissionController.GetPermission)
		api.PUT("/permissions/:id", permissionController.UpdatePermission)
		api.DELETE("/permissions/:id", permissionController.DeletePermission)

		// Role Routes
		api.POST("/roles", roleController.CreateRole)
		api.GET("/roles", roleController.GetRoles)
		api.GET("/roles/:id", roleController.GetRole)
		api.PUT("/roles/:id", roleController.UpdateRole)
		api.DELETE("/roles/:id", roleController.DeleteRole)
		api.POST("/roles-permissions/:id", roleController.AppendAllPermissions)

		// User Routes
		api.POST("/users", userController.CreateUser)
		api.GET("/users", userController.GetUsers)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)
	}

	return r
}
