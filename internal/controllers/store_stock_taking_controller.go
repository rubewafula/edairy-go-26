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

type StoreStockTakingController struct {
	service *services.StoreStockTakingService
}

func NewStoreStockTakingController() *StoreStockTakingController {
	return &StoreStockTakingController{
		service: services.NewStoreStockTakingService(),
	}
}

func (c *StoreStockTakingController) CreateStockTaking(ctx *gin.Context) {
	var req dtos.CreateStoreStockTakingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("StockTaking: Could not bind json: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("StockTaking: Could not validate json: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	items, err := c.service.CreateStockTaking(req, userID)
	if err != nil {
		log.Printf("StockTaking: Could not create stock taking: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully recorded " + strconv.Itoa(len(items)) + " stock taking entries", "stock_take_no": req.StockTakeNo})
}

func (c *StoreStockTakingController) GetStockTakings(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetStockTakings(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *StoreStockTakingController) GetStockTaking(ctx *gin.Context) {
	result, err := c.service.GetStockTaking(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Stock taking record not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
