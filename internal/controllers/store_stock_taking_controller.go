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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	stockTaking, err := c.service.CreateStockTaking(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetStockTaking(utils.Uint64ToString(stockTaking.ID))
	ctx.JSON(http.StatusCreated, response)
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

func (c *StoreStockTakingController) UpdateStockTaking(ctx *gin.Context) {
	var req dtos.UpdateStoreStockTakingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateStockTaking(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Stock taking record updated successfully"})
}

func (c *StoreStockTakingController) DeleteStockTaking(ctx *gin.Context) {
	if err := c.service.DeleteStockTaking(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Stock taking record deleted successfully"})
}
