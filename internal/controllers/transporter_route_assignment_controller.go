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

type TransporterRouteAssignmentController struct {
	service *services.TransporterRouteAssignmentService
}

func NewTransporterRouteAssignmentController() *TransporterRouteAssignmentController {
	return &TransporterRouteAssignmentController{
		service: services.NewTransporterRouteAssignmentService(),
	}
}

func (c *TransporterRouteAssignmentController) CreateAssignment(ctx *gin.Context) {
	var req dtos.CreateTransporterRouteAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterRouteAssignmentController.CreateAssignment] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterRouteAssignmentController.CreateAssignment] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	assignment, err := c.service.CreateAssignment(req)
	if err != nil {
		log.Printf("[TransporterRouteAssignmentController.CreateAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transporter route assignment"})
		return
	}
	ctx.JSON(http.StatusCreated, assignment)
}

func (c *TransporterRouteAssignmentController) GetAssignments(ctx *gin.Context) {
	assignments, total, err := c.service.GetAssignments() // Now returns dtos.TransporterRouteAssignmentResponse
	if err != nil {
		log.Printf("[TransporterRouteAssignmentController.GetAssignments] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transporter route assignments"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": assignments, "total": total})
}

func (c *TransporterRouteAssignmentController) GetAssignment(ctx *gin.Context) {
	assignment, err := c.service.GetAssignment(ctx.Param("id")) // Now returns dtos.TransporterRouteAssignmentResponse
	if err != nil {
		log.Printf("[TransporterRouteAssignmentController.GetAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter route assignment not found"})
		return
	}
	ctx.JSON(http.StatusOK, assignment)
}

func (c *TransporterRouteAssignmentController) UpdateAssignment(ctx *gin.Context) {
	var req dtos.UpdateTransporterRouteAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterRouteAssignmentController.UpdateAssignment] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterRouteAssignmentController.UpdateAssignment] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAssignment(ctx.Param("id"), req); err != nil {
		log.Printf("[TransporterRouteAssignmentController.UpdateAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transporter route assignment"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Route assignment updated successfully"})
}

func (c *TransporterRouteAssignmentController) DeleteAssignment(ctx *gin.Context) {
	if err := c.service.DeleteAssignment(ctx.Param("id")); err != nil {
		log.Printf("[TransporterRouteAssignmentController.DeleteAssignment] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transporter route assignment"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Assignment deleted successfully"})
}
