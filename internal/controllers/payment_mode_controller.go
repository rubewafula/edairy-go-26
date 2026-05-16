package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type PaymentModeController struct {
	service *services.PaymentModeService
}

func NewPaymentModeController() *PaymentModeController {
	return &PaymentModeController{
		service: services.NewPaymentModeService(),
	}
}

func (c *PaymentModeController) CreatePaymentMode(ctx *gin.Context) {
	var req dtos.CreatePaymentModeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	paymentMode, err := c.service.CreatePaymentMode(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, paymentMode)
}

func (c *PaymentModeController) GetPaymentModes(ctx *gin.Context) {
	paymentModes, total, err := c.service.GetPaymentModes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": paymentModes, "total": total})
}

func (c *PaymentModeController) GetPaymentMode(ctx *gin.Context) {
	paymentMode, err := c.service.GetPaymentMode(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment mode not found"})
		return
	}
	ctx.JSON(http.StatusOK, paymentMode)
}

func (c *PaymentModeController) UpdatePaymentMode(ctx *gin.Context) {
	var req dtos.UpdatePaymentModeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdatePaymentMode(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Payment mode updated successfully"})
}

func (c *PaymentModeController) DeletePaymentMode(ctx *gin.Context) {
	if err := c.service.DeletePaymentMode(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Payment mode deleted successfully"})
}
