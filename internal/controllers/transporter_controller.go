package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type TransporterController struct {
	service *services.TransporterService
}

func NewTransporterController() *TransporterController {
	return &TransporterController{
		service: services.NewTransporterService(),
	}
}

func (c *TransporterController) CreateTransporter(ctx *gin.Context) {
	var req dtos.CreateTransporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	transporter, err := c.service.CreateTransporter(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Convert to Response DTO for PascalCase
	response, _ := c.service.GetTransporter(utils.Uint64ToString(transporter.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterController) GetTransporters(ctx *gin.Context) {
	transporters, total, err := c.service.GetTransporters() // Now returns dtos.TransporterResponse
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": transporters, "total": total})
}

func (c *TransporterController) GetTransporter(ctx *gin.Context) {
	transporter, err := c.service.GetTransporter(ctx.Param("id")) // Now returns dtos.TransporterResponse
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Transporter not found"})
		return
	}
	ctx.JSON(http.StatusOK, transporter)
}

func (c *TransporterController) UpdateTransporter(ctx *gin.Context) {
	var req dtos.UpdateTransporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateTransporter(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter updated successfully"})
}

func (c *TransporterController) DeleteTransporter(ctx *gin.Context) {
	if err := c.service.DeleteTransporter(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter deleted successfully"})
}
