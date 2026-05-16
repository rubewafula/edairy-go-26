package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Supplied item not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SuppliedItemController) UpdateSuppliedItem(ctx *gin.Context) {
	var req dtos.UpdateSuppliedItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateSuppliedItem(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Supplied item updated successfully"})
}

func (c *SuppliedItemController) DeleteSuppliedItem(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeleteSuppliedItem(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Supplied item deleted successfully"})
}

func (c *SuppliedItemController) GetItemsBySupply(ctx *gin.Context) {
	items, err := c.service.GetSuppliedItemsBySupply(ctx.Param("supply_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}
