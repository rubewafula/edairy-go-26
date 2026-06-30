package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type PayrollDashboardController struct {
	Service *services.PayrollDashboardService
}

func NewPayrollDashboardController() *PayrollDashboardController {
	return &PayrollDashboardController{Service: services.NewPayrollDashboardService()}
}

// GetDashboard
// @Summary Get payroll dashboard statistics
// @Tags Payroll Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.PayrollDashboardResponse
// @Router /api/payroll-dashboard [get]
func (h *PayrollDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
