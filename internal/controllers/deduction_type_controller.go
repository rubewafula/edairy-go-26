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

type DeductionTypeController struct {
	service *services.DeductionTypeService
}

func NewDeductionTypeController() *DeductionTypeController {
	return &DeductionTypeController{
		service: services.NewDeductionTypeService(),
	}
}

func (c *DeductionTypeController) CreateDeductionType(ctx *gin.Context) {
	var req dtos.CreateDeductionTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	deductionType, err := c.service.CreateDeductionType(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, deductionType)
}

func (c *DeductionTypeController) GetDeductionTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetDeductionTypes(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": results, "Total": total})
}

func (c *DeductionTypeController) GetDeductionType(ctx *gin.Context) {
	deductionType, err := c.service.GetDeductionType(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Deduction type not found"})
		return
	}
	ctx.JSON(http.StatusOK, deductionType)
}

func (c *DeductionTypeController) UpdateDeductionType(ctx *gin.Context) {
	var req dtos.UpdateDeductionTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateDeductionType(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Deduction type updated successfully"})
}

func (c *DeductionTypeController) DeleteDeductionType(ctx *gin.Context) {
	if err := c.service.DeleteDeductionType(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Deduction type deleted successfully"})
}
