package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type SMSCampaignController struct {
	service *services.SMSCampaignService
}

func NewSMSCampaignController() *SMSCampaignController {
	return &SMSCampaignController{service: services.NewSMSCampaignService()}
}

func (c *SMSCampaignController) CreateCampaign(ctx *gin.Context) {
	var req dtos.CreateSMSCampaignRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	campaign, err := c.service.CreateCampaign(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, campaign)
}

func (c *SMSCampaignController) GetCampaigns(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetCampaigns(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SMSCampaignController) GetCampaign(ctx *gin.Context) {
	result, err := c.service.GetCampaign(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS campaign not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SMSCampaignController) UpdateCampaign(ctx *gin.Context) {
	var req dtos.UpdateSMSCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateCampaign(ctx.Param("id"), req, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS campaign not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS campaign updated successfully"})
}

func (c *SMSCampaignController) DeleteCampaign(ctx *gin.Context) {
	if err := c.service.DeleteCampaign(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS campaign not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
