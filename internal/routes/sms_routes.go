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
	api.GET("/sms-contacts", smsController.GetContacts)
	api.GET("/sms-contacts/:id", smsController.GetContact)
	api.PUT("/sms-contacts/:id", smsController.UpdateContact)
	api.DELETE("/sms-contacts/:id", smsController.DeleteContact)
	api.GET("/sms-groups/:id/contacts", smsController.GetContactsByGroup)

	api.POST("/sms-send", smsController.SendMessage)
	api.GET("/sms-queue", smsController.GetQueue)
	api.GET("/sms-providers", smsController.GetProviders)
	api.POST("/sms-providers", smsController.CreateProvider)
	api.POST("/sms-templates", smsController.CreateTemplate)
	api.GET("/sms-templates", smsController.GetTemplates)

	// SMS Message Routes
	api.GET("/sms-messages", smsController.GetMessages)
	api.GET("/sms-messages/:id", smsController.GetMessage)
	api.POST("/sms-messages", smsController.CreateMessage)
	api.PUT("/sms-messages/:id", smsController.UpdateMessage)
	api.DELETE("/sms-messages/:id", smsController.DeleteMessage)

	// SMS Campaign Routes
	api.POST("/sms-campaigns", smsCampaignController.CreateCampaign)
	api.GET("/sms-campaigns", smsCampaignController.GetCampaigns)
	api.GET("/sms-campaigns/:id", smsCampaignController.GetCampaign)
	api.PUT("/sms-campaigns/:id", smsCampaignController.UpdateCampaign)
	api.DELETE("/sms-campaigns/:id", smsCampaignController.DeleteCampaign)
	api.GET("/sms-campaigns/:id/recipients", smsCampaignController.GetRecipientsByCampaign)

	// SMS Campaign Recipient Routes
	api.POST("/sms-campaign-recipients", smsCampaignController.CreateRecipient)
	api.GET("/sms-campaign-recipients", smsCampaignController.GetAllRecipients)
	api.GET("/sms-campaign-recipients/:id", smsCampaignController.GetRecipient)
	api.PUT("/sms-campaign-recipients/:id", smsCampaignController.UpdateRecipient)
	api.DELETE("/sms-campaign-recipients/:id", smsCampaignController.DeleteRecipient)
}
