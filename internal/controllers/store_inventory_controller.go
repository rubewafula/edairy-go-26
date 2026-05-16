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

type StoreInventoryController struct {
	service *services.StoreInventoryService
}

func NewStoreInventoryController() *StoreInventoryController {
	return &StoreInventoryController{
		service: services.NewStoreInventoryService(),
	}
}

func (c *StoreInventoryController) CreateInventory(ctx *gin.Context) {
	var req dtos.CreateStoreInventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	inventory, err := c.service.CreateInventory(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetInventory(utils.Uint64ToString(inventory.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *StoreInventoryController) GetInventories(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetInventories(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreInventoryController) GetInventory(ctx *gin.Context) {
	inventory, err := c.service.GetInventory(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}
	ctx.JSON(http.StatusOK, inventory)
}

func (c *StoreInventoryController) UpdateInventory(ctx *gin.Context) {
	var req dtos.UpdateStoreInventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateInventory(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Inventory updated successfully"})
}

func (c *StoreInventoryController) DeleteInventory(ctx *gin.Context) {
	if err := c.service.DeleteInventory(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Inventory deleted successfully"})
}
