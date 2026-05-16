package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type DailyMilkVarianceController struct {
	service *services.DailyMilkVarianceService
}

func NewDailyMilkVarianceController() *DailyMilkVarianceController {
	return &DailyMilkVarianceController{
		service: services.NewDailyMilkVarianceService(),
	}
}

func (c *DailyMilkVarianceController) GetDailyVariances(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	variances, total, err := c.service.GetDailyVariances(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": variances, "total": total})
}
