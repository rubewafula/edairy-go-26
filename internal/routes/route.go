package routes

import (
	"os"

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

	// WebSocket Route moved above global middleware to prevent CORS interference during handshake
	r.GET("/ws", ws.ServeWS(hub))

	// CORS Middleware to allow cross-origin requests
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
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
