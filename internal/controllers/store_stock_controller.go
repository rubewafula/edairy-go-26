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

type StoreStockController struct {
	service *services.StoreStockService
}

func NewStoreStockController() *StoreStockController {
	return &StoreStockController{
		service: services.NewStoreStockService(),
	}
}

func (c *StoreStockController) CreateStock(ctx *gin.Context) {
	var req dtos.CreateStoreStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreStockController.CreateStock] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreStockController.CreateStock] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	stock, err := c.service.CreateStock(req)
	if err != nil {
		log.Printf("[StoreStockController.CreateStock] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store stock"})
		return
	}

	response, _ := c.service.GetStock(utils.Uint64ToString(stock.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *StoreStockController) GetStocks(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetStocks(page, limit)
	if err != nil {
		log.Printf("[StoreStockController.GetStocks] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve store stocks"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreStockController) GetStock(ctx *gin.Context) {
	stock, err := c.service.GetStock(ctx.Param("id"))
	if err != nil {
		log.Printf("[StoreStockController.GetStock] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Store stock entry not found"})
		return
	}
	ctx.JSON(http.StatusOK, stock)
}

func (c *StoreStockController) UpdateStock(ctx *gin.Context) {
	var req dtos.UpdateStoreStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreStockController.UpdateStock] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreStockController.UpdateStock] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateStock(ctx.Param("id"), req); err != nil {
		log.Printf("[StoreStockController.UpdateStock] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store stock"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Stock updated successfully"})
}

func (c *StoreStockController) DeleteStock(ctx *gin.Context) {
	if err := c.service.DeleteStock(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Stock deleted successfully"})
}
