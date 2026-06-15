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
	"gorm.io/gorm"
)

type SupplierBankAccountController struct {
	service *services.SupplierBankAccountService
}

func NewSupplierBankAccountController() *SupplierBankAccountController {
	return &SupplierBankAccountController{
		service: services.NewSupplierBankAccountService(),
	}
}

func (c *SupplierBankAccountController) CreateBankAccount(ctx *gin.Context) {
	var req dtos.CreateSupplierBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierBankAccountController.CreateBankAccount] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierBankAccountController.CreateBankAccount] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	account, err := c.service.CreateBankAccount(req, userID)
	if err != nil {
		log.Printf("[SupplierBankAccountController.CreateBankAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier bank account"})
		return
	}
	response, _ := c.service.GetBankAccount(utils.Uint64ToString(account.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *SupplierBankAccountController) GetBankAccounts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetBankAccounts(page, limit)
	if err != nil {
		log.Printf("[SupplierBankAccountController.GetBankAccounts] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier bank accounts"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplierBankAccountController) GetBankAccount(ctx *gin.Context) {
	result, err := c.service.GetBankAccount(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier bank account not found"})
			return
		}
		log.Printf("[SupplierBankAccountController.GetBankAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier bank account"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplierBankAccountController) UpdateBankAccount(ctx *gin.Context) {
	var req dtos.UpdateSupplierBankAccountRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierBankAccountController.UpdateBankAccount] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierBankAccountController.UpdateBankAccount] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateBankAccount(id, req, userID); err != nil {
		log.Printf("[SupplierBankAccountController.UpdateBankAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier bank account"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier bank account updated successfully"})
}

func (c *SupplierBankAccountController) DeleteBankAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteBankAccount(id, userID); err != nil {
		log.Printf("[SupplierBankAccountController.DeleteBankAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier bank account"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier bank account deleted successfully"})
}
