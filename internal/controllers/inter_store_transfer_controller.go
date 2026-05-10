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

type InterStoreTransferController struct {
	service *services.InterStoreTransferService
}

func NewInterStoreTransferController() *InterStoreTransferController {
	return &InterStoreTransferController{
		service: services.NewInterStoreTransferService(),
	}
}

func (c *InterStoreTransferController) CreateTransfer(ctx *gin.Context) {
	var req dtos.CreateInterStoreTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	transfer, err := c.service.CreateTransfer(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetTransfer(utils.Uint64ToString(transfer.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *InterStoreTransferController) GetTransfers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetTransfers(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": results, "Total": total})
}

func (c *InterStoreTransferController) GetTransfer(ctx *gin.Context) {
	result, err := c.service.GetTransfer(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Transfer not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *InterStoreTransferController) UpdateTransfer(ctx *gin.Context) {
	var req dtos.UpdateInterStoreTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateTransfer(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transfer updated successfully"})
}

func (c *InterStoreTransferController) DeleteTransfer(ctx *gin.Context) {
	if err := c.service.DeleteTransfer(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transfer deleted successfully"})
}
