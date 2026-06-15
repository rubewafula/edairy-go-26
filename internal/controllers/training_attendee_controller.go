package controllers

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type TrainingAttendeeController struct {
	service *services.TrainingAttendeeService
}

func NewTrainingAttendeeController() *TrainingAttendeeController {
	return &TrainingAttendeeController{
		service: services.NewTrainingAttendeeService(),
	}
}

func (c *TrainingAttendeeController) CreateAttendee(ctx *gin.Context) {
	var req dtos.CreateTrainingAttendeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TrainingAttendeeController.CreateAttendee] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TrainingAttendeeController.CreateAttendee] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	attendee, err := c.service.CreateAttendee(req)
	if err != nil {
		log.Printf("[TrainingAttendeeController.CreateAttendee] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create training attendee"})
		return
	}
	ctx.JSON(http.StatusCreated, attendee)
}

func (c *TrainingAttendeeController) GetAttendees(ctx *gin.Context) {
	attendees, total, err := c.service.GetAttendees()
	if err != nil {
		log.Printf("[TrainingAttendeeController.GetAttendees] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve training attendees"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendees, "total": total})
}

func (c *TrainingAttendeeController) GetAttendee(ctx *gin.Context) {
	attendee, err := c.service.GetAttendee(ctx.Param("id"))
	if err != nil {
		log.Printf("[TrainingAttendeeController.GetAttendee] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Training attendee not found"})
		return
	}
	ctx.JSON(http.StatusOK, attendee)
}

func (c *TrainingAttendeeController) UpdateAttendee(ctx *gin.Context) {
	var req dtos.UpdateTrainingAttendeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TrainingAttendeeController.UpdateAttendee] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TrainingAttendeeController.UpdateAttendee] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAttendee(ctx.Param("id"), req); err != nil {
		log.Printf("[TrainingAttendeeController.UpdateAttendee] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update training attendee"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Attendee updated successfully"})
}

func (c *TrainingAttendeeController) DeleteAttendee(ctx *gin.Context) {
	if err := c.service.DeleteAttendee(ctx.Param("id")); err != nil {
		log.Printf("[TrainingAttendeeController.DeleteAttendee] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete training attendee"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Attendee deleted successfully"})
}
