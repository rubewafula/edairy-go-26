package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type TransporterDashboardController struct {
	Service *services.TransporterDashboardService
}

func NewTransporterDashboardController() *TransporterDashboardController {
	return &TransporterDashboardController{
		Service: services.NewTransporterDashboardService(),
	}
}

// GetDashboard
// @Summary Get transporter dashboard statistics
// @Tags Transporter Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.TransporterDashboardResponse
// @Router /api/transporter-dashboard [get]
func (h *TransporterDashboardController) GetDashboard(c *gin.Context) {
	data := h.Service.GetDashboard()
	c.JSON(http.StatusOK, data)
}
