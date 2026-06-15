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
	"gorm.io/gorm"
)

type SupplierQuoteController struct {
	service *services.SupplierQuoteService
}

func NewSupplierQuoteController() *SupplierQuoteController {
	return &SupplierQuoteController{
		service: services.NewSupplierQuoteService(),
	}
}

func (c *SupplierQuoteController) CreateQuote(ctx *gin.Context) {
	var req dtos.CreateSupplierQuoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierQuoteController.CreateQuote] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")
	quote, err := c.service.CreateQuote(req, userID)
	if err != nil {
		log.Printf("[SupplierQuoteController.CreateQuote] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier quote"})
		return
	}
	ctx.JSON(http.StatusCreated, quote)
}

func (c *SupplierQuoteController) GetQuotes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetQuotes(page, limit)
	if err != nil {
		log.Printf("[SupplierQuoteController.GetQuotes] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier quotes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplierQuoteController) CreateQuoteItem(ctx *gin.Context) {
	var req dtos.CreateSupplierQuoteItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierQuoteController.CreateQuoteItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierQuoteController.CreateQuoteItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	item, err := c.service.CreateQuoteItem(req, userID)
	if err != nil {
		log.Printf("[SupplierQuoteController.CreateQuoteItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier quote item"})
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *SupplierQuoteController) GetQuoteItems(ctx *gin.Context) {
	items, err := c.service.GetQuoteItems(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplierQuoteController.GetQuoteItems] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier quote items"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

func (c *SupplierQuoteController) GetQuoteItem(ctx *gin.Context) {
	item, err := c.service.GetQuoteItem(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier quote item not found"})
			return
		}
		log.Printf("[SupplierQuoteController.GetQuoteItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier quote item"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *SupplierQuoteController) UpdateQuoteItem(ctx *gin.Context) {
	var req dtos.UpdateSupplierQuoteItemRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierQuoteController.UpdateQuoteItem] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierQuoteController.UpdateQuoteItem] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	err := c.service.UpdateQuoteItem(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier quote item not found"})
			return
		}
		log.Printf("[SupplierQuoteController.UpdateQuoteItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier quote item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier quote item updated successfully"})
}

func (c *SupplierQuoteController) DeleteQuoteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteQuoteItem(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier quote item not found"})
			return
		}
		log.Printf("[SupplierQuoteController.DeleteQuoteItem] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier quote item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier quote item deleted successfully"})
}
