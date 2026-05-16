package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CompanyTransporterController struct {
	service *services.CompanyTransporterService
}

func NewCompanyTransporterController() *CompanyTransporterController {
	return &CompanyTransporterController{
		service: services.NewCompanyTransporterService(),
	}
}

func (c *CompanyTransporterController) GetCompanyTransporters(ctx *gin.Context) {
	companies, total, err := c.service.GetCompanyTransporters()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": companies, "total": total})
}

func (c *CompanyTransporterController) GetCompanyTransporter(ctx *gin.Context) {
	company, err := c.service.GetCompanyTransporter(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Company transporter details not found"})
		return
	}
	ctx.JSON(http.StatusOK, company)
}

func (c *CompanyTransporterController) UpdateCompanyTransporter(ctx *gin.Context) {
	var req dtos.UpdateCompanyTransporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCompanyTransporter(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Company transporter details updated successfully"})
}
