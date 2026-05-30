package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerEmployeePayrollRoutes(api *gin.RouterGroup) {
	employeePayrollController := controllers.NewEmployeePayrollController()
	statutoryDeductionConfigController := controllers.NewStatutoryDeductionConfigurationController()

	// Employee Payroll Routes
	api.POST("/employee-payrolls", employeePayrollController.CreatePayroll)
	api.GET("/employee-payrolls", employeePayrollController.GetEmployeePayrolls)
	api.GET("/employee-payrolls/:id", employeePayrollController.GetEmployeePayroll)
	api.PUT("/employee-payrolls/:id/confirm", employeePayrollController.ConfirmPayroll)
	api.PUT("/employee-payrolls/:id/approve", employeePayrollController.ApprovePayroll)
	api.DELETE("/employee-payrolls/:id", employeePayrollController.DeleteEmployeePayroll)

	// Statutory Deduction Configuration Routes
	api.POST("/statutory-deduction-configurations", statutoryDeductionConfigController.CreateConfiguration)
	api.GET("/statutory-deduction-configurations", statutoryDeductionConfigController.GetConfigurations)
	api.GET("/statutory-deduction-configurations/:id", statutoryDeductionConfigController.GetConfiguration)
	api.PUT("/statutory-deduction-configurations/:id", statutoryDeductionConfigController.UpdateConfiguration)
	api.DELETE("/statutory-deduction-configurations/:id", statutoryDeductionConfigController.DeleteConfiguration)
}
