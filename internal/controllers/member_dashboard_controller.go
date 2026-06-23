package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type MemberDashboardController struct {
	Service *services.MemberDashboardService
}

func NewMemberDashboardController() *MemberDashboardController {
	return &MemberDashboardController{
		Service: services.NewMemberDashboardService(),
	}
}

// GetDashboard
// @Summary Get member dashboard statistics
// @Tags Member Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dtos.MemberDashboardResponse
// @Router /api/member-dashboard [get]
func (h *MemberDashboardController) GetDashboard(c *gin.Context) {
	data := h.Service.GetDashboard()
	c.JSON(http.StatusOK, data)
}
