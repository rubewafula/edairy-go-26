package routes

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
	"github.com/rubewafula/edairy-go-26/internal/middleware"
	ws "github.com/rubewafula/edairy-go-26/internal/ws-socket"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rubewafula/edairy-go-26/docs"
)

func SetupRouter() *gin.Engine {
	// Initialize Socket.IO Server
	hub := ws.NewHub()
	go hub.Run()

	ws.InitHub(hub)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// DO NOT use "*" with AllowCredentials: true.
		// This function dynamically mirrors the incoming origin safely.
		AllowOriginFunc: func(origin string) bool {
			allowed := map[string]bool{
				"https://arithi.edairy.africa":     true,
				"https://api.arithi.edairy.africa": true,
				"https://edairy.africa":            true,
				"http://localhost:5173":            true,
			}
			return allowed[origin] // In production, check against an allowed list
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// WebSocket Route moved above global middleware to prevent CORS interference during handshake
	r.GET("/ws", ws.ServeWS(hub))

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
	registerRBACRoutes(api)
	registerMemberRoutes(api)
	registerEmployeeRoutes(api)
	registerBankingRoutes(api)
	registerMilkRoutes(api)
	registerInventoryRoutes(api)
	registerCustomerRoutes(api)
	registerEmployeePayrollRoutes(api)
	registerUINotificationRoutes(api)
	registerSupplierRoutes(api)
	registerTransporterRoutes(api)
	registerLoanRoutes(api)
	registerShareRoutes(api)
	registerOrgRoutes(api)
	registerAssetRoutes(api)
	registerTrainingRoutes(api)
	registerLivestockRoutes(api)
	registerSMSRoutes(api)
	registerRouteRoutes(api)
	registerMasterDataRoutes(api)
	registerLocationRoutes(api)
	registerHRRoutes(api)
	registerFinanceRoutes(api)

	// Auth & Dashboard
	authController := controllers.NewAuthController()
	adminDashboardController := controllers.NewAdminDashboardController()
	api.POST("/change-password", authController.ChangePassword)
	api.GET("/admin-dashboard", adminDashboardController.GetDashboard)

}
