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

	userID := ctx.MustGet("user_id").(uint64)

	rate, err := c.service.CreateCustomerMilkRate(req, userID)
	if err != nil {
		log.Printf("[CustomerMilkRateController.CreateRate] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetCustomerMilkRate(utils.Uint64ToString(rate.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *CustomerMilkRateController) GetRates(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	rates, total, err := c.service.GetCustomerMilkRates(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rates, "total": total})
}

func (c *CustomerMilkRateController) GetRate(ctx *gin.Context) {
	rate, err := c.service.GetCustomerMilkRate(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer milk rate not found"})
		return
	}
	ctx.JSON(http.StatusOK, rate)
}

func (c *CustomerMilkRateController) UpdateRate(ctx *gin.Context) {
	var req dtos.UpdateCustomerMilkRateRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	if err := c.service.UpdateCustomerMilkRate(id, req, userID); err != nil {
		log.Printf("[CustomerMilkRateController.UpdateRate] Error updating rate %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer milk rate updated successfully"})
}

func (c *CustomerMilkRateController) DeleteRate(ctx *gin.Context) {
	if err := c.service.DeleteCustomerMilkRate(ctx.Param("id")); err != nil {
		log.Printf("[CustomerMilkRateController.DeleteRate] Error deleting rate %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer milk rate deleted successfully"})
}
