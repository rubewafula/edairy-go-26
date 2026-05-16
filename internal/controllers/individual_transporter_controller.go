package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type IndividualTransporterController struct {
	service *services.IndividualTransporterService
}

func NewIndividualTransporterController() *IndividualTransporterController {
	return &IndividualTransporterController{
		service: services.NewIndividualTransporterService(),
	}
}

func (c *IndividualTransporterController) GetIndividualTransporters(ctx *gin.Context) {
	individuals, total, err := c.service.GetIndividualTransporters()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": individuals, "total": total})
}

func (c *IndividualTransporterController) GetIndividualTransporter(ctx *gin.Context) {
	individual, err := c.service.GetIndividualTransporter(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Individual transporter details not found"})
		return
	}
	ctx.JSON(http.StatusOK, individual)
}

func (c *IndividualTransporterController) UpdateIndividualTransporter(ctx *gin.Context) {
	var req dtos.UpdateIndividualTransporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateIndividualTransporter(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Individual transporter details updated successfully"})
}
