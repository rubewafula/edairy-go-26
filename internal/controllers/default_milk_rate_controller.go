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

type DefaultMilkRateController struct {
	service *services.DefaultMilkRateService
}

func NewDefaultMilkRateController() *DefaultMilkRateController {
	return &DefaultMilkRateController{
		service: services.NewDefaultMilkRateService(),
	}
}

func (c *DefaultMilkRateController) CreateRate(ctx *gin.Context) {
	var req dtos.CreateDefaultMilkRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	rate, err := c.service.CreateDefaultMilkRate(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetDefaultMilkRate(utils.Uint64ToString(rate.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *DefaultMilkRateController) GetRates(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	rates, total, err := c.service.GetDefaultMilkRates(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": rates, "Total": total})
}

func (c *DefaultMilkRateController) GetRate(ctx *gin.Context) {
	rate, err := c.service.GetDefaultMilkRate(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Default milk rate not found"})
		return
	}
	ctx.JSON(http.StatusOK, rate)
}

func (c *DefaultMilkRateController) UpdateRate(ctx *gin.Context) {
	var req dtos.UpdateDefaultMilkRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateDefaultMilkRate(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Default milk rate updated successfully"})
}

func (c *DefaultMilkRateController) DeleteRate(ctx *gin.Context) {
	if err := c.service.DeleteDefaultMilkRate(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Default milk rate deleted successfully"})
}
