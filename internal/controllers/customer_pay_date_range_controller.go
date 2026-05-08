package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CustomerPayDateRangeController struct {
	service *services.CustomerPayDateRangeService
}

func NewCustomerPayDateRangeController() *CustomerPayDateRangeController {
	return &CustomerPayDateRangeController{
		service: services.NewCustomerPayDateRangeService(),
	}
}

func (c *CustomerPayDateRangeController) CreateRange(ctx *gin.Context) {
	var req dtos.CreateCustomerPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	dateRange, err := c.service.CreateCustomerPayDateRange(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dateRange)
}

func (c *CustomerPayDateRangeController) GetRanges(ctx *gin.Context) {
	ranges, total, err := c.service.GetCustomerPayDateRanges()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": ranges, "total": total})
}

func (c *CustomerPayDateRangeController) GetRange(ctx *gin.Context) {
	dateRange, err := c.service.GetCustomerPayDateRange(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Range not found"})
		return
	}
	ctx.JSON(http.StatusOK, dateRange)
}

func (c *CustomerPayDateRangeController) UpdateRange(ctx *gin.Context) {
	var req dtos.UpdateCustomerPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCustomerPayDateRange(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Range updated successfully"})
}

func (c *CustomerPayDateRangeController) DeleteRange(ctx *gin.Context) {
	if err := c.service.DeleteCustomerPayDateRange(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Range deleted successfully"})
}
