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

type LivestockProductionController struct {
	service *services.LivestockProductionService
}

func NewLivestockProductionController() *LivestockProductionController {
	return &LivestockProductionController{
		service: services.NewLivestockProductionService(),
	}
}

func (c *LivestockProductionController) Create(ctx *gin.Context) {
	var req dtos.CreateLivestockProductionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.MustGet("userID").(uint64)
	record, err := c.service.CreateProduction(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, record)
}

func (c *LivestockProductionController) List(ctx *gin.Context) {
	livestockID := ctx.Query("livestock_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Support 'size' parameter from client logs
	if size := ctx.Query("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			limit = s
		}
	}

	results, total, err := c.service.GetProductionRecords(livestockID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *LivestockProductionController) Get(ctx *gin.Context) {
	record, err := c.service.GetProductionRecord(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	ctx.JSON(http.StatusOK, record)
}

func (c *LivestockProductionController) Update(ctx *gin.Context) {
	var req dtos.UpdateLivestockProductionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.MustGet("userID").(uint64)
	if err := c.service.UpdateProductionRecord(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func (c *LivestockProductionController) Delete(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint64)
	if err := c.service.DeleteProductionRecord(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

func (c *LivestockProductionController) Stats(ctx *gin.Context) {
	stats, err := c.service.GetProductionStats(ctx.Param("livestock_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}
