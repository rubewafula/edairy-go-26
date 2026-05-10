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

type StoreSaleController struct {
	service *services.StoreSaleService
}

func NewStoreSaleController() *StoreSaleController {
	return &StoreSaleController{
		service: services.NewStoreSaleService(),
	}
}

func (c *StoreSaleController) CreateSale(ctx *gin.Context) {
	var req dtos.CreateStoreSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	sale, err := c.service.CreateSale(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetSale(utils.Uint64ToString(sale.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *StoreSaleController) GetSales(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetSales(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": results, "Total": total})
}

func (c *StoreSaleController) GetSale(ctx *gin.Context) {
	result, err := c.service.GetSale(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Store sale not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *StoreSaleController) UpdateSale(ctx *gin.Context) {
	var req dtos.UpdateStoreSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.UpdateSale(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Store sale updated successfully"})
}

func (c *StoreSaleController) DeleteSale(ctx *gin.Context) {
	if err := c.service.DeleteSale(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Store sale deleted successfully"})
}
