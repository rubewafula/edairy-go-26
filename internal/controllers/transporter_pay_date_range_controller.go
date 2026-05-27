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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.Create(req, userID)
	if err != nil {
		log.Printf("[TransporterPayDateRangeController.Create] Error: %v", err)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *TransporterPayDateRangeController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateTransporterPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.Update(id, req, userID); err != nil {
		log.Printf("[TransporterPayDateRangeController.Update] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter pay date range updated successfully"})
}

func (c *TransporterPayDateRangeController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		log.Printf("[TransporterPayDateRangeController.Delete] Error deleting range %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter pay date range deleted successfully"})
}
