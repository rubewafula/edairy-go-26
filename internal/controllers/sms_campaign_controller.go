package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type SMSCampaignController struct {
	service *services.SMSCampaignService
}

func NewSMSCampaignController() *SMSCampaignController {
	return &SMSCampaignController{service: services.NewSMSCampaignService()}
}

func (c *SMSCampaignController) CreateCampaign(ctx *gin.Context) {
	var req dtos.CreateSMSCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (c *SMSCampaignController) GetRecipients(ctx *gin.Context) {
	results, err := c.service.GetRecipients(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}
