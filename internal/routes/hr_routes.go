package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerHRRoutes(api *gin.RouterGroup) {
	departmentController := controllers.NewDepartmentController()
	jobPositionController := controllers.NewJobPositionController()

	// Department Routes
	api.POST("/departments", departmentController.CreateDepartment)
	api.GET("/departments", departmentController.GetDepartments)
	api.GET("/departments/:id", departmentController.GetDepartment)
	api.PUT("/departments/:id", departmentController.UpdateDepartment)
	api.DELETE("/departments/:id", departmentController.DeleteDepartment)

	// Job Position Routes
	api.POST("/job-positions", jobPositionController.CreateJobPosition)
	api.GET("/job-positions", jobPositionController.GetJobPositions)
	api.GET("/job-positions/:id", jobPositionController.GetJobPosition)
	api.PUT("/job-positions/:id", jobPositionController.UpdateJobPosition)
	api.DELETE("/job-positions/:id", jobPositionController.DeleteJobPosition)
}
