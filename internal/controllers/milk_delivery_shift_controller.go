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

type MilkDeliveryShiftController struct {
	service *services.MilkDeliveryShiftService
}

func NewMilkDeliveryShiftController() *MilkDeliveryShiftController {
	return &MilkDeliveryShiftController{
		service: services.NewMilkDeliveryShiftService(),
	}
}

func (c *MilkDeliveryShiftController) CreateShift(ctx *gin.Context) {
	var req dtos.CreateMilkDeliveryShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	shift, err := c.service.CreateShift(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, shift)
}

func (c *MilkDeliveryShiftController) GetShifts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	shifts, total, err := c.service.GetShifts(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": shifts, "Total": total})
}

func (c *MilkDeliveryShiftController) GetShift(ctx *gin.Context) {
	shift, err := c.service.GetShift(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Shift not found"})
		return
	}
	ctx.JSON(http.StatusOK, shift)
}

func (c *MilkDeliveryShiftController) UpdateShift(ctx *gin.Context) {
	var req dtos.UpdateMilkDeliveryShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShift(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Shift updated successfully"})
}

func (c *MilkDeliveryShiftController) DeleteShift(ctx *gin.Context) {
	if err := c.service.DeleteShift(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Shift deleted successfully"})
}
