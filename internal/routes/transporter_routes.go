package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerTransporterRoutes(api *gin.RouterGroup) {
	transporterController := controllers.NewTransporterController()
	individualTransporterController := controllers.NewIndividualTransporterController()
	companyTransporterController := controllers.NewCompanyTransporterController()
	driverController := controllers.NewTransporterDriverController()
	vehicleController := controllers.NewTransporterVehicleController()
	routeAssignmentController := controllers.NewTransporterRouteAssignmentController()
	driverAssignmentController := controllers.NewTransporterDriverAssignmentController()
	transporterBankAccountController := controllers.NewTransporterBankAccountController()
	transporterBenefitController := controllers.NewTransporterBenefitController()
	subRouteController := controllers.NewSubRouteController()
	routeCenterController := controllers.NewRouteCenterController()
	transportRateController := controllers.NewTransportRateController()
	transporterPayDateRangeController := controllers.NewTransporterPayDateRangeController()
	transporterPayrollController := controllers.NewTransporterPayrollController()

	api.POST("/transporters", transporterController.CreateTransporter)
	api.GET("/transporters", transporterController.GetTransporters)
	api.GET("/transporters/:id", transporterController.GetTransporter)
	api.PUT("/transporters/:id", transporterController.UpdateTransporter)
	api.DELETE("/transporters/:id", transporterController.DeleteTransporter)
	api.POST("/transporters/import", transporterController.ImportTransporters)
	api.GET("/transporters/export", transporterController.ExportTransporters)
	api.GET("/transporters/export/download/:filename", transporterController.DownloadExportFile)
	api.GET("/transporters/import-errors/:id", transporterController.GetImportErrors)

	// Individual Transporter Routes
	api.GET("/individual-transporters", individualTransporterController.GetIndividualTransporters)
	api.GET("/individual-transporters/:id", individualTransporterController.GetIndividualTransporter)
	api.PUT("/individual-transporters/:id", individualTransporterController.UpdateIndividualTransporter)

	// Company Transporter Routes
	api.GET("/company-transporters", companyTransporterController.GetCompanyTransporters)
	api.GET("/company-transporters/:id", companyTransporterController.GetCompanyTransporter)
	api.PUT("/company-transporters/:id", companyTransporterController.UpdateCompanyTransporter)

	// Transporter Driver Routes
	api.POST("/transporter-drivers", driverController.CreateDriver)
	api.GET("/transporter-drivers", driverController.GetDrivers)
	api.GET("/transporter-drivers/:id", driverController.GetDriver)
	api.PUT("/transporter-drivers/:id", driverController.UpdateDriver)
	api.DELETE("/transporter-drivers/:id", driverController.DeleteDriver)

	// Transporter Vehicle Routes
	api.POST("/transporter-vehicles", vehicleController.CreateVehicle)
	api.GET("/transporter-vehicles", vehicleController.GetVehicles)
	api.GET("/transporter-vehicles/:id", vehicleController.GetVehicle)
	api.PUT("/transporter-vehicles/:id", vehicleController.UpdateVehicle)
	api.DELETE("/transporter-vehicles/:id", vehicleController.DeleteVehicle)

	// Transporter Route Assignment Routes
	api.POST("/transporter-route-assignments", routeAssignmentController.CreateAssignment)
	api.GET("/transporter-route-assignments", routeAssignmentController.GetAssignments)
	api.GET("/transporter-route-assignments/:id", routeAssignmentController.GetAssignment)
	api.PUT("/transporter-route-assignments/:id", routeAssignmentController.UpdateAssignment)
	api.DELETE("/transporter-route-assignments/:id", routeAssignmentController.DeleteAssignment)

	// Transporter Driver Assignment Routes
	api.POST("/transporter-driver-assignments", driverAssignmentController.CreateAssignment)
	api.GET("/transporter-driver-assignments", driverAssignmentController.GetAssignments)
	api.GET("/transporter-driver-assignments/:id", driverAssignmentController.GetAssignment)
	api.PUT("/transporter-driver-assignments/:id", driverAssignmentController.UpdateAssignment)
	api.DELETE("/transporter-driver-assignments/:id", driverAssignmentController.DeleteAssignment)

	// Transporter Bank Account Routes
	api.POST("/transporter-bank-accounts", transporterBankAccountController.CreateAccount)
	api.GET("/transporter-bank-accounts", transporterBankAccountController.GetAccounts)
	api.GET("/transporter-bank-accounts/:id", transporterBankAccountController.GetAccount)
	api.PUT("/transporter-bank-accounts/:id", transporterBankAccountController.UpdateAccount)
	api.DELETE("/transporter-bank-accounts/:id", transporterBankAccountController.DeleteAccount)

	// Transporter Benefit Routes
	api.POST("/transporter-benefits", transporterBenefitController.CreateBenefit)
	api.GET("/transporter-benefits", transporterBenefitController.GetBenefits)
	api.GET("/transporter-benefits/:id", transporterBenefitController.GetBenefit)
	api.PUT("/transporter-benefits/:id", transporterBenefitController.UpdateBenefit)
	api.DELETE("/transporter-benefits/:id", transporterBenefitController.DeleteBenefit)

	// SubRoute Routes
	api.POST("/sub-routes", subRouteController.CreateSubRoute)
	api.GET("/sub-routes", subRouteController.GetSubRoutes)
	api.GET("/sub-routes/:id", subRouteController.GetSubRoute)
	api.PUT("/sub-routes/:id", subRouteController.UpdateSubRoute)
	api.DELETE("/sub-routes/:id", subRouteController.DeleteSubRoute)

	// Route Center Routes
	api.POST("/route-centers", routeCenterController.CreateCenter)
	api.GET("/route-centers", routeCenterController.GetCenters)
	api.GET("/route-centers/:id", routeCenterController.GetCenter)
	api.PUT("/route-centers/:id", routeCenterController.UpdateCenter)
	api.DELETE("/route-centers/:id", routeCenterController.DeleteCenter)

	// Transport Rate Routes
	api.POST("/transport-rates", transportRateController.CreateRate)
	api.GET("/transport-rates", transportRateController.GetRates)
	api.GET("/transport-rates/:id", transportRateController.GetRate)
	api.PUT("/transport-rates/:id", transportRateController.UpdateRate)
	api.DELETE("/transport-rates/:id", transportRateController.DeleteRate)

	// Transporter Pay Date Range Routes
	api.POST("/transporter-pay-date-ranges", transporterPayDateRangeController.Create)
	api.GET("/transporter-pay-date-ranges", transporterPayDateRangeController.List)
	api.GET("/transporter-pay-date-ranges/:id", transporterPayDateRangeController.Get)
	api.PUT("/transporter-pay-date-ranges/:id", transporterPayDateRangeController.Update)
	api.DELETE("/transporter-pay-date-ranges/:id", transporterPayDateRangeController.Delete)

	// Transporter Payroll Routes
	api.POST("/transporter-payrolls", transporterPayrollController.Create)
	api.GET("/transporter-payrolls", transporterPayrollController.List)
	api.GET("/transporter-payrolls/:id", transporterPayrollController.Get)
	api.PUT("/transporter-payrolls/:id/confirm", transporterPayrollController.Confirm)
	api.PUT("/transporter-payrolls/:id/approve", transporterPayrollController.Approve)
}
