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

type InterStoreTransferItemController struct {
	service *services.InterStoreTransferItemService
}

func NewInterStoreTransferItemController() *InterStoreTransferItemController {
	return &InterStoreTransferItemController{
		service: services.NewInterStoreTransferItemService(),
	}
}

func (c *InterStoreTransferItemController) CreateTransferItem(ctx *gin.Context) {
	var req dtos.CreateInterStoreTransferItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	item, err := c.service.CreateTransferItem(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetTransferItem(utils.Uint64ToString(item.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *InterStoreTransferItemController) GetTransferItems(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetTransferItems(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *InterStoreTransferItemController) GetTransferItem(ctx *gin.Context) {
	result, err := c.service.GetTransferItem(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Transfer item not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *InterStoreTransferItemController) UpdateTransferItem(ctx *gin.Context) {
	var req dtos.UpdateInterStoreTransferItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateTransferItem(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transfer item updated successfully"})
}

func (c *InterStoreTransferItemController) DeleteTransferItem(ctx *gin.Context) {
	if err := c.service.DeleteTransferItem(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transfer item deleted successfully"})
}
