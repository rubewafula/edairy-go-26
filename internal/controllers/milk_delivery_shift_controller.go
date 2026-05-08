package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
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
	shifts, total, err := c.service.GetShifts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": shifts, "total": total})
}

func (c *MilkDeliveryShiftController) GetShift(ctx *gin.Context) {
	shift, err := c.service.GetShift(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Shift not found"})
		return
	}
	ctx.JSON(http.StatusOK, shift)
}

func (c *MilkDeliveryShiftController) UpdateShift(ctx *gin.Context) {
	var req dtos.UpdateMilkDeliveryShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShift(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Shift updated successfully"})
}

func (c *MilkDeliveryShiftController) DeleteShift(ctx *gin.Context) {
	if err := c.service.DeleteShift(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Shift deleted successfully"})
}
