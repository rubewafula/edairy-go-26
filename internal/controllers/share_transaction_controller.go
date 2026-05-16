package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ShareTransactionController struct {
	service *services.ShareTransactionService
}

func NewShareTransactionController() *ShareTransactionController {
	return &ShareTransactionController{
		service: services.NewShareTransactionService(),
	}
}

func (c *ShareTransactionController) CreateShareTransaction(ctx *gin.Context) {
	var req dtos.CreateShareTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	transaction, err := c.service.CreateShareTransaction(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, transaction)
}

func (c *ShareTransactionController) GetShareTransactions(ctx *gin.Context) {
	transactions, total, err := c.service.GetShareTransactions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": transactions, "total": total})
}

func (c *ShareTransactionController) GetShareTransaction(ctx *gin.Context) {
	transaction, err := c.service.GetShareTransaction(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Share transaction not found"})
		return
	}
	ctx.JSON(http.StatusOK, transaction)
}

func (c *ShareTransactionController) UpdateShareTransaction(ctx *gin.Context) {
	var req dtos.UpdateShareTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShareTransaction(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share transaction updated successfully"})
}

func (c *ShareTransactionController) DeleteShareTransaction(ctx *gin.Context) {
	if err := c.service.DeleteShareTransaction(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share transaction deleted successfully"})
}
