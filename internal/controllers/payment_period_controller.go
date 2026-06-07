package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type PaymentPeriodController struct {
	service *services.PaymentPeriodService
}

func NewPaymentPeriodController() *PaymentPeriodController {
	return &PaymentPeriodController{
		service: services.NewPaymentPeriodService(),
	}
}

func (ctrl *PaymentPeriodController) Create(c *gin.Context) {
	var req dtos.CreatePaymentPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := c.GetUint64("user_id")
	res, err := ctrl.service.Create(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (ctrl *PaymentPeriodController) List(c *gin.Context) {
	res, err := ctrl.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (ctrl *PaymentPeriodController) Get(c *gin.Context) {
	res, err := ctrl.service.Get(c.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment period not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *PaymentPeriodController) Update(c *gin.Context) {
	var req dtos.UpdatePaymentPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := c.GetUint64("user_id")
	if err := ctrl.service.Update(c.Param("id"), req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment period updated successfully"})
}

func (ctrl *PaymentPeriodController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	if err := ctrl.service.Delete(c.Param("id"), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment period deleted successfully"})
}
