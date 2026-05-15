package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type CustomerBillingController struct {
	service *services.CustomerBillingService
}

func NewCustomerBillingController() *CustomerBillingController {
	return &CustomerBillingController{
		service: services.NewCustomerBillingService(),
	}
}

func (c *CustomerBillingController) GetBillings(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetBillings(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *CustomerBillingController) GetBilling(ctx *gin.Context) {
	result, err := c.service.GetBilling(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Billing not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *CustomerBillingController) GetBillingItems(ctx *gin.Context) {
	items, err := c.service.GetBillingItems(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}
