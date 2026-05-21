package controllers

import (
	"log"
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
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("Create Transporter: ShouldBind: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("Create Transporter: Validation Error: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	req.PassportPhoto, _ = ctx.FormFile("passport_photo")
	req.IDFrontPhoto, _ = ctx.FormFile("id_front_photo")
	req.IDBackPhoto, _ = ctx.FormFile("id_back_photo")
	req.CertificateOfIncorporation, _ = ctx.FormFile("certificate_of_incorporation")

	transporter, err := c.service.CreateTransporter(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transporter)
}

func (c *TransporterController) GetTransporters(ctx *gin.Context) {
	transporters, total, err := c.service.GetTransporters() // Now returns dtos.TransporterResponse
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": transporters, "total": total})
}

func (c *TransporterController) GetTransporter(ctx *gin.Context) {
	transporter, err := c.service.GetTransporter(ctx.Param("id")) // Now returns dtos.TransporterResponse
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter not found"})
		return
	}
	ctx.JSON(http.StatusOK, transporter)
}

func (c *TransporterController) UpdateTransporter(ctx *gin.Context) {
	var req dtos.UpdateTransporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateTransporter(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter updated successfully"})
}

func (c *TransporterController) DeleteTransporter(ctx *gin.Context) {
	if err := c.service.DeleteTransporter(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter deleted successfully"})
}
