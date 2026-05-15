package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type EmployeeBenefitController struct {
	service *services.EmployeeBenefitService
}

func NewEmployeeBenefitController() *EmployeeBenefitController {
	return &EmployeeBenefitController{
		service: services.NewEmployeeBenefitService(),
	}
}

func (c *EmployeeBenefitController) CreateEmployeeBenefit(ctx *gin.Context) {
	var req dtos.CreateEmployeeBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	benefit, err := c.service.CreateEmployeeBenefit(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, benefit)
}

func (c *EmployeeBenefitController) GetEmployeeBenefits(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetEmployeeBenefits(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *EmployeeBenefitController) GetEmployeeBenefit(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetEmployeeBenefit(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee benefit not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmployeeBenefitController) UpdateEmployeeBenefit(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeeBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateEmployeeBenefit(id, req, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee benefit not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee benefit updated successfully"})
}

func (c *EmployeeBenefitController) DeleteEmployeeBenefit(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteEmployeeBenefit(id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee benefit not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee benefit deleted successfully"})
}
