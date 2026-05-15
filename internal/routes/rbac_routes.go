package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerRBACRoutes(api *gin.RouterGroup) {
	permissionController := controllers.NewPermissionController()
	roleController := controllers.NewRoleController()
	userController := controllers.NewUserController()

	api.POST("/permissions", permissionController.CreatePermission)
	api.GET("/permissions", permissionController.GetPermissions)
	api.GET("/permissions/:id", permissionController.GetPermission)
	api.PUT("/permissions/:id", permissionController.UpdatePermission)
	api.DELETE("/permissions/:id", permissionController.DeletePermission)

	api.POST("/roles", roleController.CreateRole)
	api.GET("/roles", roleController.GetRoles)
	api.GET("/roles/:id", roleController.GetRole)
	api.PUT("/roles/:id", roleController.UpdateRole)
	api.DELETE("/roles/:id", roleController.DeleteRole)
	api.POST("/roles-permissions/:id", roleController.AppendAllPermissions)

	api.POST("/users", userController.CreateUser)
	api.GET("/users", userController.GetUsers)
	api.GET("/users/:id", userController.GetUser)
	api.PUT("/users/:id", userController.UpdateUser)
	api.DELETE("/users/:id", userController.DeleteUser)
}
