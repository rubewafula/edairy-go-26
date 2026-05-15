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

	jwtSecret := os.Getenv("JWT_SECRET")
	authMiddleware := middleware.AuthMiddleware([]byte(jwtSecret))

	registerPublicRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	api.Use(authMiddleware)
	api.Use(middleware.RequireAutoPermission())
	{
		registerAuthenticatedRoutes(api)
	}
	return r
}

func registerPublicRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()
	r.POST("/api/signup", authController.Signup)
	r.GET("/api/verify", authController.Verify)
	r.POST("/api/login", authController.Login)
	r.POST("/api/forgot-password", authController.ForgotPassword)
	r.POST("/api/reset-password", authController.ResetPassword)
}

func registerAuthenticatedRoutes(api *gin.RouterGroup) {
	authController := controllers.NewAuthController()
	memberController := controllers.NewMemberController()
	employeeController := controllers.NewEmployeeController()
	bankController := controllers.NewBankController()
	bankBranchController := controllers.NewBankBranchController()
	memberBankAccountController := controllers.NewMemberBankAccountController()
	memberDependantController := controllers.NewMemberDependantController()
	memberTypeController := controllers.NewMemberTypeController()
	milkDeliveryShiftController := controllers.NewMilkDeliveryShiftController()
	coolerMilkCollectionController := controllers.NewCoolerMilkCollectionController()
	milkRejectController := controllers.NewMilkRejectController()
	milkJournalController := controllers.NewMilkJournalController()
	milkJournalEntryController := controllers.NewMilkJournalEntryController()
	milkDeliveryController := controllers.NewMilkDeliveryController()
	milkLocalSaleController := controllers.NewMilkLocalSaleController()
	dailyMilkVarianceController := controllers.NewDailyMilkVarianceController()
	transporterController := controllers.NewTransporterController()
	individualTransporterController := controllers.NewIndividualTransporterController()
	companyTransporterController := controllers.NewCompanyTransporterController()
	driverController := controllers.NewTransporterDriverController()
	vehicleController := controllers.NewTransporterVehicleController()
	transportRateController := controllers.NewTransportRateController()
	routeController := controllers.NewRouteController()
	subRouteController := controllers.NewSubRouteController()
	routeCenterController := controllers.NewRouteCenterController()
	productGradeController := controllers.NewProductGradeController()
	storeController := controllers.NewStoreController()
	storeItemController := controllers.NewStoreItemController()
	storeInventoryController := controllers.NewStoreInventoryController()
	storeStockController := controllers.NewStoreStockController()
	storeSaleController := controllers.NewStoreSaleController()
	interStoreTransferController := controllers.NewInterStoreTransferController()
	customerController := controllers.NewCustomerController()
	customerClassController := controllers.NewCustomerClassController()
	customerTypeController := controllers.NewCustomerTypeController()
	supplierController := controllers.NewSupplierController()
	supplyController := controllers.NewSupplyController()
	purchaseOrderController := controllers.NewPurchaseOrderController()
	smsController := controllers.NewSMSController()
	smsCampaignController := controllers.NewSMSCampaignController()
	loanController := controllers.NewLoanController()
	loanManagementController := controllers.NewLoanManagementController()
	organizationAddressController := controllers.NewOrganizationAddressController()
	organizationMemberController := controllers.NewOrganizationMemberController()
	organizationBankController := controllers.NewOrganizationBankController()
	organizationWalletController := controllers.NewOrganizationWalletController()
	walletTypeController := controllers.NewWalletTypeController()
	shareTypeController := controllers.NewShareTypeController()
	shareAccountController := controllers.NewShareAccountController()
	sharePaymentController := controllers.NewSharePaymentController()
	dividendDeclarationController := controllers.NewDividendDeclarationController()
	cattleBreedController := controllers.NewCattleBreedController()
	trainingController := controllers.NewTrainingController()
	trainingSessionController := controllers.NewTrainingSessionController()
	exchangeVisitController := controllers.NewExchangeVisitController()
	paymentModeController := controllers.NewPaymentModeController()
	assetCategoryController := controllers.NewAssetCategoryController()
	assetController := controllers.NewAssetController()
	roleController := controllers.NewRoleController()
	permissionController := controllers.NewPermissionController()
	userController := controllers.NewUserController()
	adminDashboardController := controllers.NewAdminDashboardController()
	trainingAttendeeController := controllers.NewTrainingAttendeeController()
	shareTransferController := controllers.NewShareTransferController()
	supplierContactController := controllers.NewSupplierContactController()
	supplierDocumentController := controllers.NewSupplierDocumentController()
	supplierBankAccountController := controllers.NewSupplierBankAccountController()
	supplierQuoteController := controllers.NewSupplierQuoteController()
	suppliedItemController := controllers.NewSuppliedItemController()
	supplyRejectController := controllers.NewSupplyRejectController()
	loanOrganizationProfileController := controllers.NewLoanOrganizationProfileController()
	organizationLeadershipController := controllers.NewOrganizationLeadershipController()
	organizationDocumentController := controllers.NewOrganizationDocumentController()
	organizationKybCommentController := controllers.NewOrganizationKybCommentController()
	shareDividendController := controllers.NewShareDividendController()
	shareTransactionController := controllers.NewShareTransactionController()
	exchangeVisitAttendeeController := controllers.NewExchangeVisitAttendeeController()
	interStoreTransferItemController := controllers.NewInterStoreTransferItemController()
	canMovementController := controllers.NewCanMovementController()
	routeAssignmentController := controllers.NewTransporterRouteAssignmentController()
	driverAssignmentController := controllers.NewTransporterDriverAssignmentController()
	transporterBankAccountController := controllers.NewTransporterBankAccountController()
	transporterBenefitController := controllers.NewTransporterBenefitController()
	defaultMilkRateController := controllers.NewDefaultMilkRateController()
	milkCanController := controllers.NewMilkCanController()
	deductionTypeController := controllers.NewDeductionTypeController()
	deductionPricingRuleController := controllers.NewDeductionPricingRuleController()
	itemCategoryController := controllers.NewItemCategoryController()
	storeItemUnitController := controllers.NewStoreItemUnitController()
	storeStockTakingController := controllers.NewStoreStockTakingController()
	storeStockMovementController := controllers.NewStoreStockMovementController()
	storeStockMovementTypeController := controllers.NewStoreStockMovementTypeController()
	storeSaleItemController := controllers.NewStoreSaleItemController()
	customerBillingController := controllers.NewCustomerBillingController()
	customerInvoiceController := controllers.NewCustomerInvoiceController()
	customerPaymentController := controllers.NewCustomerPaymentController()
	customerMilkRateController := controllers.NewCustomerMilkRateController()
	customerPayDateRangeController := controllers.NewCustomerPayDateRangeController()
	supplierCategoryController := controllers.NewSupplierCategoryController()
	assetAssignmentController := controllers.NewAssetAssignmentController()
	assetDepreciationController := controllers.NewAssetDepreciationController()

	// =====================================================
	// AUTH & DASHBOARD MODULE
	// =====================================================

	// Auth Routes
	api.POST("/change-password", authController.ChangePassword)

	// Dashboard Routes
	api.GET("/admin-dashboard", adminDashboardController.GetDashboard)

	// Training Session Attendee Routes
	api.POST("/training-session-attendees", trainingAttendeeController.CreateAttendee)
	api.GET("/training-session-attendees", trainingAttendeeController.GetAttendees)
	api.GET("/training-session-attendees/:id", trainingAttendeeController.GetAttendee)
	api.PUT("/training-session-attendees/:id", trainingAttendeeController.UpdateAttendee)
	api.DELETE("/training-session-attendees/:id", trainingAttendeeController.DeleteAttendee)

	// =====================================================
	// RBAC & USER MANAGEMENT MODULE
	// =====================================================

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

	// =====================================================
	// MEMBER MANAGEMENT MODULE
	// =====================================================

	// Member Routes
	api.POST("/member/create", memberController.CreateMember)
	api.GET("/members", memberController.GetMembers)
	api.GET("/member/:id", memberController.GetMember)
	api.PUT("/member/:id", memberController.UpdateMember)
	api.DELETE("/member/:id", memberController.DeleteMember)

	// Member Type Routes
	api.POST("/member-types", memberTypeController.CreateMemberType)
	api.GET("/member-types", memberTypeController.GetMemberTypes)
	api.GET("/member-types/:id", memberTypeController.GetMemberType)
	api.PUT("/member-types/:id", memberTypeController.UpdateMemberType)
	api.DELETE("/member-types/:id", memberTypeController.DeleteMemberType)

	// Member Bank Account Routes
	api.POST("/member-bank-accounts", memberBankAccountController.CreateAccount)
	api.GET("/member-bank-accounts", memberBankAccountController.GetAccounts)
	api.GET("/member-bank-accounts/:id", memberBankAccountController.GetAccount)
	api.PUT("/member-bank-accounts/:id", memberBankAccountController.UpdateAccount)
	api.DELETE("/member-bank-accounts/:id", memberBankAccountController.DeleteAccount)

	// Member Dependant Routes
	api.POST("/member-dependants", memberDependantController.CreateDependant)
	api.GET("/member-dependants", memberDependantController.GetDependants)
	api.GET("/member-dependants/:id", memberDependantController.GetDependant)
	api.PUT("/member-dependants/:id", memberDependantController.UpdateDependant)
	api.DELETE("/member-dependants/:id", memberDependantController.DeleteDependant)

	// =====================================================
	// EMPLOYEE MANAGEMENT MODULE
	// =====================================================

	api.POST("/employees", employeeController.CreateEmployee)
	api.GET("/employees", employeeController.GetEmployees)
	api.GET("/employees/:id", employeeController.GetEmployee)
	api.PUT("/employees/:id", employeeController.UpdateEmployee)
	api.DELETE("/employees/:id", employeeController.DeleteEmployee)
	api.POST("/employees/:id/salaries", employeeController.CreateSalary)
	api.POST("/employees/:id/bank-accounts", employeeController.CreateBankAccount)

	// =====================================================
	// BANKING MODULE
	// =====================================================

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

	// =====================================================
	// MILK COLLECTION & DELIVERY MODULE
	// =====================================================

	// Milk Delivery Shift Routes
	api.POST("/milk-delivery-shifts", milkDeliveryShiftController.CreateShift)
	api.GET("/milk-delivery-shifts", milkDeliveryShiftController.GetShifts)
	api.GET("/milk-delivery-shifts/:id", milkDeliveryShiftController.GetShift)
	api.PUT("/milk-delivery-shifts/:id", milkDeliveryShiftController.UpdateShift)
	api.DELETE("/milk-delivery-shifts/:id", milkDeliveryShiftController.DeleteShift)

	// Cooler Milk Collection Routes
	api.POST("/cooler-milk-collections", coolerMilkCollectionController.CreateCollection)
	api.GET("/cooler-milk-collections", coolerMilkCollectionController.GetCollections)
	api.GET("/cooler-milk-collections/:id", coolerMilkCollectionController.GetCollection)
	api.PUT("/cooler-milk-collections/:id", coolerMilkCollectionController.UpdateCollection)
	api.DELETE("/cooler-milk-collections/:id", coolerMilkCollectionController.DeleteCollection)

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

	// Stray Milk Collections
	api.GET("/stray-milk-collections", milkJournalEntryController.GetStrayEntries)

	// Milk Reject Routes
	api.POST("/milk-rejects", milkRejectController.CreateReject)
	api.GET("/milk-rejects", milkRejectController.GetRejects)
	api.GET("/milk-rejects/:id", milkRejectController.GetReject)
	api.DELETE("/milk-rejects/:id", milkRejectController.DeleteReject)

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

	// =====================================================
	// TRANSPORT & LOGISTICS MODULE
	// =====================================================

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

	// Route Routes
	api.POST("/routes", routeController.CreateRoute)
	api.GET("/routes", routeController.GetRoutes)
	api.GET("/routes/:id", routeController.GetRoute)
	api.PUT("/routes/:id", routeController.UpdateRoute)
	api.DELETE("/routes/:id", routeController.DeleteRoute)

	// Sub Route Routes
	api.POST("/sub-routes", subRouteController.CreateSubRoute)
	api.GET("/sub-routes", subRouteController.GetSubRoutes)
	api.GET("/sub-routes/:id", subRouteController.GetSubRoute)
	api.PUT("/sub-routes/:id", subRouteController.UpdateSubRoute)
	api.DELETE("/sub-routes/:id", subRouteController.DeleteSubRoute)

	// Route Center Routes
	api.POST("/route-centers", routeCenterController.CreateCenter)
	api.GET("/route-centers", routeCenterController.GetCenters)
	api.GET("/route-centers/:id", routeCenterController.GetCenter)
	api.PUT("/route-centers/:id", routeCenterController.UpdateCenter)
	api.DELETE("/route-centers/:id", routeCenterController.DeleteCenter)

	// Transport Rate Routes
	api.POST("/transport-rates", transportRateController.CreateRate)
	api.GET("/transport-rates", transportRateController.GetRates)
	api.GET("/transport-rates/:id", transportRateController.GetRate)
	api.PUT("/transport-rates/:id", transportRateController.UpdateRate)
	api.DELETE("/transport-rates/:id", transportRateController.DeleteRate)

	// =====================================================
	// INVENTORY & STORE MODULE
	// =====================================================

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

	// Store Sale Routes
	api.POST("/store-sales", storeSaleController.CreateSale)
	api.GET("/store-sales", storeSaleController.GetSales)
	api.GET("/store-sales/:id", storeSaleController.GetSale)
	api.PUT("/store-sales/:id", storeSaleController.UpdateSale)
	api.DELETE("/store-sales/:id", storeSaleController.DeleteSale)

	// Inter Store Transfer Routes
	api.POST("/inter-store-transfers", interStoreTransferController.CreateTransfer)
	api.GET("/inter-store-transfers", interStoreTransferController.GetTransfers)
	api.GET("/inter-store-transfers/:id", interStoreTransferController.GetTransfer)
	api.PUT("/inter-store-transfers/:id", interStoreTransferController.UpdateTransfer)
	api.DELETE("/inter-store-transfers/:id", interStoreTransferController.DeleteTransfer)

	// =====================================================
	// CUSTOMER MANAGEMENT MODULE
	// =====================================================

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

	// Customer Type Routes
	api.POST("/customer-types", customerTypeController.CreateCustomerType)
	api.GET("/customer-types", customerTypeController.GetCustomerTypes)
	api.GET("/customer-types/:id", customerTypeController.GetCustomerType)
	api.PUT("/customer-types/:id", customerTypeController.UpdateCustomerType)
	api.DELETE("/customer-types/:id", customerTypeController.DeleteCustomerType)

	// =====================================================
	// SUPPLIER & PROCUREMENT MODULE
	// =====================================================

	// Supplier Routes
	api.POST("/suppliers", supplierController.CreateSupplier)
	api.GET("/suppliers", supplierController.GetSuppliers)
	api.GET("/suppliers/:id", supplierController.GetSupplier)

	// Purchase Order Routes
	api.POST("/purchase-orders", purchaseOrderController.CreatePO)
	api.GET("/purchase-orders", purchaseOrderController.GetPOs)
	api.GET("/purchase-orders/:id", purchaseOrderController.GetPO)
	api.PUT("/purchase-orders/:id", purchaseOrderController.UpdatePO)
	api.DELETE("/purchase-orders/:id", purchaseOrderController.DeletePO)

	// Supply Routes
	api.POST("/supplies", supplyController.CreateSupply)
	api.GET("/supplies", supplyController.GetSupplies)

	// =====================================================
	// LOAN MANAGEMENT MODULE
	// =====================================================

	// Loan Routes
	api.POST("/loans", loanController.CreateLoan)
	api.GET("/loans", loanController.GetLoans)
	api.GET("/loans/:id", loanController.GetLoan)
	api.PUT("/loans/:id", loanController.UpdateLoan)
	api.DELETE("/loans/:id", loanController.DeleteLoan)

	// Loan Management Routes
	api.POST("/loan-accounts", loanManagementController.CreateAccount)
	api.GET("/loan-accounts", loanManagementController.GetAccounts)
	api.POST("/loan-transactions", loanManagementController.CreateTransaction)

	// =====================================================
	// SHARES & DIVIDENDS MODULE
	// =====================================================

	// Share Type Routes
	api.POST("/share-types", shareTypeController.CreateShareType)
	api.GET("/share-types", shareTypeController.GetShareTypes)

	// Share Transfer Routes
	api.POST("/share-transfers", shareTransferController.CreateShareTransfer)
	api.GET("/share-transfers", shareTransferController.GetShareTransfers)
	api.GET("/share-transfers/:id", shareTransferController.GetShareTransfer)
	api.PUT("/share-transfers/:id", shareTransferController.UpdateShareTransfer)
	api.DELETE("/share-transfers/:id", shareTransferController.DeleteShareTransfer)

	// Share Dividend Routes
	api.POST("/share-dividends", shareDividendController.CreateDividend)
	api.GET("/share-dividends", shareDividendController.GetDividends)
	api.GET("/share-dividends/:id", shareDividendController.GetDividend)
	api.PUT("/share-dividends/:id", shareDividendController.UpdateDividend)
	api.DELETE("/share-dividends/:id", shareDividendController.DeleteDividend)

	// Share Account Routes
	api.POST("/share-accounts", shareAccountController.CreateAccount)
	api.GET("/share-accounts", shareAccountController.GetAccounts)

	// Share Payment Routes
	api.POST("/share-payments", sharePaymentController.CreateSharePayment)
	api.GET("/share-payments", sharePaymentController.GetSharePayments)

	// Dividend Declaration Routes
	api.POST("/dividend-declarations", dividendDeclarationController.CreateDeclaration)
	api.GET("/dividend-declarations", dividendDeclarationController.GetDeclarations)

	// =====================================================
	// ORGANIZATION MANAGEMENT MODULE
	// =====================================================

	// Organization Member Routes
	api.POST("/organization-members", organizationMemberController.CreateMember)
	api.GET("/organization-members", organizationMemberController.GetMembers)

	// Organization Address Routes
	api.POST("/organization-addresses", organizationAddressController.CreateAddress)
	api.GET("/organization-addresses", organizationAddressController.GetAddresses)

	// Organization Bank Routes
	api.POST("/organization-banks", organizationBankController.CreateBank)
	api.GET("/organization-banks", organizationBankController.GetBanks)

	// Organization Wallet Routes
	api.POST("/organization-wallets", organizationWalletController.CreateWallet)
	api.GET("/organization-wallets", organizationWalletController.GetWallets)

	// =====================================================
	// ASSET MANAGEMENT MODULE
	// =====================================================

	// Asset Category Routes
	api.POST("/asset-categories", assetCategoryController.CreateCategory)
	api.GET("/asset-categories", assetCategoryController.GetCategories)
	api.GET("/asset-categories/:id", assetCategoryController.GetCategory)
	api.PUT("/asset-categories/:id", assetCategoryController.UpdateCategory)
	api.DELETE("/asset-categories/:id", assetCategoryController.DeleteCategory)

	// Fixed Asset Routes
	api.POST("/fixed-assets", assetController.CreateAsset)
	api.GET("/fixed-assets", assetController.GetAssets)
	api.GET("/fixed-assets/:id", assetController.GetAsset)
	api.PUT("/fixed-assets/:id", assetController.UpdateAsset)
	api.DELETE("/fixed-assets/:id", assetController.DeleteAsset)

	// =====================================================
	// TRAINING & EDUCATION MODULE
	// =====================================================

	// Training Routes
	api.POST("/trainings", trainingController.CreateTraining)
	api.GET("/trainings", trainingController.GetTrainings)

	// Training Session Routes
	api.POST("/training-sessions", trainingSessionController.CreateSession)
	api.GET("/training-sessions", trainingSessionController.GetSessions)

	// Exchange Visit Routes
	api.POST("/exchange-visits", exchangeVisitController.CreateVisit)
	api.GET("/exchange-visits", exchangeVisitController.GetVisits)

	// =====================================================
	// SMS & COMMUNICATION MODULE
	// =====================================================

	// SMS Routes
	api.POST("/sms-groups", smsController.CreateGroup)
	api.GET("/sms-groups", smsController.GetGroups)
	api.POST("/sms-contacts", smsController.CreateContact)
	api.POST("/sms-send", smsController.SendMessage)
	api.GET("/sms-queue", smsController.GetQueue)

	// SMS Campaign Routes
	api.POST("/sms-campaigns", smsCampaignController.CreateCampaign)
	api.GET("/sms-campaigns", smsCampaignController.GetCampaigns)

	// =====================================================
	// MASTER DATA MODULE
	// =====================================================

	// Wallet Type Routes
	api.POST("/wallet-types", walletTypeController.CreateWalletType)
	api.GET("/wallet-types", walletTypeController.GetWalletTypes)

	// Cattle Breed Routes
	api.POST("/cattle-breeds", cattleBreedController.CreateCattleBreed)
	api.GET("/cattle-breeds", cattleBreedController.GetCattleBreeds)

	// Product Grade Routes
	api.POST("/product-grades", productGradeController.CreateGrade)
	api.GET("/product-grades", productGradeController.GetGrades)

	// Payment Mode Routes
	api.POST("/payment-modes", paymentModeController.CreatePaymentMode)
	api.GET("/payment-modes", paymentModeController.GetPaymentModes)

	// =====================================================
	// ADDITIONAL MODULE ROUTES FROM DEFINITION
	// =====================================================

	// Supplier Sub-entity Routes
	api.POST("/suppliers/:id/contacts", supplierController.CreateContact)
	api.GET("/suppliers/:id/contacts", supplierController.GetSupplierContacts)
	api.POST("/suppliers/:id/bank-accounts", supplierController.CreateBankAccount)
	api.GET("/suppliers/:id/bank-accounts", supplierController.GetSupplierBankAccounts)

	// Supplier Contact Routes
	api.POST("/supplier-contacts", supplierContactController.CreateContact)
	api.GET("/supplier-contacts", supplierContactController.GetContacts)
	api.GET("/supplier-contacts/:id", supplierContactController.GetContact)
	api.PUT("/supplier-contacts/:id", supplierContactController.UpdateContact)
	api.DELETE("/supplier-contacts/:id", supplierContactController.DeleteContact)

	// Supplier Document Routes
	api.POST("/supplier-documents", supplierDocumentController.CreateDocument)
	api.GET("/supplier-documents", supplierDocumentController.GetDocuments)
	api.GET("/supplier-documents/:id", supplierDocumentController.GetDocument)
	api.PUT("/supplier-documents/:id", supplierDocumentController.UpdateDocument)
	api.DELETE("/supplier-documents/:id", supplierDocumentController.DeleteDocument)
	api.PATCH("/supplier-documents/:id/verify", supplierDocumentController.VerifyDocument)

	// Supplier Bank Account Routes
	api.POST("/supplier-bank-accounts", supplierBankAccountController.CreateBankAccount)
	api.GET("/supplier-bank-accounts", supplierBankAccountController.GetBankAccounts)
	api.GET("/supplier-bank-accounts/:id", supplierBankAccountController.GetBankAccount)
	api.PUT("/supplier-bank-accounts/:id", supplierBankAccountController.UpdateBankAccount)
	api.DELETE("/supplier-bank-accounts/:id", supplierBankAccountController.DeleteBankAccount)

	// Supplier Quote Routes
	api.POST("/supplier-quotes", supplierQuoteController.CreateQuote)
	api.GET("/supplier-quotes", supplierQuoteController.GetQuotes)
	api.POST("/supplier-quotes/:id/items", supplierQuoteController.CreateQuoteItem)
	api.GET("/supplier-quotes/:id/items", supplierQuoteController.GetQuoteItems)
	api.GET("/supplier-quote-items/:id", supplierQuoteController.GetQuoteItem)
	api.PUT("/supplier-quote-items/:id", supplierQuoteController.UpdateQuoteItem)
	api.DELETE("/supplier-quote-items/:id", supplierQuoteController.DeleteQuoteItem)

	// Supply Reject Routes
	api.POST("/supply-rejects", supplyRejectController.CreateReject)
	api.GET("/supply-rejects", supplyRejectController.GetRejects)
	api.GET("/supplies/:id/rejects", supplyRejectController.GetRejectsBySupply)
	api.GET("/supply-rejects/:id", supplyRejectController.GetReject)
	api.PUT("/supply-rejects/:id", supplyRejectController.UpdateReject)
	api.DELETE("/supply-rejects/:id", supplyRejectController.DeleteReject)

	// Supplied Item Routes
	api.GET("/supplied-items/:id", suppliedItemController.GetSuppliedItem)
	api.PUT("/supplied-items/:id", suppliedItemController.UpdateSuppliedItem)
	api.DELETE("/supplied-items/:id", suppliedItemController.DeleteSuppliedItem)

	// Purchase Requisition Routes
	api.POST("/purchase-requisitions", purchaseOrderController.CreateRequisition)
	api.GET("/purchase-requisitions", purchaseOrderController.GetRequisitions)
	api.GET("/purchase-requisitions/:id", purchaseOrderController.GetRequisition)
	api.PUT("/purchase-requisitions/:id", purchaseOrderController.UpdateRequisition)
	api.DELETE("/purchase-requisitions/:id", purchaseOrderController.DeleteRequisition)
	api.POST("/purchase-requisition-items", purchaseOrderController.CreateRequisitionItem)
	api.GET("/purchase-requisitions/:id/items", purchaseOrderController.GetRequisitionItems)
	api.GET("/purchase-requisition-items/:id", purchaseOrderController.GetRequisitionItem)
	api.PUT("/purchase-requisition-items/:id", purchaseOrderController.UpdateRequisitionItem)
	api.DELETE("/purchase-requisition-items/:id", purchaseOrderController.DeleteRequisitionItem)

	// Expanded Organization Management
	api.GET("/organization-members/:id", organizationMemberController.GetMember)
	api.PUT("/organization-members/:id", organizationMemberController.UpdateMember)
	api.DELETE("/organization-members/:id", organizationMemberController.DeleteMember)

	api.GET("/organization-banks/:id", organizationBankController.GetBank)
	api.PUT("/organization-banks/:id", organizationBankController.UpdateBank)
	api.DELETE("/organization-banks/:id", organizationBankController.DeleteBank)

	api.GET("/organization-addresses/:id", organizationAddressController.GetAddress)
	api.PUT("/organization-addresses/:id", organizationAddressController.UpdateAddress)
	api.DELETE("/organization-addresses/:id", organizationAddressController.DeleteAddress)

	api.GET("/organization-wallets/:id", organizationWalletController.GetWallet)
	api.PUT("/organization-wallets/:id", organizationWalletController.UpdateWallet)
	api.DELETE("/organization-wallets/:id", organizationWalletController.DeleteWallet)

	api.POST("/organization-leaderships", organizationLeadershipController.CreateLeadership)
	api.GET("/organization-leaderships", organizationLeadershipController.GetLeaderships)
	api.GET("/organization-leaderships/:id", organizationLeadershipController.GetLeadership)
	api.PUT("/organization-leaderships/:id", organizationLeadershipController.UpdateLeadership)
	api.DELETE("/organization-leaderships/:id", organizationLeadershipController.DeleteLeadership)
	api.GET("/organization-leaderships/national-id/:id_no", organizationLeadershipController.GetLeadershipByNationalID)

	api.POST("/organization-documents", organizationDocumentController.CreateDocument)
	api.GET("/organization-documents", organizationDocumentController.GetDocuments)
	api.GET("/organization-documents/:id", organizationDocumentController.GetDocument)
	api.PUT("/organization-documents/:id", organizationDocumentController.UpdateDocument)
	api.DELETE("/organization-documents/:id", organizationDocumentController.DeleteDocument)
	api.GET("/organization-documents/astra/:astra_id", organizationDocumentController.GetDocumentsByAstraID)

	api.POST("/organization-kyb-comments", organizationKybCommentController.CreateComment)
	api.GET("/organization-kyb-comments", organizationKybCommentController.GetComments)
	api.GET("/organization-kyb-comments/:id", organizationKybCommentController.GetComment)
	api.PUT("/organization-kyb-comments/:id", organizationKybCommentController.UpdateComment)
	api.DELETE("/organization-kyb-comments/:id", organizationKybCommentController.DeleteComment)
	api.GET("/organization-kyb-comments/iteration/:iteration", organizationKybCommentController.GetCommentsByIteration)

	// Expanded Loan Management
	api.POST("/loan-callbacks", loanManagementController.CreateCallback)
	api.GET("/loan-callbacks", loanManagementController.GetCallbacks)
	api.GET("/loan-callbacks/:id", loanManagementController.GetCallback)
	api.PUT("/loan-callbacks/:id", loanManagementController.UpdateCallback)
	api.DELETE("/loan-callbacks/:id", loanManagementController.DeleteCallback)

	api.POST("/member-loans", loanManagementController.CreateMemberLoan)
	api.GET("/member-loans", loanManagementController.GetMemberLoans)
	api.PUT("/member-loans/:id", loanManagementController.UpdateMemberLoan)

	api.POST("/loan-origination-logs", loanManagementController.CreateOriginationLog)

	// Loan Organization Profile Routes
	api.POST("/loan-organization-profiles", loanOrganizationProfileController.CreateProfile)
	api.GET("/loan-organization-profiles", loanOrganizationProfileController.GetProfiles)
	api.GET("/loan-organization-profiles/:id", loanOrganizationProfileController.GetProfile)
	api.PUT("/loan-organization-profiles/:id", loanOrganizationProfileController.UpdateProfile)
	api.DELETE("/loan-organization-profiles/:id", loanOrganizationProfileController.DeleteProfile)

	// Share Transaction Routes
	api.POST("/share-transactions", shareTransactionController.CreateShareTransaction)
	api.GET("/share-transactions", shareTransactionController.GetShareTransactions)
	api.GET("/share-transactions/:id", shareTransactionController.GetShareTransaction)
	api.PUT("/share-transactions/:id", shareTransactionController.UpdateShareTransaction)
	api.DELETE("/share-transactions/:id", shareTransactionController.DeleteShareTransaction)

	// Exchange Visit Attendee Routes
	api.POST("/exchange-visit-attendees", exchangeVisitAttendeeController.CreateAttendee)
	api.GET("/exchange-visit-attendees", exchangeVisitAttendeeController.GetAttendees)
	api.GET("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.GetAttendee)
	api.PUT("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.UpdateAttendee)
	api.DELETE("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.DeleteAttendee)

	// Logistics & Transporter Management
	api.POST("/can-movements", canMovementController.CreateMovement)
	api.GET("/can-movements", canMovementController.GetMovements)
	api.GET("/can-movements/:id", canMovementController.GetMovement)
	api.PUT("/can-movements/:id", canMovementController.UpdateMovement)
	api.DELETE("/can-movements/:id", canMovementController.DeleteMovement)

	api.POST("/transporter-route-assignments", routeAssignmentController.CreateAssignment)
	api.GET("/transporter-route-assignments", routeAssignmentController.GetAssignments)
	api.GET("/transporter-route-assignments/:id", routeAssignmentController.GetAssignment)
	api.PUT("/transporter-route-assignments/:id", routeAssignmentController.UpdateAssignment)
	api.DELETE("/transporter-route-assignments/:id", routeAssignmentController.DeleteAssignment)

	api.POST("/transporter-driver-assignments", driverAssignmentController.CreateAssignment)
	api.GET("/transporter-driver-assignments", driverAssignmentController.GetAssignments)
	api.GET("/transporter-driver-assignments/:id", driverAssignmentController.GetAssignment)
	api.PUT("/transporter-driver-assignments/:id", driverAssignmentController.UpdateAssignment)
	api.DELETE("/transporter-driver-assignments/:id", driverAssignmentController.DeleteAssignment)

	api.POST("/transporter-bank-accounts", transporterBankAccountController.CreateAccount)
	api.GET("/transporter-bank-accounts", transporterBankAccountController.GetAccounts)
	api.GET("/transporter-bank-accounts/:id", transporterBankAccountController.GetAccount)
	api.PUT("/transporter-bank-accounts/:id", transporterBankAccountController.UpdateAccount)
	api.DELETE("/transporter-bank-accounts/:id", transporterBankAccountController.DeleteAccount)

	api.POST("/transporter-benefits", transporterBenefitController.CreateBenefit)
	api.GET("/transporter-benefits", transporterBenefitController.GetBenefits)
	api.GET("/transporter-benefits/:id", transporterBenefitController.GetBenefit)
	api.PUT("/transporter-benefits/:id", transporterBenefitController.UpdateBenefit)
	api.DELETE("/transporter-benefits/:id", transporterBenefitController.DeleteBenefit)

	// Milk & Rate Management
	api.POST("/default-milk-rates", defaultMilkRateController.CreateRate)
	api.GET("/default-milk-rates", defaultMilkRateController.GetRates)
	api.GET("/default-milk-rates/:id", defaultMilkRateController.GetRate)
	api.PUT("/default-milk-rates/:id", defaultMilkRateController.UpdateRate)
	api.DELETE("/default-milk-rates/:id", defaultMilkRateController.DeleteRate)

	api.POST("/milk-cans", milkCanController.CreateMilkCan)
	api.GET("/milk-cans", milkCanController.GetMilkCans)
	api.GET("/milk-cans/:id", milkCanController.GetMilkCan)
	api.PUT("/milk-cans/:id", milkCanController.UpdateMilkCan)
	api.DELETE("/milk-cans/:id", milkCanController.DeleteMilkCan)

	// Inventory & Store Details
	api.POST("/item-categories", itemCategoryController.CreateCategory)
	api.GET("/item-categories", itemCategoryController.GetCategories)
	api.GET("/item-categories/:id", itemCategoryController.GetCategory)
	api.PUT("/item-categories/:id", itemCategoryController.UpdateCategory)
	api.DELETE("/item-categories/:id", itemCategoryController.DeleteCategory)

	api.POST("/store-item-units", storeItemUnitController.CreateUnit)
	api.GET("/store-item-units", storeItemUnitController.GetUnits)
	api.GET("/store-item-units/:id", storeItemUnitController.GetUnit)
	api.PUT("/store-item-units/:id", storeItemUnitController.UpdateUnit)
	api.DELETE("/store-item-units/:id", storeItemUnitController.DeleteUnit)

	api.POST("/store-stock-takings", storeStockTakingController.CreateStockTaking)
	api.GET("/store-stock-takings", storeStockTakingController.GetStockTakings)
	api.GET("/store-stock-takings/:id", storeStockTakingController.GetStockTaking)
	api.PUT("/store-stock-takings/:id", storeStockTakingController.UpdateStockTaking)
	api.DELETE("/store-stock-takings/:id", storeStockTakingController.DeleteStockTaking)

	api.POST("/store-stock-movements", storeStockMovementController.CreateMovement)
	api.GET("/store-stock-movements", storeStockMovementController.GetMovements)
	api.GET("/store-stock-movements/:id", storeStockMovementController.GetMovement)
	api.PUT("/store-stock-movements/:id", storeStockMovementController.UpdateMovement)
	api.DELETE("/store-stock-movements/:id", storeStockMovementController.DeleteMovement)

	api.POST("/store-stock-movement-types", storeStockMovementTypeController.CreateMovementType)
	api.GET("/store-stock-movement-types", storeStockMovementTypeController.GetMovementTypes)
	api.GET("/store-stock-movement-types/:id", storeStockMovementTypeController.GetMovementType)
	api.PUT("/store-stock-movement-types/:id", storeStockMovementTypeController.UpdateMovementType)
	api.DELETE("/store-stock-movement-types/:id", storeStockMovementTypeController.DeleteMovementType)

	api.POST("/inter-store-transfer-items", interStoreTransferItemController.CreateTransferItem)
	api.GET("/inter-store-transfer-items", interStoreTransferItemController.GetTransferItems)
	api.GET("/inter-store-transfer-items/:id", interStoreTransferItemController.GetTransferItem)
	api.PUT("/inter-store-transfer-items/:id", interStoreTransferItemController.UpdateTransferItem)
	api.DELETE("/inter-store-transfer-items/:id", interStoreTransferItemController.DeleteTransferItem)

	api.POST("/store-sale-items", storeSaleItemController.CreateSaleItem)
	api.GET("/store-sale-items", storeSaleItemController.GetSaleItems)
	api.GET("/store-sale-items/:id", storeSaleItemController.GetSaleItem)
	api.PUT("/store-sale-items/:id", storeSaleItemController.UpdateSaleItem)
	api.DELETE("/store-sale-items/:id", storeSaleItemController.DeleteSaleItem)

	// Customer Billing & Financials
	api.GET("/customer-billings", customerBillingController.GetBillings)
	api.GET("/customer-billings/:id", customerBillingController.GetBilling)
	api.GET("/customer-billings/:id/items", customerBillingController.GetBillingItems)

	api.POST("/customer-invoices", customerInvoiceController.CreateInvoice)
	api.GET("/customer-invoices", customerInvoiceController.GetInvoices)
	api.GET("/customer-invoices/:id", customerInvoiceController.GetInvoice)
	api.DELETE("/customer-invoices/:id", customerInvoiceController.DeleteInvoice)

	api.POST("/customer-payments", customerPaymentController.CreatePayment)
	api.GET("/customer-payments", customerPaymentController.GetPayments)
	api.GET("/customer-payments/:id", customerPaymentController.GetPayment)

	api.POST("/customer-milk-rates", customerMilkRateController.CreateRate)
	api.GET("/customer-milk-rates", customerMilkRateController.GetRates)
	api.GET("/customer-milk-rates/:id", customerMilkRateController.GetRate)
	api.PUT("/customer-milk-rates/:id", customerMilkRateController.UpdateRate)
	api.DELETE("/customer-milk-rates/:id", customerMilkRateController.DeleteRate)

	api.POST("/customer-pay-date-ranges", customerPayDateRangeController.CreateCustomerPayDateRange)
	api.GET("/customer-pay-date-ranges", customerPayDateRangeController.GetCustomerPayDateRanges)
	api.GET("/customer-pay-date-ranges/:id", customerPayDateRangeController.GetCustomerPayDateRange)
	api.PUT("/customer-pay-date-ranges/:id", customerPayDateRangeController.UpdateCustomerPayDateRange)
	api.DELETE("/customer-pay-date-ranges/:id", customerPayDateRangeController.DeleteCustomerPayDateRange)

	// Deductions & Suppliers
	api.POST("/deduction-types", deductionTypeController.CreateDeductionType)
	api.GET("/deduction-types", deductionTypeController.GetDeductionTypes)
	api.GET("/deduction-types/:id", deductionTypeController.GetDeductionType)
	api.PUT("/deduction-types/:id", deductionTypeController.UpdateDeductionType)
	api.DELETE("/deduction-types/:id", deductionTypeController.DeleteDeductionType)

	api.POST("/deduction-pricing-rules", deductionPricingRuleController.CreateRule)
	api.GET("/deduction-pricing-rules", deductionPricingRuleController.GetRules)
	api.GET("/deduction-pricing-rules/:id", deductionPricingRuleController.GetRule)
	api.PUT("/deduction-pricing-rules/:id", deductionPricingRuleController.UpdateRule)
	api.DELETE("/deduction-pricing-rules/:id", deductionPricingRuleController.DeleteRule)

	api.POST("/supplier-categories", supplierCategoryController.CreateCategory)
	api.GET("/supplier-categories", supplierCategoryController.GetCategories)
	api.GET("/supplier-categories/:id", supplierCategoryController.GetCategory)
	api.PUT("/supplier-categories/:id", supplierCategoryController.UpdateCategory)
	api.DELETE("/supplier-categories/:id", supplierCategoryController.DeleteCategory)

	// Fixed Asset Management
	api.POST("/asset-assignments", assetAssignmentController.CreateAssignment)
	api.GET("/asset-assignments", assetAssignmentController.GetAssignments)
	api.GET("/asset-assignments/:id", assetAssignmentController.GetAssignment)
	api.PUT("/asset-assignments/:id", assetAssignmentController.UpdateAssignment)
	api.DELETE("/asset-assignments/:id", assetAssignmentController.DeleteAssignment)

	api.POST("/asset-depreciation-entries", assetDepreciationController.CreateEntry)
	api.GET("/asset-depreciation-entries", assetDepreciationController.GetEntries)
	api.GET("/asset-depreciation-entries/:id", assetDepreciationController.GetEntry)
	api.DELETE("/asset-depreciation-entries/:id", assetDepreciationController.DeleteEntry)

}
