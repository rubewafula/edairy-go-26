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
	api.GET("/sms-groups/:id", smsController.GetGroup)
	api.PUT("/sms-groups/:id", smsController.UpdateGroup)
	api.DELETE("/sms-groups/:id", smsController.DeleteGroup)
	api.GET("/sms-groups/:id/contacts", smsController.GetContactsByGroup)

	api.POST("/sms-contacts", smsController.CreateContact)
	api.GET("/sms-contacts", smsController.GetContacts)
	api.GET("/sms-contacts/:id", smsController.GetContact)
	api.PUT("/sms-contacts/:id", smsController.UpdateContact)
	api.DELETE("/sms-contacts/:id", smsController.DeleteContact)

	//sms providers
	api.GET("/sms-providers", smsController.GetProviders)
	api.GET("/sms-providers/:id", smsController.GetProvider)
	api.PUT("/sms-providers/:id", smsController.UpdateProvider)
	api.POST("/sms-providers", smsController.CreateProvider)
	api.DELETE("/sms-providers/:id", smsController.DeleteProvider)

	api.POST("/sms-templates", smsController.CreateTemplate)
	api.GET("/sms-templates", smsController.GetTemplates)
	api.PUT("/sms-templates/:id", smsController.UpdateTemplate)
	api.GET("/sms-templates/:id", smsController.GetTemplate)
	api.DELETE("/sms-templates/:id", smsController.DeleteTemplate)

	// SMS Message Routes
	api.GET("/sms-messages", smsController.GetMessages)
	api.GET("/sms-messages/:id", smsController.GetMessage)
	api.POST("/sms-messages", smsController.CreateMessage)
	api.PUT("/sms-messages/:id", smsController.UpdateMessage)
	api.DELETE("/sms-messages/:id", smsController.DeleteMessage)

	//sms providers
	api.GET("/sms-outboxes", smsController.GetAllOutboxes)
	api.GET("/sms-outboxes/:id", smsController.GetOutbox)
	api.PUT("/sms-outboxes/:id", smsController.UpdateOutbox)
	api.POST("/sms-outboxes", smsController.CreateOutbox)
	api.DELETE("/sms-outboxes/:id", smsController.DeleteOutbox)

	// SMS Campaign Routes
	api.POST("/sms-campaigns", smsCampaignController.CreateCampaign)
	api.GET("/sms-campaigns", smsCampaignController.GetCampaigns)
	api.GET("/sms-campaigns/:id", smsCampaignController.GetCampaign)
	api.PUT("/sms-campaigns/:id", smsCampaignController.UpdateCampaign)
	api.DELETE("/sms-campaigns/:id", smsCampaignController.DeleteCampaign)

	// In-App Configuration Routes
	api.GET("/sms-inapp-configurations", smsController.GetInAppConfigs)
	api.GET("/sms-inapp-configurations/:id", smsController.GetInAppConfig)
	api.POST("/sms-inapp-configurations", smsController.CreateInAppConfig)
	api.PUT("/sms-inapp-configurations/:id", smsController.UpdateInAppConfig)
	api.DELETE("/sms-inapp-configurations/:id", smsController.DeleteInAppConfig)

}
