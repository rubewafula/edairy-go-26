package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

// registerMemberRoutes registers the API routes for member-related functionalities.
func registerMemberRoutes(api *gin.RouterGroup) {
	memberController := controllers.NewMemberController()
	memberBankAccountController := controllers.NewMemberBankAccountController()
	memberDependantController := controllers.NewMemberDependantController()
	memberTypeController := controllers.NewMemberTypeController()

	// Member Registration Routes
	api.POST("/members", memberController.CreateMember)
	api.GET("/members", memberController.GetMembers)
	api.GET("/members/:id", memberController.GetMember)
	api.PUT("/members/:id", memberController.UpdateMember)
	api.DELETE("/members/:id", memberController.DeleteMember)

	// Member Type Routes
	api.POST("/member-types", memberTypeController.CreateMemberType)
	api.GET("/member-types", memberTypeController.GetMemberTypes)
	api.GET("/member-types/:id", memberTypeController.GetMemberType)
	api.PUT("/member-types/:id", memberTypeController.UpdateMemberType)
	api.DELETE("/member-types/:id", memberTypeController.DeleteMemberType)

	api.POST("/member-bank-accounts", memberBankAccountController.CreateAccount)
	api.GET("/member-bank-accounts", memberBankAccountController.GetAccounts)
	api.GET("/member-bank-accounts/:id", memberBankAccountController.GetAccount)
	api.PUT("/member-bank-accounts/:id", memberBankAccountController.UpdateAccount)
	api.DELETE("/member-bank-accounts/:id", memberBankAccountController.DeleteAccount)

	api.POST("/member-dependants", memberDependantController.CreateDependant)
	api.GET("/member-dependants", memberDependantController.GetDependants)
	api.GET("/member-dependants/:id", memberDependantController.GetDependant)
	api.PUT("/member-dependants/:id", memberDependantController.UpdateDependant)
	api.DELETE("/member-dependants/:id", memberDependantController.DeleteDependant)

}
