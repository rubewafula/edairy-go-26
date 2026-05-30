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
	employeeLeaveApplicationController := controllers.NewEmployeeLeaveApplicationController()
	employeeLeaveAssignmentController := controllers.NewEmployeeLeaveAssignmentController()
	documentTypeController := controllers.NewDocumentTypeController()
	employeeReliefController := controllers.NewEmployeeReliefController()
	deductionTypeController := controllers.NewEmployeeDeductionTypeController()
	employeeQualificationController := controllers.NewEmployeeQualificationController()
	employeeBankAccountController := controllers.NewEmployeeBankAccountController()
	employeeSalaryController := controllers.NewEmployeeSalaryController()
	employeeContractController := controllers.NewEmployeeContractController()
	employeeExitDetailController := controllers.NewEmployeeExitDetailController()
	employeeDependantController := controllers.NewEmployeeDependantController()
	employeeLeaveDetailController := controllers.NewEmployeeLeaveDetailController()
	employeePayrollBenefitController := controllers.NewEmployeePayrollBenefitController()
	employeePayrollDeductionController := controllers.NewEmployeePayrollDeductionController()
	employeePayrollReliefController := controllers.NewEmployeePayrollReliefController()
	employeePayslipController := controllers.NewEmployeePayslipController()
	employeeProfessionalTitleController := controllers.NewEmployeeProfessionalTitleController()
	employeeTerminationCategoryController := controllers.NewEmployeeTerminationCategoryController()

	api.POST("/employees", employeeController.CreateEmployee)
	api.GET("/employees", employeeController.GetEmployees)
	api.GET("/employees/:id", employeeController.GetEmployee)
	api.PUT("/employees/:id", employeeController.UpdateEmployee)
	api.DELETE("/employees/:id", employeeController.DeleteEmployee)

	// Employee Salary Routes (now using dedicated controller)
	api.POST("/employee-salaries", employeeSalaryController.Create)
	api.GET("/employee-salaries", employeeSalaryController.List)
	api.GET("/employee-salaries/:id", employeeSalaryController.Get)
	api.PUT("/employee-salaries/:id", employeeSalaryController.Update)
	api.DELETE("/employee-salaries/:id", employeeSalaryController.Delete)

	// Employee Bank Account Routes
	api.POST("/employee-bank-accounts", employeeBankAccountController.Create)
	api.GET("/employee-bank-accounts", employeeBankAccountController.List)
	api.GET("/employee-bank-accounts/:id", employeeBankAccountController.Get)
	api.PUT("/employee-bank-accounts/:id", employeeBankAccountController.Update)
	api.DELETE("/employee-bank-accounts/:id", employeeBankAccountController.Delete)

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

	// Document Types
	api.POST("/document-types", documentTypeController.Create)
	api.GET("/document-types", documentTypeController.List)
	api.GET("/document-types/:id", documentTypeController.Get)
	api.PUT("/document-types/:id", documentTypeController.Update)
	api.DELETE("/document-types/:id", documentTypeController.Delete)

	// Employee Leave Types
	api.POST("/employee-leave-types", employeeLeaveTypeController.CreateEmployeeLeaveType)
	api.GET("/employee-leave-types", employeeLeaveTypeController.GetEmployeeLeaveTypes)
	api.GET("/employee-leave-types/:id", employeeLeaveTypeController.GetEmployeeLeaveType)
	api.PUT("/employee-leave-types/:id", employeeLeaveTypeController.UpdateEmployeeLeaveType)
	api.DELETE("/employee-leave-types/:id", employeeLeaveTypeController.DeleteEmployeeLeaveType)

	// Employee Leave Applications
	api.POST("/employee-leave-applications", employeeLeaveApplicationController.Create)
	api.GET("/employee-leave-applications", employeeLeaveApplicationController.List)
	api.GET("/employee-leave-applications/:id", employeeLeaveApplicationController.Get)
	api.PUT("/employee-leave-applications/:id", employeeLeaveApplicationController.Update)
	api.DELETE("/employee-leave-applications/:id", employeeLeaveApplicationController.Delete)

	// Employee Leave Assignments
	api.POST("/employee-leave-assignments", employeeLeaveAssignmentController.Create)
	api.GET("/employee-leave-assignments", employeeLeaveAssignmentController.List)
	api.GET("/employee-leave-assignments/:id", employeeLeaveAssignmentController.Get)
	api.PUT("/employee-leave-assignments/:id", employeeLeaveAssignmentController.Update)
	api.DELETE("/employee-leave-assignments/:id", employeeLeaveAssignmentController.Delete)

	// Employee Reliefs
	api.POST("/employee-reliefs", employeeReliefController.Create)
	api.GET("/employee-reliefs", employeeReliefController.List)
	api.GET("/employee-reliefs/:id", employeeReliefController.Get)
	api.PUT("/employee-reliefs/:id", employeeReliefController.Update)
	api.DELETE("/employee-reliefs/:id", employeeReliefController.Delete)

	// Employee Deduction Types
	api.POST("/employee-deduction-types", deductionTypeController.CreateDeductionType)
	api.GET("/employee-deduction-types", deductionTypeController.GetDeductionTypes)
	api.GET("/employee-deduction-types/:id", deductionTypeController.GetDeductionType)
	api.PUT("/employee-deduction-types/:id", deductionTypeController.UpdateDeductionType)
	api.DELETE("/employee-deduction-types/:id", deductionTypeController.DeleteDeductionType)

	// Employee Qualifiactions
	api.POST("/employee-qualifications", employeeQualificationController.Create)
	api.GET("/employee-qualifications", employeeQualificationController.List)
	api.GET("/employee-qualifications/:id", employeeQualificationController.Get)
	api.PUT("/employee-qualifications/:id", employeeQualificationController.Update)
	api.DELETE("/employee-qualifications/:id", employeeQualificationController.Delete)

	// Employee Contract Routes
	api.POST("/employee-contracts", employeeContractController.Create)
	api.GET("/employee-contracts", employeeContractController.List)
	api.GET("/employee-contracts/:id", employeeContractController.Get)
	api.PUT("/employee-contracts/:id", employeeContractController.Update)
	api.DELETE("/employee-contracts/:id", employeeContractController.Delete)

	// Employee Dependant Routes
	api.POST("/employee-dependants", employeeDependantController.CreateEmployeeDependant)
	api.GET("/employee-dependants", employeeDependantController.GetEmployeeDependants)
	api.GET("/employee-dependants/:id", employeeDependantController.GetEmployeeDependant)
	api.PUT("/employee-dependants/:id", employeeDependantController.UpdateEmployeeDependant)
	api.DELETE("/employee-dependants/:id", employeeDependantController.DeleteEmployeeDependant)

	// Employee Exit Detail Routes
	api.POST("/employee-exit-details", employeeExitDetailController.Create)
	api.GET("/employee-exit-details", employeeExitDetailController.List)
	api.GET("/employee-exit-details/:id", employeeExitDetailController.Get)
	api.PUT("/employee-exit-details/:id", employeeExitDetailController.Update)
	api.DELETE("/employee-exit-details/:id", employeeExitDetailController.Delete)

	// Employee Leave Detail Routes
	api.POST("/employee-leave-details", employeeLeaveDetailController.Create)
	api.GET("/employee-leave-details", employeeLeaveDetailController.List)
	api.GET("/employee-leave-details/:id", employeeLeaveDetailController.Get)
	api.PUT("/employee-leave-details/:id", employeeLeaveDetailController.Update)
	api.DELETE("/employee-leave-details/:id", employeeLeaveDetailController.Delete)

	// Employee Payroll Benefits
	api.POST("/employee-payroll-benefits", employeePayrollBenefitController.Create)
	api.GET("/employee-payroll-benefits", employeePayrollBenefitController.List)
	api.GET("/employee-payroll-benefits/:id", employeePayrollBenefitController.Get)
	api.PUT("/employee-payroll-benefits/:id", employeePayrollBenefitController.Update)
	api.DELETE("/employee-payroll-benefits/:id", employeePayrollBenefitController.Delete)

	// Employee Payroll Deductions
	api.POST("/employee-payroll-deductions", employeePayrollDeductionController.Create)
	api.GET("/employee-payroll-deductions", employeePayrollDeductionController.List)
	api.GET("/employee-payroll-deductions/:id", employeePayrollDeductionController.Get)
	api.PUT("/employee-payroll-deductions/:id", employeePayrollDeductionController.Update)
	api.DELETE("/employee-payroll-deductions/:id", employeePayrollDeductionController.Delete)

	// Employee Payroll Reliefs
	api.POST("/employee-payroll-reliefs", employeePayrollReliefController.Create)
	api.GET("/employee-payroll-reliefs", employeePayrollReliefController.List)
	api.GET("/employee-payroll-reliefs/:id", employeePayrollReliefController.Get)
	api.PUT("/employee-payroll-reliefs/:id", employeePayrollReliefController.Update)
	api.DELETE("/employee-payroll-reliefs/:id", employeePayrollReliefController.Delete)

	// Employee Payslips
	api.POST("/employee-payslips", employeePayslipController.Create)
	api.GET("/employee-payslips", employeePayslipController.List)
	api.GET("/employee-payslips/:id", employeePayslipController.Get)
	api.PUT("/employee-payslips/:id", employeePayslipController.Update)
	api.DELETE("/employee-payslips/:id", employeePayslipController.Delete)

	// Employee Professional Titles
	api.POST("/employee-professional-titles", employeeProfessionalTitleController.Create)
	api.GET("/employee-professional-titles", employeeProfessionalTitleController.List)
	api.GET("/employee-professional-titles/:id", employeeProfessionalTitleController.Get)
	api.PUT("/employee-professional-titles/:id", employeeProfessionalTitleController.Update)
	api.DELETE("/employee-professional-titles/:id", employeeProfessionalTitleController.Delete)

	// Employee Termination Categories
	api.POST("/employee-termination-categories", employeeTerminationCategoryController.Create)
	api.GET("/employee-termination-categories", employeeTerminationCategoryController.List)
	api.GET("/employee-termination-categories/:id", employeeTerminationCategoryController.Get)
	api.PUT("/employee-termination-categories/:id", employeeTerminationCategoryController.Update)
	api.DELETE("/employee-termination-categories/:id", employeeTerminationCategoryController.Delete)
}
