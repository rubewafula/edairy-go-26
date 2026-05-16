package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	account, err := c.service.CreateAccount(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetAccount(utils.Uint64ToString(account.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterBankAccountController) GetAccounts(ctx *gin.Context) {
	accounts, total, err := c.service.GetAccounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": accounts, "total": total})
}

func (c *TransporterBankAccountController) GetAccount(ctx *gin.Context) {
	account, err := c.service.GetAccount(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Transporter bank account not found"})
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (c *TransporterBankAccountController) UpdateAccount(ctx *gin.Context) {
	var req dtos.UpdateTransporterBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := c.service.UpdateAccount(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter bank account updated successfully"})
}

func (c *TransporterBankAccountController) DeleteAccount(ctx *gin.Context) {
	if err := c.service.DeleteAccount(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter bank account deleted successfully"})
}
