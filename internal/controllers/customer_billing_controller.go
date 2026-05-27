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
	// No request body needed as it will automatically find the next unprocessed pay date range
	userID := ctx.GetUint64("user_id")

	// Proceed to generate all pending billings up to now
	if err := c.service.CreateBilling(userID); err != nil {
		log.Printf("[CustomerBillingController.CreateBilling] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Customer billing generation started in the background"})
}

// New controller methods for Confirm and Approve
func (c *CustomerBillingController) ConfirmBilling(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid billing ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ConfirmBilling(id, userID); err != nil {
		log.Printf("[CustomerBillingController.ConfirmBilling] Error confirming billing %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer billing confirmed successfully"})
}

func (c *CustomerBillingController) ApproveBilling(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid billing ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ApproveBilling(id, userID); err != nil {
		log.Printf("[CustomerBillingController.ApproveBilling] Error approving billing %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer billing approved successfully"})
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
