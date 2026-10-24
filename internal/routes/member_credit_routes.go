package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerMemberCreditRoutes(api *gin.RouterGroup) {
	h := controllers.NewMemberCreditController()
	api.GET("/member-credit-limit", h.GetMemberCreditLimit)
	api.GET("/members/:id/credit-limit", h.GetCreditLimit)
	api.GET("/members/:id/credit-limit-logs", h.ListCreditLimitLogs)
	api.POST("/members/:id/credit-limit/recalculate", h.RecalculateCreditLimit)
	api.POST("/member-credit-limits/recalculate-all", h.RecalculateAll)
	api.POST("/member-credit-limits/reconcile", h.Reconcile)
}
