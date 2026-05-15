package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerEmployeeRoutes(api *gin.RouterGroup) {
	employeeController := controllers.NewEmployeeController()
	employeeBenefitController := controllers.NewEmployeeBenefitController()
	employeeDocumentController := controllers.NewEmployeeDocumentController()
	employeeLeaveTypeController := controllers.NewEmployeeLeaveTypeController()
	employeePayrollController := controllers.NewEmployeePayrollController()

	api.POST("/employees", employeeController.CreateEmployee)
	api.GET("/employees", employeeController.GetEmployees)
	api.GET("/employees/:id", employeeController.GetEmployee)
	api.PUT("/employees/:id", employeeController.UpdateEmployee)
	api.DELETE("/employees/:id", employeeController.DeleteEmployee)

	api.POST("/employees/:id/salaries", employeeController.CreateSalary)
	api.POST("/employees/:id/bank-accounts", employeeController.CreateBankAccount)

	// Employee Benefits
	api.POST("/employee-benefits", employeeBenefitController.CreateEmployeeBenefit)
	api.GET("/employee-benefits", employeeBenefitController.GetEmployeeBenefits)
	api.GET("/employee-benefits/:id", employeeBenefitController.GetEmployeeBenefit)
	api.PUT("/employee-benefits/:id", employeeBenefitController.UpdateEmployeeBenefit)
	api.DELETE("/employee-benefits/:id", employeeBenefitController.DeleteEmployeeBenefit)

	// Employee Documents
	api.POST("/employee-documents", employeeDocumentController.CreateEmployeeDocument)
	api.GET("/employee-documents", employeeDocumentController.GetEmployeeDocuments)
	api.GET("/employee-documents/:id", employeeDocumentController.GetEmployeeDocument)
	api.PUT("/employee-documents/:id", employeeDocumentController.UpdateEmployeeDocument)
	api.DELETE("/employee-documents/:id", employeeDocumentController.DeleteEmployeeDocument)

	// Employee Leave Types
	api.POST("/employee-leave-types", employeeLeaveTypeController.CreateEmployeeLeaveType)
	api.GET("/employee-leave-types", employeeLeaveTypeController.GetEmployeeLeaveTypes)
	api.GET("/employee-leave-types/:id", employeeLeaveTypeController.GetEmployeeLeaveType)
	api.PUT("/employee-leave-types/:id", employeeLeaveTypeController.UpdateEmployeeLeaveType)
	api.DELETE("/employee-leave-types/:id", employeeLeaveTypeController.DeleteEmployeeLeaveType)

	// Employee Payrolls
	api.POST("/employee-payrolls", employeePayrollController.CreateEmployeePayroll)
	api.GET("/employee-payrolls", employeePayrollController.GetEmployeePayrolls)
	api.GET("/employee-payrolls/:id", employeePayrollController.GetEmployeePayroll)
	api.PUT("/employee-payrolls/:id", employeePayrollController.UpdateEmployeePayroll)
	api.DELETE("/employee-payrolls/:id", employeePayrollController.DeleteEmployeePayroll)
}
