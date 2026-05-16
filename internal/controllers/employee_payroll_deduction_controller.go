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

type EmployeePayrollDeductionController struct {
	service *services.EmployeePayrollDeductionService
}

func NewEmployeePayrollDeductionController() *EmployeePayrollDeductionController {
	return &EmployeePayrollDeductionController{
		service: services.NewEmployeePayrollDeductionService(),
	}
}

func (c *EmployeePayrollDeductionController) Create(ctx *gin.Context) {
	var req dtos.CreateEmployeePayrollDeductionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreatePayrollDeduction(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *EmployeePayrollDeductionController) List(ctx *gin.Context) {
	employeeID := ctx.Query("employee_id")
	payrollID := ctx.Query("payroll_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetPayrollDeductions(employeeID, payrollID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (c *EmployeePayrollDeductionController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetPayrollDeduction(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Payroll deduction not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *EmployeePayrollDeductionController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeePayrollDeductionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdatePayrollDeduction(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payroll deduction updated successfully"})
}

func (c *EmployeePayrollDeductionController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeletePayrollDeduction(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payroll deduction deleted successfully"})
}
