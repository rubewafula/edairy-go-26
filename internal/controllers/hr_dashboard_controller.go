package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type HRDashboardController struct {
	Service *services.HRDashboardService
}

func NewHRDashboardController() *HRDashboardController {
	return &HRDashboardController{Service: services.NewHRDashboardService()}
}

// GetDashboard
// @Summary Get human resources dashboard statistics
// @Tags HR Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.HRDashboardResponse
// @Router /api/hr-dashboard [get]
func (h *HRDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
