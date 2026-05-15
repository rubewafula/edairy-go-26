package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerSupplierRoutes(api *gin.RouterGroup) {
	supplierCategoryController := controllers.NewSupplierCategoryController()
	supplierController := controllers.NewSupplierController()
	supplierContactController := controllers.NewSupplierContactController()
	supplierDocumentController := controllers.NewSupplierDocumentController()
	supplierQuoteController := controllers.NewSupplierQuoteController()
	supplyController := controllers.NewSupplyController()
	supplierBankAccountController := controllers.NewSupplierBankAccountController()
	suppliedItemController := controllers.NewSuppliedItemController()
	supplyRejectController := controllers.NewSupplyRejectController()
	purchaseOrderController := controllers.NewPurchaseOrderController()

	// Supplier Category Routes
	api.POST("/supplier-categories", supplierCategoryController.CreateCategory)
	api.GET("/supplier-categories", supplierCategoryController.GetCategories)
	api.GET("/supplier-categories/:id", supplierCategoryController.GetCategory)
	api.PUT("/supplier-categories/:id", supplierCategoryController.UpdateCategory)
	api.DELETE("/supplier-categories/:id", supplierCategoryController.DeleteCategory)

	// Supplier Routes
	api.POST("/suppliers", supplierController.CreateSupplier)
	api.GET("/suppliers", supplierController.GetSuppliers)
	api.GET("/suppliers/:id", supplierController.GetSupplier)
	api.POST("/suppliers/:id/contacts", supplierController.CreateContact)
	api.GET("/suppliers/:id/contacts", supplierController.GetSupplierContacts)
	api.POST("/suppliers/:id/bank-accounts", supplierController.CreateBankAccount)
	api.GET("/suppliers/:id/bank-accounts", supplierController.GetSupplierBankAccounts)

	// Supplier Contact Routes
	api.POST("/supplier-contacts", supplierContactController.CreateContact)
	api.GET("/supplier-contacts", supplierContactController.GetContacts)
	api.GET("/supplier-contacts/:id", supplierContactController.GetContact)
	api.PUT("/supplier-contacts/:id", supplierContactController.UpdateContact)
	api.DELETE("/supplier-contacts/:id", supplierContactController.DeleteContact)

	// Supplier Document Routes
	api.POST("/supplier-documents", supplierDocumentController.CreateDocument)
	api.GET("/supplier-documents", supplierDocumentController.GetDocuments)
	api.GET("/supplier-documents/:id", supplierDocumentController.GetDocument)
	api.PUT("/supplier-documents/:id", supplierDocumentController.UpdateDocument)
	api.DELETE("/supplier-documents/:id", supplierDocumentController.DeleteDocument)
	api.PATCH("/supplier-documents/:id/verify", supplierDocumentController.VerifyDocument)

	// Supplier Bank Account Routes
	api.POST("/supplier-bank-accounts", supplierBankAccountController.CreateBankAccount)
	api.GET("/supplier-bank-accounts", supplierBankAccountController.GetBankAccounts)
	//api.GET("/supplier-bank-accounts/:id", supplierBankAccountController.GetAccount)
	api.PUT("/supplier-bank-accounts/:id", supplierBankAccountController.UpdateBankAccount)
	api.DELETE("/supplier-bank-accounts/:id", supplierBankAccountController.DeleteBankAccount)

	// Supplier Quote Routes
	api.POST("/supplier-quotes", supplierQuoteController.CreateQuote)
	api.GET("/supplier-quotes", supplierQuoteController.GetQuotes)
	api.POST("/supplier-quotes/:id/items", supplierQuoteController.CreateQuoteItem)
	api.GET("/supplier-quotes/:id/items", supplierQuoteController.GetQuoteItems)
	api.GET("/supplier-quote-items/:id", supplierQuoteController.GetQuoteItem)
	api.PUT("/supplier-quote-items/:id", supplierQuoteController.UpdateQuoteItem)
	api.DELETE("/supplier-quote-items/:id", supplierQuoteController.DeleteQuoteItem)

	// Supply Routes
	api.POST("/supplies", supplyController.CreateSupply)
	api.GET("/supplies", supplyController.GetSupplies)
	api.GET("/supplies/:id/items", supplyController.GetSuppliedItems)

	// Supplied Item Routes
	api.GET("/supplied-items/:id", suppliedItemController.GetSuppliedItem)
	api.PUT("/supplied-items/:id", suppliedItemController.UpdateSuppliedItem)
	api.DELETE("/supplied-items/:id", suppliedItemController.DeleteSuppliedItem)
	//api.GET("/supplies/:supply_id/items-list", suppliedItemController.GetItemsBySupply) // This was commented out

	// Supply Reject Routes
	api.POST("/supply-rejects", supplyRejectController.CreateReject)
	api.GET("/supply-rejects", supplyRejectController.GetRejects)
	api.GET("/supplies/:id/rejects", supplyRejectController.GetRejectsBySupply)
	api.GET("/supply-rejects/:id", supplyRejectController.GetReject)
	api.PUT("/supply-rejects/:id", supplyRejectController.UpdateReject)
	api.DELETE("/supply-rejects/:id", supplyRejectController.DeleteReject)

	// Purchase Order Routes
	api.POST("/purchase-orders", purchaseOrderController.CreatePO)
	api.GET("/purchase-orders", purchaseOrderController.GetPOs)
	api.GET("/purchase-orders/:id", purchaseOrderController.GetPO)
	api.PUT("/purchase-orders/:id", purchaseOrderController.UpdatePO)
	api.DELETE("/purchase-orders/:id", purchaseOrderController.DeletePO)
	api.GET("/purchase-orders/:id/items", purchaseOrderController.GetPOItems)
	api.GET("/purchase-order-items/:id", purchaseOrderController.GetPOItem)
	api.PUT("/purchase-order-items/:id", purchaseOrderController.UpdatePOItem)
	api.DELETE("/purchase-order-items/:id", purchaseOrderController.DeletePOItem)

	// Purchase Requisition Routes
	api.POST("/purchase-requisitions", purchaseOrderController.CreateRequisition)
	api.GET("/purchase-requisitions", purchaseOrderController.GetRequisitions)
	api.GET("/purchase-requisitions/:id", purchaseOrderController.GetRequisition)
	api.PUT("/purchase-requisitions/:id", purchaseOrderController.UpdateRequisition)
	api.DELETE("/purchase-requisitions/:id", purchaseOrderController.DeleteRequisition)
	api.POST("/purchase-requisition-items", purchaseOrderController.CreateRequisitionItem)
	api.GET("/purchase-requisitions/:id/items", purchaseOrderController.GetRequisitionItems)
	api.GET("/purchase-requisition-items/:id", purchaseOrderController.GetRequisitionItem)
	api.PUT("/purchase-requisition-items/:id", purchaseOrderController.UpdateRequisitionItem)
	api.DELETE("/purchase-requisition-items/:id", purchaseOrderController.DeleteRequisitionItem)
}
