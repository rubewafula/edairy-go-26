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

type StoreSaleItemController struct {
	service *services.StoreSaleItemService
}

func NewStoreSaleItemController() *StoreSaleItemController {
	return &StoreSaleItemController{
		service: services.NewStoreSaleItemService(),
	}
}

func (c *StoreSaleItemController) CreateSaleItem(ctx *gin.Context) {
	var req dtos.CreateStoreSaleItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	item, err := c.service.CreateSaleItem(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetSaleItem(utils.Uint64ToString(item.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *StoreSaleItemController) GetSaleItems(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetSaleItems(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreSaleItemController) GetSaleItem(ctx *gin.Context) {
	result, err := c.service.GetSaleItem(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Sale item not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *StoreSaleItemController) UpdateSaleItem(ctx *gin.Context) {
	var req dtos.UpdateStoreSaleItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateSaleItem(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Sale item updated successfully"})
}

func (c *StoreSaleItemController) DeleteSaleItem(ctx *gin.Context) {
	if err := c.service.DeleteSaleItem(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Sale item deleted successfully"})
}
