package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type SupplierController struct {
	service *services.SupplierService
}

func NewSupplierController() *SupplierController {
	return &SupplierController{
		service: services.NewSupplierService(),
	}
}

func (c *SupplierController) CreateSupplier(ctx *gin.Context) {
	var req dtos.CreateSupplierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	supplier, err := c.service.CreateSupplier(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, supplier)
}

func (c *SupplierController) GetSuppliers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetSuppliers(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplierController) GetSupplier(ctx *gin.Context) {
	result, err := c.service.GetSupplier(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplierController) CreateContact(ctx *gin.Context) {
	var req dtos.CreateSupplierContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	_, err := c.service.CreateContact(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convenience method usually just returns success as the specialized
	// SupplierContactController handles full enriched retrieval.
	ctx.JSON(http.StatusCreated, gin.H{"message": "Supplier contact created successfully"})
}

func (c *SupplierController) GetSupplierContacts(ctx *gin.Context) {
	contacts, err := c.service.GetSupplierContacts(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": contacts})
}

func (c *SupplierController) CreateBankAccount(ctx *gin.Context) {
	var req dtos.CreateSupplierBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	account, err := c.service.CreateBankAccount(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

func (c *SupplierController) GetSupplierBankAccounts(ctx *gin.Context) {
	accounts, err := c.service.GetSupplierBankAccounts(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": accounts})
}
