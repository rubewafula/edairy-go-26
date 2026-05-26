package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerCustomerRoutes(api *gin.RouterGroup) {
	customerController := controllers.NewCustomerController()
	customerTypeController := controllers.NewCustomerTypeController()
	customerBillingController := controllers.NewCustomerBillingController()
	customerInvoiceController := controllers.NewCustomerInvoiceController()
	customerPaymentController := controllers.NewCustomerPaymentController()
	customerMilkRateController := controllers.NewCustomerMilkRateController()
	customerPayDateRangeController := controllers.NewCustomerPayDateRangeController()
	memberPayrollController := controllers.NewMemberPayrollController()

	// Customer Routes
	api.POST("/customers", customerController.CreateCustomer)
	api.GET("/customers", customerController.GetCustomers)
	api.GET("/customers/:id", customerController.GetCustomer)
	api.PUT("/customers/:id", customerController.UpdateCustomer)
	api.DELETE("/customers/:id", customerController.DeleteCustomer)

	// Customer Type Routes
	api.POST("/customer-types", customerTypeController.CreateType)
	api.GET("/customer-types", customerTypeController.GetTypes)
	api.GET("/customer-types/:id", customerTypeController.GetType)
	api.PUT("/customer-types/:id", customerTypeController.UpdateType)
	api.DELETE("/customer-types/:id", customerTypeController.DeleteType)

	// Customer Billing Routes
	api.GET("/customer-billings", customerBillingController.GetBillings)
	api.GET("/customer-billings/:id", customerBillingController.GetBilling)
	api.GET("/customer-billings/:id/items", customerBillingController.GetBillingItems)

	// Customer Invoice Routes
	api.POST("/customer-invoices", customerInvoiceController.CreateInvoice)
	api.GET("/customer-invoices", customerInvoiceController.GetInvoices)
	api.GET("/customer-invoices/:id", customerInvoiceController.GetInvoice)
	api.DELETE("/customer-invoices/:id", customerInvoiceController.DeleteInvoice)

	// Customer Payment Routes
	api.POST("/customer-payments", customerPaymentController.CreatePayment)
	api.GET("/customer-payments", customerPaymentController.GetPayments)
	api.GET("/customer-payments/:id", customerPaymentController.GetPayment)

	// Customer Milk Rate Routes
	api.POST("/customer-milk-rates", customerMilkRateController.CreateRate)
	api.GET("/customer-milk-rates", customerMilkRateController.GetRates)
	api.GET("/customer-milk-rates/:id", customerMilkRateController.GetRate)
	api.PUT("/customer-milk-rates/:id", customerMilkRateController.UpdateRate)
	api.DELETE("/customer-milk-rates/:id", customerMilkRateController.DeleteRate)

	// Customer Pay Date Range Routes
	api.POST("/customer-pay-date-ranges", customerPayDateRangeController.CreateCustomerPayDateRange)
	api.GET("/customer-pay-date-ranges", customerPayDateRangeController.GetCustomerPayDateRanges)
	api.GET("/customer-pay-date-ranges/:id", customerPayDateRangeController.GetCustomerPayDateRange)
	api.PUT("/customer-pay-date-ranges/:id", customerPayDateRangeController.UpdateCustomerPayDateRange)
	api.DELETE("/customer-pay-date-ranges/:id", customerPayDateRangeController.DeleteCustomerPayDateRange)

	// Member Payroll Routes
	api.POST("/member-payrolls", memberPayrollController.Create)
	api.GET("/member-payrolls", memberPayrollController.List)
	api.GET("/member-payrolls/:id", memberPayrollController.Get)
	api.PUT("/member-payrolls/:id", memberPayrollController.Update)
	api.DELETE("/member-payrolls/:id", memberPayrollController.Delete)
}
