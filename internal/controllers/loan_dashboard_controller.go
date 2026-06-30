package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type LoanDashboardController struct {
	Service *services.LoanDashboardService
}

func NewLoanDashboardController() *LoanDashboardController {
	return &LoanDashboardController{Service: services.NewLoanDashboardService()}
}

// GetDashboard
// @Summary Get loan dashboard statistics
// @Tags Loan Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.LoanDashboardResponse
// @Router /api/loan-dashboard [get]
func (h *LoanDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
