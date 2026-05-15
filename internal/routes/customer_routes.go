package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerCustomerRoutes(api *gin.RouterGroup) {
	customerController := controllers.NewCustomerController()
	customerClassController := controllers.NewCustomerClassController()
	customerTypeController := controllers.NewCustomerTypeController()
	customerBillingController := controllers.NewCustomerBillingController()
	customerInvoiceController := controllers.NewCustomerInvoiceController()
	customerPaymentController := controllers.NewCustomerPaymentController()
	customerMilkRateController := controllers.NewCustomerMilkRateController()
	customerPayDateRangeController := controllers.NewCustomerPayDateRangeController()

	// Customer Routes
	api.POST("/customers", customerController.CreateCustomer)
	api.GET("/customers", customerController.GetCustomers)
	api.GET("/customers/:id", customerController.GetCustomer)
	api.PUT("/customers/:id", customerController.UpdateCustomer)
	api.DELETE("/customers/:id", customerController.DeleteCustomer)

	// Customer Class Routes
	api.POST("/customer-classes", customerClassController.CreateClass)
	api.GET("/customer-classes", customerClassController.GetClasses)
	api.GET("/customer-classes/:id", customerClassController.GetClass)
	api.PUT("/customer-classes/:id", customerClassController.UpdateClass)
	api.DELETE("/customer-classes/:id", customerClassController.DeleteClass)

	// Customer Type Routes
	api.POST("/customer-types", customerTypeController.CreateCustomerType)
	api.GET("/customer-types", customerTypeController.GetCustomerTypes)
	api.GET("/customer-types/:id", customerTypeController.GetCustomerType)
	api.PUT("/customer-types/:id", customerTypeController.UpdateCustomerType)
	api.DELETE("/customer-types/:id", customerTypeController.DeleteCustomerType)

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
}
