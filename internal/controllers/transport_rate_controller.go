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

type TransportRateController struct {
	service *services.TransportRateService
}

func NewTransportRateController() *TransportRateController {
	return &TransportRateController{
		service: services.NewTransportRateService(),
	}
}

func (c *TransportRateController) CreateRate(ctx *gin.Context) {
	var req dtos.CreateTransportRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	rate, err := c.service.CreateTransportRate(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, rate)
}

func (c *TransportRateController) GetRates(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	transporterNo := ctx.Query("transporter_no")
	routeID := ctx.Query("route_id")
	memberNo := ctx.Query("member_no")

	rates, total, err := c.service.GetTransportRates(page, limit, transporterNo, routeID, memberNo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rates, "total": total})
}

func (c *TransportRateController) GetRate(ctx *gin.Context) {
	rate, err := c.service.GetTransportRate(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transport rate not found"})
		return
	}
	ctx.JSON(http.StatusOK, rate)
}

func (c *TransportRateController) UpdateRate(ctx *gin.Context) {
	var req dtos.UpdateTransportRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateTransportRate(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transport rate updated successfully"})
}

func (c *TransportRateController) DeleteRate(ctx *gin.Context) {
	if err := c.service.DeleteTransportRate(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transport rate deleted successfully"})
}
