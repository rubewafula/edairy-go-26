package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type ProduceDashboardController struct {
	Service *services.ProduceDashboardService
}

func NewProduceDashboardController() *ProduceDashboardController {
	return &ProduceDashboardController{
		Service: services.NewProduceDashboardService(),
	}
}

// GetDashboard
// @Summary Get produce dashboard statistics
// @Tags Produce Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.ProduceDashboardResponse
// @Router /api/produce-dashboard [get]
func (h *ProduceDashboardController) GetDashboard(c *gin.Context) {
	data := h.Service.GetDashboard()
	c.JSON(http.StatusOK, data)
}
