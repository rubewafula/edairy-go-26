package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CustomerInvoiceController struct {
	service *services.CustomerInvoiceService
}

func NewCustomerInvoiceController() *CustomerInvoiceController {
	return &CustomerInvoiceController{
		service: services.NewCustomerInvoiceService(),
	}
}

func (c *CustomerInvoiceController) CreateInvoice(ctx *gin.Context) {
	var req dtos.CreateCustomerInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	invoice, err := c.service.CreateInvoice(req, userID)
	if err != nil {
		log.Printf("[CustomerInvoiceController.CreateInvoice] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetInvoice(utils.Uint64ToString(invoice.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *CustomerInvoiceController) GetInvoices(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetInvoices(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *CustomerInvoiceController) GetInvoice(ctx *gin.Context) {
	result, err := c.service.GetInvoice(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *CustomerInvoiceController) DeleteInvoice(ctx *gin.Context) {
	if err := c.service.DeleteInvoice(ctx.Param("id")); err != nil {
		log.Printf("[CustomerInvoiceController.DeleteInvoice] Error deleting invoice %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
