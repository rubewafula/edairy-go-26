package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/controllers"
)

func registerTrainingRoutes(api *gin.RouterGroup) {
	trainingController := controllers.NewTrainingController()
	trainingSessionController := controllers.NewTrainingSessionController()
	trainingAttendeeController := controllers.NewTrainingAttendeeController()
	exchangeVisitController := controllers.NewExchangeVisitController()
	exchangeVisitAttendeeController := controllers.NewExchangeVisitAttendeeController()

	// Training Routes
	api.POST("/trainings", trainingController.CreateTraining)
	api.GET("/trainings", trainingController.GetTrainings)
	api.GET("/trainings/:id", trainingController.GetTraining)
	api.PUT("/trainings/:id", trainingController.UpdateTraining)
	api.DELETE("/trainings/:id", trainingController.DeleteTraining)

	// Training Session Routes
	api.POST("/training-sessions", trainingSessionController.CreateSession)
	api.GET("/training-sessions", trainingSessionController.GetSessions)
	api.GET("/training-sessions/:id", trainingSessionController.GetSession)
	api.PUT("/training-sessions/:id", trainingSessionController.UpdateSession)
	api.DELETE("/training-sessions/:id", trainingSessionController.DeleteSession)

	// Training Session Attendee Routes
	api.POST("/training-session-attendees", trainingAttendeeController.CreateAttendee)
	api.GET("/training-session-attendees", trainingAttendeeController.GetAttendees)
	api.GET("/training-session-attendees/:id", trainingAttendeeController.GetAttendee)
	api.PUT("/training-session-attendees/:id", trainingAttendeeController.UpdateAttendee)
	api.DELETE("/training-session-attendees/:id", trainingAttendeeController.DeleteAttendee)

	// Exchange Visit Routes
	api.POST("/exchange-visits", exchangeVisitController.CreateVisit)
	api.GET("/exchange-visits", exchangeVisitController.GetVisits)
	api.GET("/exchange-visits/:id", exchangeVisitController.GetVisit)
	api.PUT("/exchange-visits/:id", exchangeVisitController.UpdateVisit)
	api.DELETE("/exchange-visits/:id", exchangeVisitController.DeleteVisit)

	// Exchange Visit Attendee Routes
	api.POST("/exchange-visit-attendees", exchangeVisitAttendeeController.CreateAttendee)
	api.GET("/exchange-visit-attendees", exchangeVisitAttendeeController.GetAttendees)
	api.GET("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.GetAttendee)
	api.PUT("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.UpdateAttendee)
	api.DELETE("/exchange-visit-attendees/:id", exchangeVisitAttendeeController.DeleteAttendee)
}
