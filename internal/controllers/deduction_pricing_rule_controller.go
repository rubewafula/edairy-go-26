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

type DeductionPricingRuleController struct {
	service *services.DeductionPricingRuleService
}

func NewDeductionPricingRuleController() *DeductionPricingRuleController {
	return &DeductionPricingRuleController{
		service: services.NewDeductionPricingRuleService(),
	}
}

func (c *DeductionPricingRuleController) CreateRule(ctx *gin.Context) {
	var req dtos.CreateDeductionPricingRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	rule, err := c.service.CreateRule(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response, _ := c.service.GetRule(utils.Uint64ToString(rule.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *DeductionPricingRuleController) GetRules(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	rules, total, err := c.service.GetRules(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": rules, "Total": total})
}

func (c *DeductionPricingRuleController) GetRule(ctx *gin.Context) {
	rule, err := c.service.GetRule(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Pricing rule not found"})
		return
	}
	ctx.JSON(http.StatusOK, rule)
}

func (c *DeductionPricingRuleController) UpdateRule(ctx *gin.Context) {
	var req dtos.UpdateDeductionPricingRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateRule(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Pricing rule updated successfully"})
}

func (c *DeductionPricingRuleController) DeleteRule(ctx *gin.Context) {
	if err := c.service.DeleteRule(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Pricing rule deleted successfully"})
}
