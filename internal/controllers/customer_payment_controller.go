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

type CustomerPaymentController struct {
	service *services.CustomerPaymentService
}

func NewCustomerPaymentController() *CustomerPaymentController {
	return &CustomerPaymentController{
		service: services.NewCustomerPaymentService(),
	}
}

func (c *CustomerPaymentController) CreatePayment(ctx *gin.Context) {
	var req dtos.CreateCustomerPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	payment, err := c.service.CreatePayment(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetPayment(utils.Uint64ToString(payment.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *CustomerPaymentController) GetPayments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetPayments(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *CustomerPaymentController) GetPayment(ctx *gin.Context) {
	result, err := c.service.GetPayment(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
