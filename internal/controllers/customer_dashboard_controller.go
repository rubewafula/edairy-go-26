package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type CustomerDashboardController struct {
	Service *services.CustomerDashboardService
}

func NewCustomerDashboardController() *CustomerDashboardController {
	return &CustomerDashboardController{Service: services.NewCustomerDashboardService()}
}

// GetDashboard
// @Summary Get customer dashboard statistics
// @Tags Customer Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.CustomerDashboardResponse
// @Router /api/customer-dashboard [get]
func (h *CustomerDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
