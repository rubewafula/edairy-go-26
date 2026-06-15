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
		log.Printf("[StoreStockMovementController.CreateMovement] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreStockMovementController.CreateMovement] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	movement, err := c.service.CreateMovement(req, userID)
	if err != nil {
		log.Printf("[StoreStockMovementController.CreateMovement] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": movement})
}

func (c *StoreStockMovementController) GetMovements(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetMovements(page, limit)
	if err != nil {
		log.Printf("[StoreStockMovementController.GetMovements] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stock movements"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreStockMovementController) GetMovement(ctx *gin.Context) {
	result, err := c.service.GetMovement(ctx.Param("id"))
	if err != nil {
		log.Printf("[StoreStockMovementController.GetMovement] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stock movement record not found"}) // Specific not found is fine
		return
	}
	ctx.JSON(http.StatusOK, result)
}
