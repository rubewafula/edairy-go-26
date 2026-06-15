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

type StoreItemUnitController struct {
	service *services.StoreItemUnitService
}

func NewStoreItemUnitController() *StoreItemUnitController {
	return &StoreItemUnitController{
		service: services.NewStoreItemUnitService(),
	}
}

func (c *StoreItemUnitController) CreateUnit(ctx *gin.Context) {
	var req dtos.CreateStoreItemUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreItemUnitController.CreateUnit] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreItemUnitController.CreateUnit] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	unit, err := c.service.CreateUnit(req, userID)
	if err != nil {
		log.Printf("[StoreItemUnitController.CreateUnit] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store item unit"})
		return
	}
	ctx.JSON(http.StatusCreated, unit)
}

func (c *StoreItemUnitController) GetUnits(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetUnits(page, limit)
	if err != nil {
		log.Printf("[StoreItemUnitController.GetUnits] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve store item units"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreItemUnitController) GetUnit(ctx *gin.Context) {
	unit, err := c.service.GetUnit(ctx.Param("id"))
	if err != nil {
		log.Printf("[StoreItemUnitController.GetUnit] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Store item unit not found"})
		return
	}
	ctx.JSON(http.StatusOK, unit)
}

func (c *StoreItemUnitController) UpdateUnit(ctx *gin.Context) {
	var req dtos.UpdateStoreItemUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreItemUnitController.UpdateUnit] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreItemUnitController.UpdateUnit] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateUnit(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[StoreItemUnitController.UpdateUnit] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store item unit"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Unit updated successfully"})
}

func (c *StoreItemUnitController) DeleteUnit(ctx *gin.Context) {
	if err := c.service.DeleteUnit(ctx.Param("id")); err != nil {
		log.Printf("[StoreItemUnitController.DeleteUnit] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete store item unit"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Unit deleted successfully"})
}
