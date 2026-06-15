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
	"gorm.io/gorm"
)

type TransporterPayDateRangeController struct {
	service *services.TransporterPayDateRangeService
}

func NewTransporterPayDateRangeController() *TransporterPayDateRangeController {
	return &TransporterPayDateRangeController{
		service: services.NewTransporterPayDateRangeService(),
	}
}

func (c *TransporterPayDateRangeController) Create(ctx *gin.Context) {
	var req dtos.CreateTransporterPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterPayDateRangeController.Create] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterPayDateRangeController.Create] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.Create(req, userID)
	if err != nil {
		log.Printf("[TransporterPayDateRangeController.Create] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *TransporterPayDateRangeController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	res, total, err := c.service.List(page, limit)
	if err != nil {
		log.Printf("[TransporterPayDateRangeController.List] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pay date ranges"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *TransporterPayDateRangeController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter pay date range not found"})
			return
		}
		log.Printf("[TransporterPayDateRangeController.Get] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pay date range"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *TransporterPayDateRangeController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateTransporterPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterPayDateRangeController.Update] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.Update(id, req, userID); err != nil {
		log.Printf("[TransporterPayDateRangeController.Update] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter pay date range updated successfully"})
}

func (c *TransporterPayDateRangeController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		log.Printf("[TransporterPayDateRangeController.Delete] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter pay date range deleted successfully"})
}
