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
		log.Printf("[StoreSaleItemController.CreateSaleItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreSaleItemController.CreateSaleItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	item, err := c.service.CreateSaleItem(req, userID)
	if err != nil {
		log.Printf("[StoreSaleItemController.CreateSaleItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store sale item"})
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
		log.Printf("[StoreSaleItemController.GetSaleItems] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve store sale items"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreSaleItemController) GetSaleItem(ctx *gin.Context) {
	result, err := c.service.GetSaleItem(ctx.Param("id"))
	if err != nil {
		log.Printf("[StoreSaleItemController.GetSaleItem] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Store sale item not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *StoreSaleItemController) UpdateSaleItem(ctx *gin.Context) {
	var req dtos.UpdateStoreSaleItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreSaleItemController.UpdateSaleItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreSaleItemController.UpdateSaleItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateSaleItem(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[StoreSaleItemController.UpdateSaleItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store sale item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Sale item updated successfully"})
}

func (c *StoreSaleItemController) DeleteSaleItem(ctx *gin.Context) {
	if err := c.service.DeleteSaleItem(ctx.Param("id")); err != nil {
		log.Printf("[StoreSaleItemController.DeleteSaleItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete store sale item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Sale item deleted successfully"})
}
