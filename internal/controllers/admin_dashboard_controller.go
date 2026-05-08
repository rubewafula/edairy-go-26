package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type AdminDashboardController struct {
	Service *services.AdminDashboardService
}

func NewAdminDashboardController() *AdminDashboardController {
	return &AdminDashboardController{
		Service: services.NewAdminDashboardService(),
	}
}

// GetDashboard
// @Summary Get dashboard statistics
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.DashboardResponse
// @Router /api/dashboard [get]
func (h *AdminDashboardController) GetDashboard(c *gin.Context) {

	data := h.Service.GetDashboard()

	c.JSON(http.StatusOK, data)

}
