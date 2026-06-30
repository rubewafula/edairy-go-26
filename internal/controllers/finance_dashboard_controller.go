package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type FinanceDashboardController struct {
	Service *services.FinanceDashboardService
}

func NewFinanceDashboardController() *FinanceDashboardController {
	return &FinanceDashboardController{Service: services.NewFinanceDashboardService()}
}

// GetDashboard
// @Summary Get finance dashboard statistics
// @Tags Finance Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.FinanceDashboardResponse
// @Router /api/finance-dashboard [get]
func (h *FinanceDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
