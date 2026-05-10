package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type TransporterBenefitController struct {
	service *services.TransporterBenefitService
}

func NewTransporterBenefitController() *TransporterBenefitController {
	return &TransporterBenefitController{
		service: services.NewTransporterBenefitService(),
	}
}

func (c *TransporterBenefitController) CreateBenefit(ctx *gin.Context) {
	var req dtos.CreateTransporterBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	benefit, err := c.service.CreateBenefit(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetBenefit(utils.Uint64ToString(benefit.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterBenefitController) GetBenefits(ctx *gin.Context) {
	benefits, total, err := c.service.GetBenefits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": benefits, "Total": total})
}

func (c *TransporterBenefitController) GetBenefit(ctx *gin.Context) {
	benefit, err := c.service.GetBenefit(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Transporter benefit not found"})
		return
	}
	ctx.JSON(http.StatusOK, benefit)
}

func (c *TransporterBenefitController) UpdateBenefit(ctx *gin.Context) {
	var req dtos.UpdateTransporterBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := c.service.UpdateBenefit(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter benefit updated successfully"})
}

func (c *TransporterBenefitController) DeleteBenefit(ctx *gin.Context) {
	if err := c.service.DeleteBenefit(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Transporter benefit deleted successfully"})
}
