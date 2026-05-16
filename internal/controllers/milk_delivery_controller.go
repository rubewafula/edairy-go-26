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

type MilkDeliveryController struct {
	service *services.MilkDeliveryService
}

func NewMilkDeliveryController() *MilkDeliveryController {
	return &MilkDeliveryController{
		service: services.NewMilkDeliveryService(),
	}
}

func (c *MilkDeliveryController) CreateDelivery(ctx *gin.Context) {
	var req dtos.CreateMilkDeliveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	delivery, err := c.service.CreateDelivery(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetDelivery(utils.Uint64ToString(delivery.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *MilkDeliveryController) GetDeliveries(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	deliveries, total, err := c.service.GetDeliveries(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": deliveries, "total": total})
}

func (c *MilkDeliveryController) GetDelivery(ctx *gin.Context) {
	delivery, err := c.service.GetDelivery(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery entry not found"})
		return
	}
	ctx.JSON(http.StatusOK, delivery)
}

func (c *MilkDeliveryController) UpdateDelivery(ctx *gin.Context) {
	var req dtos.UpdateMilkDeliveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateDelivery(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Delivery entry updated successfully"})
}

func (c *MilkDeliveryController) DeleteDelivery(ctx *gin.Context) {
	if err := c.service.DeleteDelivery(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Delivery entry deleted successfully"})
}
