package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type StoreDashboardController struct {
	Service *services.StoreDashboardService
}

func NewStoreDashboardController() *StoreDashboardController {
	return &StoreDashboardController{
		Service: services.NewStoreDashboardService(),
	}
}

// GetDashboard
// @Summary Get store dashboard statistics
// @Tags Store Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.StoreDashboardResponse
// @Router /api/store-dashboard [get]
func (h *StoreDashboardController) GetDashboard(c *gin.Context) {
	data := h.Service.GetDashboard()
	c.JSON(http.StatusOK, data)
}
