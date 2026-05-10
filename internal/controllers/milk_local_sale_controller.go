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

type MilkLocalSaleController struct {
	service *services.MilkLocalSaleService
}

func NewMilkLocalSaleController() *MilkLocalSaleController {
	return &MilkLocalSaleController{
		service: services.NewMilkLocalSaleService(),
	}
}

func (c *MilkLocalSaleController) CreateMilkLocalSale(ctx *gin.Context) {
	var req dtos.CreateMilkLocalSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	sale, err := c.service.CreateMilkLocalSale(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetMilkLocalSale(utils.Uint64ToString(sale.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *MilkLocalSaleController) GetMilkLocalSales(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	sales, total, err := c.service.GetMilkLocalSales(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": sales, "Total": total})
}

func (c *MilkLocalSaleController) GetMilkLocalSale(ctx *gin.Context) {
	sale, err := c.service.GetMilkLocalSale(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Local sale not found"})
		return
	}
	ctx.JSON(http.StatusOK, sale)
}

func (c *MilkLocalSaleController) UpdateMilkLocalSale(ctx *gin.Context) {
	var req dtos.UpdateMilkLocalSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMilkLocalSale(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Local sale updated successfully"})
}

func (c *MilkLocalSaleController) DeleteMilkLocalSale(ctx *gin.Context) {
	if err := c.service.DeleteMilkLocalSale(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Local sale deleted successfully"})
}
