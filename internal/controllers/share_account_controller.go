package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ShareAccountController struct {
	service *services.ShareAccountService
}

func NewShareAccountController() *ShareAccountController {
	return &ShareAccountController{
		service: services.NewShareAccountService(),
	}
}

func (c *ShareAccountController) CreateAccount(ctx *gin.Context) {
	var req dtos.CreateShareAccountRequest
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
	ctx.JSON(http.StatusCreated, account)
}

func (c *ShareAccountController) GetAccounts(ctx *gin.Context) {
	accounts, total, err := c.service.GetShareAccounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": accounts, "Total": total})
}

func (c *ShareAccountController) GetAccount(ctx *gin.Context) {
	account, err := c.service.GetShareAccount(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Share account not found"})
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (c *ShareAccountController) UpdateAccount(ctx *gin.Context) {
	var req dtos.UpdateShareAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAccount(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share account updated successfully"})
}

func (c *ShareAccountController) DeleteAccount(ctx *gin.Context) {
	if err := c.service.DeleteAccount(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share account deleted successfully"})
}
