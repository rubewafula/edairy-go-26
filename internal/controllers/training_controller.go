package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type TrainingController struct {
	service *services.TrainingService
}

func NewTrainingController() *TrainingController {
	return &TrainingController{
		service: services.NewTrainingService(),
	}
}

func (c *TrainingController) CreateTraining(ctx *gin.Context) {
	var req dtos.CreateTrainingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	training, err := c.service.CreateTraining(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, training)
}

func (c *TrainingController) GetTrainings(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	venue := ctx.Query("venue")
	topic := ctx.Query("topic")
	facilitator := ctx.Query("facilitator")

	trainings, total, err := c.service.GetTrainings(page, limit, venue, topic, facilitator)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": trainings, "total": total})
}

func (c *TrainingController) GetTraining(ctx *gin.Context) {
	training, err := c.service.GetTraining(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Training not found"})
		return
	}
	ctx.JSON(http.StatusOK, training)
}

func (c *TrainingController) UpdateTraining(ctx *gin.Context) {
	var req dtos.UpdateTrainingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateTraining(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Training updated successfully"})
}

func (c *TrainingController) DeleteTraining(ctx *gin.Context) {
	if err := c.service.DeleteTraining(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Training deleted successfully"})
}
