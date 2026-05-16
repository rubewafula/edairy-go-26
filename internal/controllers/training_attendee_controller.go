package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	attendee, err := c.service.CreateAttendee(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, attendee)
}

func (c *TrainingAttendeeController) GetAttendees(ctx *gin.Context) {
	attendees, total, err := c.service.GetAttendees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendees, "total": total})
}

func (c *TrainingAttendeeController) GetAttendee(ctx *gin.Context) {
	attendee, err := c.service.GetAttendee(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Attendee not found"})
		return
	}
	ctx.JSON(http.StatusOK, attendee)
}

func (c *TrainingAttendeeController) UpdateAttendee(ctx *gin.Context) {
	var req dtos.UpdateTrainingAttendeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAttendee(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Attendee updated successfully"})
}

func (c *TrainingAttendeeController) DeleteAttendee(ctx *gin.Context) {
	if err := c.service.DeleteAttendee(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Attendee deleted successfully"})
}
