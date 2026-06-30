package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

// registerMemberRoutes registers the API routes for member-related functionalities.
func registerMemberRoutes(api *gin.RouterGroup) {
	memberDashboardController := controllers.NewMemberDashboardController()
	payrollDashboardController := controllers.NewPayrollDashboardController()
	memberController := controllers.NewMemberController()
	memberBankAccountController := controllers.NewMemberBankAccountController()
	memberDependantController := controllers.NewMemberDependantController()
	memberTypeController := controllers.NewMemberTypeController()
	memberNextOfKinController := controllers.NewMemberNextOfKinController()

	memberPayrollController := controllers.NewMemberPayrollController()
	memberPayDateRangeController := controllers.NewMemberPayDateRangeController()
	memberPayslipController := controllers.NewMemberPayslipController()
	memberPayrollDeductionController := controllers.NewMemberPayrollDeductionController()

	api.GET("/member-dashboard", memberDashboardController.GetDashboard)
	api.GET("/payroll-dashboard", payrollDashboardController.GetDashboard)

	// Member Registration Routes
	api.POST("/members", memberController.CreateMember)
	api.GET("/members", memberController.GetMembers)
	api.GET("/members/export", memberController.ExportMembers)
	api.GET("/members/agm-reports/export", memberController.ExportAGMReport)
	api.GET("/members/export/download/:filename", memberController.DownloadExportFile)
	api.GET("/members/agm-report/download/:filename", memberController.DownloadExportFile)
	api.GET("/members/:id", memberController.GetMember)
	api.PUT("/members/:id", memberController.UpdateMember)
	api.DELETE("/members/:id", memberController.DeleteMember)
	api.PUT("/members/suspend/:id", memberController.SuspendMember)
	api.POST("/members/import", memberController.ImportMembers)
	api.GET("/members/import-errors/:importid", memberController.GetMemberImportErrors)

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

	api.POST("/member-next-of-kins", memberNextOfKinController.CreateMemberNextOfKin)
	api.GET("/member-next-of-kins", memberNextOfKinController.GetMemberNextOfKins)
	api.GET("/member-next-of-kins/:id", memberNextOfKinController.GetMemberNextOfKin)
	api.PUT("/member-next-of-kins/:id", memberNextOfKinController.UpdateMemberNextOfKin)
	api.DELETE("/member-next-of-kins/:id", memberNextOfKinController.DeleteMemberNextOfKin)

	// Member Pay Date Range Routes
	api.POST("/member-pay-date-ranges", memberPayDateRangeController.Create)
	api.GET("/member-pay-date-ranges", memberPayDateRangeController.List)
	api.GET("/member-pay-date-ranges/next", memberPayDateRangeController.GetNextRange)
	api.GET("/member-pay-date-ranges/:id", memberPayDateRangeController.Get)
	api.PUT("/member-pay-date-ranges/:id", memberPayDateRangeController.Update)
	api.DELETE("/member-pay-date-ranges/:id", memberPayDateRangeController.Delete)

	// Member Payroll Routes
	api.POST("/member-payrolls", memberPayrollController.Create)
	api.GET("/member-payrolls", memberPayrollController.List)
	api.GET("/member-payrolls/:id", memberPayrollController.Get)
	api.PUT("/member-payrolls/:id/confirm", memberPayrollController.Confirm)
	api.PUT("/member-payrolls/:id/approve", memberPayrollController.Approve)
	api.GET("/member-payrolls/generation-errors/:payrollID", memberPayrollController.GetGenerationErrors)
	api.GET("/member-payrolls/approval-errors/:payrollID", memberPayrollController.GetApprovalErrors)

	// Member Payslip Routes
	api.GET("/member-payslips", memberPayslipController.GetPayslips)
	api.GET("/member-payslips/statements/:payslip_id/:member_id", memberPayslipController.ExportStatements)
	api.GET("/member-payslips/export", memberPayslipController.ExportPayslips)
	api.GET("/member-payslips/export/download/:filename", memberPayslipController.DownloadExportFile)
	api.GET("/member-payslips/:id", memberPayslipController.GetPayslip)

	// Member Payroll Deduction Routes
	api.GET("/member-payroll-deductions", memberPayrollDeductionController.GetMemberPayrollDeductions)
	api.GET("/member-payroll-deductions/export", memberPayrollDeductionController.ExportDeductions)
	api.GET("/member-payroll-deductions/export/download/:filename", memberPayrollDeductionController.DownloadExportFile)
	api.GET("/member-payroll-deductions/:id", memberPayrollDeductionController.GetMemberPayrollDeduction)

}
