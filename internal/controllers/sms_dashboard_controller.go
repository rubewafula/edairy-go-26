package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type SMSDashboardController struct {
	Service *services.SMSDashboardService
}

func NewSMSDashboardController() *SMSDashboardController {
	return &SMSDashboardController{Service: services.NewSMSDashboardService()}
}

// GetDashboard
// @Summary Get SMS dashboard statistics
// @Tags SMS Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.SMSDashboardResponse
// @Router /api/sms-dashboard [get]
func (h *SMSDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
