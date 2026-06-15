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

type SupplyRejectController struct {
	service *services.SupplyRejectService
}

func NewSupplyRejectController() *SupplyRejectController {
	return &SupplyRejectController{
		service: services.NewSupplyRejectService(),
	}
}

func (c *SupplyRejectController) CreateReject(ctx *gin.Context) {
	var req dtos.CreateSupplyRejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplyRejectController.CreateReject] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplyRejectController.CreateReject] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	reject, err := c.service.CreateReject(req, userID)
	if err != nil {
		log.Printf("[SupplyRejectController.CreateReject] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supply reject record"})
		return
	}
	ctx.JSON(http.StatusCreated, reject)
}

func (c *SupplyRejectController) GetRejects(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetRejects(page, limit)
	if err != nil {
		log.Printf("[SupplyRejectController.GetRejects] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supply rejects"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplyRejectController) GetRejectsBySupply(ctx *gin.Context) {
	results, err := c.service.GetRejectsBySupply(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplyRejectController.GetRejectsBySupply] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supply rejects by supply ID"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func (c *SupplyRejectController) GetReject(ctx *gin.Context) {
	result, err := c.service.GetReject(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplyRejectController.GetReject] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Supply reject record not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplyRejectController) UpdateReject(ctx *gin.Context) {
	var req dtos.UpdateSupplyRejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateReject(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Reject record updated successfully"})
}

func (c *SupplyRejectController) DeleteReject(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteReject(ctx.Param("id"), userID); err != nil {
		log.Printf("[SupplyRejectController.DeleteReject] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supply reject record"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Reject record deleted successfully"})
}
