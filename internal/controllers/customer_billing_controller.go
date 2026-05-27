package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type CustomerBillingController struct {
	service *services.CustomerBillingService
}

func NewCustomerBillingController() *CustomerBillingController {
	return &CustomerBillingController{
		service: services.NewCustomerBillingService(),
	}
}

func (c *CustomerBillingController) CreateBilling(ctx *gin.Context) {
	var req struct {
		PayDateRangeID uint64 `json:"pay_date_range_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.CreateBilling(req.PayDateRangeID, userID); err != nil {
		log.Printf("[CustomerBillingController.CreateBilling] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Customer billings generated successfully"})
}

func (c *CustomerBillingController) GetBillings(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetBillings(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *CustomerBillingController) GetBilling(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetBilling(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer billing not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *CustomerBillingController) GetBillingItems(ctx *gin.Context) {
	id := ctx.Param("id")
	items, err := c.service.GetBillingItems(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

func (c *CustomerBillingController) DeleteBilling(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteBilling(id); err != nil {
		log.Printf("[CustomerBillingController.DeleteBilling] Error deleting billing %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer billing deleted successfully"})
}
