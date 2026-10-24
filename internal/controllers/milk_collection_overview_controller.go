package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type MilkCollectionOverviewController struct {
	Service *services.MilkCollectionOverviewService
}

func NewMilkCollectionOverviewController() *MilkCollectionOverviewController {
	return &MilkCollectionOverviewController{
		Service: services.NewMilkCollectionOverviewService(),
	}
}

// GetOverview
// @Summary Member milk collection overview (mobile dashboard)
// @Description Daily litres collected from members over the last N days (today inclusive), with comparison to the previous N days.
// @Tags Produce Dashboard
// @Security BearerAuth
// @Produce json
// @Param days query int false "Number of days (default 7, max 90)"
// @Success 200 {object} dtos.MilkCollectionOverviewResponse
// @Router /api/milk-collection-overview [get]
func (h *MilkCollectionOverviewController) GetOverview(c *gin.Context) {
	days := queryIntParam(c, services.MilkCollectionOverviewDefaultDays, "days", "Days")
	scope := services.ResolveMilkOverviewScope(c.GetUint64("user_id"))
	c.JSON(http.StatusOK, h.Service.GetOverview(days, scope))
}
