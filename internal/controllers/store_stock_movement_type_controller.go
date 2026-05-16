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

type StoreStockMovementTypeController struct {
	service *services.StoreStockMovementTypeService
}

func NewStoreStockMovementTypeController() *StoreStockMovementTypeController {
	return &StoreStockMovementTypeController{
		service: services.NewStoreStockMovementTypeService(),
	}
}

func (c *StoreStockMovementTypeController) CreateMovementType(ctx *gin.Context) {
	var req dtos.CreateStoreStockMovementTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	movementType, err := c.service.CreateMovementType(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, movementType)
}

func (c *StoreStockMovementTypeController) GetMovementTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetMovementTypes(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreStockMovementTypeController) GetMovementType(ctx *gin.Context) {
	movementType, err := c.service.GetMovementType(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Movement type not found"})
		return
	}
	ctx.JSON(http.StatusOK, movementType)
}

func (c *StoreStockMovementTypeController) UpdateMovementType(ctx *gin.Context) {
	var req dtos.UpdateStoreStockMovementTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMovementType(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Movement type updated successfully"})
}

func (c *StoreStockMovementTypeController) DeleteMovementType(ctx *gin.Context) {
	if err := c.service.DeleteMovementType(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Movement type deleted successfully"})
}
