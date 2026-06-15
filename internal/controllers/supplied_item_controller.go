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

type SuppliedItemController struct {
	service *services.SuppliedItemService
}

func NewSuppliedItemController() *SuppliedItemController {
	return &SuppliedItemController{
		service: services.NewSuppliedItemService(),
	}
}

func (c *SuppliedItemController) GetSuppliedItem(ctx *gin.Context) {
	result, err := c.service.GetSuppliedItem(ctx.Param("id"))
	if err != nil {
		log.Printf("[SuppliedItemController.GetSuppliedItem] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplied item not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SuppliedItemController) UpdateSuppliedItem(ctx *gin.Context) {
	var req dtos.UpdateSuppliedItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SuppliedItemController.UpdateSuppliedItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SuppliedItemController.UpdateSuppliedItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateSuppliedItem(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[SuppliedItemController.UpdateSuppliedItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplied item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Supplied item updated successfully"})
}

func (c *SuppliedItemController) DeleteSuppliedItem(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeleteSuppliedItem(ctx.Param("id"), userID); err != nil {
		log.Printf("[SuppliedItemController.DeleteSuppliedItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplied item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Supplied item deleted successfully"})
}

func (c *SuppliedItemController) GetItemsBySupply(ctx *gin.Context) {
	items, err := c.service.GetSuppliedItemsBySupply(ctx.Param("supply_id"))
	if err != nil {
		log.Printf("[SuppliedItemController.GetItemsBySupply] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplied items by supply ID"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}
