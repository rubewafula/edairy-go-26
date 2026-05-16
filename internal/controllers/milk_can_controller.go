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

type MilkCanController struct {
	service *services.MilkCanService
}

func NewMilkCanController() *MilkCanController {
	return &MilkCanController{
		service: services.NewMilkCanService(),
	}
}

func (c *MilkCanController) CreateMilkCan(ctx *gin.Context) {
	var req dtos.CreateMilkCanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	milkCan, err := c.service.CreateMilkCan(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetMilkCan(utils.Uint64ToString(milkCan.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *MilkCanController) GetMilkCans(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	milkCans, total, err := c.service.GetMilkCans(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": milkCans, "total": total})
}

func (c *MilkCanController) GetMilkCan(ctx *gin.Context) {
	milkCan, err := c.service.GetMilkCan(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Milk Can not found"})
		return
	}
	ctx.JSON(http.StatusOK, milkCan)
}

func (c *MilkCanController) UpdateMilkCan(ctx *gin.Context) {
	var req dtos.UpdateMilkCanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMilkCan(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Milk Can updated successfully"})
}

func (c *MilkCanController) DeleteMilkCan(ctx *gin.Context) {
	if err := c.service.DeleteMilkCan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Milk Can deleted successfully"})
}
