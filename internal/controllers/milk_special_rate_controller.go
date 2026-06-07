package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type MilkSpecialRateController struct {
	service *services.MilkSpecialRateService
}

func NewMilkSpecialRateController() *MilkSpecialRateController {
	return &MilkSpecialRateController{
		service: services.NewMilkSpecialRateService(),
	}
}

func (c *MilkSpecialRateController) Create(ctx *gin.Context) {
	var req dtos.CreateMilkSpecialRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.Create(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *MilkSpecialRateController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	payDateRangeID := ctx.Query("pay_date_range_id")

	res, total, err := c.service.List(payDateRangeID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *MilkSpecialRateController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Rate not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *MilkSpecialRateController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateMilkSpecialRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.Update(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Rate updated successfully"})
}

func (c *MilkSpecialRateController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Rate deleted successfully"})
}
