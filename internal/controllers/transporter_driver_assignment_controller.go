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

type TransporterDriverAssignmentController struct {
	service *services.TransporterDriverAssignmentService
}

func NewTransporterDriverAssignmentController() *TransporterDriverAssignmentController {
	return &TransporterDriverAssignmentController{
		service: services.NewTransporterDriverAssignmentService(),
	}
}

func (c *TransporterDriverAssignmentController) CreateAssignment(ctx *gin.Context) {
	var req dtos.CreateTransporterDriverAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterDriverAssignmentController.CreateAssignment] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterDriverAssignmentController.CreateAssignment] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	assignment, err := c.service.CreateAssignment(req)
	if err != nil {
		log.Printf("[TransporterDriverAssignmentController.CreateAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transporter driver assignment"})
		return
	}
	ctx.JSON(http.StatusCreated, assignment)
}

func (c *TransporterDriverAssignmentController) GetAssignments(ctx *gin.Context) {
	assignments, total, err := c.service.GetAssignments() // Now returns dtos.TransporterDriverAssignmentResponse
	if err != nil {
		log.Printf("[TransporterDriverAssignmentController.GetAssignments] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transporter driver assignments"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": assignments, "total": total})
}

func (c *TransporterDriverAssignmentController) GetAssignment(ctx *gin.Context) {
	assignment, err := c.service.GetAssignment(ctx.Param("id")) // Now returns dtos.TransporterDriverAssignmentResponse
	if err != nil {
		log.Printf("[TransporterDriverAssignmentController.GetAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter driver assignment not found"})
		return
	}
	ctx.JSON(http.StatusOK, assignment)
}

func (c *TransporterDriverAssignmentController) UpdateAssignment(ctx *gin.Context) {
	var req dtos.UpdateTransporterDriverAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterDriverAssignmentController.UpdateAssignment] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterDriverAssignmentController.UpdateAssignment] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAssignment(ctx.Param("id"), req); err != nil {
		log.Printf("[TransporterDriverAssignmentController.UpdateAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transporter driver assignment"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Driver assignment updated successfully"})
}

func (c *TransporterDriverAssignmentController) DeleteAssignment(ctx *gin.Context) {
	if err := c.service.DeleteAssignment(ctx.Param("id")); err != nil {
		log.Printf("[TransporterDriverAssignmentController.DeleteAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transporter driver assignment"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Assignment deleted successfully"})
}
