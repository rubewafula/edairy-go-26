package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CustomerMilkRateController struct {
	service *services.CustomerMilkRateService
}

func NewCustomerMilkRateController() *CustomerMilkRateController {
	return &CustomerMilkRateController{
		service: services.NewCustomerMilkRateService(),
	}
}

func (c *CustomerMilkRateController) CreateRate(ctx *gin.Context) {
	var req dtos.CreateCustomerMilkRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	rate, err := c.service.CreateCustomerMilkRate(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, rate)
}

func (c *CustomerMilkRateController) GetRates(ctx *gin.Context) {
	rates, total, err := c.service.GetCustomerMilkRates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rates, "total": total})
}

func (c *CustomerMilkRateController) GetRate(ctx *gin.Context) {
	rate, err := c.service.GetCustomerMilkRate(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Rate not found"})
		return
	}
	ctx.JSON(http.StatusOK, rate)
}

func (c *CustomerMilkRateController) UpdateRate(ctx *gin.Context) {
	var req dtos.UpdateCustomerMilkRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCustomerMilkRate(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Rate updated successfully"})
}

func (c *CustomerMilkRateController) DeleteRate(ctx *gin.Context) {
	if err := c.service.DeleteCustomerMilkRate(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Rate deleted successfully"})
}
