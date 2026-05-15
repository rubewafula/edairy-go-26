package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerSMSRoutes(api *gin.RouterGroup) {
	smsController := controllers.NewSMSController()
	smsCampaignController := controllers.NewSMSCampaignController()

	// SMS Routes
	api.POST("/sms-groups", smsController.CreateGroup)
	api.GET("/sms-groups", smsController.GetGroups)
	api.POST("/sms-contacts", smsController.CreateContact)
	api.POST("/sms-send", smsController.SendMessage)
	api.GET("/sms-queue", smsController.GetQueue)
	api.GET("/sms-providers", smsController.GetProviders)
	api.POST("/sms-providers", smsController.CreateProvider)
	api.POST("/sms-templates", smsController.CreateTemplate)
	api.GET("/sms-templates", smsController.GetTemplates)

	// SMS Campaign Routes
	api.POST("/sms-campaigns", smsCampaignController.CreateCampaign)
	api.GET("/sms-campaigns", smsCampaignController.GetCampaigns)
	api.GET("/sms-campaigns/:id/recipients", smsCampaignController.GetRecipients)
}
