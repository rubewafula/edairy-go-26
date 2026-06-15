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

type CashTransactionController struct {
	service *services.CashTransactionService
}

func NewCashTransactionController() *CashTransactionController {
	return &CashTransactionController{
		service: services.NewCashTransactionService(),
	}
}

func (c *CashTransactionController) Create(ctx *gin.Context) {
	var req dtos.CreateCashTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[CashTransactionController.Create] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[CashTransactionController.Create] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.Create(req, userID)
	if err != nil {
		log.Printf("[CashTransactionController.Create] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cash transaction"})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *CashTransactionController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.List(page, limit)
	if err != nil {
		log.Printf("[CashTransactionController.List] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cash transactions"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *CashTransactionController) Get(ctx *gin.Context) {
	res, err := c.service.Get(ctx.Param("id"))
	if err != nil {
		log.Printf("[CashTransactionController.Get] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cash transaction not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *CashTransactionController) Update(ctx *gin.Context) {
	var req dtos.CreateCashTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[CashTransactionController.Update] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[CashTransactionController.Update] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.Update(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[CashTransactionController.Update] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cash transaction"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cash transaction updated successfully"})
}

func (c *CashTransactionController) Delete(ctx *gin.Context) {
	if err := c.service.Delete(ctx.Param("id")); err != nil {
		log.Printf("[CashTransactionController.Delete] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cash transaction"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cash transaction deleted successfully"})
}
