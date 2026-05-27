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

type TransporterBenefitController struct {
	service *services.TransporterBenefitService
}

func NewTransporterBenefitController() *TransporterBenefitController {
	return &TransporterBenefitController{
		service: services.NewTransporterBenefitService(),
	}
}

func (c *TransporterBenefitController) CreateBenefit(ctx *gin.Context) {
	var req dtos.CreateTransporterBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	benefit, err := c.service.CreateBenefit(req, userID)
	if err != nil {
		log.Printf("[TransporterBenefitController.CreateBenefit] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetBenefit(utils.Uint64ToString(benefit.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterBenefitController) GetBenefits(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	routeID := ctx.Query("route_id")

	results, total, err := c.service.GetBenefits(routeID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *TransporterBenefitController) GetBenefit(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetBenefit(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter benefit not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *TransporterBenefitController) UpdateBenefit(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateTransporterBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateBenefit(id, req, userID); err != nil {
		log.Printf("[TransporterBenefitController.UpdateBenefit] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter benefit updated successfully"})
}

func (c *TransporterBenefitController) DeleteBenefit(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteBenefit(id); err != nil {
		log.Printf("[TransporterBenefitController.DeleteBenefit] Error deleting benefit %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter benefit deleted successfully"})
}
