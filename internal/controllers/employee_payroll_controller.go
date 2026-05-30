package controllers

import (
	"log"
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

// CreatePayroll godoc
// @Summary Initiate employee payroll generation
// @Description Triggers the background process to generate employee payrolls for all pending pay date ranges.
// @Tags Employee Payrolls
// @Accept json
// @Produce json
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls [post]
func (c *EmployeePayrollController) CreatePayroll(ctx *gin.Context) {
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

	if err := c.service.CreatePayroll(req, userID); err != nil {
		log.Printf("[EmployeePayrollController.CreatePayroll] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Employee payroll generation started in the background"})
}

// ConfirmPayroll godoc
// @Summary Confirm an employee payroll
// @Description Confirms a generated employee payroll, moving it from draft/incomplete to confirmed status.
// @Tags Employee Payrolls
// @Accept json
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 200 {object} models.EmployeePayroll
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id}/confirm [put]
func (c *EmployeePayrollController) ConfirmPayroll(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	payroll, err := c.service.ConfirmPayroll(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		log.Printf("[EmployeePayrollController.ConfirmPayroll] Error confirming payroll %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, payroll)
}

// ApprovePayroll godoc
// @Summary Approve an employee payroll
// @Description Approves a confirmed employee payroll, triggering GL postings in the background.
// @Tags Employee Payrolls
// @Accept json
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 202 {object} models.EmployeePayroll
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id}/approve [put]
func (c *EmployeePayrollController) ApprovePayroll(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	payroll, err := c.service.ApprovePayroll(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		log.Printf("[EmployeePayrollController.ApprovePayroll] Error approving payroll %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, payroll)
}

// GetEmployeePayrolls godoc
// @Summary Get all employee payrolls
// @Description Retrieve a list of all employee payrolls
// @Tags Employee Payrolls
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} dtos.EmployeePayrollResponse
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls [get]
func (c *EmployeePayrollController) GetEmployeePayrolls(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetEmployeePayrolls(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

// GetEmployeePayroll godoc
// @Summary Get a single employee payroll by ID
// @Description Retrieve an employee payroll by its ID
// @Tags Employee Payrolls
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 200 {object} dtos.EmployeePayrollResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id} [get]
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

// DeleteEmployeePayroll godoc
// @Summary Delete an employee payroll
// @Description Delete an employee payroll by its ID. Only draft or incomplete payrolls can be deleted.
// @Tags Employee Payrolls
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id} [delete]
func (c *EmployeePayrollController) DeleteEmployeePayroll(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteEmployeePayroll(utils.Uint64ToString(id), userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		log.Printf("[EmployeePayrollController.DeleteEmployeePayroll] Error deleting payroll %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee payroll deleted successfully"})
}
