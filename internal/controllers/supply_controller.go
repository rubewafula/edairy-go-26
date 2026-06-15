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

type SupplyController struct {
	service *services.SupplyService
}

func NewSupplyController() *SupplyController {
	return &SupplyController{
		service: services.NewSupplyService(),
	}
}

func (c *SupplyController) CreateSupply(ctx *gin.Context) {
	var req dtos.CreateSupplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplyController.CreateSupply] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplyController.CreateSupply] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	supply, err := c.service.CreateSupply(req, userID)
	if err != nil {
		log.Printf("[SupplyController.CreateSupply] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supply record"})
		return
	}
	ctx.JSON(http.StatusCreated, supply)
}

func (c *SupplyController) GetSupplies(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetSupplies(page, limit)
	if err != nil {
		log.Printf("[SupplyController.GetSupplies] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplies"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplyController) GetSuppliedItems(ctx *gin.Context) {
	items, err := c.service.GetSuppliedItems(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplyController.GetSuppliedItems] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplied items"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}
