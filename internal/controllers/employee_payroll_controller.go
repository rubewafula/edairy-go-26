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

type EmployeePayrollController struct {
	service *services.EmployeePayrollService
}

func NewEmployeePayrollController() *EmployeePayrollController {
	return &EmployeePayrollController{
		service: services.NewEmployeePayrollService(),
	}
}

func (c *EmployeePayrollController) CreateEmployeePayroll(ctx *gin.Context) {
	var req dtos.CreateEmployeePayrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	payroll, err := c.service.CreateEmployeePayroll(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, payroll)
}

func (c *EmployeePayrollController) GetEmployeePayrolls(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetEmployeePayrolls(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *EmployeePayrollController) GetEmployeePayroll(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetEmployeePayroll(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmployeePayrollController) UpdateEmployeePayroll(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeePayrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateEmployeePayroll(id, req, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee payroll updated successfully"})
}

func (c *EmployeePayrollController) DeleteEmployeePayroll(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteEmployeePayroll(id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee payroll deleted successfully"})
}
