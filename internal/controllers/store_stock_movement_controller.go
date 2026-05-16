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

type StoreStockMovementController struct {
	service *services.StoreStockMovementService
}

func NewStoreStockMovementController() *StoreStockMovementController {
	return &StoreStockMovementController{
		service: services.NewStoreStockMovementService(),
	}
}

func (c *StoreStockMovementController) CreateMovement(ctx *gin.Context) {
	var req dtos.CreateStoreStockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	movement, err := c.service.CreateMovement(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetMovement(utils.Uint64ToString(movement.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *StoreStockMovementController) GetMovements(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetMovements(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreStockMovementController) GetMovement(ctx *gin.Context) {
	result, err := c.service.GetMovement(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stock movement record not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *StoreStockMovementController) UpdateMovement(ctx *gin.Context) {
	var req dtos.UpdateStoreStockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMovement(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Stock movement record updated successfully"})
}

func (c *StoreStockMovementController) DeleteMovement(ctx *gin.Context) {
	if err := c.service.DeleteMovement(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Stock movement record deleted successfully"})
}
