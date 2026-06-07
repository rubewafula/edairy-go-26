package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerUserStoreAssignmentRoutes(api *gin.RouterGroup) {
	userStoreAssignmentController := controllers.NewUserStoreAssignmentController()

	api.POST("/user-store-assignments", userStoreAssignmentController.CreateAssignment)
	api.GET("/user-store-assignments", userStoreAssignmentController.GetAssignments)
	api.GET("/user-store-assignments/:id", userStoreAssignmentController.GetAssignment)
	api.PUT("/user-store-assignments/:id", userStoreAssignmentController.UpdateAssignment)
	api.DELETE("/user-store-assignments/:id", userStoreAssignmentController.DeleteAssignment)
}
