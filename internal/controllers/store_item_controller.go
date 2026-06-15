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

type StoreItemController struct {
	service *services.StoreItemService
}

func NewStoreItemController() *StoreItemController {
	return &StoreItemController{
		service: services.NewStoreItemService(),
	}
}

func (c *StoreItemController) CreateItem(ctx *gin.Context) {
	var req dtos.CreateStoreItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreItemController.CreateItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreItemController.CreateItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	item, err := c.service.CreateStoreItem(req)
	if err != nil {
		log.Printf("[StoreItemController.CreateItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store item"})
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *StoreItemController) GetItems(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	items, total, err := c.service.GetStoreItems(page, limit)
	if err != nil {
		log.Printf("[StoreItemController.GetItems] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve store items"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (c *StoreItemController) GetItem(ctx *gin.Context) {
	item, err := c.service.GetStoreItem(ctx.Param("id"))
	if err != nil {
		log.Printf("[StoreItemController.GetItem] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Store item not found"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *StoreItemController) UpdateItem(ctx *gin.Context) {
	var req dtos.UpdateStoreItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreItemController.UpdateItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreItemController.UpdateItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateStoreItem(ctx.Param("id"), req); err != nil {
		log.Printf("[StoreItemController.UpdateItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Item updated successfully"})
}

func (c *StoreItemController) DeleteItem(ctx *gin.Context) {
	if err := c.service.DeleteStoreItem(ctx.Param("id")); err != nil {
		log.Printf("[StoreItemController.DeleteItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete store item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Item deleted successfully"})
}
