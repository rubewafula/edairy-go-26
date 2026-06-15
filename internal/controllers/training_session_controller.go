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

type TrainingSessionController struct {
	service *services.TrainingSessionService
}

func NewTrainingSessionController() *TrainingSessionController {
	return &TrainingSessionController{
		service: services.NewTrainingSessionService(),
	}
}

func (c *TrainingSessionController) CreateSession(ctx *gin.Context) {
	var req dtos.CreateTrainingSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TrainingSessionController.CreateSession] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TrainingSessionController.CreateSession] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	session, err := c.service.CreateSession(req)
	if err != nil {
		log.Printf("[TrainingSessionController.CreateSession] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create training session"})
		return
	}
	ctx.JSON(http.StatusCreated, session)
}

func (c *TrainingSessionController) GetSessions(ctx *gin.Context) {
	sessions, total, err := c.service.GetSessions()
	if err != nil {
		log.Printf("[TrainingSessionController.GetSessions] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve training sessions"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": sessions, "total": total})
}

func (c *TrainingSessionController) GetSession(ctx *gin.Context) {
	session, err := c.service.GetSession(ctx.Param("id"))
	if err != nil {
		log.Printf("[TrainingSessionController.GetSession] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Training session not found"})
		return
	}
	ctx.JSON(http.StatusOK, session)
}

func (c *TrainingSessionController) UpdateSession(ctx *gin.Context) {
	var req dtos.UpdateTrainingSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TrainingSessionController.UpdateSession] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TrainingSessionController.UpdateSession] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateSession(ctx.Param("id"), req); err != nil {
		log.Printf("[TrainingSessionController.UpdateSession] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update training session"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Training session updated successfully"})
}

func (c *TrainingSessionController) DeleteSession(ctx *gin.Context) {
	if err := c.service.DeleteSession(ctx.Param("id")); err != nil {
		log.Printf("[TrainingSessionController.DeleteSession] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete training session"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Training session deleted successfully"})
}
