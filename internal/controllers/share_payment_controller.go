package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type SharePaymentController struct {
	service *services.SharePaymentService
}

func NewSharePaymentController() *SharePaymentController {
	return &SharePaymentController{
		service: services.NewSharePaymentService(),
	}
}

func (c *SharePaymentController) CreateSharePayment(ctx *gin.Context) {
	var req dtos.CreateSharePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	payment, err := c.service.CreateSharePayment(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, payment)
}

func (c *SharePaymentController) GetSharePayments(ctx *gin.Context) {
	payments, total, err := c.service.GetSharePayments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": payments, "total": total})
}

func (c *SharePaymentController) GetSharePayment(ctx *gin.Context) {
	payment, err := c.service.GetSharePayment(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Share payment not found"})
		return
	}
	ctx.JSON(http.StatusOK, payment)
}

func (c *SharePaymentController) UpdateSharePayment(ctx *gin.Context) {
	var req dtos.UpdateSharePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateSharePayment(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share payment updated successfully"})
}

func (c *SharePaymentController) DeleteSharePayment(ctx *gin.Context) {
	if err := c.service.DeleteSharePayment(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share payment deleted successfully"})
}
