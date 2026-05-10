package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	assignment, err := c.service.CreateAssignment(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, assignment)
}

func (c *TransporterDriverAssignmentController) GetAssignments(ctx *gin.Context) {
	assignments, total, err := c.service.GetAssignments() // Now returns dtos.TransporterDriverAssignmentResponse
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": assignments, "Total": total})
}

func (c *TransporterDriverAssignmentController) GetAssignment(ctx *gin.Context) {
	assignment, err := c.service.GetAssignment(ctx.Param("id")) // Now returns dtos.TransporterDriverAssignmentResponse
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Assignment not found"})
		return
	}
	ctx.JSON(http.StatusOK, assignment)
}

func (c *TransporterDriverAssignmentController) UpdateAssignment(ctx *gin.Context) {
	var req dtos.UpdateTransporterDriverAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAssignment(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Driver assignment updated successfully"})
}

func (c *TransporterDriverAssignmentController) DeleteAssignment(ctx *gin.Context) {
	if err := c.service.DeleteAssignment(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Assignment deleted successfully"})
}
