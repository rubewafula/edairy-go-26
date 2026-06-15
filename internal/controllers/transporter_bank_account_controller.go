package controllers

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type TransporterBankAccountController struct {
	service *services.TransporterBankAccountService
}

func NewTransporterBankAccountController() *TransporterBankAccountController {
	return &TransporterBankAccountController{
		service: services.NewTransporterBankAccountService(),
	}
}

func (c *TransporterBankAccountController) CreateAccount(ctx *gin.Context) {
	var req dtos.CreateTransporterBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterBankAccountController.CreateAccount] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterBankAccountController.CreateAccount] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	account, err := c.service.CreateAccount(req)
	if err != nil {
		log.Printf("[TransporterBankAccountController.CreateAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transporter bank account"})
		return
	}

	response, _ := c.service.GetAccount(utils.Uint64ToString(account.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterBankAccountController) GetAccounts(ctx *gin.Context) {
	accounts, total, err := c.service.GetAccounts()
	if err != nil {
		log.Printf("[TransporterBankAccountController.GetAccounts] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transporter bank accounts"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": accounts, "total": total})
}

func (c *TransporterBankAccountController) GetAccount(ctx *gin.Context) {
	account, err := c.service.GetAccount(ctx.Param("id"))
	if err != nil {
		log.Printf("[TransporterBankAccountController.GetAccount] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter bank account not found"})
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (c *TransporterBankAccountController) UpdateAccount(ctx *gin.Context) {
	var req dtos.UpdateTransporterBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterBankAccountController.UpdateAccount] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := c.service.UpdateAccount(ctx.Param("id"), req); err != nil {
		log.Printf("[TransporterBankAccountController.UpdateAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transporter bank account"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter bank account updated successfully"})
}

func (c *TransporterBankAccountController) DeleteAccount(ctx *gin.Context) {
	if err := c.service.DeleteAccount(ctx.Param("id")); err != nil {
		log.Printf("[TransporterBankAccountController.DeleteAccount] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transporter bank account"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter bank account deleted successfully"})
}
