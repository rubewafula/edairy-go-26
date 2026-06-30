package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type SupplierDashboardController struct {
	Service *services.SupplierDashboardService
}

func NewSupplierDashboardController() *SupplierDashboardController {
	return &SupplierDashboardController{Service: services.NewSupplierDashboardService()}
}

// GetDashboard
// @Summary Get supplier dashboard statistics
// @Tags Supplier Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.SupplierDashboardResponse
// @Router /api/supplier-dashboard [get]
func (h *SupplierDashboardController) GetDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, h.Service.GetDashboard())
}
